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
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/internal/idp"
	"github.com/thunder-id/thunderid/internal/system/resourcedependency"
	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
	"github.com/thunder-id/thunderid/tests/mocks/idp/idpmock"
)

// HandlerTestSuite covers the shared HTTP plumbing (decode, status mapping, empty-id and
// list/delete handlers), exercised through Google as the representative vendor.
type HandlerTestSuite struct {
	suite.Suite
	handler *handler
	mockIDP *idpmock.IDPServiceInterfaceMock
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (s *HandlerTestSuite) SetupTest() {
	s.handler, s.mockIDP = newConnectionTestHandler(s.T())
}

func (s *HandlerTestSuite) TestListConnections() {
	s.mockIDP.On("GetIdentityProviderList", mock.Anything).Return([]idp.BasicIDPDTO{
		{ID: "1", Type: providers.IDPTypeGoogle},
		{ID: "2", Type: providers.IDPTypeGoogle},
		{ID: "3", Type: providers.IDPTypeOIDC},
	}, (*tidcommon.ServiceError)(nil))

	req := httptest.NewRequest(http.MethodGet, "/connections", nil)
	rr := httptest.NewRecorder()
	s.handler.handleListConnections(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	var resp connectionListResponse
	s.Require().NoError(json.NewDecoder(rr.Body).Decode(&resp))

	byType := make(map[string]connectionTypeSummary, len(resp.Connections))
	for _, c := range resp.Connections {
		byType[c.Type] = c
	}
	s.Len(resp.Connections, len(idpBackedVendors))
	s.Equal(2, byType["google"].InstanceCount)
	s.True(byType["google"].Configured)
	s.Equal(1, byType["oidc"].InstanceCount)
	s.Equal(0, byType["github"].InstanceCount)
	s.False(byType["github"].Configured)
}

func (s *HandlerTestSuite) TestListConnectionsServiceError() {
	s.mockIDP.On("GetIdentityProviderList", mock.Anything).
		Return(([]idp.BasicIDPDTO)(nil), &tidcommon.InternalServerError)

	req := httptest.NewRequest(http.MethodGet, "/connections", nil)
	rr := httptest.NewRecorder()
	s.handler.handleListConnections(rr, req)

	s.Equal(http.StatusInternalServerError, rr.Code)
}

func (s *HandlerTestSuite) googleBody() []byte {
	body, _ := json.Marshal(googleConnectionRequest{
		Name: "g", ClientID: "c", ClientSecret: "s", RedirectURI: "https://app/cb",
	})
	return body
}

func (s *HandlerTestSuite) TestCreateInvalidBody() {
	req := httptest.NewRequest(http.MethodPost, "/connections/google", bytes.NewReader([]byte("{bad")))
	rr := httptest.NewRecorder()
	createHandler(s.handler, googleToIDPDTO, googleFromIDPDTO)(rr, req)
	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *HandlerTestSuite) TestCreateServiceErrorConflict() {
	s.mockIDP.On("CreateIdentityProvider", mock.Anything, mock.Anything).
		Return((*providers.IDPDTO)(nil), &idp.ErrorIDPAlreadyExists)

	req := httptest.NewRequest(http.MethodPost, "/connections/google", bytes.NewReader(s.googleBody()))
	rr := httptest.NewRecorder()
	createHandler(s.handler, googleToIDPDTO, googleFromIDPDTO)(rr, req)
	s.Equal(http.StatusConflict, rr.Code)
}

func (s *HandlerTestSuite) TestGetEmptyID() {
	req := httptest.NewRequest(http.MethodGet, "/connections/google/", nil)
	rr := httptest.NewRecorder()
	getHandler(s.handler, providers.IDPTypeGoogle, googleFromIDPDTO)(rr, req)
	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *HandlerTestSuite) TestGetServiceErrorNotFound() {
	s.mockIDP.On("GetIdentityProvider", mock.Anything, "missing").
		Return((*providers.IDPDTO)(nil), &idp.ErrorIDPNotFound)

	req := httptest.NewRequest(http.MethodGet, "/connections/google/missing", nil)
	req.SetPathValue("id", "missing")
	rr := httptest.NewRecorder()
	getHandler(s.handler, providers.IDPTypeGoogle, googleFromIDPDTO)(rr, req)
	s.Equal(http.StatusNotFound, rr.Code)
}

func (s *HandlerTestSuite) TestUpdateEmptyID() {
	req := httptest.NewRequest(http.MethodPut, "/connections/google/", bytes.NewReader(s.googleBody()))
	rr := httptest.NewRecorder()
	updateHandler(s.handler, providers.IDPTypeGoogle, googleToIDPDTO, googleFromIDPDTO)(rr, req)
	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *HandlerTestSuite) TestUpdateInvalidBody() {
	req := httptest.NewRequest(http.MethodPut, "/connections/google/g-1", bytes.NewReader([]byte("{bad")))
	req.SetPathValue("id", "g-1")
	rr := httptest.NewRecorder()
	updateHandler(s.handler, providers.IDPTypeGoogle, googleToIDPDTO, googleFromIDPDTO)(rr, req)
	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *HandlerTestSuite) TestListInstancesServiceError() {
	s.mockIDP.On("GetIdentityProviderList", mock.Anything).
		Return(([]idp.BasicIDPDTO)(nil), &tidcommon.InternalServerError)

	req := httptest.NewRequest(http.MethodGet, "/connections/google", nil)
	rr := httptest.NewRecorder()
	s.handler.listInstances(providers.IDPTypeGoogle)(rr, req)
	s.Equal(http.StatusInternalServerError, rr.Code)
}

func (s *HandlerTestSuite) TestDeleteEmptyID() {
	req := httptest.NewRequest(http.MethodDelete, "/connections/google/", nil)
	rr := httptest.NewRecorder()
	s.handler.deleteInstance(providers.IDPTypeGoogle)(rr, req)
	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *HandlerTestSuite) TestUsagesEmptyID() {
	req := httptest.NewRequest(http.MethodGet, "/connections/google//usages", nil)
	rr := httptest.NewRecorder()
	s.handler.usagesInstance(providers.IDPTypeGoogle)(rr, req)
	s.Equal(http.StatusBadRequest, rr.Code)
	s.mockIDP.AssertNotCalled(s.T(), "GetIDPUsages", mock.Anything, mock.Anything)
}

func (s *HandlerTestSuite) TestUsagesServiceError() {
	s.mockIDP.On("GetIdentityProvider", mock.Anything, "missing").
		Return((*providers.IDPDTO)(nil), &idp.ErrorIDPNotFound)

	req := httptest.NewRequest(http.MethodGet, "/connections/google/missing/usages", nil)
	req.SetPathValue("id", "missing")
	rr := httptest.NewRecorder()
	s.handler.usagesInstance(providers.IDPTypeGoogle)(rr, req)
	s.Equal(http.StatusNotFound, rr.Code)
}

func (s *HandlerTestSuite) TestUsagesSuccess() {
	total := 1
	usages := &resourcedependency.DependenciesResponse{
		TotalResults: &total,
		Count:        1,
		Summary:      map[string]int{"flow": 1},
		Usages: []resourcedependency.ResourceDependency{
			{ResourceType: "flow", ID: "flow-1", DisplayName: "Login Flow", BehaviorOnDelete: "restrict"},
		},
	}
	s.mockIDP.On("GetIdentityProvider", mock.Anything, "g-1").
		Return(&providers.IDPDTO{ID: "g-1", Type: providers.IDPTypeGoogle}, (*tidcommon.ServiceError)(nil))
	s.mockIDP.On("GetIDPUsages", mock.Anything, "g-1").Return(usages, (*tidcommon.ServiceError)(nil))

	req := httptest.NewRequest(http.MethodGet, "/connections/google/g-1/usages", nil)
	req.SetPathValue("id", "g-1")
	rr := httptest.NewRecorder()
	s.handler.usagesInstance(providers.IDPTypeGoogle)(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	var resp resourcedependency.DependenciesResponse
	s.Require().NoError(json.NewDecoder(rr.Body).Decode(&resp))
	s.Require().NotNil(resp.TotalResults)
	s.Equal(1, *resp.TotalResults)
	s.Require().Len(resp.Usages, 1)
	s.Equal("flow-1", resp.Usages[0].ID)
	s.Equal("restrict", resp.Usages[0].BehaviorOnDelete)
}
