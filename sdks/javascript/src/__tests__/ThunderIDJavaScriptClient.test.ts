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

import {describe, expect, it, vi, beforeEach} from 'vitest';
import ThunderIDJavaScriptClient from '../ThunderIDJavaScriptClient';
import type {Storage} from '../models/store';

vi.mock('../IsomorphicCrypto', () => ({
  IsomorphicCrypto: class MockIsomorphicCrypto {
    constructor(_cryptoUtils: unknown) {}
  },
}));

vi.mock('../utils/AuthenticationHelper', () => ({
  default: class MockAuthenticationHelper {
    constructor(_storage: unknown, _crypto: unknown) {}
  },
}));

class MemoryStore implements Storage {
  private store = new Map<string, string>();

  async getData(key: string): Promise<string> {
    return this.store.get(key) ?? null!;
  }

  async setData(key: string, value: string): Promise<void> {
    this.store.set(key, value);
  }

  async removeData(key: string): Promise<void> {
    this.store.delete(key);
  }
}

async function getStoredConfig(client: ThunderIDJavaScriptClient): Promise<Record<string, any>> {
  return (client as any).storageManager.getConfigData();
}

describe('ThunderIDJavaScriptClient', () => {
  let store: MemoryStore;

  beforeEach(() => {
    vi.clearAllMocks();
    store = new MemoryStore();
  });

  describe('initialize()', () => {
    it('should apply DEFAULT_CONFIG baseline when no overrides are provided', async () => {
      const client = new ThunderIDJavaScriptClient(store, {} as any);

      await client.initialize({baseUrl: 'https://example.com', clientId: 'test-client'} as any);

      const config = await getStoredConfig(client);

      expect(config['enablePKCE']).toBe(true);
      expect(config['sendCookiesInRequests']).toBe(true);
      expect(config['tokenValidation']['idToken']['clockTolerance']).toBe(300);
      expect(config['tokenValidation']['idToken']['validate']).toBe(true);
      expect(config['tokenValidation']['idToken']['validateIssuer']).toBe(true);
    });

    it('should deep-merge partial tokenValidation, preserving sibling defaults', async () => {
      const client = new ThunderIDJavaScriptClient(store, {} as any);

      await client.initialize({
        baseUrl: 'https://example.com',
        clientId: 'test-client',
        tokenValidation: {idToken: {validate: false}},
      } as any);

      const config = await getStoredConfig(client);

      expect(config['tokenValidation']['idToken']['validate']).toBe(false);
      expect(config['tokenValidation']['idToken']['clockTolerance']).toBe(300);
      expect(config['tokenValidation']['idToken']['validateIssuer']).toBe(true);
    });

    it('should allow individual tokenValidation fields to be overridden independently', async () => {
      const client = new ThunderIDJavaScriptClient(store, {} as any);

      await client.initialize({
        baseUrl: 'https://example.com',
        clientId: 'test-client',
        tokenValidation: {idToken: {clockTolerance: 60}},
      } as any);

      const config = await getStoredConfig(client);

      expect(config['tokenValidation']['idToken']['clockTolerance']).toBe(60);
      expect(config['tokenValidation']['idToken']['validate']).toBe(true);
      expect(config['tokenValidation']['idToken']['validateIssuer']).toBe(true);
    });

    it('should set explicit fields (applicationId, scope) at highest precedence', async () => {
      const client = new ThunderIDJavaScriptClient(store, {} as any);

      await client.initialize({
        applicationId: 'app-123',
        baseUrl: 'https://example.com',
        clientId: 'test-client',
        scopes: ['openid', 'profile'],
      } as any);

      const config = await getStoredConfig(client);

      expect(config['applicationId']).toBe('app-123');
      expect(config['scope']).toContain('openid');
    });
  });
});
