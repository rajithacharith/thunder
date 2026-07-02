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

import {useQuery, type UseQueryResult} from '@tanstack/react-query';
import {useConfig} from '@thunderid/contexts';
import {useThunderID} from '@thunderid/react';
import SettingsQueryKeys from '../constants/settings-query-keys';
import type {CorsConfigResponse} from '../models/responses';

/**
 * Fetches the CORS server-config section.
 *
 * @returns TanStack Query result containing the CORS config layers.
 *
 * @public
 */
export default function useGetCorsConfig(): UseQueryResult<CorsConfigResponse> {
  const {http} = useThunderID();
  const {getServerUrl} = useConfig();

  return useQuery<CorsConfigResponse>({
    queryKey: [SettingsQueryKeys.SERVER_CONFIG, SettingsQueryKeys.CORS],
    queryFn: async (): Promise<CorsConfigResponse> => {
      const serverUrl: string = getServerUrl();

      const response: {
        data: CorsConfigResponse;
      } = await http.request({
        url: `${serverUrl}/server-config/cors`,
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      } as unknown as Parameters<typeof http.request>[0]);

      return response.data;
    },
  });
}
