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

import {useConfig} from '@thunderid/contexts';
import {Box, Button, CircularProgress, Stack, Typography} from '@wso2/oxygen-ui';
import {CheckCircle, FileCode, FileText, FolderOpen, RefreshCw, XCircle} from '@wso2/oxygen-ui-icons-react';
import type {JSX} from 'react';
import {useState} from 'react';
import {useTranslation} from 'react-i18next';
import useImportConfiguration from '../../import-export/api/useImportConfiguration';
import type {ImportResponse} from '../../import-export/models/import-configuration';

const IMPORTED_KEY = 'thunderid-wayfinder-config-imported';

type Status = 'idle' | 'selected' | 'importing' | 'success' | 'alreadyDone' | 'error';

interface FoundFiles {
  folderName: string;
  yamlFile: File;
  envFile: File | null;
}

function parseEnvFile(content: string): Record<string, string> {
  const vars: Record<string, string> = {};
  for (const line of content.split(/\r?\n/)) {
    const trimmed = line.trim();
    if (!trimmed || trimmed.startsWith('#')) continue;
    const eq = trimmed.indexOf('=');
    if (eq === -1) continue;
    vars[trimmed.slice(0, eq).trim()] = trimmed.slice(eq + 1).trim();
  }
  return vars;
}

async function pickConfigFiles(): Promise<FoundFiles> {
  // showDirectoryPicker only reads metadata for files we explicitly request —
  // no enumeration of the whole folder tree.
  const w = window as unknown as {showDirectoryPicker: () => Promise<FileSystemDirectoryHandle>};
  const dirHandle = await w.showDirectoryPicker();

  const configDirHandle = await dirHandle.getDirectoryHandle('thunderid-config');

  let yamlFileHandle: FileSystemFileHandle | null = null;
  for (const name of ['thunderid-config.yaml', 'thunderid-config.yml']) {
    try {
      yamlFileHandle = await configDirHandle.getFileHandle(name);
      break;
    } catch {
      // try next
    }
  }

  if (!yamlFileHandle) {
    throw new Error('yaml_not_found');
  }

  let envFileHandle: FileSystemFileHandle | null = null;
  try {
    envFileHandle = await configDirHandle.getFileHandle('thunderid.env');
  } catch {
    // env file is optional
  }

  const yamlFile = await yamlFileHandle.getFile();
  const envFile = envFileHandle ? await envFileHandle.getFile() : null;

  return {folderName: dirHandle.name, yamlFile, envFile};
}

function formatImportedDate(ts: string): string {
  const date = new Date(Number(ts));
  if (isNaN(date.getTime())) return '';
  return date.toLocaleDateString(undefined, {day: 'numeric', month: 'short', year: 'numeric'});
}

interface WayfinderFolderImportProps {
  onSuccess?: () => void;
}

