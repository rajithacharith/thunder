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

import {useLogger} from '@thunderid/logger';
import {Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions, Button, Alert} from '@wso2/oxygen-ui';
import {useState, type JSX} from 'react';
import {useTranslation} from 'react-i18next';
import useRegenerateFlowSecret from '../api/useRegenerateFlowSecret';

/**
 * Props for the {@link RegenerateFlowSecretDialog} component.
 */
export interface RegenerateFlowSecretDialogProps {
  /**
   * Whether the dialog is open
   */
  open: boolean;
  /**
   * The ID of the application whose Flow Secret will be regenerated
   */
  applicationId: string | null;
  /**
   * Callback when the dialog should be closed
   */
  onClose: () => void;
  /**
   * Callback when the Flow Secret is successfully regenerated with the new Flow Secret
   */
  onSuccess?: (newFlowSecret: string) => void;
  /**
   * Callback when the regeneration fails
   */
  onError?: (message: string) => void;
}

/**
 * Dialog component for confirming Flow Secret regeneration.
 *
 * Warns users that regenerating the Flow Secret immediately invalidates the current one, which will
 * break any server-side flow initiation until the new secret is deployed.
 *
 * @param props - Component props
 * @returns The regenerate Flow Secret confirmation dialog
 */
export default function RegenerateFlowSecretDialog({
  open,
  applicationId,
  onClose,
  onSuccess = undefined,
  onError = undefined,
}: RegenerateFlowSecretDialogProps): JSX.Element {
  const {t} = useTranslation();
  const logger = useLogger('RegenerateFlowSecretDialog');
  const [error, setError] = useState<string | null>(null);
  const regenerateFlowSecret = useRegenerateFlowSecret();

  const handleCancel = (): void => {
    setError(null);
    onClose();
  };

  const handleConfirm = (): void => {
    if (!applicationId) {
      setError(t('applications:regenerateFlowSecret.dialog.error'));
      return;
    }

    setError(null);
    logger.info('Regenerating application Flow Secret', {applicationId});

    regenerateFlowSecret.mutate(
      {applicationId},
      {
        onSuccess: ({flowSecret}) => {
          logger.info('Application Flow Secret regenerated successfully.', {applicationId});
          onClose();
          onSuccess?.(flowSecret);
        },
        onError: (err) => {
          const errorMessage = err instanceof Error ? err.message : t('applications:regenerateFlowSecret.dialog.error');
          logger.error('Failed to regenerate Flow Secret', {
            applicationId,
            errorMessage,
            errorName: err instanceof Error ? err.name : 'UnknownError',
          });
          setError(errorMessage);
          onError?.(errorMessage);
        },
      },
    );
  };

  return (
    <Dialog open={open} onClose={handleCancel} maxWidth="sm" fullWidth>
      <DialogTitle>{t('applications:regenerateFlowSecret.dialog.title')}</DialogTitle>
      <DialogContent>
        <DialogContentText sx={{mb: 2}}>{t('applications:regenerateFlowSecret.dialog.message')}</DialogContentText>
        <Alert severity="warning" sx={{mb: 2}}>
          {t('applications:regenerateFlowSecret.dialog.disclaimer')}
        </Alert>
        {error && (
          <Alert severity="error" sx={{mt: 2}}>
            {error}
          </Alert>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={handleCancel} disabled={regenerateFlowSecret.isPending}>
          {t('common:actions.cancel')}
        </Button>
        <Button
          onClick={handleConfirm}
          color="error"
          variant="contained"
          disabled={regenerateFlowSecret.isPending || !applicationId}
        >
          {regenerateFlowSecret.isPending
            ? t('applications:regenerateFlowSecret.dialog.regenerating')
            : t('applications:regenerateFlowSecret.dialog.confirmButton')}
        </Button>
      </DialogActions>
    </Dialog>
  );
}
