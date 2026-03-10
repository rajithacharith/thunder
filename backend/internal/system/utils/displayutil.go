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

package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// ExtractDisplayValue extracts a string value from JSON attributes using a dot-notation path.
// Returns an empty string if the path is empty, attributes are nil/empty, or the value cannot be found.
// Non-string values are converted to their string representation.
func ExtractDisplayValue(attributes json.RawMessage, attrPath string) string {
	if len(attributes) == 0 || attrPath == "" {
		return ""
	}

	var data map[string]interface{}
	if err := json.Unmarshal(attributes, &data); err != nil {
		return ""
	}

	parts := strings.Split(attrPath, ".")
	var current interface{} = data

	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return ""
		}
		current = m[part]
		if current == nil {
			return ""
		}
	}

	switch v := current.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%v", v)
	default:
		return ""
	}
}

// DisplayAttributeLookupFunc is a function type for looking up display attribute paths by user type names.
type DisplayAttributeLookupFunc func(ctx context.Context, typeNames []string) (map[string]string, error)

// ResolveDisplayAttributePaths collects unique user types from the provided type-name pairs and
// resolves their display attribute paths using the given lookup function.
// Returns nil if there are no types to resolve or if the lookup fails.
func ResolveDisplayAttributePaths(
	ctx context.Context, userTypes []string, lookupFn DisplayAttributeLookupFunc,
) map[string]string {
	if lookupFn == nil || len(userTypes) == 0 {
		return nil
	}

	seen := make(map[string]struct{})
	var uniqueTypes []string
	for _, t := range userTypes {
		if t != "" {
			if _, ok := seen[t]; !ok {
				seen[t] = struct{}{}
				uniqueTypes = append(uniqueTypes, t)
			}
		}
	}

	if len(uniqueTypes) == 0 {
		return nil
	}

	displayPaths, err := lookupFn(ctx, uniqueTypes)
	if err != nil {
		return nil
	}

	return displayPaths
}

// ResolveUserDisplay resolves a display value for a user from their attributes using
// a schema-configured display attribute path. Falls back to the user ID if no display
// attribute is configured or extraction fails.
func ResolveUserDisplay(id, userType string, attributes json.RawMessage, displayAttrPaths map[string]string) string {
	if displayAttrPaths != nil && userType != "" {
		if path, ok := displayAttrPaths[userType]; ok && path != "" {
			if val := ExtractDisplayValue(attributes, path); val != "" {
				return val
			}
		}
	}

	return id
}
