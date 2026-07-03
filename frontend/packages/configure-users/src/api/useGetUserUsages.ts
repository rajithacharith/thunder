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
import UserQueryKeys from '../constants/user-query-keys';
import type {UserUsagesResponse} from '../models/users';

/**
 * Custom hook to fetch resources that reference a user, such as agents that list the user
 * as their owner. Used to populate the pre-delete confirmation dialog.
 *
 * @param userId - The unique identifier of the user
 * @param enabled - Whether the query should run (default true)
 * @returns TanStack Query result with user usages data
 */
export default function useGetUserUsages(userId: string | null, enabled = true): UseQueryResult<UserUsagesResponse> {
  const {http} = useThunderID();
  const {getServerUrl} = useConfig();

  return useQuery<UserUsagesResponse>({
    queryKey: [UserQueryKeys.USER_USAGES, userId],
    queryFn: async (): Promise<UserUsagesResponse> => {
      const serverUrl: string = getServerUrl();

      const response: {data: UserUsagesResponse} = await http.request({
        url: `${serverUrl}/users/${encodeURIComponent(userId!)}/usages`,
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      } as unknown as Parameters<typeof http.request>[0]);

      return response.data;
    },
    enabled: Boolean(userId) && enabled,
  });
}
