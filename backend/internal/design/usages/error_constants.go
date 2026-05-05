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

package usages

import (
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	"github.com/asgardeo/thunder/internal/system/i18n/core"
)

// Client errors for design usages operations.
var (
	// ErrorInvalidUsageType is returned when the 'type' query parameter is missing or unrecognised.
	ErrorInvalidUsageType = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "DSU-1001",
		Error: core.I18nMessage{
			Key:          "design.usages.error.invalid_type",
			DefaultValue: "Invalid request format",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "design.usages.error.invalid_type_description",
			DefaultValue: "The 'type' query parameter is required and must be one of: THEME, LAYOUT, FLOW",
		},
	}

	// ErrorMissingResourceID is returned when the 'id' query parameter is absent.
	ErrorMissingResourceID = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "DSU-1002",
		Error: core.I18nMessage{
			Key:          "design.usages.error.missing_id",
			DefaultValue: "Invalid request format",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "design.usages.error.missing_id_description",
			DefaultValue: "The 'id' query parameter is required",
		},
	}

	// ErrorUnsupportedUsageType is returned when the type is valid in structure but not yet supported.
	ErrorUnsupportedUsageType = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "DSU-1003",
		Error: core.I18nMessage{
			Key:          "design.usages.error.unsupported_type",
			DefaultValue: "Unsupported usage type",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "design.usages.error.unsupported_type_description",
			DefaultValue: "The specified usage type is not supported. Supported types are: THEME, LAYOUT, FLOW",
		},
	}

	// ErrorResourceNotFound is returned when the referenced design resource does not exist.
	ErrorResourceNotFound = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "DSU-1004",
		Error: core.I18nMessage{
			Key:          "design.usages.error.resource_not_found",
			DefaultValue: "Resource not found",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "design.usages.error.resource_not_found_description",
			DefaultValue: "The design resource with the specified id does not exist",
		},
	}
)
