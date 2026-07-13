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

import {SettingsCard} from '@thunderid/components';
import {Box, Button, FormControl, FormLabel, IconButton, Stack, TextField, Tooltip, Typography} from '@wso2/oxygen-ui';
import {Plus, Trash} from '@wso2/oxygen-ui-icons-react';
import {useEffect, useRef, useState} from 'react';
import {useTranslation} from 'react-i18next';
import type {AttestationConfig} from '../../../models/oauth';

/**
 * Props for the {@link AttestationSection} component.
 */
interface AttestationSectionProps {
  /**
   * The current attestation config (from the OAuth config).
   * null or undefined means no attestation is configured.
   */
  attestation?: AttestationConfig | null;
  /**
   * Called when the user changes any attestation field.
   * Passes null when no attestation fields are set.
   */
  onAttestationChange: (attestation: AttestationConfig | null) => void;
  /**
   * Whether inputs should be disabled (e.g. read-only resource).
   */
  disabled?: boolean;
}

/**
 * Section component for configuring Google Play Integrity attestation for Android mobile clients.
 *
 * A mobile application that configures attestation may initiate an authentication flow directly by
 * presenting a Play Integrity token, which the server verifies against the registered package name
 * and signing certificate digests. The service account credentials are write-only: they are never
 * returned by the API, so leaving the field blank when editing preserves the stored value.
 *
 * @param props - Component props
 * @returns Attestation configuration UI within a SettingsCard
 */
