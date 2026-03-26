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

import Editor from '@monaco-editor/react';
import {isValidStylesheetUrl, isInsecureStylesheetUrl} from '@thunder/shared-design';
import type {Stylesheet, UrlStylesheet} from '@thunder/shared-design';
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Box,
  Button,
  Chip,
  Dialog,
  DialogContent,
  DialogTitle,
  FormControl,
  FormLabel,
  IconButton,
  Stack,
  Switch,
  TextField,
  Tooltip,
  Typography,
  useColorScheme,
} from '@wso2/oxygen-ui';
import {Plus, Trash, ChevronUp, ChevronDown, Maximize, Edit, X} from '@wso2/oxygen-ui-icons-react';
import {forwardRef, useCallback, useEffect, useImperativeHandle, useRef, useState, type JSX} from 'react';

// Shared DOM node for Monaco overflow widgets (context menu, suggest, etc.)
// Appended to <body> so they are never clipped by parent overflow.
let sharedOverflowNode: HTMLDivElement | null = null;
function getOverflowWidgetsDomNode(): HTMLDivElement {
  if (!sharedOverflowNode) {
    sharedOverflowNode = document.createElement('div');
    sharedOverflowNode.className = 'monaco-editor';
    sharedOverflowNode.style.zIndex = '9999';
    document.body.appendChild(sharedOverflowNode);
  }
  return sharedOverflowNode;
}

// ── Inline CSS Monaco editor with debounce ──────────────────────────────────

interface InlineCSSFieldProps {
  id: string;
  content: string;
  colorMode: 'light' | 'dark';
  onChange: (content: string) => void;
  /** Called on mount/update so the parent can flush this field's pending debounce. */
  registerFlush: (flush: (() => void) | null) => void;
}

const EDITOR_OPTIONS = {
  minimap: {enabled: false},
  scrollBeyondLastLine: false,
  automaticLayout: true,
  fontSize: 12,
  lineHeight: 18,
  tabSize: 2,
  wordWrap: 'on' as const,
  folding: false,
  fixedOverflowWidgets: true,
  overflowWidgetsDomNode: getOverflowWidgetsDomNode(),
  scrollbar: {verticalScrollbarSize: 2, horizontalScrollbarSize: 2, useShadows: false},
  overviewRulerLanes: 0,
  hideCursorInOverviewRuler: true,
  overviewRulerBorder: false,
  renderLineHighlight: 'none' as const,
  padding: {top: 4, bottom: 4},
  quickSuggestions: true,
  suggestOnTriggerCharacters: true,
  suggest: {showProperties: true, showValues: true, showColors: true, showKeywords: true},
};

