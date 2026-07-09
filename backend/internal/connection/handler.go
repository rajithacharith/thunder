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

package connection

import (
	"net/http"
	"strings"

	"github.com/thunder-id/thunderid/internal/idp"
	sysutils "github.com/thunder-id/thunderid/internal/system/utils"
	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// handler serves the connection HTTP endpoints. Each vendor file (google.go, ...) defines
// its own typed request/response structs, its toDTO/fromDTO mappers, and thin handler
// methods that delegate the request plumbing to the generic helpers below.
type handler struct {
	svc *service
}

// newHandler creates a new connection handler.
func newHandler(svc *service) *handler {
	return &handler{svc: svc}
}

// createConnection decodes a typed request, maps it to an IdP DTO via the vendor's mapper,
// delegates creation, and writes the encoded response.
func createConnection[Req any, Resp any](h *handler, w http.ResponseWriter, r *http.Request,
	toDTO func(Req) (*providers.IDPDTO, error), fromDTO func(providers.IDPDTO) (Resp, error)) {
	ctx := r.Context()
	req, err := sysutils.DecodeJSONBody[Req](r)
	if err != nil {
		writeInvalidBody(ctx, w)
		return
	}
	dto, err := toDTO(*req)
	if err != nil {
		writeServiceError(ctx, w, &tidcommon.InternalServerError)
		return
	}
	created, svcErr := h.svc.create(ctx, dto)
	if svcErr != nil {
		writeServiceError(ctx, w, svcErr)
		return
	}
	resp, err := fromDTO(*created)
	if err != nil {
		writeServiceError(ctx, w, &tidcommon.InternalServerError)
		return
	}
	sysutils.WriteSuccessResponse(ctx, w, http.StatusCreated, resp)
}

// getConnection fetches an instance of the given type and writes the encoded response.
func getConnection[Resp any](h *handler, w http.ResponseWriter, r *http.Request,
	idpType providers.IDPType, fromDTO func(providers.IDPDTO) (Resp, error)) {
	ctx := r.Context()
	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		writeServiceError(ctx, w, &idp.ErrorInvalidIDPID)
		return
	}
	dto, svcErr := h.svc.getByType(ctx, idpType, id)
	if svcErr != nil {
		writeServiceError(ctx, w, svcErr)
		return
	}
	resp, err := fromDTO(*dto)
	if err != nil {
		writeServiceError(ctx, w, &tidcommon.InternalServerError)
		return
	}
	sysutils.WriteSuccessResponse(ctx, w, http.StatusOK, resp)
}

// updateConnection decodes a typed request, maps it, delegates the update (which preserves
// any secret the request omits), and writes the encoded response.
func updateConnection[Req any, Resp any](h *handler, w http.ResponseWriter, r *http.Request,
	idpType providers.IDPType, toDTO func(Req) (*providers.IDPDTO, error),
	fromDTO func(providers.IDPDTO) (Resp, error)) {
	ctx := r.Context()
	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		writeServiceError(ctx, w, &idp.ErrorInvalidIDPID)
		return
	}
	req, err := sysutils.DecodeJSONBody[Req](r)
	if err != nil {
		writeInvalidBody(ctx, w)
		return
	}
	dto, err := toDTO(*req)
	if err != nil {
		writeServiceError(ctx, w, &tidcommon.InternalServerError)
		return
	}
	updated, svcErr := h.svc.update(ctx, idpType, id, dto)
	if svcErr != nil {
		writeServiceError(ctx, w, svcErr)
		return
	}
	resp, err := fromDTO(*updated)
	if err != nil {
		writeServiceError(ctx, w, &tidcommon.InternalServerError)
		return
	}
	sysutils.WriteSuccessResponse(ctx, w, http.StatusOK, resp)
}

// createHandler binds a vendor's mappers to createConnection, yielding a registerable handler.
func createHandler[Req any, Resp any](h *handler, toDTO func(Req) (*providers.IDPDTO, error),
	fromDTO func(providers.IDPDTO) (Resp, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		createConnection(h, w, r, toDTO, fromDTO)
	}
}

// getHandler binds a vendor's type and mapper to getConnection, yielding a registerable handler.
func getHandler[Resp any](h *handler, idpType providers.IDPType,
	fromDTO func(providers.IDPDTO) (Resp, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getConnection(h, w, r, idpType, fromDTO)
	}
}

// updateHandler binds a vendor's type and mappers to updateConnection, yielding a registerable handler.
func updateHandler[Req any, Resp any](h *handler, idpType providers.IDPType,
	toDTO func(Req) (*providers.IDPDTO, error),
	fromDTO func(providers.IDPDTO) (Resp, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		updateConnection(h, w, r, idpType, toDTO, fromDTO)
	}
}

// listInstances returns a handler that lists the configured instances of a connection type.
func (h *handler) listInstances(idpType providers.IDPType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		instances, svcErr := h.svc.listByType(ctx, idpType)
		if svcErr != nil {
			writeServiceError(ctx, w, svcErr)
			return
		}
		summaries := make([]connectionInstanceSummary, 0, len(instances))
		for _, instance := range instances {
			summaries = append(summaries, connectionInstanceSummary{
				ID:          instance.ID,
				Name:        instance.Name,
				Description: instance.Description,
			})
		}
		sysutils.WriteSuccessResponse(ctx, w, http.StatusOK, summaries)
	}
}

// deleteInstance returns a handler that deletes an instance of a connection type.
func (h *handler) deleteInstance(idpType providers.IDPType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.PathValue("id")
		if strings.TrimSpace(id) == "" {
			writeServiceError(ctx, w, &idp.ErrorInvalidIDPID)
			return
		}
		if svcErr := h.svc.deleteByType(ctx, idpType, id); svcErr != nil {
			writeServiceError(ctx, w, svcErr)
			return
		}
		sysutils.WriteSuccessResponse(ctx, w, http.StatusNoContent, nil)
	}
}

// usagesInstance returns a handler that lists the resources referencing an instance of a
// connection type. Drives the pre-delete confirmation dialog.
func (h *handler) usagesInstance(idpType providers.IDPType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.PathValue("id")
		if strings.TrimSpace(id) == "" {
			writeServiceError(ctx, w, &idp.ErrorInvalidIDPID)
			return
		}
		usages, svcErr := h.svc.usagesByType(ctx, idpType, id)
		if svcErr != nil {
			writeServiceError(ctx, w, svcErr)
			return
		}
		sysutils.WriteSuccessResponse(ctx, w, http.StatusOK, usages)
	}
}

// handleListConnections handles GET /connections, returning the available connection types
// with their configured status and instance count.
func (h *handler) handleListConnections(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	counts, svcErr := h.svc.typeCounts(ctx)
	if svcErr != nil {
		writeServiceError(ctx, w, svcErr)
		return
	}

	summaries := make([]connectionTypeSummary, 0, len(idpBackedVendors))
	for _, vendor := range idpBackedVendors {
		count := counts[vendor.idpType]
		summaries = append(summaries, connectionTypeSummary{
			Type:          vendor.name,
			Configured:    count > 0,
			InstanceCount: count,
		})
	}

	sysutils.WriteSuccessResponse(ctx, w, http.StatusOK, connectionListResponse{Connections: summaries})
}
