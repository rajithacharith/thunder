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

import {waitFor} from '@testing-library/react';
import {renderHook} from '@thunderid/test-utils';
import {describe, it, expect, beforeEach, afterEach, vi} from 'vitest';
import type {CorsConfigResponse} from '../../models/responses';

const mockHttpRequest = vi.fn();
vi.mock('@thunderid/react', () => ({
  useThunderID: () => ({
    http: {request: mockHttpRequest},
  }),
}));

const mockGetServerUrl = vi.fn<() => string>(() => 'https://localhost:8090');
vi.mock('@thunderid/contexts', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@thunderid/contexts')>();
  return {
    ...actual,
    useConfig: () => ({getServerUrl: mockGetServerUrl}),
  };
});

const {default: useUpdateCorsConfig} = await import('../useUpdateCorsConfig');

describe('useUpdateCorsConfig', () => {
  const mockResponse: CorsConfigResponse = {
    readOnly: {allowedOrigins: ['https://localhost:5190']},
    writable: {allowedOrigins: ['https://app.example.com']},
    merged: {allowedOrigins: ['https://localhost:5190', 'https://app.example.com']},
  };

  beforeEach(() => {
    mockHttpRequest.mockReset();
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it('PUTs the writable value to /server-config/cors', async () => {
    mockHttpRequest.mockResolvedValue({data: mockResponse});
    const {result} = renderHook(() => useUpdateCorsConfig());

    result.current.mutate({data: {allowedOrigins: ['https://app.example.com']}});

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(mockHttpRequest).toHaveBeenCalledWith(
      expect.objectContaining({
        url: 'https://localhost:8090/server-config/cors',
        method: 'PUT',
        data: {allowedOrigins: ['https://app.example.com']},
      }),
    );
  });

  it('writes the returned config into the query cache on success', async () => {
    mockHttpRequest.mockResolvedValue({data: mockResponse});
    const {result, queryClient} = renderHook(() => useUpdateCorsConfig());

    result.current.mutate({data: {allowedOrigins: ['https://app.example.com']}});

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(queryClient.getQueryData(['server-config', 'cors'])).toEqual(mockResponse);
  });

  it('surfaces update errors', async () => {
    mockHttpRequest.mockRejectedValue(new Error('nope'));
    const {result} = renderHook(() => useUpdateCorsConfig());

    result.current.mutate({data: {allowedOrigins: []}});

    await waitFor(() => {
      expect(result.current.error?.message).toBe('nope');
    });
  });
});
