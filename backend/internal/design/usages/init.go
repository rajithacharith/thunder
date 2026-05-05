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

	"github.com/asgardeo/thunder/internal/system/middleware"
)

// Initialize creates and wires the design usages service and registers its HTTP routes.
func Initialize(
	mux *http.ServeMux,
	resolver ApplicationUsageResolver,
	existenceChecker ResourceExistenceChecker,
) DesignUsageServiceInterface {
	svc := newDesignUsageService(resolver, existenceChecker)
	h := newDesignUsageHandler(svc)
	registerRoutes(mux, h)
	return svc
}

func registerRoutes(mux *http.ServeMux, h *designUsageHandler) {
	opts := middleware.CORSOptions{
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   middleware.DefaultAllowedHeaders,
		AllowCredentials: true,
		MaxAge:           600,
	}
	mux.HandleFunc(middleware.WithCORS("GET /design/usages", h.HandleUsagesRequest, opts))
	mux.HandleFunc(middleware.WithCORS("OPTIONS /design/usages", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}, opts))
}
