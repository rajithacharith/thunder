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
	"fmt"
	"strings"
	"time"
)

const dbTimeFormat = "2006-01-02 15:04:05.999999999"

// ParseDBTimeField parses a time value read from a database column.
// Accepts a time.Time (returned as-is) or a string in SQLite datetime format,
// normalising to UTC. Falls back to ISO 8601 if the primary format does not match.
func ParseDBTimeField(field interface{}, fieldName string) (time.Time, error) {
	switch v := field.(type) {
	case string:
		parts := strings.SplitN(v, " ", 3)
		trimmed := v
		if len(parts) >= 2 {
			trimmed = parts[0] + " " + parts[1]
		}
		if t, err := time.Parse(dbTimeFormat, trimmed); err == nil {
			return t.UTC(), nil
		}
		t, err := time.Parse("2006-01-02T15:04:05Z07:00", v)
		if err != nil {
			return time.Time{}, fmt.Errorf("error parsing %s: %w", fieldName, err)
		}
		return t.UTC(), nil
	case time.Time:
		return v.UTC(), nil
	default:
		return time.Time{}, fmt.Errorf("unexpected type for %s", fieldName)
	}
}
