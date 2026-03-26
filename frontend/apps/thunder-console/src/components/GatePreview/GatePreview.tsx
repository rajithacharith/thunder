/**
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import type {EmbeddedFlowComponent} from '@asgardeo/react';
import createCache from '@emotion/cache';
import {CacheProvider} from '@emotion/react';
import {FlowComponentRenderer, DesignProvider, useDesign, AuthPageLayout, AuthCardLayout} from '@thunder/shared-design';
import type {Theme, ColorSchemeOption, DesignResolveResponse, Stylesheet} from '@thunder/shared-design';
import {useTemplateLiteralResolver} from '@thunder/shared-hooks';
import {TemplateLiteralType} from '@thunder/utils';
import {AcrylicOrangeTheme, Box, CircularProgress, ThemeProvider, Typography, useColorScheme} from '@wso2/oxygen-ui';
import {useCallback, useEffect, useLayoutEffect, useMemo, useRef, useState, type JSX, type ReactNode} from 'react';
import {createPortal} from 'react-dom';
import {useTranslation} from 'react-i18next';
import buildPreviewMock from './mocks/buildPreviewMock';
import ElementInspector from '../../features/design/components/layouts/ElementInspector';
import PreviewToolbar from '../../features/design/components/PreviewToolbar';
import {VIEWPORT_WIDTHS, VIEWPORT_HEIGHTS} from '../../features/design/components/viewportConstants';

// ── Constants ────────────────────────────────────────────────────────────────

const ZOOM_STEPS = [25, 50, 75, 100, 125, 150];

/** Minimum width (px) the content needs so the 450px sign-in card + padding renders without clipping. */
const MIN_CONTENT_WIDTH = 520;

/** Minimum height (px) the content needs so a typical sign-in form renders without clipping. */
const MIN_CONTENT_HEIGHT = 700;

/** No-op handlers for preview mode — the form is purely visual. */
const noopSubmit = (): void => { /* no-op */ };
const noopInputChange = (): void => { /* no-op */ };


/**
 * Initial HTML written into the preview iframe. Sets up the full height chain
 * so AuthPageLayout's minHeight: 100% resolves correctly.
 */
const IFRAME_INITIAL_HTML = [
  '<!DOCTYPE html><html style="height:100%"><head>',
  '<link rel="preconnect" href="https://fonts.googleapis.com">',
  '<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>',
  '<link href="https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&display=swap" rel="stylesheet">',
  '<style>body{margin:0;height:100%;font-family:\'Inter\',sans-serif}#root,#root>*{height:100%}</style>',
  '</head><body><div id="root"></div></body></html>',
].join('');

// ── Helper components ────────────────────────────────────────────────────────

/** Syncs the nested ThemeProvider's mode with the preview's colorScheme prop. */
function ColorSchemeSync({mode}: {mode: 'light' | 'dark'}): null {
  const {setMode} = useColorScheme();
  useEffect(() => {
    setMode(mode);
  }, [mode, setMode]);
  return null;
}

/**
 * Wraps children in a ThemeProvider scoped to the preview iframe.
 *
 * Key: `colorSchemeNode` must point to the iframe's `<html>` element so MUI
 * sets `data-color-scheme` inside the iframe (not on the parent document).
 * Without this, the CSS-vars selectors like `[data-color-scheme="dark"]`
 * never match and the theme doesn't switch.
 */
function PreviewThemeProvider({
  colorScheme,
  colorSchemeNode = undefined,
  children,
}: {
  colorScheme: 'light' | 'dark';
  colorSchemeNode?: HTMLElement | null;
  children: ReactNode;
}): JSX.Element {
  const {theme} = useDesign(AcrylicOrangeTheme as Theme);

  // MUI's ThemeProvider supports CSS-vars-specific props (colorSchemeNode,
  // disableNestedContext, storageManager) at runtime, but the TypeScript types
  // only expose them when the module-augmentation `CssThemeVariables` is set to
  // `{ enabled: true }`.  We need these props to isolate the preview iframe's
  // color-scheme attribute and prevent localStorage conflicts.
  const cssVarsProps = {
    storageManager: null,
    disableNestedContext: true,
    ...(colorSchemeNode ? {colorSchemeNode} : {}),
  } as Record<string, unknown>;

  return (
    <ThemeProvider
      theme={theme ?? AcrylicOrangeTheme}
      defaultMode={colorScheme}
      {...cssVarsProps}
    >
      <ColorSchemeSync mode={colorScheme} />
      {children}
    </ThemeProvider>
  );
}

