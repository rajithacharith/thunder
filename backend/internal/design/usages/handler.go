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
	"net/http"
	"strings"

	"github.com/asgardeo/thunder/internal/system/error/apierror"
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	"github.com/asgardeo/thunder/internal/system/log"
	"github.com/asgardeo/thunder/internal/system/utils"
)

const handlerLogger = "DesignUsageHandler"

// designUsageHandler handles HTTP requests for the design usages endpoint.
type designUsageHandler struct {
	usageService DesignUsageServiceInterface
	logger       *log.Logger
}

// newDesignUsageHandler creates a new designUsageHandler.
func newDesignUsageHandler(usageService DesignUsageServiceInterface) *designUsageHandler {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, handlerLogger))
	return &designUsageHandler{
		usageService: usageService,
		logger:       logger,
	}
}

// HandleUsagesRequest handles GET /design/usages?type=THEME|LAYOUT|FLOW&id={resourceID}.
func (h *designUsageHandler) HandleUsagesRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resourceType := DesignUsageType(strings.ToUpper(r.URL.Query().Get("type")))
	id := r.URL.Query().Get("id")

	response, svcErr := h.usageService.GetResourceUsages(ctx, resourceType, id)
	if svcErr != nil {
		h.handleError(w, svcErr)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, response)

	h.logger.Debug("Successfully resolved design resource usages",
		log.String("type", string(resourceType)),
		log.String("id", id))
}

// handleError maps service errors to HTTP status codes.
func (h *designUsageHandler) handleError(w http.ResponseWriter, svcErr *serviceerror.ServiceError) {
	statusCode := http.StatusInternalServerError
	if svcErr.Type == serviceerror.ClientErrorType {
		switch svcErr.Code {
		case ErrorInvalidUsageType.Code,
			ErrorMissingResourceID.Code,
			ErrorUnsupportedUsageType.Code:
			statusCode = http.StatusBadRequest
		case ErrorResourceNotFound.Code:
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusBadRequest
		}
	}

	errResp := apierror.ErrorResponse{
		Code:        svcErr.Code,
		Message:     svcErr.Error,
		Description: svcErr.ErrorDescription,
	}

	utils.WriteErrorResponse(w, statusCode, errResp)
}
