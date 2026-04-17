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

package manager

import "github.com/asgardeo/thunder/internal/system/error/serviceerror"

var (
	// ErrorAuthenticationFailed is returned when the underlying provider rejects the authentication
	// attempt due to a client-side reason (e.g. invalid credentials).
	ErrorAuthenticationFailed = serviceerror.ServiceError{
		Type:             serviceerror.ClientErrorType,
		Code:             "AUTHN-MGR-1001",
		Error:            "Authentication failed",
		ErrorDescription: "The authentication attempt failed",
	}

	// ErrorGetAttributesClientError is returned when the underlying provider rejects the
	// attribute fetch due to a client-side reason (e.g. invalid or expired token).
	ErrorGetAttributesClientError = serviceerror.ServiceError{
		Type:             serviceerror.ClientErrorType,
		Code:             "AUTHN-MGR-1004",
		Error:            "Failed to get attributes",
		ErrorDescription: "The attribute fetch was rejected by the provider",
	}

	// ErrorUserNotFound is returned when the underlying provider indicates no user was found
	// matching the provided identifiers.
	ErrorUserNotFound = serviceerror.ServiceError{
		Type:             serviceerror.ClientErrorType,
		Code:             "AUTHN-MGR-1007",
		Error:            "User not found",
		ErrorDescription: "No user found matching the provided identifiers",
	}

	// ErrorInvalidRequest is returned when the underlying provider rejects the authentication
	// request as invalid.
	ErrorInvalidRequest = serviceerror.ServiceError{
		Type:             serviceerror.ClientErrorType,
		Code:             "AUTHN-MGR-1008",
		Error:            "Invalid request",
		ErrorDescription: "The authentication request is invalid",
	}
)
