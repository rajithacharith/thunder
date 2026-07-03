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

package flowmgt

import (
	"regexp"

	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// handleFormatRegex matches valid handle format:
// - starts with lowercase letter or digit
// - contains only lowercase letters, digits, underscores, or dashes
// - ends with lowercase letter or digit
var handleFormatRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*[a-z0-9]$|^[a-z0-9]$`)

// isValidHandleFormat validates that the handle follows the required format.
func isValidHandleFormat(handle string) bool {
	return handleFormatRegex.MatchString(handle)
}

// isValidFlowType checks if the provided flow type is valid.
func isValidFlowType(flowType providers.FlowType) bool {
	for _, valid := range providers.ValidFlowTypes {
		if flowType == valid {
			return true
		}
	}
	return false
}
