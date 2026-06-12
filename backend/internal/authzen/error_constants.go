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

package authzen

import (
	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
	"github.com/thunder-id/thunderid/internal/system/i18n/core"
)

var (
	// ErrorInvalidRequestFormat is returned when the request JSON is malformed.
	ErrorInvalidRequestFormat = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "AZN-1001",
		Error: core.I18nMessage{
			Key:          "error.authzen.invalid_request_format",
			DefaultValue: "Invalid request format",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.authzen.invalid_request_format_description",
			DefaultValue: "The request body is malformed or contains invalid data",
		},
	}
	// ErrorMissingSubject is returned when subject id is not provided.
	ErrorMissingSubject = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "AZN-1002",
		Error: core.I18nMessage{
			Key:          "error.authzen.missing_subject",
			DefaultValue: "Missing subject",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.authzen.missing_subject_description",
			DefaultValue: "Subject id is required",
		},
	}
	// ErrorMissingResource is returned when resource type is not provided.
	ErrorMissingResource = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "AZN-1003",
		Error: core.I18nMessage{
			Key:          "error.authzen.missing_resource",
			DefaultValue: "Missing resource",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.authzen.missing_resource_description",
			DefaultValue: "Resource type is required",
		},
	}
	// ErrorMissingResourceID is returned when resource id is not provided.
	ErrorMissingResourceID = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "AZN-1004",
		Error: core.I18nMessage{
			Key:          "error.authzen.missing_resource_id",
			DefaultValue: "Missing resource id",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.authzen.missing_resource_id_description",
			DefaultValue: "Resource id is required",
		},
	}
	// ErrorMissingAction is returned when action name is not provided.
	ErrorMissingAction = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "AZN-1005",
		Error: core.I18nMessage{
			Key:          "error.authzen.missing_action",
			DefaultValue: "Missing action",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.authzen.missing_action_description",
			DefaultValue: "Action name is required",
		},
	}
	// ErrorMissingEvaluations is returned when batch request has no evaluations.
	ErrorMissingEvaluations = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "AZN-1006",
		Error: core.I18nMessage{
			Key:          "error.authzen.missing_evaluations",
			DefaultValue: "Missing evaluations",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.authzen.missing_evaluations_description",
			DefaultValue: "At least one evaluation is required",
		},
	}
	// ErrorInvalidSubject is returned when subject type does not match subject id.
	ErrorInvalidSubject = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "AZN-1007",
		Error: core.I18nMessage{
			Key:          "error.authzen.invalid_subject",
			DefaultValue: "Invalid subject",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.authzen.invalid_subject_description",
			DefaultValue: "Subject id does not match subject type",
		},
	}
	// ErrorInvalidAction is returned when action is not registered on the resource server.
	ErrorInvalidAction = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "AZN-1008",
		Error: core.I18nMessage{
			Key:          "error.authzen.invalid_action",
			DefaultValue: "Invalid action",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.authzen.invalid_action_description",
			DefaultValue: "Action name is not registered on the resource server",
		},
	}
	// ErrorInvalidResource is returned when resource type does not resolve to a resource server.
	ErrorInvalidResource = serviceerror.ServiceError{
		Type: serviceerror.ClientErrorType,
		Code: "AZN-1009",
		Error: core.I18nMessage{
			Key:          "error.authzen.invalid_resource",
			DefaultValue: "Invalid resource",
		},
		ErrorDescription: core.I18nMessage{
			Key:          "error.authzen.invalid_resource_description",
			DefaultValue: "Resource type is not registered as a resource server",
		},
	}
)
