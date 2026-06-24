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

package definition

import (
	"errors"
	"net/http"

	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
	"github.com/thunder-id/thunderid/internal/system/i18n/core"
)

// Internal sentinel errors for the composite presentation-definition store.
var (
	// ErrDefinitionIsImmutable is returned when trying to modify or delete an
	// immutable (file-based) presentation definition.
	ErrDefinitionIsImmutable = errors.New("presentation definition is immutable")

	// ErrResultLimitExceededInCompositeMode is the internal sentinel error returned
	// when composite store results exceed the configured limit.
	ErrResultLimitExceededInCompositeMode = errors.New("result limit exceeded in composite mode")

	// ErrDefinitionDataCorrupted is returned when declarative store data is malformed.
	ErrDefinitionDataCorrupted = errors.New("presentation definition data is corrupted")
)

// Client-facing API errors for the presentation-definition management endpoints.
var (
	// ErrorDefinitionInvalidRequest indicates a malformed create/update request.
	ErrorDefinitionInvalidRequest = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "VP-2001",
		Error: core.I18nMessage{
			Key:          "error.vp.definition_invalid_request",
			DefaultValue: "Invalid request",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.vp.definition_invalid_request_description",
			DefaultValue: "The presentation definition request is missing required fields or is malformed",
		},
	}

	// ErrorDefinitionNotFound indicates the presentation definition does not exist.
	ErrorDefinitionNotFound = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "VP-2002",
		Error: core.I18nMessage{
			Key:          "error.vp.definition_not_found",
			DefaultValue: "Presentation definition not found",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.vp.definition_not_found_description",
			DefaultValue: "No presentation definition exists for the supplied identifier",
		},
	}

	// ErrorDefinitionAlreadyExists indicates the handle is already in use.
	ErrorDefinitionAlreadyExists = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "VP-2003",
		Error: core.I18nMessage{
			Key:          "error.vp.definition_already_exists",
			DefaultValue: "Presentation definition already exists",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.vp.definition_already_exists_description",
			DefaultValue: "A presentation definition with the supplied handle already exists",
		},
	}

	// ErrorDefinitionUnsupportedFormat indicates an unsupported credential format.
	ErrorDefinitionUnsupportedFormat = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "VP-2004",
		Error: core.I18nMessage{
			Key:          "error.vp.definition_unsupported_format",
			DefaultValue: "Unsupported credential format",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.vp.definition_unsupported_format_description",
			DefaultValue: "Only the dc+sd-jwt credential format is supported",
		},
	}

	// ErrorDefinitionImmutable indicates the presentation definition is declarative
	// (file-based) and cannot be modified or deleted via the management API.
	ErrorDefinitionImmutable = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "VP-2005",
		Error: core.I18nMessage{
			Key:          "error.vp.definition_immutable",
			DefaultValue: "Presentation definition is immutable",
		},
		ErrorDescription: core.I18nMessage{
			Key: "error.vp.definition_immutable_description",
			DefaultValue: "The presentation definition is defined in declarative configuration " +
				"and cannot be modified or deleted",
		},
	}

	// ErrorDefinitionResultLimitExceeded indicates the merged composite-store result
	// set exceeds the supported maximum.
	ErrorDefinitionResultLimitExceeded = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "VP-2006",
		Error: core.I18nMessage{
			Key:          "error.vp.definition_result_limit_exceeded",
			DefaultValue: "Result limit exceeded",
		},
		ErrorDescription: core.I18nMessage{
			Key: "error.vp.definition_result_limit_exceeded_description",
			DefaultValue: "The number of presentation definitions exceeds the supported limit in " +
				"hybrid mode. Use search for larger datasets",
		},
	}

	// ErrorDefinitionInvalidOU indicates the organization unit is missing or does not exist.
	ErrorDefinitionInvalidOU = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "VP-2007",
		Error: core.I18nMessage{
			Key:          "error.vp.definition_invalid_ou",
			DefaultValue: "Invalid organization unit",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.vp.definition_invalid_ou_description",
			DefaultValue: "A valid organization unit (ouId or ouHandle) is required",
		},
	}
)

// definitionClientErrorStatus maps a client-facing definition error to its HTTP status.
func definitionClientErrorStatus(code string) int {
	switch code {
	case ErrorDefinitionNotFound.Code:
		return http.StatusNotFound
	case ErrorDefinitionAlreadyExists.Code, ErrorDefinitionImmutable.Code:
		return http.StatusConflict
	default:
		return http.StatusBadRequest
	}
}
