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

package layoutmgt

import "github.com/asgardeo/thunder/internal/system/error/serviceerror"

var (
	// ErrorInvalidLayoutData is returned when invalid layout data is provided.
	ErrorInvalidLayoutData = serviceerror.ServiceError{
		Type:             serviceerror.ClientErrorType,
		Code:             "LAY-1001",
		Error:            "Invalid layout data",
		ErrorDescription: "The provided layout data is invalid",
	}

	// ErrorInvalidLayoutID is returned when an invalid layout ID is provided.
	ErrorInvalidLayoutID = serviceerror.ServiceError{
		Type:             serviceerror.ClientErrorType,
		Code:             "LAY-1002",
		Error:            "Invalid layout ID",
		ErrorDescription: "The provided layout ID is invalid",
	}

	// ErrorLayoutNotFound is returned when a layout is not found.
	ErrorLayoutNotFound = serviceerror.ServiceError{
		Type:             serviceerror.ClientErrorType,
		Code:             "LAY-1003",
		Error:            "Layout not found",
		ErrorDescription: "The requested layout configuration was not found",
	}

	// ErrorLayoutAlreadyExists is returned when trying to create a layout that already exists.
	ErrorLayoutAlreadyExists = serviceerror.ServiceError{
		Type:             serviceerror.ClientErrorType,
		Code:             "LAY-1004",
		Error:            "Layout already exists",
		ErrorDescription: "A layout with the same ID already exists",
	}

	// ErrorMissingDisplayName is returned when display name is not provided.
	ErrorMissingDisplayName = serviceerror.ServiceError{
		Type:             serviceerror.ClientErrorType,
		Code:             "LAY-1005",
		Error:            "Missing display name",
		ErrorDescription: "Display name is required",
	}

	// ErrorMissingLayout is returned when layout field is not provided.
	ErrorMissingLayout = serviceerror.ServiceError{
		Type:             serviceerror.ClientErrorType,
		Code:             "LAY-1006",
		Error:            "Missing layout",
		ErrorDescription: "Layout field is required",
	}

	// ErrorInvalidLayoutFormat is returned when layout JSON is invalid.
	ErrorInvalidLayoutFormat = serviceerror.ServiceError{
		Type:             serviceerror.ClientErrorType,
		Code:             "LAY-1007",
		Error:            "Invalid layout format",
		ErrorDescription: "Layout must be a valid JSON object",
	}

	// ErrorLayoutInUse is returned when trying to delete a layout that is being used by applications.
	ErrorLayoutInUse = serviceerror.ServiceError{
		Type:             serviceerror.ClientErrorType,
		Code:             "LAY-1008",
		Error:            "Layout in use",
		ErrorDescription: "Cannot delete layout that is currently associated with one or more applications",
	}

	// ErrorInvalidLimitValue is returned when limit validation fails in service layer.
	ErrorInvalidLimitValue = serviceerror.ServiceError{
		Type:             serviceerror.ClientErrorType,
		Code:             "LAY-1009",
		Error:            "Invalid limit",
		ErrorDescription: "Limit value is out of valid range",
	}

	// ErrorInvalidOffsetValue is returned when offset validation fails in service layer.
	ErrorInvalidOffsetValue = serviceerror.ServiceError{
		Type:             serviceerror.ClientErrorType,
		Code:             "LAY-1010",
		Error:            "Invalid offset",
		ErrorDescription: "Offset must be non-negative",
	}

	// ErrorInvalidLimitParam is returned when limit parameter cannot be parsed.
	ErrorInvalidLimitParam = serviceerror.ServiceError{
		Type:             serviceerror.ClientErrorType,
		Code:             "LAY-1011",
		Error:            "Invalid limit",
		ErrorDescription: "Limit must be a valid integer",
	}

	// ErrorInvalidOffsetParam is returned when offset parameter cannot be parsed.
	ErrorInvalidOffsetParam = serviceerror.ServiceError{
		Type:             serviceerror.ClientErrorType,
		Code:             "LAY-1012",
		Error:            "Invalid offset",
		ErrorDescription: "Offset must be a valid integer",
	}
)
