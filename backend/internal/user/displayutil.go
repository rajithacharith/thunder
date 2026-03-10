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

package user

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/asgardeo/thunder/internal/system/log"
	"github.com/asgardeo/thunder/internal/system/utils"
	"github.com/asgardeo/thunder/internal/userschema"
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

// ResolveDisplayAttributePaths collects unique user types and resolves their display
// attribute paths from the user schema service.
// Returns nil if there are no types to resolve or if the lookup fails.
func ResolveDisplayAttributePaths(
	ctx context.Context, userTypes []string, schemaService userschema.UserSchemaServiceInterface,
	logger *log.Logger,
) map[string]string {
	if schemaService == nil || len(userTypes) == 0 {
		return nil
	}

	uniqueTypes := utils.UniqueNonEmptyStrings(userTypes)
	if len(uniqueTypes) == 0 {
		return nil
	}

	displayPaths, svcErr := schemaService.GetDisplayAttributesByNames(ctx, uniqueTypes)
	if svcErr != nil {
		if logger != nil {
			logger.Warn("Failed to resolve display attribute paths, skipping display resolution",
				log.Any("error", svcErr))
		}
		return nil
	}

	return displayPaths
}

// ResolveUserDisplay resolves a display value for a user from their attributes using
// a schema-configured display attribute path. Falls back to the user ID if no display
// attribute is configured or extraction fails.
func ResolveUserDisplay(id, userType string, attributes json.RawMessage, displayAttrPaths map[string]string) string {
	if displayAttrPaths == nil || userType == "" {
		return id
	}

	path, ok := displayAttrPaths[userType]
	if !ok || path == "" {
		return id
	}

	value := ExtractDisplayValue(attributes, path)
	if value == "" {
		return id
	}

	return value
}