function InlineCSSField({id, content, colorMode, onChange, registerFlush}: InlineCSSFieldProps): JSX.Element {
  const [localContent, setLocalContent] = useState(content);
  const [expanded, setExpanded] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const localContentRef = useRef(localContent);
  const onChangeRef = useRef(onChange);

  useEffect(() => {
    localContentRef.current = localContent;
    onChangeRef.current = onChange;
  });

  const [prevContent, setPrevContent] = useState(content);
  if (prevContent !== content) {
    setPrevContent(content);
    setLocalContent(content);
  }

  // Register a flush callback so the parent can synchronously commit pending edits.
  const flush = useCallback(() => {
    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
      debounceRef.current = null;
      onChangeRef.current(localContentRef.current);
    }
  }, []);

  useEffect(() => {
    registerFlush(flush);
    return () => registerFlush(null);
  }, [registerFlush, flush]);

  // Cleanup debounce on unmount
  useEffect(
    () => () => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
    },
    [],
  );

  const handleEditorChange = (raw: string | undefined): void => {
    const text = raw ?? '';
    setLocalContent(text);
    if (debounceRef.current) clearTimeout(debounceRef.current);
    debounceRef.current = setTimeout(() => onChange(text), 400);
  };

  return (
    <>
      <Box sx={{border: '1px solid', borderColor: 'divider', borderRadius: 1}}>
        <Editor
          height="220px"
          language="css"
          theme={colorMode === 'dark' ? 'vs-dark' : 'vs'}
          value={localContent}
          onChange={handleEditorChange}
          options={{
            ...EDITOR_OPTIONS,
            lineNumbers: 'on',
            glyphMargin: false,
            lineDecorationsWidth: 4,
            lineNumbersMinChars: 3,
          }}
        />
        <Stack
          direction="row"
          alignItems="center"
          justifyContent="flex-end"
          sx={{
            borderTop: '1px solid',
            borderColor: 'divider',
            px: 0.5,
            py: 0.25,
            bgcolor: colorMode === 'dark' ? 'grey.900' : 'grey.50',
          }}
        >
          <Tooltip title="Open in full editor">
            <IconButton size="small" onClick={() => setExpanded(true)} sx={{p: 0.25}}>
              <Maximize size={14} />
            </IconButton>
          </Tooltip>
        </Stack>
      </Box>

      <Dialog
        open={expanded}
        onClose={() => setExpanded(false)}
        maxWidth="md"
        fullWidth
        slotProps={{paper: {sx: {height: '80vh'}}}}
      >
        <DialogTitle sx={{display: 'flex', alignItems: 'center', py: 1, px: 2}}>
          <Typography component="span" variant="subtitle2" sx={{flex: 1, fontFamily: 'monospace'}}>
            {id}
          </Typography>
          <IconButton size="small" onClick={() => setExpanded(false)}>
            <X size={16} />
          </IconButton>
        </DialogTitle>
        <DialogContent sx={{p: 0}}>
          <Editor
            height="100%"
            language="css"
            theme={colorMode === 'dark' ? 'vs-dark' : 'vs'}
            value={localContent}
            onChange={handleEditorChange}
            options={{...EDITOR_OPTIONS, lineNumbers: 'on', fixedOverflowWidgets: false, overflowWidgetsDomNode: undefined}}
          />
        </DialogContent>
      </Dialog>
    </>
  );
}

// ── URL stylesheet field ────────────────────────────────────────────────────

interface UrlFieldProps {
  sheet: UrlStylesheet;
  onUpdate: (patch: Partial<UrlStylesheet>) => void;
}

function UrlField({sheet, onUpdate}: UrlFieldProps): JSX.Element {
  const hasError = !!sheet.href && !isValidStylesheetUrl(sheet.href);
  const isInsecure = Boolean(sheet.href) && !hasError && isInsecureStylesheetUrl(sheet.href);

  let helperText: string | undefined;
  if (hasError) helperText = 'URL must be a valid http:// or https:// address';
  else if (isInsecure) helperText = 'Using HTTP is insecure. Consider using HTTPS instead.';

  return (
    <FormControl fullWidth>
      <FormLabel>URL</FormLabel>
      <TextField
        size="small"
        value={sheet.href}
        onChange={(e) => onUpdate({href: e.target.value})}
        fullWidth
        error={hasError}
        color={isInsecure ? 'warning' : undefined}
        focused={isInsecure ?? undefined}
        helperText={helperText}
        slotProps={{
          input: {sx: {fontSize: '0.8rem', fontFamily: 'monospace'}},
          formHelperText: isInsecure ? {sx: {color: 'warning.main'}} : undefined,
        }}
      />
    </FormControl>
  );
}

// ── Editable title ──────────────────────────────────────────────────────────

interface EditableTitleProps {
  value: string;
  onChange: (value: string) => void;
}

