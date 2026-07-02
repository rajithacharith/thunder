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

import {useMutation, useQueryClient, type UseMutationResult} from '@tanstack/react-query';
import {useConfig, useToast} from '@thunderid/contexts';
import {useThunderID} from '@thunderid/react';
import {getErrorMessage} from '@thunderid/utils';
import {useTranslation} from 'react-i18next';
import SettingsQueryKeys from '../constants/settings-query-keys';
import type {CorsConfigResponse, CorsValue} from '../models/responses';

/**
 * Variables for the {@link useUpdateCorsConfig} mutation.
 *
 * @public
 */
export interface UpdateCorsConfigVariables {
  /**
   * The writable CORS value to persist.
   */
  data: CorsValue;
}

/**
 * Updates the writable layer of the CORS server-config section.
 *
 * @returns TanStack Query mutation for updating allowed origins.
 *
 * @public
 */
export default function useUpdateCorsConfig(): UseMutationResult<CorsConfigResponse, Error, UpdateCorsConfigVariables> {
  const {http} = useThunderID();
  const {getServerUrl} = useConfig();
  const queryClient: ReturnType<typeof useQueryClient> = useQueryClient();
  const {t} = useTranslation();
  const {showToast} = useToast();

  return useMutation<CorsConfigResponse, Error, UpdateCorsConfigVariables>({
    mutationFn: async ({data}: UpdateCorsConfigVariables): Promise<CorsConfigResponse> => {
      const serverUrl: string = getServerUrl();

      const response: {
        data: CorsConfigResponse;
      } = await http.request({
        url: `${serverUrl}/server-config/cors`,
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        data,
      } as unknown as Parameters<typeof http.request>[0]);

      return response.data;
    },
    onSuccess: (data) => {
      // PUT returns the full config, so keep the read model in sync without a refetch.
      queryClient.setQueryData([SettingsQueryKeys.SERVER_CONFIG, SettingsQueryKeys.CORS], data);
      showToast(t('settings:cors.save.success'), 'success');
    },
    onError: (error) => {
      showToast(getErrorMessage(error, t, 'settings:cors.save.error'), 'error');
    },
  });
}
