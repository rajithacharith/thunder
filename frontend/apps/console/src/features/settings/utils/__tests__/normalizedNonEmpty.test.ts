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

import {describe, it, expect} from 'vitest';
import normalizedNonEmpty from '../normalizedNonEmpty';

describe('normalizedNonEmpty', () => {
  it('normalizes each value (lowercase + trailing slash) and preserves order', () => {
    expect(normalizedNonEmpty(['HTTPS://Example.COM/', 'https://app.io'])).toEqual([
      'https://example.com',
      'https://app.io',
    ]);
  });

  it('drops empty and whitespace-only entries', () => {
    expect(normalizedNonEmpty(['', '   ', 'https://app.io'])).toEqual(['https://app.io']);
  });

  it('returns an empty array when there are no non-empty values', () => {
    expect(normalizedNonEmpty(['', '  '])).toEqual([]);
  });
});
