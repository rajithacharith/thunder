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

package template

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
)

type TemplateServiceTestSuite struct {
	suite.Suite
	mockStore *templateStoreInterfaceMock
	service   TemplateServiceInterface
}

func TestTemplateServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TemplateServiceTestSuite))
}

func (suite *TemplateServiceTestSuite) SetupTest() {
	suite.mockStore = newTemplateStoreInterfaceMock(suite.T())
	suite.service = newTemplateService(suite.mockStore)
}

func (suite *TemplateServiceTestSuite) TestGetTemplateByScenario() {
	dto := &TemplateDTO{ID: "test-1", Scenario: ScenarioUserInvite}
	suite.mockStore.On("GetTemplateByScenario", mock.Anything, ScenarioUserInvite).Return(dto, nil)

	res, err := suite.service.GetTemplateByScenario(context.Background(), ScenarioUserInvite)
	suite.Nil(err)
	suite.Equal("test-1", res.ID)
}

func (suite *TemplateServiceTestSuite) TestRender() {
	dto := &TemplateDTO{
		ID:          "1",
		Scenario:    ScenarioUserInvite,
		Subject:     "Test Invite",
		ContentType: "text/html",
		Body:        "Link: {{ctx(inviteLink)}}",
	}
	suite.mockStore.On("GetTemplateByScenario", mock.Anything, ScenarioUserInvite).Return(dto, nil)

	res, err := suite.service.Render(context.Background(), ScenarioUserInvite,
		TemplateData{"inviteLink": "http://example.com"})
	suite.Nil(err)
	suite.Equal("Test Invite", res.Subject)
	suite.Equal("Link: http://example.com", res.Body)
	suite.True(res.IsHTML)
}

func (suite *TemplateServiceTestSuite) TestRender_NotFound() {
	suite.mockStore.On("GetTemplateByScenario", mock.Anything, ScenarioUserInvite).Return(nil, errTemplateNotFound)

	res, err := suite.service.Render(context.Background(), ScenarioUserInvite, TemplateData{})
	suite.NotNil(err)
	suite.Equal(&ErrorTemplateNotFound, err)
	suite.Nil(res)
}

func (suite *TemplateServiceTestSuite) TestRender_StoreError() {
	storeErr := errors.New("store error")
	suite.mockStore.On("GetTemplateByScenario", mock.Anything, ScenarioUserInvite).Return(nil, storeErr)

	res, err := suite.service.Render(context.Background(), ScenarioUserInvite, TemplateData{})
	suite.NotNil(err)
	suite.Equal(&serviceerror.InternalServerErrorWithI18n, err)
	suite.Nil(res)
}

func (suite *TemplateServiceTestSuite) TestRender_UnknownPlaceholder() {
	dto := &TemplateDTO{
		ID:          "1",
		Scenario:    ScenarioUserInvite,
		Subject:     "Test",
		ContentType: "text/html",
		Body:        "Unknown: {{ctx(unknownKey)}}",
	}
	suite.mockStore.On("GetTemplateByScenario", mock.Anything, ScenarioUserInvite).Return(dto, nil)

	res, err := suite.service.Render(context.Background(), ScenarioUserInvite, TemplateData{})
	suite.Nil(err)
	suite.Equal("Unknown: {{ctx(unknownKey)}}", res.Body)
}