export default function WayfinderFolderImport({onSuccess = undefined}: WayfinderFolderImportProps): JSX.Element {
  const {t} = useTranslation(['common']);
  const {config} = useConfig();
  const productName = config.brand.product_name;
  const {mutateAsync: importConfig} = useImportConfiguration();

  const [status, setStatus] = useState<Status>(() => {
    const ts = localStorage.getItem(IMPORTED_KEY);
    return ts ? 'alreadyDone' : 'idle';
  });
  const [found, setFound] = useState<FoundFiles | null>(null);
  const [result, setResult] = useState<ImportResponse | null>(null);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const previousTs = localStorage.getItem(IMPORTED_KEY);

  const handleSelectFolder = async (): Promise<void> => {
    try {
      const files = await pickConfigFiles();
      setFound(files);
      setErrorMsg(null);
      setStatus('selected');
    } catch (err) {
      if (err instanceof DOMException && err.name === 'AbortError') return;
      setErrorMsg(t('common:welcome.wayfinderFolderImport.errors.cannotReadFolder'));
      setStatus('error');
      setFound(null);
    }
  };

  const handleImport = async (): Promise<void> => {
    if (!found) return;
    setStatus('importing');
    setErrorMsg(null);

    try {
      const content = await found.yamlFile.text();
      const variables = found.envFile ? parseEnvFile(await found.envFile.text()) : undefined;

      const response = await importConfig({content, variables, options: {upsert: true}});
      setResult(response);
      if (response.summary.failed > 0) {
        setErrorMsg(t('common:welcome.wayfinderFolderImport.errors.partialFailure', {count: response.summary.failed}));
        setStatus('error');
      } else {
        localStorage.setItem(IMPORTED_KEY, Date.now().toString());
        setStatus('success');
        onSuccess?.();
      }
    } catch {
      setErrorMsg(t('common:welcome.wayfinderFolderImport.errors.importFailed'));
      setStatus('error');
    }
  };

  const handleReset = (): void => {
    setStatus('idle');
    setFound(null);
    setResult(null);
    setErrorMsg(null);
  };

  if (status === 'importing') {
    return (
      <Stack direction="row" spacing={1.5} alignItems="center">
        <CircularProgress size={18} />
        <Typography variant="body2" color="text.secondary">
          {t('common:welcome.wayfinderFolderImport.status.importing')}
        </Typography>
      </Stack>
    );
  }

  if (status === 'alreadyDone') {
    return (
      <Stack spacing={1}>
        <Stack direction="row" spacing={1} alignItems="center">
          <CheckCircle size={20} style={{color: 'var(--oxygen-palette-success-main)'}} />
          <Typography variant="body2" fontWeight={600} color="success.main">
            {t('common:welcome.wayfinderFolderImport.status.alreadyDone', {productName})}
          </Typography>
        </Stack>
        {previousTs && (
          <Typography variant="caption" color="text.secondary">
            {t('common:welcome.wayfinderFolderImport.status.lastImported', {date: formatImportedDate(previousTs)})}
          </Typography>
        )}
        <Box>
          <Button
            variant="text"
            size="small"
            startIcon={<RefreshCw size={13} />}
            sx={{pl: 0, color: 'text.secondary', fontSize: '0.75rem'}}
            onClick={handleReset}
          >
            {t('common:welcome.wayfinderFolderImport.actions.reImport')}
          </Button>
        </Box>
      </Stack>
    );
  }

  if (status === 'success' && result) {
    return (
      <Stack spacing={1}>
        <Stack direction="row" spacing={1} alignItems="center">
          <CheckCircle size={20} style={{color: 'var(--oxygen-palette-success-main)'}} />
          <Typography variant="body2" fontWeight={600} color="success.main">
            {t('common:welcome.wayfinderFolderImport.status.success', {productName})}
          </Typography>
        </Stack>
        <Typography variant="caption" color="text.secondary">
          {t('common:welcome.wayfinderFolderImport.status.resourcesImported', {count: result.summary.imported})}
        </Typography>
      </Stack>
    );
  }

  return (
    <Stack spacing={1.5}>
      {status === 'idle' && (
        <Box>
          <Button
            variant="outlined"
            size="small"
            startIcon={<FolderOpen size={16} />}
            onClick={() => void handleSelectFolder()}
          >
            {t('common:welcome.wayfinderFolderImport.actions.selectFolder')}
          </Button>
        </Box>
      )}

      {status === 'selected' && found && (
        <Stack spacing={1.5}>
          <Box sx={{border: '1px solid', borderColor: 'divider', borderRadius: 1.5, px: 2, py: 1.5}}>
            <Stack direction="row" spacing={1} alignItems="center" sx={{mb: 1}}>
              <FolderOpen size={16} style={{opacity: 0.5, flexShrink: 0}} />
              <Typography variant="body2" fontWeight={600} noWrap sx={{flex: 1}}>
                {found.folderName}/
              </Typography>
              <Button
                size="small"
                variant="text"
                sx={{fontSize: '0.7rem', flexShrink: 0}}
                onClick={() => void handleSelectFolder()}
              >
                {t('common:welcome.wayfinderFolderImport.actions.change')}
              </Button>
            </Stack>
            <Stack spacing={0.5} sx={{pl: 0.5}}>
              <Stack direction="row" spacing={1} alignItems="center">
                <FileCode size={13} style={{opacity: 0.5, flexShrink: 0}} />
                <Typography variant="caption" color="text.secondary" noWrap>
                  thunderid-config/{found.yamlFile.name}
                </Typography>
              </Stack>
              <Stack direction="row" spacing={1} alignItems="center">
                <FileText size={13} style={{opacity: 0.5, flexShrink: 0}} />
                <Typography variant="caption" color="text.secondary" noWrap>
                  {found.envFile
                    ? `thunderid-config/${found.envFile.name}`
                    : t('common:welcome.wayfinderFolderImport.status.envNotFound')}
                </Typography>
              </Stack>
            </Stack>
          </Box>
          <Box>
            <Button variant="contained" size="small" onClick={() => void handleImport()}>
              {t('common:welcome.wayfinderFolderImport.actions.importConfig')}
            </Button>
          </Box>
        </Stack>
      )}

      {status === 'error' && (
        <Stack spacing={1}>
          <Stack direction="row" spacing={1} alignItems="center">
            <XCircle size={16} style={{color: 'var(--oxygen-palette-error-main)', flexShrink: 0}} />
            <Typography variant="caption" color="error.main">
              {errorMsg}
            </Typography>
          </Stack>
          {result?.results
            .filter((r) => r.status === 'failed')
            .map((r) => (
              <Typography
                key={`${r.resourceType}-${r.resourceId ?? r.resourceName ?? r.message}`}
                variant="caption"
                color="error.main"
                sx={{pl: 3, display: 'block'}}
              >
                {r.resourceType}
                {r.resourceName ? ` · ${r.resourceName}` : ''}: {r.message}
              </Typography>
            ))}
          <Box>
            <Button
              variant="outlined"
              size="small"
              startIcon={<RefreshCw size={13} />}
              onClick={() => {
                handleReset();
                void handleSelectFolder();
              }}
            >
              {t('common:welcome.wayfinderFolderImport.actions.reSelectFolder')}
            </Button>
          </Box>
        </Stack>
      )}
    </Stack>
  );
}
