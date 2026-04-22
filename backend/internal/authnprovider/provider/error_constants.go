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

package provider

import (
	authnprovidercm "github.com/asgardeo/thunder/internal/authnprovider/common"
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	"github.com/asgardeo/thunder/internal/system/i18n/core"
)

var errorSystemError = serviceerror.ServiceError{
	Type: serviceerror.ServerErrorType,
	Code: authnprovidercm.ErrorCodeSystemError,
	Error: core.I18nMessage{
		Key:          "error.authnproviderservice.system_error",
		DefaultValue: "System error",
	},
	ErrorDescription: core.I18nMessage{
		Key:          "error.authnproviderservice.system_error_description",
		DefaultValue: "An internal server error occurred",
	},
}
