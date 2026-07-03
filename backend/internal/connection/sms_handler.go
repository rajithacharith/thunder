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

	"github.com/thunder-id/thunderid/internal/notification"
	ncommon "github.com/thunder-id/thunderid/internal/notification/common"
	sysutils "github.com/thunder-id/thunderid/internal/system/utils"
	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
)

// createSMSConnection decodes a typed request, maps it to a notification-sender DTO via the
// vendor's mapper, delegates creation, and writes the encoded response.
func createSMSConnection[Req any, Resp any](h *handler, w http.ResponseWriter, r *http.Request,
	toDTO func(Req) (*ncommon.NotificationSenderDTO, error),
	fromDTO func(ncommon.NotificationSenderDTO) (Resp, error)) {
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
	created, svcErr := h.svc.createSMS(ctx, *dto)
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

// getSMSConnection fetches a message sender of the given provider and writes the encoded response.
func getSMSConnection[Resp any](h *handler, w http.ResponseWriter, r *http.Request,
	provider ncommon.MessageProviderType, fromDTO func(ncommon.NotificationSenderDTO) (Resp, error)) {
	ctx := r.Context()
	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		writeServiceError(ctx, w, &notification.ErrorInvalidSenderID)
		return
	}
	dto, svcErr := h.svc.getSMSByProvider(ctx, provider, id)
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

// updateSMSConnection decodes a typed request, maps it, delegates the update (which preserves
// any secret the request omits), and writes the encoded response.
func updateSMSConnection[Req any, Resp any](h *handler, w http.ResponseWriter, r *http.Request,
	provider ncommon.MessageProviderType, toDTO func(Req) (*ncommon.NotificationSenderDTO, error),
	fromDTO func(ncommon.NotificationSenderDTO) (Resp, error)) {
	ctx := r.Context()
	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		writeServiceError(ctx, w, &notification.ErrorInvalidSenderID)
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
	updated, svcErr := h.svc.updateSMS(ctx, provider, id, *dto)
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

// createSMSHandler binds a vendor's mappers to createSMSConnection, yielding a registerable handler.
func createSMSHandler[Req any, Resp any](h *handler,
	toDTO func(Req) (*ncommon.NotificationSenderDTO, error),
	fromDTO func(ncommon.NotificationSenderDTO) (Resp, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		createSMSConnection(h, w, r, toDTO, fromDTO)
	}
}

// getSMSHandler binds a vendor's provider and mapper to getSMSConnection, yielding a handler.
func getSMSHandler[Resp any](h *handler, provider ncommon.MessageProviderType,
	fromDTO func(ncommon.NotificationSenderDTO) (Resp, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getSMSConnection(h, w, r, provider, fromDTO)
	}
}

// updateSMSHandler binds a vendor's provider and mappers to updateSMSConnection.
func updateSMSHandler[Req any, Resp any](h *handler, provider ncommon.MessageProviderType,
	toDTO func(Req) (*ncommon.NotificationSenderDTO, error),
	fromDTO func(ncommon.NotificationSenderDTO) (Resp, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		updateSMSConnection(h, w, r, provider, toDTO, fromDTO)
	}
}

// listSMSInstances returns a handler that lists the configured senders of a message provider.
func (h *handler) listSMSInstances(provider ncommon.MessageProviderType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		instances, svcErr := h.svc.listSMSByProvider(ctx, provider)
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

// deleteSMSInstance returns a handler that deletes a sender of a message provider.
func (h *handler) deleteSMSInstance(provider ncommon.MessageProviderType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.PathValue("id")
		if strings.TrimSpace(id) == "" {
			writeServiceError(ctx, w, &notification.ErrorInvalidSenderID)
			return
		}
		if svcErr := h.svc.deleteSMSByProvider(ctx, provider, id); svcErr != nil {
			writeServiceError(ctx, w, svcErr)
			return
		}
		sysutils.WriteSuccessResponse(ctx, w, http.StatusNoContent, nil)
	}
}