/**
 * Wraps preview content with an emotion CacheProvider that injects MUI styles
 * into the iframe's <head>, and injects custom stylesheets into the iframe document.
 */
function IframeContent({
  iframeDoc,
  colorScheme,
  theme,
  stylesheets,
  pageBackground,
  mock,
  inspectorEnabled,
  onSelectSelector = undefined,
}: {
  iframeDoc: Document;
  colorScheme: 'light' | 'dark';
  theme: Theme | undefined;
  stylesheets: Stylesheet[];
  pageBackground: string | undefined;
  mock: EmbeddedFlowComponent[];
  inspectorEnabled: boolean;
  onSelectSelector?: (selector: string) => void;
}): JSX.Element {
  // Resolve {{t(key)}} templates using the app's i18n catalog.
  const {t} = useTranslation();
  const {resolveAll} = useTemplateLiteralResolver();
  const previewResolve = useMemo(
    () => (template: string | undefined): string | undefined =>
      resolveAll(template, {[TemplateLiteralType.TRANSLATION]: t}),
    [resolveAll, t],
  );

  // Create an emotion cache that injects styles into the iframe's <head>.
  const cache = useMemo(() => createCache({key: 'preview', container: iframeDoc.head}), [iframeDoc]);

  // Inject custom stylesheets into the iframe document (not the parent).
  const serializedSheets = JSON.stringify(stylesheets);
  useEffect(() => {
    const parsed: Stylesheet[] = JSON.parse(serializedSheets) as Stylesheet[];
    const injectedIds: string[] = [];

    parsed.forEach((sheet) => {
      const elementId = `thunder-preview-${sheet.id}`;
      iframeDoc.getElementById(elementId)?.remove();

      if (sheet.disabled) return;

      if (sheet.type === 'inline') {
        const style = iframeDoc.createElement('style');
        style.id = elementId;
        style.textContent = sheet.content;
        iframeDoc.head.appendChild(style);
        injectedIds.push(elementId);
      } else if (sheet.type === 'url') {
        const link = iframeDoc.createElement('link');
        link.id = elementId;
        link.rel = 'stylesheet';
        link.href = sheet.href;
        iframeDoc.head.appendChild(link);
        injectedIds.push(elementId);
      }
    });

    return () => {
      injectedIds.forEach((id) => iframeDoc.getElementById(id)?.remove());
    };
  }, [iframeDoc, serializedSheets]);

  return (
    <CacheProvider value={cache}>
      <DesignProvider
        shouldResolveDesignInternally={false}
        design={theme ? ({theme} as DesignResolveResponse) : undefined}
      >
        <PreviewThemeProvider colorScheme={colorScheme} colorSchemeNode={iframeDoc.documentElement}>
          <ElementInspector enabled={inspectorEnabled} onSelectSelector={onSelectSelector}>
            <AuthPageLayout variant="SignIn" background={pageBackground}>
              <AuthCardLayout variant="SignInBox" showLogo={false}>
                <Box sx={{display: 'flex', flexDirection: 'column', gap: 2}}>
                  {mock.map((component, index) => (
                    <FlowComponentRenderer
                      key={component.id ?? index}
                      component={component}
                      index={index}
                      values={{}}
                      isLoading={false}
                      resolve={previewResolve}
                      onInputChange={noopInputChange}
                      onSubmit={noopSubmit}
                    />
                  ))}
                </Box>
              </AuthCardLayout>
            </AuthPageLayout>
          </ElementInspector>
        </PreviewThemeProvider>
      </DesignProvider>
    </CacheProvider>
  );
}

// ── Types & Props ────────────────────────────────────────────────────────────

export type Viewport = 'desktop' | 'tablet' | 'mobile';