export default function AttestationSection({
  attestation = undefined,
  onAttestationChange,
  disabled = false,
}: AttestationSectionProps) {
  const {t} = useTranslation();

  const android = attestation?.android;
  // The fields are backed by local state (seeded from props) so editing is never gated on the
  // parent's config round-trip. Credentials are write-only and never seeded from props.
  const [packageName, setPackageName] = useState<string>(android?.packageName ?? '');
  const [digests, setDigests] = useState<string[]>(android?.certificateSha256Digests ?? []);
  const [credentials, setCredentials] = useState<string>('');

  // Identity of the incoming config (package name + digests). The effect below resyncs local state
  // when the attestation prop is replaced externally — e.g. the application reloads, or the config
  // is cleared — while ignoring the echo of this component's own emissions (tracked via the ref).
  const identityKey = JSON.stringify({
    packageName: android?.packageName ?? '',
    digests: android?.certificateSha256Digests ?? [],
  });
  const lastSyncedKeyRef = useRef<string>(identityKey);

  useEffect(() => {
    if (identityKey === lastSyncedKeyRef.current) {
      return;
    }
    lastSyncedKeyRef.current = identityKey;
    setPackageName(android?.packageName ?? '');
    setDigests(android?.certificateSha256Digests ?? []);
    // Credentials are write-only; an external config change resets the editable field to blank.
    setCredentials('');
    // identityKey is the canonical trigger; android is read for the values it encodes.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [identityKey]);

  const emit = (pkg: string, digs: string[], creds: string) => {
    const cleanedDigests = digs.map((d) => d.trim()).filter((d) => d !== '');
    const cleanedPackageName = pkg.trim();

    // Record the identity being emitted so the resync effect ignores the resulting prop echo and
    // preserves the user's in-progress edits.
    lastSyncedKeyRef.current = JSON.stringify({packageName: cleanedPackageName, digests: cleanedDigests});

    if (cleanedPackageName === '' && cleanedDigests.length === 0 && creds === '') {
      onAttestationChange(null);
      return;
    }

    const androidConfig: NonNullable<AttestationConfig['android']> = {};
    if (cleanedPackageName !== '') {
      androidConfig.packageName = cleanedPackageName;
    }
    if (cleanedDigests.length > 0) {
      androidConfig.certificateSha256Digests = cleanedDigests;
    }
    if (creds !== '') {
      androidConfig.serviceAccountCredentials = creds;
    }
    onAttestationChange({android: androidConfig});
  };

  const handlePackageNameChange = (value: string) => {
    setPackageName(value);
    emit(value, digests, credentials);
  };

  const handleCredentialsChange = (value: string) => {
    setCredentials(value);
    emit(packageName, digests, value);
  };

  const commitDigests = (nextDigests: string[]) => {
    setDigests(nextDigests);
    emit(packageName, nextDigests, credentials);
  };

  const handleAddDigest = () => {
    setDigests((prev) => [...prev, '']);
  };

  const handleDigestChange = (index: number, value: string) => {
    setDigests((prev) => prev.map((d, i) => (i === index ? value : d)));
  };

  const handleRemoveDigest = (index: number) => {
    commitDigests(digests.filter((_, i) => i !== index));
  };

  return (
    <SettingsCard
      title={t('applications:edit.advanced.labels.attestation')}
      description={t('applications:edit.advanced.attestation.intro')}
    >
      <Stack spacing={2}>
        <FormControl fullWidth>
          <FormLabel htmlFor="attestation-package-name">
            {t('applications:edit.advanced.attestation.labels.packageName')}
          </FormLabel>
          <TextField
            id="attestation-package-name"
            fullWidth
            value={packageName}
            onChange={(e) => handlePackageNameChange(e.target.value)}
            placeholder={t('applications:edit.advanced.attestation.placeholder.packageName')}
            helperText={t('applications:edit.advanced.attestation.hint.packageName')}
            disabled={disabled}
          />
        </FormControl>

        <FormControl fullWidth>
          <FormLabel htmlFor="attestation-digests-section">
            {t('applications:edit.advanced.attestation.labels.certificateSha256Digests')}
          </FormLabel>
          <Typography variant="caption" color="text.secondary" sx={{display: 'block', mb: 2}}>
            {t('applications:edit.advanced.attestation.hint.certificateSha256Digests')}
          </Typography>
          <Stack spacing={2} id="attestation-digests-section">
            {digests.map((digest, index) => (
              // IMPORTANT: Do not remove the suppression since it affects functionality.
              // eslint-disable-next-line react/no-array-index-key
              <Stack key={index} direction="row" spacing={1} alignItems="flex-start">
                <FormControl fullWidth sx={{flex: 1}}>
                  <TextField
                    fullWidth
                    id={`attestation-digest-${index}-input`}
                    // Each repeated field needs a unique accessible name; the shared FormLabel points
                    // at the surrounding Stack, not the individual inputs.
                    aria-label={`${t(
                      'applications:edit.advanced.attestation.labels.certificateSha256Digests',
                    )} ${index + 1}`}
                    value={digest}
                    onChange={(e) => handleDigestChange(index, e.target.value)}
                    onBlur={() => commitDigests(digests)}
                    placeholder={t('applications:edit.advanced.attestation.placeholder.certificateSha256Digest')}
                    disabled={disabled}
                  />
                </FormControl>
                <Tooltip title={t('common:actions.delete')}>
                  <IconButton onClick={() => handleRemoveDigest(index)} color="error" sx={{mt: 1}} disabled={disabled}>
                    <Trash size={20} />
                  </IconButton>
                </Tooltip>
              </Stack>
            ))}
            <Box>
              <Button
                variant="outlined"
                startIcon={<Plus />}
                onClick={handleAddDigest}
                size="small"
                disabled={disabled}
              >
                {t('applications:edit.advanced.attestation.addDigest')}
              </Button>
            </Box>
          </Stack>
        </FormControl>

        <FormControl fullWidth>
          <FormLabel htmlFor="attestation-service-account">
            {t('applications:edit.advanced.attestation.labels.serviceAccountCredentials')}
          </FormLabel>
          <TextField
            id="attestation-service-account"
            fullWidth
            multiline
            rows={4}
            value={credentials}
            onChange={(e) => handleCredentialsChange(e.target.value)}
            placeholder={t('applications:edit.advanced.attestation.placeholder.serviceAccountCredentials')}
            helperText={t('applications:edit.advanced.attestation.hint.serviceAccountCredentials')}
            disabled={disabled}
          />
        </FormControl>
      </Stack>
    </SettingsCard>
  );
}
