/*
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
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

package userschema

import (
	"errors"

	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	"github.com/asgardeo/thunder/internal/system/i18n/core"
)

// Client errors for user schema management operations.
var (
	// ErrorInvalidRequestFormat is the error returned when the request format is invalid.
	ErrorInvalidRequestFormat = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "USRS-1001",
		Error: core.I18nMessage{
			Key:          "error.userschemaservice.invalid_request_format",
			DefaultValue: "Invalid request format",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.userschemaservice.invalid_request_format_description",
			DefaultValue: "The request body is malformed or contains invalid data",
		},
	}
	// ErrorUserSchemaNotFound is the error returned when a user schema is not found.
	ErrorUserSchemaNotFound = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "USRS-1002",
		Error: core.I18nMessage{
			Key:          "error.userschemaservice.user_schema_not_found",
			DefaultValue: "User schema not found",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.userschemaservice.user_schema_not_found_description",
			DefaultValue: "The user schema with the specified id does not exist",
		},
	}
	// ErrorUserSchemaNameConflict is the error returned when user schema name already exists.
	ErrorUserSchemaNameConflict = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "USRS-1003",
		Error: core.I18nMessage{
			Key:          "error.userschemaservice.user_schema_name_conflict",
			DefaultValue: "User schema name conflict",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.userschemaservice.user_schema_name_conflict_description",
			DefaultValue: "A user schema with the same name already exists",
		},
	}
	// ErrorInvalidUserSchemaRequest is the error returned when user schema request is invalid.
	ErrorInvalidUserSchemaRequest = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "USRS-1004",
		Error: core.I18nMessage{
			Key:          "error.userschemaservice.invalid_user_schema_request",
			DefaultValue: "Invalid user schema request",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.userschemaservice.invalid_user_schema_request_description",
			DefaultValue: "The user schema request contains invalid or missing required fields",
		},
	}
	// ErrorInvalidLimit is the error returned when limit parameter is invalid.
	ErrorInvalidLimit = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "USRS-1005",
		Error: core.I18nMessage{
			Key:          "error.userschemaservice.invalid_limit_parameter",
			DefaultValue: "Invalid pagination parameter",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.userschemaservice.invalid_limit_parameter_description",
			DefaultValue: "The limit parameter must be a positive integer",
		},
	}
	// ErrorInvalidOffset is the error returned when offset parameter is invalid.
	ErrorInvalidOffset = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "USRS-1006",
		Error: core.I18nMessage{
			Key:          "error.userschemaservice.invalid_offset_parameter",
			DefaultValue: "Invalid pagination parameter",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.userschemaservice.invalid_offset_parameter_description",
			DefaultValue: "The offset parameter must be a non-negative integer",
		},
	}
	// ErrorUserValidationFailed is the error returned when user attributes do not conform to the schema.
	ErrorUserValidationFailed = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "USRS-1007",
		Error: core.I18nMessage{
			Key:          "error.userschemaservice.user_validation_failed",
			DefaultValue: "User validation failed",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.userschemaservice.user_validation_failed_description",
			DefaultValue: "User attributes do not conform to the required schema",
		},
	}
	// ErrorCannotModifyDeclarativeResource is the error returned when trying to modify a declarative resource.
	ErrorCannotModifyDeclarativeResource = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "USRS-1008",
		Error: core.I18nMessage{
			Key:          "error.userschemaservice.cannot_modify_declarative_resource",
			DefaultValue: "Cannot modify declarative resource",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.userschemaservice.cannot_modify_declarative_resource_description",
			DefaultValue: "The user schema is declarative and cannot be modified or deleted",
		},
	}
	// ErrorResultLimitExceededInCompositeMode is the error returned when
	// the result limit is exceeded in composite mode.
	ErrorResultLimitExceededInCompositeMode = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "USRS-1009",
		Error: core.I18nMessage{
			Key:          "error.userschemaservice.result_limit_exceeded",
			DefaultValue: "Result limit exceeded",
		},
		ErrorDescription: core.I18nMessage{
			Key: "error.userschemaservice.result_limit_exceeded_description",
			DefaultValue: "The combined result set from both file-based and database " +
				"stores exceeds the maximum limit. Please refine your query to return " +
				"fewer results.",
		},
	}
	// ErrorConsentSyncFailed is the error returned when user schema changes failed to sync with the consent service.
	ErrorConsentSyncFailed = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "USRS-1010",
		Error: core.I18nMessage{
			Key:          "error.userschemaservice.consent_synchronization_failed",
			DefaultValue: "Consent synchronization failed",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.userschemaservice.consent_synchronization_failed_description",
			DefaultValue: "Failed to synchronize consent configurations for the user schema",
		},
	}
	// ErrorInvalidDisplayAttribute is the error returned when the display attribute
	// does not reference a valid top-level attribute in the schema.
	ErrorInvalidDisplayAttribute = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "USRS-1011",
		Error: core.I18nMessage{
			Key:          "error.userschemaservice.invalid_display_attribute",
			DefaultValue: "Invalid display attribute",
		},
		ErrorDescription: core.I18nMessage{
			Key: "error.userschemaservice.invalid_display_attribute_description",
			DefaultValue: "Display attribute must reference an attribute defined in the schema " +
				"(use dot notation for nested attributes, e.g. 'address.city')",
		},
	}
	// ErrorNonDisplayableAttribute is the error returned when the display attribute
	// references an attribute with a non-displayable type (e.g. object or array).
	ErrorNonDisplayableAttribute = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "USRS-1012",
		Error: core.I18nMessage{
			Key:          "error.userschemaservice.non_displayable_attribute_type",
			DefaultValue: "Non-displayable attribute type",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.userschemaservice.non_displayable_attribute_type_description",
			DefaultValue: "Display attribute must reference a string or number type",
		},
	}
	// ErrorCredentialDisplayAttribute is the error returned when the display attribute
	// references an attribute marked as a credential.
	ErrorCredentialDisplayAttribute = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "USRS-1013",
		Error: core.I18nMessage{
			Key:          "error.userschemaservice.credential_attribute_not_allowed_as_display",
			DefaultValue: "Credential attribute not allowed as display",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.userschemaservice.credential_attribute_not_allowed_as_display_description",
			DefaultValue: "Display attribute must not reference a credential attribute",
		},
	}
)

// Server errors for user schema management operations.
var (
	// ErrorInternalServerError is the error returned when an internal server error occurs.
	ErrorInternalServerError = serviceerror.ServiceError{
		Type: serviceerror.ServerErrorType,
		Code: "USRS-5000",
		Error: core.I18nMessage{
			Key:          "error.userschemaservice.internal_server_error",
			DefaultValue: "Internal server error",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.userschemaservice.internal_server_error_description",
			DefaultValue: "An unexpected error occurred while processing the request",
		},
	}
)

// Error variables for user schema operations.
var (
	// ErrUserSchemaNotFound is returned when the user schema is not found in the system.
	ErrUserSchemaNotFound = errors.New("user schema not found")

	// ErrUserSchemaAlreadyExists is returned when a user schema with the same name already exists.
	ErrUserSchemaAlreadyExists = errors.New("user schema already exists")

	// ErrInvalidSchemaDefinition is returned when the schema definition is invalid.
	ErrInvalidSchemaDefinition = errors.New("invalid schema definition")

	// errResultLimitExceededInCompositeMode is returned when the result limit is exceeded in composite mode.
	errResultLimitExceededInCompositeMode = errors.New("result limit exceeded in composite mode")
)
