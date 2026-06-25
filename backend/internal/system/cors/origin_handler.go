/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package cors

import "encoding/json"

// OriginHandler decodes, validates, and merges CORS origin config. It satisfies a consumer's
// value-handler interface structurally, so the cors package stays decoupled from the configuration store.
type OriginHandler struct{}

// Decode parses raw JSON origin entries into typed OriginEntries. Empty input yields an empty (non-nil)
// list so an unset layer serializes as [] rather than null, consistent with the merged layer.
func (OriginHandler) Decode(raw json.RawMessage) (any, error) {
	if len(raw) == 0 {
		return OriginEntries{}, nil
	}
	var entries OriginEntries
	if err := json.Unmarshal(raw, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// Validate checks that incoming is a valid set of origin entries. The readOnly and writable layers are
// unused for CORS — an origin list is valid on its own.
func (OriginHandler) Validate(incoming, _, _ any) error {
	entries, _ := incoming.(OriginEntries)
	return Validate(entries)
}

// Merge returns the union of the readOnly and writable origin lists, de-duplicated, preserving order
// (readOnly entries first). Empty or absent layers contribute nothing.
func (OriginHandler) Merge(readOnly, writable any) any {
	seen := make(map[string]struct{})
	out := make(OriginEntries, 0)
	for _, layer := range []any{readOnly, writable} {
		entries, _ := layer.(OriginEntries)
		for _, e := range entries {
			key := entryKey(e)
			if _, dup := seen[key]; dup {
				continue
			}
			seen[key] = struct{}{}
			out = append(out, e)
		}
	}
	return out
}
