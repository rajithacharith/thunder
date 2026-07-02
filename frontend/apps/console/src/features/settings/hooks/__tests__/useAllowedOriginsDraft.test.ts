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

import {act} from '@testing-library/react';
import {renderHook} from '@thunderid/test-utils';
import {describe, it, expect} from 'vitest';
import type {AllowedOrigin, CorsConfigResponse} from '../../models/responses';
import useAllowedOriginsDraft from '../useAllowedOriginsDraft';

function makeData(readOnly: AllowedOrigin[], writable: AllowedOrigin[]): CorsConfigResponse {
  return {
    readOnly: {allowedOrigins: readOnly},
    writable: {allowedOrigins: writable},
    merged: {allowedOrigins: []},
  };
}

describe('useAllowedOriginsDraft', () => {
  it('loads writable entries as editable rows (strings as-is, regex as its pattern)', () => {
    const {result} = renderHook(() =>
      useAllowedOriginsDraft(makeData([], ['https://app.acme.com', {regex: '^https://[a-z]+\\.acme\\.io$'}])),
    );
    expect(result.current.draft).toEqual(['https://app.acme.com', '^https://[a-z]+\\.acme\\.io$']);
    expect(result.current.dirty).toBe(false);
    expect(result.current.hasErrors).toBe(false);
  });

  it('adding an empty row is not dirty; typing a value makes it dirty', () => {
    const {result} = renderHook(() => useAllowedOriginsDraft(makeData([], ['https://app.acme.com'])));
    act(() => result.current.addRow());
    expect(result.current.dirty).toBe(false);
    act(() => result.current.changeRow(1, 'https://new.example.com'));
    expect(result.current.dirty).toBe(true);
  });

  it('removing a row marks the draft dirty', () => {
    const {result} = renderHook(() =>
      useAllowedOriginsDraft(makeData([], ['https://app.acme.com', 'https://other.acme.com'])),
    );
    act(() => result.current.removeRow(0));
    expect(result.current.dirty).toBe(true);
    expect(result.current.draft).toEqual(['https://other.acme.com']);
  });

  it('reset clears local edits and reverts to the saved value', () => {
    const {result} = renderHook(() => useAllowedOriginsDraft(makeData([], ['https://app.acme.com'])));
    act(() => result.current.changeRow(0, 'https://changed.example.com'));
    expect(result.current.dirty).toBe(true);
    act(() => result.current.reset());
    expect(result.current.dirty).toBe(false);
    expect(result.current.draft).toEqual(['https://app.acme.com']);
  });

  it('normalizes a row on blur (lowercase + trailing slash), preserving an explicit port', () => {
    const {result} = renderHook(() => useAllowedOriginsDraft(makeData([], ['https://app.acme.com'])));
    act(() => result.current.changeRow(0, 'HTTPS://Example.COM:443/'));
    act(() => result.current.blurRow(0));
    expect(result.current.draft[0]).toBe('https://example.com:443');
    expect(result.current.hasErrors).toBe(false);
  });

  it('flags a row that is neither a valid origin nor a compilable regex', () => {
    const {result} = renderHook(() => useAllowedOriginsDraft(makeData([], ['(bad'])));
    let ok = true;
    act(() => {
      ok = result.current.validateAll();
    });
    expect(ok).toBe(false);
    expect(result.current.errors[0]).toBeTruthy();
  });

  it('accepts both a valid origin and a valid regex row', () => {
    const {result} = renderHook(() =>
      useAllowedOriginsDraft(makeData([], ['https://ok.example.com', {regex: '^https://.*\\.ok\\.io$'}])),
    );
    let ok = false;
    act(() => {
      ok = result.current.validateAll();
    });
    expect(ok).toBe(true);
    expect(result.current.hasErrors).toBe(false);
  });

  it('flags duplicates within the draft and clears them when a counterpart row is removed', () => {
    const {result} = renderHook(() =>
      useAllowedOriginsDraft(makeData([], ['https://dup.example.com', 'https://dup.example.com'])),
    );
    act(() => {
      result.current.validateAll();
    });
    expect(result.current.errors[0]).toBeTruthy();
    expect(result.current.errors[1]).toBeTruthy();

    act(() => result.current.removeRow(1));
    expect(result.current.errors[0]).toBeUndefined();
  });

  it('clears a duplicate error when the counterpart is edited to a unique value', () => {
    const {result} = renderHook(() =>
      useAllowedOriginsDraft(makeData([], ['https://dup.example.com', 'https://dup.example.com'])),
    );
    act(() => {
      result.current.validateAll();
    });
    expect(result.current.errors[0]).toBeTruthy();

    act(() => result.current.changeRow(1, 'https://unique.example.com'));
    expect(result.current.errors[0]).toBeUndefined();
  });

  it('treats a default port as distinct from the port-less origin (no false duplicate)', () => {
    const {result} = renderHook(() =>
      useAllowedOriginsDraft(makeData([], ['https://example.com', 'https://example.com:443'])),
    );
    act(() => {
      result.current.validateAll();
    });
    expect(result.current.hasErrors).toBe(false);
  });

  it('flags a custom origin that duplicates a read-only origin', () => {
    const {result} = renderHook(() =>
      useAllowedOriginsDraft(makeData(['https://console.example.com'], ['https://console.example.com'])),
    );
    act(() => {
      result.current.validateAll();
    });
    expect(result.current.errors[0]).toBeTruthy();
  });

  it('flags a custom regex that duplicates a read-only regex', () => {
    const {result} = renderHook(() =>
      useAllowedOriginsDraft(
        makeData([{regex: '^https://[a-z]+\\.acme\\.io$'}], [{regex: '^https://[a-z]+\\.acme\\.io$'}]),
      ),
    );
    act(() => {
      result.current.validateAll();
    });
    expect(result.current.errors[0]).toBeTruthy();
  });

  describe('buildPayload', () => {
    it('classifies origins as strings and non-origins as regex entries, dropping empty rows', () => {
      const {result} = renderHook(() =>
        useAllowedOriginsDraft(makeData([], ['https://app.example.com', {regex: '^https://[a-z]+\\.example\\.com$'}])),
      );
      // Add a trailing empty row that Save must drop.
      act(() => result.current.addRow());

      expect(result.current.buildPayload()).toEqual({
        allowedOrigins: ['https://app.example.com', {regex: '^https://[a-z]+\\.example\\.com$'}],
      });
    });

    it('preserves an explicit default port in the saved string entry', () => {
      const {result} = renderHook(() => useAllowedOriginsDraft(makeData([], ['https://example.com:443'])));
      expect(result.current.buildPayload()).toEqual({allowedOrigins: ['https://example.com:443']});
    });

    it('round-trips a loaded regex entry back to a {regex} entry', () => {
      const {result} = renderHook(() => useAllowedOriginsDraft(makeData([], [{regex: '^https://x\\.io$'}])));
      expect(result.current.draft).toEqual(['^https://x\\.io$']);
      expect(result.current.buildPayload()).toEqual({allowedOrigins: [{regex: '^https://x\\.io$'}]});
    });

    it('keeps the "null" literal as a string entry', () => {
      const {result} = renderHook(() => useAllowedOriginsDraft(makeData([], ['null'])));
      expect(result.current.buildPayload()).toEqual({allowedOrigins: ['null']});
    });
  });
});