export interface GatePreviewProps {
  /** The theme to render. Null shows a loading spinner; undefined shows an empty prompt. */
  theme: Theme | null | undefined;
  displayName?: string;
  showToolbar?: boolean;
  viewport?: {
    width: string | number;
    height: string | number;
  };
  colorScheme?: ColorSchemeOption;
  /** When true, the preview tracks the host app's color scheme instead of the toolbar toggle. */
  syncColorSchemeWithSystem?: boolean;
  mock?: EmbeddedFlowComponent[];
  /** Optional page background CSS value (color, gradient, or image). Overrides theme background when set. */
  pageBackground?: string;
  /** Custom stylesheets to inject into the isolated preview iframe. */
  stylesheets?: Stylesheet[];
  /** When true, enables the element inspector overlay inside the preview. */
  inspectorEnabled?: boolean;
  /** Callback when a CSS selector is picked via the inspector. */
  onSelectSelector?: (selector: string) => void;
}

// ── Main component ───────────────────────────────────────────────────────────

export default function GatePreview({
  theme,
  displayName = '',
  showToolbar = true,
  viewport = undefined,
  mock = buildPreviewMock(),
  colorScheme = undefined,
  syncColorSchemeWithSystem = false,
  pageBackground = undefined,
  stylesheets = [],
  inspectorEnabled = false,
  onSelectSelector = undefined,
}: GatePreviewProps): JSX.Element {
  const {mode, systemMode} = useColorScheme();
  const [previewColorScheme, setPreviewColorScheme] = useState<'light' | 'dark' | 'system'>('light');
  const [viewportState, setViewport] = useState<Viewport>('desktop');
  const [zoom, setZoom] = useState(75);
  const canvasRef = useRef<HTMLDivElement>(null);
  const iframeRef = useRef<HTMLIFrameElement | null>(null);
  const dimensionsRef = useRef<HTMLSpanElement>(null);
  const [iframeDoc, setIframeDoc] = useState<Document | null>(null);

  // Callback ref: initializes the iframe document whenever the <iframe> mounts.
  // This handles the case where theme starts as null (loading spinner), so the
  // iframe doesn't exist on first render — the callback fires when it appears.
  // We skip re-initialization if #root already exists (React Strict Mode calls
  // the callback ref twice; re-writing the document would destroy the portal
  // target without triggering a re-render since the doc reference is the same).
  const iframeCallbackRef = useCallback((iframe: HTMLIFrameElement | null) => {
    iframeRef.current = iframe;
    if (!iframe) return;
    const doc = iframe.contentDocument;
    if (!doc) return;
    if (doc.getElementById('root')) {
      setIframeDoc(doc);
      return;
    }
    doc.open();
    doc.write(IFRAME_INITIAL_HTML);
    doc.close();
    setIframeDoc(doc);
  }, []);

  const resolvedSystemMode: 'light' | 'dark' = (mode === 'system' ? systemMode : mode) === 'dark' ? 'dark' : 'light';
  const activeScheme = colorScheme !== 'system' ? colorScheme : undefined;
  let effectiveScheme: 'light' | 'dark';
  if (activeScheme) {
    effectiveScheme = activeScheme;
  } else if (syncColorSchemeWithSystem) {
    effectiveScheme = resolvedSystemMode;
  } else if (previewColorScheme !== 'system') {
    effectiveScheme = previewColorScheme;
  } else {
    effectiveScheme = resolvedSystemMode;
  }

  const zoomIdx = ZOOM_STEPS.indexOf(zoom);

  // Imperatively size & scale the iframe to fit the canvas — no React state, no re-renders.
  useLayoutEffect(() => {
    const canvas = canvasRef.current;
    const iframe = iframeRef.current;
    if (!canvas || !iframe) return undefined;

    const update = (): void => {
      const cw = canvas.clientWidth;
      const ch = canvas.clientHeight;
      if (!cw || !ch) return;

      const userScale = zoom / 100;
      // Scale down to fit both dimensions so the card never clips.
      const fitScaleW = Math.min(1, cw / MIN_CONTENT_WIDTH);
      const fitScaleH = Math.min(1, ch / MIN_CONTENT_HEIGHT);
      const fitScale = Math.min(fitScaleW, fitScaleH);
      const totalScale = fitScale * userScale;

      // Inverse-scale: render iframe at (canvas / totalScale) so after
      // transform: scale(totalScale) it visually fills the canvas exactly.
      const iframeW = Math.round(cw / totalScale);
      const iframeH = Math.round(ch / totalScale);
      iframe.style.width = `${iframeW}px`;
      iframe.style.height = `${iframeH}px`;
      iframe.style.transform = `scale(${totalScale})`;

      // Update dimensions label without triggering a React re-render.
      if (dimensionsRef.current) {
        dimensionsRef.current.textContent = `${iframeW} × ${iframeH}`;
      }
    };

    const observer = new ResizeObserver(update);
    observer.observe(canvas);
    update();

    return () => observer.disconnect();
  }, [zoom, iframeDoc]);

  if (theme === null) {
    return (
      <Box sx={{height: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center'}}>
        <CircularProgress size={32} />
      </Box>
    );
  }

  return (
    <Box sx={{height: '100%', display: 'flex', flexDirection: 'column'}}>
      {/* Toolbar */}
      {showToolbar && (
        <Box sx={{display: 'flex', justifyContent: 'center', py: 1.5, flexShrink: 0}}>
          <PreviewToolbar
            viewport={viewportState}
            setViewport={setViewport}
            previewColorScheme={previewColorScheme}
            setPreviewColorScheme={setPreviewColorScheme}
            zoom={zoom}
            setZoom={setZoom}
            zoomIdx={zoomIdx}
          />
        </Box>
      )}

      {/* Viewport container */}
      <Box
        sx={{flex: 1, overflow: 'hidden', display: 'flex', justifyContent: 'center', alignItems: 'flex-start', p: 2}}
      >
        <Box
          sx={{
            backgroundColor: 'background.paper',
            borderRadius: 1,
            width: viewport?.width ?? VIEWPORT_WIDTHS[viewportState],
            height: viewport?.height ?? VIEWPORT_HEIGHTS[viewportState],
            transition: 'width 0.2s ease, height 0.2s ease',
            display: 'flex',
            flexDirection: 'column',
          }}
        >
          {/* Browser chrome */}
          <Box
            sx={{
              px: 3,
              py: 1.5,
              borderBottom: '1px solid',
              borderColor: 'divider',
              display: 'flex',
              alignItems: 'center',
              gap: 1,
              flexShrink: 0,
            }}
          >
            <Box sx={{width: 8, height: 8, borderRadius: '50%', bgcolor: '#fc5c57'}} />
            <Box sx={{width: 8, height: 8, borderRadius: '50%', bgcolor: '#febc2e'}} />
            <Box sx={{width: 8, height: 8, borderRadius: '50%', bgcolor: '#29c840'}} />
            <Box
              sx={{
                flex: 1,
                mx: 2,
                height: 22,
                bgcolor: 'action.hover',
                borderRadius: 1,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              <Typography variant="caption" color="text.disabled" sx={{fontSize: 10}}>
                {displayName ? `${displayName} — Preview` : 'Preview'}
              </Typography>
            </Box>
          </Box>

          {/* Canvas — fills the browser chrome frame like a real viewport */}
          <Box
            ref={canvasRef}
            sx={{
              flex: 1,
              overflow: 'hidden',
              position: 'relative',
            }}
          >
            <Typography
              component="span"
              ref={dimensionsRef}
              variant="caption"
              sx={{
                position: 'absolute',
                top: 4,
                right: 6,
                zIndex: 1,
                fontSize: 9,
                fontFamily: 'monospace',
                color: 'text.disabled',
                opacity: 0.7,
                pointerEvents: 'none',
              }}
            />
            <iframe ref={iframeCallbackRef} title="Gate Preview" style={{border: 'none', transformOrigin: 'top left', position: 'absolute', top: 0, left: 0}} />
            {iframeDoc?.getElementById('root') &&
              createPortal(
                <IframeContent
                  iframeDoc={iframeDoc}
                  colorScheme={effectiveScheme}
                  theme={theme}
                  stylesheets={stylesheets}
                  pageBackground={pageBackground}
                  mock={mock}
                  inspectorEnabled={inspectorEnabled}
                  onSelectSelector={onSelectSelector}
                />,
                iframeDoc.getElementById('root')!,
              )}
          </Box>
        </Box>
      </Box>
    </Box>
  );
}
