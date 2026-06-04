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

import {Box, Button, Chip, Collapse, Stack, Typography} from '@wso2/oxygen-ui';
import {ChevronDown, ChevronUp, Download} from '@wso2/oxygen-ui-icons-react';
import type {JSX} from 'react';
import {useEffect, useMemo, useState} from 'react';
import {useTranslation} from 'react-i18next';
import useWayfinderReleases from '../api/useWayfinderReleases';
import OS_LABELS from '../constants/download-assets';
import type {DetectedPlatform} from '../models/download-assets';
import {detectPlatform, groupAssetsByOs, parseDistributionAssets, pickAssetForPlatform} from '../utils/downloadAssets';

export default function WayfinderSampleDownload({releasesUrl}: {releasesUrl: string}): JSX.Element | null {
  const {t} = useTranslation(['common']);
  const [platform, setPlatform] = useState<DetectedPlatform | null>(null);
  const [showOthers, setShowOthers] = useState(false);

  const pattern = useMemo(() => /^sample-app-wayfinder-[0-9A-Za-z.+-]+-(macos|linux|win)-(arm64|x64)\.zip$/i, []);

  useEffect(() => {
    detectPlatform()
      .then(setPlatform)
      .catch(() => undefined);
  }, []);

  const {data, isError: errored} = useWayfinderReleases(releasesUrl);

  const release = data ? (data.latestRelease ?? data.releases?.[0] ?? null) : null;
  const tag = release?.tagName ?? '';
  const assets = useMemo(() => (release ? parseDistributionAssets(release.assets, pattern) : null), [release, pattern]);

  const selected = useMemo(() => pickAssetForPlatform(assets ?? [], platform), [assets, platform]);
  const matched = selected && selected.os === platform?.os && selected.arch === platform?.arch;
  const groups = useMemo(() => groupAssetsByOs(assets ?? [], selected?.os ?? null), [assets, selected]);

  if (errored || (assets !== null && assets.length === 0)) return null;
  if (!selected) return null;

  return (
    <Box sx={{border: '1px solid', borderColor: 'divider', borderRadius: 2, overflow: 'hidden'}}>
      {/* Primary download row */}
      <Box
        sx={{display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 2, p: 2, flexWrap: 'wrap'}}
      >
        <Box>
          <Typography variant="caption" color="text.secondary" sx={{display: 'block', mb: 0.25}}>
            {matched
              ? t('common:welcome.wayfinderSampleDownload.recommendedLabel')
              : t('common:welcome.wayfinderSampleDownload.selectedLabel')}
          </Typography>
          <Typography variant="subtitle2" fontWeight={600}>
            {selected.osLabel} · {selected.archLabel}
          </Typography>
          <Stack direction="row" spacing={1} sx={{mt: 0.5, flexWrap: 'wrap'}}>
            {tag && <Chip label={tag} size="small" />}
            {selected.sizeLabel && <Chip label={selected.sizeLabel} size="small" variant="outlined" />}
          </Stack>
        </Box>
        <Button
          variant="contained"
          size="small"
          startIcon={<Download size={16} />}
          href={selected.downloadUrl}
          target="_blank"
          rel="noreferrer"
          component="a"
        >
          {t('common:welcome.wayfinderSampleDownload.downloadButton', {osLabel: selected.osLabel})}
        </Button>
      </Box>

      {/* Other platforms toggle */}
      {groups.length > 1 && (
        <>
          <Box
            role="button"
            tabIndex={0}
            onClick={() => setShowOthers((v) => !v)}
            onKeyDown={(e: React.KeyboardEvent) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                setShowOthers((v) => !v);
              }
            }}
            sx={{
              display: 'flex',
              alignItems: 'center',
              gap: 1,
              px: 2,
              py: 1,
              borderTop: '1px solid',
              borderColor: 'divider',
              cursor: 'pointer',
              bgcolor: 'action.hover',
              '&:hover': {bgcolor: 'action.selected'},
            }}
          >
            {showOthers ? <ChevronUp size={14} /> : <ChevronDown size={14} />}
            <Typography variant="caption">
              {showOthers
                ? t('common:welcome.wayfinderSampleDownload.hidePlatforms')
                : t('common:welcome.wayfinderSampleDownload.otherPlatforms')}
            </Typography>
          </Box>
          <Collapse in={showOthers}>
            <Box sx={{p: 2, display: 'flex', gap: 2, flexWrap: 'wrap'}}>
              {groups.map(({os, assets: osAssets}) => (
                <Box key={os} sx={{minWidth: 160}}>
                  <Typography variant="caption" fontWeight={600} color="text.secondary" sx={{display: 'block', mb: 1}}>
                    {OS_LABELS[os]}
                  </Typography>
                  <Stack spacing={0.5}>
                    {osAssets.map((asset) => (
                      <Box
                        key={asset.name}
                        component="a"
                        href={asset.downloadUrl}
                        target="_blank"
                        rel="noreferrer"
                        sx={{
                          display: 'flex',
                          alignItems: 'center',
                          gap: 1,
                          textDecoration: 'none',
                          color: 'inherit',
                          px: 1.5,
                          py: 0.75,
                          borderRadius: 1,
                          border: '1px solid',
                          borderColor: asset.downloadUrl === selected.downloadUrl ? 'primary.main' : 'divider',
                          '&:hover': {bgcolor: 'action.hover'},
                        }}
                      >
                        <Box sx={{flex: 1}}>
                          <Typography variant="caption" sx={{display: 'block', fontWeight: 500}}>
                            {asset.osLabel} {asset.archLabel}
                          </Typography>
                          {asset.sizeLabel && (
                            <Typography variant="caption" color="text.secondary">
                              {asset.sizeLabel}
                            </Typography>
                          )}
                        </Box>
                        <Download size={12} style={{opacity: 0.5}} />
                      </Box>
                    ))}
                  </Stack>
                </Box>
              ))}
            </Box>
          </Collapse>
        </>
      )}
    </Box>
  );
}
