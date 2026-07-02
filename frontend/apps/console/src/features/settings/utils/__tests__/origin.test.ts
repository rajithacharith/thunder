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
import {isValidOrigin, isValidRegex, normalizeOrigin} from '../origin';

describe('isValidOrigin', () => {
  it.each(['https://app.example.com', 'http://localhost:3000', 'https://example.com:8443', 'https://[::1]:8443'])(
    'accepts the valid origin %s',
    (value) => {
      expect(isValidOrigin(value)).toBe(true);
    },
  );

  it('accepts the "null" literal', () => {
    expect(isValidOrigin('null')).toBe(true);
  });

  it.each([
    '',
    'https://example.com/path',
    'https://example.com/foo/..',
    'https://example.com/.',
    'http:example.com/..',
    'https:example.com',
    'https://example.com\\path',
    'https://example.com?q=1',
    'https://example.com#frag',
    'https://user@example.com',
    'ftp://example.com',
    'https://*.example.com',
    '*',
  ])('rejects the non-origin %s', (value) => {
    expect(isValidOrigin(value)).toBe(false);
  });
});

describe('isValidRegex', () => {
  it.each(['^https://.*\\.example\\.com$', 'abc', 'https://example.com/path'])(
    'accepts the compilable pattern %s',
    (value) => {
      expect(isValidRegex(value)).toBe(true);
    },
  );

  it.each(['', '[', '(unbalanced'])('rejects the non-compilable pattern %s', (value) => {
    expect(isValidRegex(value)).toBe(false);
  });
});

describe('normalizeOrigin', () => {
  it('lowercases scheme and host and strips a trailing slash', () => {
    expect(normalizeOrigin('HTTPS://Example.COM/')).toBe('https://example.com');
  });

  it('preserves an explicit default port (unlike URL.origin)', () => {
    expect(normalizeOrigin('https://example.com:443')).toBe('https://example.com:443');
    expect(normalizeOrigin('https://example.com')).not.toBe(normalizeOrigin('https://example.com:443'));
  });

  it('trims surrounding whitespace', () => {
    expect(normalizeOrigin('  https://app.io  ')).toBe('https://app.io');
  });

  it('leaves the "null" literal and regex/invalid input unchanged', () => {
    expect(normalizeOrigin('null')).toBe('null');
    expect(normalizeOrigin('^https://x$')).toBe('^https://x$');
    expect(normalizeOrigin('')).toBe('');
  });
});