function EditableTitle({value, onChange}: EditableTitleProps): JSX.Element {
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState(value);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    setDraft(value);
  }, [value]);

  useEffect(() => {
    if (editing) {
      inputRef.current?.focus();
      inputRef.current?.select();
    }
  }, [editing]);

  const commit = (): void => {
    const trimmed = draft.trim();
    if (trimmed && trimmed !== value) onChange(trimmed);
    else setDraft(value);
    setEditing(false);
  };

  if (editing) {
    return (
      <TextField
        inputRef={inputRef}
        size="small"
        value={draft}
        onChange={(e) => setDraft(e.target.value)}
        onBlur={commit}
        onKeyDown={(e) => {
          if (e.key === 'Enter') commit();
          if (e.key === 'Escape') {
            setDraft(value);
            setEditing(false);
          }
          e.stopPropagation();
        }}
        onClick={(e) => e.stopPropagation()}
        variant="standard"
        slotProps={{input: {sx: {fontSize: '0.8rem', fontWeight: 600, fontFamily: 'monospace', py: 0}, disableUnderline: false}}}
        sx={{maxWidth: 140}}
      />
    );
  }

  return (
    <Stack direction="row" alignItems="center" gap={0.25} sx={{overflow: 'hidden', minWidth: 0}}>
      <Typography
        variant="body2"
        sx={{
          fontWeight: 600,
          fontSize: '0.8rem',
          fontFamily: 'monospace',
          overflow: 'hidden',
          textOverflow: 'ellipsis',
          whiteSpace: 'nowrap',
        }}
      >
        {value}
      </Typography>
      <Box
        component="span"
        role="button"
        tabIndex={0}
        onClick={(e) => {
          e.stopPropagation();
          setEditing(true);
        }}
        onKeyDown={(e) => {
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            e.stopPropagation();
            setEditing(true);
          }
        }}
        sx={{
          display: 'flex',
          flexShrink: 0,
          cursor: 'pointer',
          color: 'text.disabled',
          borderRadius: 0.5,
          p: 0.25,
          '&:hover': {color: 'text.secondary'},
        }}
      >
        <Edit size={12} />
      </Box>
    </Stack>
  );
}

// ── Individual stylesheet item ──────────────────────────────────────────────

interface StylesheetItemProps {
  sheet: Stylesheet;
  idx: number;
  total: number;
  expanded: boolean;
  colorMode: 'light' | 'dark';
  onToggle: () => void;
  onRemove: () => void;
  onMove: (direction: -1 | 1) => void;
  onUpdate: (patch: Partial<Stylesheet>) => void;
  registerFlush: (flush: (() => void) | null) => void;
}

function StylesheetItem({
  sheet,
  idx,
  total,
  expanded,
  colorMode,
  onToggle,
  onRemove,
  onMove,
  onUpdate,
  registerFlush,
}: StylesheetItemProps): JSX.Element {
  const isInline = sheet.type === 'inline';
  const isDisabled = !!sheet.disabled;

  return (
    <Accordion
      expanded={expanded}
      onChange={onToggle}
      disableGutters
      square
      sx={{
        backgroundColor: 'transparent',
        '&:before': {display: 'none'},
        overflow: 'visible',
        opacity: isDisabled ? 0.5 : 1,
        transition: 'opacity 0.15s ease',
      }}
    >
      <AccordionSummary
        expandIcon={<ChevronDown size={16} />}
        sx={{
          '& .MuiAccordionSummary-content': {alignItems: 'center', gap: 0.5, overflow: 'hidden', minWidth: 0},
          minHeight: 40,
          '&.Mui-expanded': {minHeight: 40},
        }}
      >
        {/* Reorder arrows */}
        <Stack
          component="span"
          onClick={(e) => e.stopPropagation()}
          sx={{flexShrink: 0, mr: 0.25}}
        >
          <Box
            component="span"
            role="button"
            tabIndex={0}
            onClick={() => { if (idx > 0) onMove(-1); }}
            onKeyDown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); if (idx > 0) onMove(-1); } }}
            sx={{display: 'flex', cursor: idx === 0 ? 'default' : 'pointer', opacity: idx === 0 ? 0.25 : 0.5, '&:hover': {opacity: idx === 0 ? 0.25 : 1}, lineHeight: 0}}
          >
            <ChevronUp size={12} />
          </Box>
          <Box
            component="span"
            role="button"
            tabIndex={0}
            onClick={() => { if (idx < total - 1) onMove(1); }}
            onKeyDown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); if (idx < total - 1) onMove(1); } }}
            sx={{display: 'flex', cursor: idx === total - 1 ? 'default' : 'pointer', opacity: idx === total - 1 ? 0.25 : 0.5, '&:hover': {opacity: idx === total - 1 ? 0.25 : 1}, lineHeight: 0}}
          >
            <ChevronDown size={12} />
          </Box>
        </Stack>

        <EditableTitle value={sheet.id} onChange={(id) => onUpdate({id})} />

        <Chip
          label={isInline ? 'Inline' : 'URL'}
          size="small"
          variant="outlined"
          sx={{fontSize: '0.6rem', height: 18, flexShrink: 0, '& .MuiChip-label': {px: 0.75}}}
        />

        <Box sx={{flex: 1}} />

        {/* Enable/disable toggle */}
        <Box component="span" onClick={(e) => e.stopPropagation()} sx={{display: 'flex', alignItems: 'center'}}>
          <Switch
            size="small"
            checked={!isDisabled}
            onChange={(e) => onUpdate({disabled: !e.target.checked})}
            sx={{transform: 'scale(0.7)'}}
          />
        </Box>

        {/* Delete */}
        <Box
          component="span"
          role="button"
          tabIndex={0}
          onClick={(e) => {
            e.stopPropagation();
            onRemove();
          }}
          onKeyDown={(e) => {
            if (e.key === 'Enter' || e.key === ' ') {
              e.preventDefault();
              e.stopPropagation();
              onRemove();
            }
          }}
          sx={{
            display: 'flex',
            p: 0.25,
            cursor: 'pointer',
            color: 'text.disabled',
            borderRadius: 0.5,
            '&:hover': {color: 'error.main'},
          }}
        >
          <Trash size={14} />
        </Box>
      </AccordionSummary>

      <AccordionDetails sx={{display: 'flex', flexDirection: 'column', gap: 1, overflow: 'visible'}}>
        {isInline ? (
          <InlineCSSField
            id={sheet.id}
            content={sheet.content}
            colorMode={colorMode}
            onChange={(content) => onUpdate({content})}
            registerFlush={registerFlush}
          />
        ) : (
          <UrlField sheet={sheet} onUpdate={onUpdate} />
        )}
      </AccordionDetails>
    </Accordion>
  );
}

