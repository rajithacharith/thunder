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

package client

import (
	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
	"github.com/thunder-id/thunderid/internal/system/i18n/core"
)

// ErrorInvalidProvider is the error returned when an invalid provider is specified.
var ErrorInvalidProvider = serviceerror.ServiceError{
	Type: serviceerror.ClientErrorType,
	Code: "MNC-1001",
	Error: core.I18nMessage{
		Key:          "error.notificationclient.unsupported_notification_provider",
		DefaultValue: "Unsupported notification provider",
	},
	ErrorDescription: core.I18nMessage{
		Key:          "error.notificationclient.unsupported_notification_provider.description",
		DefaultValue: "The requested notification provider is not supported.",
	},
}
