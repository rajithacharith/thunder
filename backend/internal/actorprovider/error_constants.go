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

package actorprovider

import (
	"errors"

	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
	"github.com/thunder-id/thunderid/internal/system/i18n/core"
)

// ErrActorNotFound indicates the requested actor or inbound-client row does not exist.
var ErrActorNotFound = errors.New("actor not found")

// ErrorActorNotFound is returned when the requested actor cannot be resolved.
var ErrorActorNotFound = serviceerror.ServiceError{
	Type: serviceerror.ClientErrorType,
	Code: "SSE-4041",
	Error: core.I18nMessage{
		Key:          "error.actor_not_found",
		DefaultValue: "Actor not found",
	},
	ErrorDescription: core.I18nMessage{
		Key:          "error.actor_not_found_description",
		DefaultValue: "The requested actor does not exist",
	},
}

// ErrorActorFetchFailed is returned when actor resolution fails unexpectedly.
var ErrorActorFetchFailed = serviceerror.ServiceError{
	Type: serviceerror.ServerErrorType,
	Code: "SSE-5002",
	Error: core.I18nMessage{
		Key:          "error.actor_fetch_failed",
		DefaultValue: "Failed to fetch actor",
	},
	ErrorDescription: core.I18nMessage{
		Key:          "error.actor_fetch_failed_description",
		DefaultValue: "An unexpected error occurred while resolving the actor",
	},
}