// ── Main component ──────────────────────────────────────────────────────────

export interface CustomCSSEditorHandle {
  /** Flush all pending debounced edits synchronously. */
  flush: () => void;
}

export interface CustomCSSEditorProps {
  stylesheets: Stylesheet[];
  onChange: (stylesheets: Stylesheet[]) => void;
}

/** Generates a short unique ID like "custom-1", "custom-2", etc. */
function nextId(stylesheets: Stylesheet[]): string {
  const existing = new Set(stylesheets.map((s) => s.id));
  let n = stylesheets.length + 1;
  while (existing.has(`custom-${n}`)) n += 1;
  return `custom-${n}`;
}

const CustomCSSEditor = forwardRef<CustomCSSEditorHandle, CustomCSSEditorProps>(
  ({stylesheets, onChange}, ref): JSX.Element => {
    const {mode, systemMode} = useColorScheme();
    const colorMode: 'light' | 'dark' = (mode === 'system' ? systemMode : mode) === 'dark' ? 'dark' : 'light';

    const [expandedIdx, setExpandedIdx] = useState<number | null>(null);

    // Track flush callbacks from InlineCSSField instances
    const flushMapRef = useRef<Map<number, () => void>>(new Map());

    useImperativeHandle(ref, () => ({
      flush: () => {
        flushMapRef.current.forEach((fn) => fn());
      },
    }));

    // Stable React keys — not tied to the editable `id` field.
    const [keyCounter, setKeyCounter] = useState(stylesheets.length);
    const nextKeyRef = useRef(keyCounter);
    const nextKey = (): number => {
      nextKeyRef.current += 1;
      setKeyCounter(nextKeyRef.current);
      return nextKeyRef.current;
    };
    const [stableKeys, setStableKeys] = useState<number[]>(() =>
      Array.from({length: stylesheets.length}, (_, i) => i + 1),
    );

    // Sync stable keys when stylesheets are replaced externally (e.g. server load).
    const [prevLength, setPrevLength] = useState(stylesheets.length);
    if (stylesheets.length !== prevLength && stableKeys.length !== stylesheets.length) {
      setPrevLength(stylesheets.length);
      const newKeys = Array.from({length: stylesheets.length}, (_, i) => keyCounter + i + 1);
      setStableKeys(newKeys);
      setKeyCounter(keyCounter + stylesheets.length);
    }

    const handleAdd = (type: 'inline' | 'url'): void => {
      const id = nextId(stylesheets);
      const newSheet: Stylesheet = type === 'inline' ? {id, type: 'inline', content: ''} : {id, type: 'url', href: ''};
      const updated = [...stylesheets, newSheet];
      onChange(updated);
      setStableKeys((prev) => [...prev, nextKey()]);
      setExpandedIdx(updated.length - 1);
    };

    const handleRemove = (idx: number): void => {
      const updated = stylesheets.filter((_, i) => i !== idx);
      onChange(updated);
      setStableKeys((prev) => prev.filter((_, i) => i !== idx));
      if (expandedIdx === idx) setExpandedIdx(null);
      else if (expandedIdx !== null && expandedIdx > idx) setExpandedIdx(expandedIdx - 1);
    };

    const handleMove = (idx: number, direction: -1 | 1): void => {
      const target = idx + direction;
      if (target < 0 || target >= stylesheets.length) return;

      // Collapse so Monaco editors are not visible during reorder.
      setExpandedIdx(null);

      const updated = [...stylesheets];
      [updated[idx], updated[target]] = [updated[target], updated[idx]];
      onChange(updated);

      // Generate fresh keys so React fully unmounts/remounts items
      // instead of trying to move DOM nodes (which corrupts Monaco).
      setStableKeys(updated.map(() => nextKey()));
    };

    const handleUpdate = (idx: number, patch: Partial<Stylesheet>): void => {
      const updated = stylesheets.map((s, i) => (i === idx ? {...s, ...patch} : s));
      onChange(updated as Stylesheet[]);
    };

    return (
      <Stack gap={1}>
        {stylesheets.length === 0 && (
          <Box sx={{py: 3, display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 1}}>
            <Typography variant="body2" color="text.disabled" sx={{fontSize: '0.8rem', textAlign: 'center'}}>
              No custom stylesheets yet.
            </Typography>
            <Typography variant="caption" color="text.disabled" sx={{fontSize: '0.7rem', textAlign: 'center'}}>
              Add an inline stylesheet or link an external CSS file to customize the appearance.
            </Typography>
          </Box>
        )}

        {stylesheets.map((sheet, idx) => (
          <StylesheetItem
            key={stableKeys[idx]}
            sheet={sheet}
            idx={idx}
            total={stylesheets.length}
            expanded={expandedIdx === idx}
            colorMode={colorMode}
            onToggle={() => setExpandedIdx(expandedIdx === idx ? null : idx)}
            onRemove={() => handleRemove(idx)}
            onMove={(dir) => handleMove(idx, dir)}
            onUpdate={(patch) => handleUpdate(idx, patch)}
            registerFlush={(flush) => {
              const key = stableKeys[idx];
              if (flush) flushMapRef.current.set(key, flush);
              else flushMapRef.current.delete(key);
            }}
          />
        ))}

        <Stack direction="row" spacing={0.75} sx={{mt: 0.25}}>
          <Button
            size="small"
            variant="text"
            startIcon={<Plus size={12} />}
            onClick={() => handleAdd('inline')}
            sx={{textTransform: 'none', fontSize: '0.7rem', color: 'text.secondary', px: 1, minWidth: 0}}
          >
            Inline
          </Button>
          <Button
            size="small"
            variant="text"
            startIcon={<Plus size={12} />}
            onClick={() => handleAdd('url')}
            sx={{textTransform: 'none', fontSize: '0.7rem', color: 'text.secondary', px: 1, minWidth: 0}}
          >
            External URL
          </Button>
        </Stack>
      </Stack>
    );
  },
);

CustomCSSEditor.displayName = 'CustomCSSEditor';

export default CustomCSSEditor;
