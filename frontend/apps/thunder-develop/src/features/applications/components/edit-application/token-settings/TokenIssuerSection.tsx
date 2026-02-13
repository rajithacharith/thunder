/**
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
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

import {TextField} from '@wso2/oxygen-ui';
import {useTranslation} from 'react-i18next';
import {Controller} from 'react-hook-form';
import type {Control, FieldErrors} from 'react-hook-form';
import SettingsCard from '../../../../../components/SettingsCard';

/**
 * Props for the {@link TokenIssuerSection} component.
 */
interface TokenIssuerSectionProps {
  /**
   * React Hook Form control object for the token config form
   */
  control: Control<{
    validityPeriod: number;
    accessTokenValidity: number;
    idTokenValidity: number;
    issuer?: string | undefined;
  }>;
  /**
   * Form validation errors
   */
  errors: FieldErrors<{
    validityPeriod: number;
    accessTokenValidity: number;
    idTokenValidity: number;
    issuer?: string | undefined;
  }>;
}

/**
 * Section component for configuring the token issuer URL.
 *
 * Displays a text input for setting the JWT `iss` (issuer) claim value.
 * Includes URL validation to ensure the issuer is a valid URL.
 *
 * @param props - Component props
 * @returns Token issuer input UI within a SettingsCard
 */
export default function TokenIssuerSection({control, errors}: TokenIssuerSectionProps) {
  const {t} = useTranslation();

  return (
    <SettingsCard
      title={t('applications:edit.token.tokenIssuer.title')}
      description={t('applications:edit.token.tokenIssuer.description')}
    >
      <Controller
        name="issuer"
        control={control}
        render={({field}) => (
          <TextField
            {...field}
            fullWidth
            label={t('applications:edit.token.labels.issuer')}
            placeholder={t('applications:edit.token.issuer.placeholder')}
            error={!!errors.issuer}
            helperText={errors.issuer?.message ?? t('applications:edit.token.issuer.hint')}
          />
        )}
      />
    </SettingsCard>
  );
}
