/*
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
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

package granthandlers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	appmodel "github.com/asgardeo/thunder/internal/application/model"
	"github.com/asgardeo/thunder/internal/authz"
	"github.com/asgardeo/thunder/internal/entityprovider"
	"github.com/asgardeo/thunder/internal/oauth/oauth2/constants"
	"github.com/asgardeo/thunder/internal/oauth/oauth2/model"
	"github.com/asgardeo/thunder/internal/oauth/oauth2/tokenservice"
	"github.com/asgardeo/thunder/internal/ou"
	"github.com/asgardeo/thunder/internal/system/config"
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	"github.com/asgardeo/thunder/tests/mocks/authzmock"
	"github.com/asgardeo/thunder/tests/mocks/entityprovidermock"
	"github.com/asgardeo/thunder/tests/mocks/jose/jwtmock"
	"github.com/asgardeo/thunder/tests/mocks/oauth/oauth2/tokenservicemock"
	"github.com/asgardeo/thunder/tests/mocks/oumock"
)

// nolint:gosec // Test token, not a real credential
const testJWTToken = "test-jwt-token-123"
const testResourceURL = "https://mcp.example.com/mcp"

type ClientCredentialsGrantHandlerTestSuite struct {
	suite.Suite
	mockJWTService   *jwtmock.JWTServiceInterfaceMock
	mockTokenBuilder *tokenservicemock.TokenBuilderInterfaceMock
	mockOUService    *oumock.OrganizationUnitServiceInterfaceMock
	mockAuthzService *authzmock.AuthorizationServiceInterfaceMock
	mockEntityProv   *entityprovidermock.EntityProviderInterfaceMock
	handler          *clientCredentialsGrantHandler
	oauthApp         *appmodel.OAuthAppConfigProcessedDTO
}

func TestClientCredentialsGrantHandlerSuite(t *testing.T) {
	suite.Run(t, new(ClientCredentialsGrantHandlerTestSuite))
}

func (suite *ClientCredentialsGrantHandlerTestSuite) SetupTest() {
	// Initialize Thunder Runtime for tests
	testConfig := &config.Config{
		JWT: config.JWTConfig{
			Issuer:         "https://test.thunder.io",
			ValidityPeriod: 3600,
		},
	}
	err := config.InitializeThunderRuntime("", testConfig)
	assert.NoError(suite.T(), err)

	suite.mockJWTService = jwtmock.NewJWTServiceInterfaceMock(suite.T())
	suite.mockTokenBuilder = tokenservicemock.NewTokenBuilderInterfaceMock(suite.T())
	suite.mockOUService = oumock.NewOrganizationUnitServiceInterfaceMock(suite.T())
	suite.mockAuthzService = authzmock.NewAuthorizationServiceInterfaceMock(suite.T())
	suite.mockEntityProv = entityprovidermock.NewEntityProviderInterfaceMock(suite.T())
	suite.handler = &clientCredentialsGrantHandler{
		tokenBuilder: suite.mockTokenBuilder,
		ouService:    suite.mockOUService,
		authzService: suite.mockAuthzService,
		entityProv:   suite.mockEntityProv,
	}
	suite.mockEntityProv.On("GetTransitiveEntityGroups", mock.Anything).
		Return([]entityprovider.EntityGroup{}, nil).Maybe()

	suite.oauthApp = &appmodel.OAuthAppConfigProcessedDTO{
		AppID:                   "app123",
		ClientID:                testClientID,
		RedirectURIs:            []string{"https://example.com/callback"},
		GrantTypes:              []constants.GrantType{constants.GrantTypeClientCredentials},
		ResponseTypes:           []constants.ResponseType{constants.ResponseTypeCode},
		TokenEndpointAuthMethod: constants.TokenEndpointAuthMethodClientSecretBasic,
	}
}

func (suite *ClientCredentialsGrantHandlerTestSuite) TestNewClientCredentialsGrantHandler() {
	handler := newClientCredentialsGrantHandler(
		suite.mockTokenBuilder, suite.mockOUService, suite.mockAuthzService, suite.mockEntityProv)
	assert.NotNil(suite.T(), handler)
	assert.Implements(suite.T(), (*GrantHandlerInterface)(nil), handler)
}

func (suite *ClientCredentialsGrantHandlerTestSuite) TestValidateGrant_Success() {
	tokenRequest := &model.TokenRequest{
		GrantType:    "client_credentials",
		ClientID:     testClientID,
		ClientSecret: "secret123",
		Scope:        "read",
	}

	result := suite.handler.ValidateGrant(context.Background(), tokenRequest, suite.oauthApp)
	assert.Nil(suite.T(), result)
}

func (suite *ClientCredentialsGrantHandlerTestSuite) TestValidateGrant_WrongGrantType() {
	tokenRequest := &model.TokenRequest{
		GrantType:    "authorization_code",
		ClientID:     testClientID,
		ClientSecret: "secret123",
	}

	result := suite.handler.ValidateGrant(context.Background(), tokenRequest, suite.oauthApp)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), constants.ErrorUnsupportedGrantType, result.Error)
	assert.Equal(suite.T(), "Unsupported grant type", result.ErrorDescription)
}

func (suite *ClientCredentialsGrantHandlerTestSuite) TestHandleGrant_Success() {
	testCases := []struct {
		name              string
		scope             string
		expectedJWTClaims map[string]interface{}
		expectedScopes    []string
	}{
		{
			name:              "WithValidScope",
			scope:             "read write",
			expectedJWTClaims: map[string]interface{}{"scope": "read write"},
			expectedScopes:    []string{"read", "write"},
		},
		{
			name:              "WithoutScope",
			scope:             "",
			expectedJWTClaims: map[string]interface{}{},
			expectedScopes:    []string{},
		},
		{
			name:              "WithWhitespaceScope",
			scope:             "   ",
			expectedJWTClaims: map[string]interface{}{},
			expectedScopes:    []string{},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Reset mocks for each test case
			suite.mockJWTService.Mock = mock.Mock{}
			suite.mockAuthzService.Mock = mock.Mock{}

			tokenRequest := &model.TokenRequest{
				GrantType:    "client_credentials",
				ClientID:     testClientID,
				ClientSecret: "secret123",
				Scope:        tc.scope,
			}

			// Mock authz service for non-OIDC scopes
			if len(tc.expectedScopes) > 0 {
				suite.mockAuthzService.On("GetAuthorizedPermissions", mock.Anything,
					authz.GetAuthorizedPermissionsRequest{
						EntityID:             suite.oauthApp.AppID,
						RequestedPermissions: tc.expectedScopes,
					}).Return(&authz.GetAuthorizedPermissionsResponse{
					AuthorizedPermissions: tc.expectedScopes,
				}, nil)
			}

			expectedToken := testJWTToken
			suite.mockTokenBuilder.On("BuildAccessToken",
				mock.MatchedBy(func(ctx *tokenservice.AccessTokenBuildContext) bool {
					return ctx.Subject == testClientID &&
						ctx.Audience == testClientID &&
						ctx.ClientID == testClientID &&
						tokenservice.JoinScopes(ctx.Scopes) == tokenservice.JoinScopes(tc.expectedScopes)
				})).Return(&model.TokenDTO{
				Token:     expectedToken,
				TokenType: constants.TokenTypeBearer,
				IssuedAt:  int64(1234567890),
				ExpiresIn: 3600,
				Scopes:    tc.expectedScopes,
				ClientID:  testClientID,
				Subject:   testClientID,
				Audience:  testClientID,
			}, nil)

			result, errResp := suite.handler.HandleGrant(context.Background(), tokenRequest, suite.oauthApp)

			assert.Nil(t, errResp)
			assert.NotNil(t, result)
			assert.Equal(t, expectedToken, result.AccessToken.Token)
			assert.Equal(t, constants.TokenTypeBearer, result.AccessToken.TokenType)
			assert.Equal(t, int64(3600), result.AccessToken.ExpiresIn)
			assert.Equal(t, tc.expectedScopes, result.AccessToken.Scopes)
			assert.Equal(t, testClientID, result.AccessToken.ClientID)

			// Verify token attributes
			assert.Equal(t, testClientID, result.AccessToken.Subject)
			assert.Equal(t, testClientID, result.AccessToken.Audience)

			suite.mockTokenBuilder.AssertExpectations(t)
		})
	}
}

func (suite *ClientCredentialsGrantHandlerTestSuite) TestHandleGrant_JWTGenerationError() {
	tokenRequest := &model.TokenRequest{
		GrantType:    "client_credentials",
		ClientID:     testClientID,
		ClientSecret: "secret123",
		Scope:        "read",
	}

	suite.mockAuthzService.On("GetAuthorizedPermissions", mock.Anything,
		authz.GetAuthorizedPermissionsRequest{
			EntityID:             suite.oauthApp.AppID,
			RequestedPermissions: []string{"read"},
		}).Return(&authz.GetAuthorizedPermissionsResponse{
		AuthorizedPermissions: []string{"read"},
	}, nil)

	suite.mockTokenBuilder.On("BuildAccessToken", mock.Anything).
		Return(nil, errors.New("JWT generation failed"))

	result, errResp := suite.handler.HandleGrant(context.Background(), tokenRequest, suite.oauthApp)

	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), errResp)
	assert.Equal(suite.T(), constants.ErrorServerError, errResp.Error)
	assert.Equal(suite.T(), "Failed to generate token", errResp.ErrorDescription)

	suite.mockTokenBuilder.AssertExpectations(suite.T())
}

func (suite *ClientCredentialsGrantHandlerTestSuite) TestHandleGrant_NilTokenAttributes() {
	tokenRequest := &model.TokenRequest{
		GrantType:    "client_credentials",
		ClientID:     testClientID,
		ClientSecret: "secret123",
		Scope:        "read",
	}

	suite.mockAuthzService.On("GetAuthorizedPermissions", mock.Anything,
		authz.GetAuthorizedPermissionsRequest{
			EntityID:             suite.oauthApp.AppID,
			RequestedPermissions: []string{"read"},
		}).Return(&authz.GetAuthorizedPermissionsResponse{
		AuthorizedPermissions: []string{"read"},
	}, nil)

	expectedToken := testJWTToken
	suite.mockTokenBuilder.On("BuildAccessToken", mock.MatchedBy(func(ctx *tokenservice.AccessTokenBuildContext) bool {
		return ctx.Subject == testClientID && ctx.Audience == testClientID &&
			tokenservice.JoinScopes(ctx.Scopes) == testScopeRead
	})).Return(&model.TokenDTO{
		Token:     expectedToken,
		TokenType: constants.TokenTypeBearer,
		IssuedAt:  int64(1234567890),
		ExpiresIn: 3600,
		Scopes:    []string{"read"},
		ClientID:  "client123",
		Subject:   testClientID,
		Audience:  testClientID,
	}, nil)

	result, errResp := suite.handler.HandleGrant(context.Background(), tokenRequest, suite.oauthApp)

	assert.Nil(suite.T(), errResp)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedToken, result.AccessToken.Token)

	// Verify token attributes
	assert.Equal(suite.T(), testClientID, result.AccessToken.Subject)
	assert.Equal(suite.T(), testClientID, result.AccessToken.Audience)

	suite.mockTokenBuilder.AssertExpectations(suite.T())
}

func (suite *ClientCredentialsGrantHandlerTestSuite) TestHandleGrant_TokenTimingValidation() {
	tokenRequest := &model.TokenRequest{
		GrantType:    "client_credentials",
		ClientID:     testClientID,
		ClientSecret: "secret123",
		Scope:        "read",
	}

	suite.mockAuthzService.On("GetAuthorizedPermissions", mock.Anything,
		authz.GetAuthorizedPermissionsRequest{
			EntityID:             suite.oauthApp.AppID,
			RequestedPermissions: []string{"read"},
		}).Return(&authz.GetAuthorizedPermissionsResponse{
		AuthorizedPermissions: []string{"read"},
	}, nil)

	expectedToken := testJWTToken
	now := time.Now().Unix()
	suite.mockTokenBuilder.On("BuildAccessToken", mock.Anything).
		Return(&model.TokenDTO{
			Token:     expectedToken,
			TokenType: constants.TokenTypeBearer,
			IssuedAt:  now,
			ExpiresIn: 3600,
			Scopes:    []string{"read"},
			ClientID:  testClientID,
		}, nil)

	startTime := time.Now().Unix()
	result, errResp := suite.handler.HandleGrant(context.Background(), tokenRequest, suite.oauthApp)
	endTime := time.Now().Unix()

	assert.Nil(suite.T(), errResp)
	assert.NotNil(suite.T(), result)

	// Verify the issued time is within reasonable bounds
	assert.GreaterOrEqual(suite.T(), result.AccessToken.IssuedAt, startTime)
	assert.LessOrEqual(suite.T(), result.AccessToken.IssuedAt, endTime)

	suite.mockTokenBuilder.AssertExpectations(suite.T())
}

func (suite *ClientCredentialsGrantHandlerTestSuite) TestHandleGrant_ClientAttributeError() {
	tokenRequest := &model.TokenRequest{
		GrantType:    "client_credentials",
		ClientID:     testClientID,
		ClientSecret: "secret123",
		Scope:        "read",
	}

	oauthAppWithOU := &appmodel.OAuthAppConfigProcessedDTO{
		AppID:                   "app123",
		ClientID:                testClientID,
		OUID:                    "ou-456",
		GrantTypes:              []constants.GrantType{constants.GrantTypeClientCredentials},
		TokenEndpointAuthMethod: constants.TokenEndpointAuthMethodClientSecretBasic,
	}

	suite.mockAuthzService.On("GetAuthorizedPermissions", mock.Anything,
		authz.GetAuthorizedPermissionsRequest{
			EntityID:             oauthAppWithOU.AppID,
			RequestedPermissions: []string{"read"},
		}).Return(&authz.GetAuthorizedPermissionsResponse{
		AuthorizedPermissions: []string{"read"},
	}, nil)

	suite.mockOUService.On("GetOrganizationUnit", context.Background(), "ou-456").Return(
		ou.OrganizationUnit{},
		&serviceerror.ServiceError{Code: "OU-0001", Error: "not found"},
	)

	result, errResp := suite.handler.HandleGrant(context.Background(), tokenRequest, oauthAppWithOU)

	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), errResp)
	assert.Equal(suite.T(), constants.ErrorServerError, errResp.Error)
	assert.Equal(suite.T(), "Failed to generate token", errResp.ErrorDescription)
}

// Resource Parameter Tests (RFC 8707) for Client Credentials Grant

func (suite *ClientCredentialsGrantHandlerTestSuite) TestHandleGrant_WithResourceParameter() {
	tokenRequest := &model.TokenRequest{
		GrantType:    "client_credentials",
		ClientID:     testClientID,
		ClientSecret: "secret123",
		Scope:        "read",
		Resource:     "https://mcp.example.com/mcp",
	}

	suite.mockAuthzService.On("GetAuthorizedPermissions", mock.Anything,
		authz.GetAuthorizedPermissionsRequest{
			EntityID:             suite.oauthApp.AppID,
			RequestedPermissions: []string{"read"},
		}).Return(&authz.GetAuthorizedPermissionsResponse{
		AuthorizedPermissions: []string{"read"},
	}, nil)

	var capturedAudience string
	suite.mockTokenBuilder.On("BuildAccessToken", mock.MatchedBy(func(ctx *tokenservice.AccessTokenBuildContext) bool {
		capturedAudience = ctx.Audience
		return ctx.Subject == testClientID && ctx.Audience == "https://mcp.example.com/mcp"
	})).Return(&model.TokenDTO{
		Token:     testJWTToken,
		TokenType: constants.TokenTypeBearer,
		IssuedAt:  int64(1234567890),
		ExpiresIn: 3600,
		Scopes:    []string{"read"},
		ClientID:  "client123",
		Subject:   testClientID,
		Audience:  "https://mcp.example.com/mcp",
	}, nil)

	result, errResp := suite.handler.HandleGrant(context.Background(), tokenRequest, suite.oauthApp)

	assert.Nil(suite.T(), errResp)
	assert.NotNil(suite.T(), result)

	// Verify resource was included in audience
	assert.Equal(suite.T(), "https://mcp.example.com/mcp", capturedAudience)

	// Verify token attributes use resource as audience
	assert.Equal(suite.T(), "https://mcp.example.com/mcp", result.AccessToken.Audience)
}

func (suite *ClientCredentialsGrantHandlerTestSuite) TestHandleGrant_WithoutResourceParameter() {
	tokenRequest := &model.TokenRequest{
		GrantType:    "client_credentials",
		ClientID:     testClientID,
		ClientSecret: "secret123",
		Scope:        "read",
	}

	suite.mockAuthzService.On("GetAuthorizedPermissions", mock.Anything,
		authz.GetAuthorizedPermissionsRequest{
			EntityID:             suite.oauthApp.AppID,
			RequestedPermissions: []string{"read"},
		}).Return(&authz.GetAuthorizedPermissionsResponse{
		AuthorizedPermissions: []string{"read"},
	}, nil)

	var capturedAudience string
	suite.mockTokenBuilder.On("BuildAccessToken", mock.MatchedBy(func(ctx *tokenservice.AccessTokenBuildContext) bool {
		capturedAudience = ctx.Audience
		return ctx.Subject == testClientID && ctx.Audience == testClientID
	})).Return(&model.TokenDTO{
		Token:     testJWTToken,
		TokenType: constants.TokenTypeBearer,
		IssuedAt:  int64(1234567890),
		ExpiresIn: 3600,
		Scopes:    []string{"read"},
		ClientID:  "client123",
		Subject:   testClientID,
		Audience:  testClientID,
	}, nil)

	result, errResp := suite.handler.HandleGrant(context.Background(), tokenRequest, suite.oauthApp)

	assert.Nil(suite.T(), errResp)
	assert.NotNil(suite.T(), result)

	// Verify default audience (client_id) when no resource parameter
	assert.Equal(suite.T(), testClientID, capturedAudience)

	// Verify token attributes use client ID as audience when no resource
	assert.Equal(suite.T(), testClientID, result.AccessToken.Audience)
}

// App Authorization Integration Tests — verify scope filtering via RBAC roles

func (suite *ClientCredentialsGrantHandlerTestSuite) TestHandleGrant_PartialScopeAuthorization() {
	tokenRequest := &model.TokenRequest{
		GrantType:    "client_credentials",
		ClientID:     testClientID,
		ClientSecret: "secret123",
		Scope:        "read write delete",
	}

	// App is only authorized for "read" and "write" via its role assignments.
	suite.mockAuthzService.On("GetAuthorizedPermissions", mock.Anything,
		authz.GetAuthorizedPermissionsRequest{
			EntityID:             suite.oauthApp.AppID,
			RequestedPermissions: []string{"read", "write", "delete"},
		}).Return(&authz.GetAuthorizedPermissionsResponse{
		AuthorizedPermissions: []string{"read", "write"},
	}, nil)

	suite.mockTokenBuilder.On("BuildAccessToken",
		mock.MatchedBy(func(ctx *tokenservice.AccessTokenBuildContext) bool {
			return tokenservice.JoinScopes(ctx.Scopes) == tokenservice.JoinScopes([]string{"read", "write"})
		})).Return(&model.TokenDTO{
		Token:     testJWTToken,
		TokenType: constants.TokenTypeBearer,
		IssuedAt:  int64(1234567890),
		ExpiresIn: 3600,
		Scopes:    []string{"read", "write"},
		ClientID:  testClientID,
	}, nil)

	result, errResp := suite.handler.HandleGrant(context.Background(), tokenRequest, suite.oauthApp)

	assert.Nil(suite.T(), errResp)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), []string{"read", "write"}, result.AccessToken.Scopes)
}

func (suite *ClientCredentialsGrantHandlerTestSuite) TestHandleGrant_NoAuthorizedScopes() {
	tokenRequest := &model.TokenRequest{
		GrantType:    "client_credentials",
		ClientID:     testClientID,
		ClientSecret: "secret123",
		Scope:        "admin:full",
	}

	// App has no role granting "admin:full".
	suite.mockAuthzService.On("GetAuthorizedPermissions", mock.Anything,
		authz.GetAuthorizedPermissionsRequest{
			EntityID:             suite.oauthApp.AppID,
			RequestedPermissions: []string{"admin:full"},
		}).Return(&authz.GetAuthorizedPermissionsResponse{
		AuthorizedPermissions: []string{},
	}, nil)

	suite.mockTokenBuilder.On("BuildAccessToken",
		mock.MatchedBy(func(ctx *tokenservice.AccessTokenBuildContext) bool {
			return len(ctx.Scopes) == 0
		})).Return(&model.TokenDTO{
		Token:     testJWTToken,
		TokenType: constants.TokenTypeBearer,
		IssuedAt:  int64(1234567890),
		ExpiresIn: 3600,
		Scopes:    []string{},
		ClientID:  testClientID,
	}, nil)

	result, errResp := suite.handler.HandleGrant(context.Background(), tokenRequest, suite.oauthApp)

	assert.Nil(suite.T(), errResp)
	assert.NotNil(suite.T(), result)
	assert.Empty(suite.T(), result.AccessToken.Scopes)
}

func (suite *ClientCredentialsGrantHandlerTestSuite) TestHandleGrant_AuthzServiceError() {
	tokenRequest := &model.TokenRequest{
		GrantType:    "client_credentials",
		ClientID:     testClientID,
		ClientSecret: "secret123",
		Scope:        "read",
	}

	suite.mockAuthzService.On("GetAuthorizedPermissions", mock.Anything,
		authz.GetAuthorizedPermissionsRequest{
			EntityID:             suite.oauthApp.AppID,
			RequestedPermissions: []string{"read"},
		}).Return((*authz.GetAuthorizedPermissionsResponse)(nil),
		&serviceerror.ServiceError{
			Code:  "AUTHZ-0001",
			Error: "authorization check failed",
		})

	result, errResp := suite.handler.HandleGrant(context.Background(), tokenRequest, suite.oauthApp)

	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), errResp)
	assert.Equal(suite.T(), constants.ErrorServerError, errResp.Error)
}

func (suite *ClientCredentialsGrantHandlerTestSuite) TestHandleGrant_EmptyScope_SkipsAuthzCall() {
	tokenRequest := &model.TokenRequest{
		GrantType:    "client_credentials",
		ClientID:     testClientID,
		ClientSecret: "secret123",
		Scope:        "",
	}

	suite.mockTokenBuilder.On("BuildAccessToken",
		mock.MatchedBy(func(ctx *tokenservice.AccessTokenBuildContext) bool {
			return len(ctx.Scopes) == 0
		})).Return(&model.TokenDTO{
		Token:     testJWTToken,
		TokenType: constants.TokenTypeBearer,
		IssuedAt:  int64(1234567890),
		ExpiresIn: 3600,
		Scopes:    []string{},
		ClientID:  testClientID,
	}, nil)

	result, errResp := suite.handler.HandleGrant(context.Background(), tokenRequest, suite.oauthApp)

	assert.Nil(suite.T(), errResp)
	assert.NotNil(suite.T(), result)
	// Verify authz service was NOT called when no scopes requested.
	suite.mockAuthzService.AssertNotCalled(suite.T(), "GetAuthorizedPermissions", mock.Anything, mock.Anything)
}
