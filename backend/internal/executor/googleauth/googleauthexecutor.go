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

// Package googleauth provides the Google OIDC authentication executor.
package googleauth

import (
	"errors"

	authngoogle "github.com/asgardeo/thunder/internal/authn/google"
	authnoauth "github.com/asgardeo/thunder/internal/authn/oauth"
	authnoidc "github.com/asgardeo/thunder/internal/authn/oidc"
	"github.com/asgardeo/thunder/internal/executor/oauth/model"
	"github.com/asgardeo/thunder/internal/executor/oidcauth"
	"github.com/asgardeo/thunder/internal/flow"
	"github.com/asgardeo/thunder/internal/idp"
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	httpservice "github.com/asgardeo/thunder/internal/system/http"
	"github.com/asgardeo/thunder/internal/system/log"
)

const loggerComponentName = "GoogleOIDCAuthExecutor"

// GoogleOIDCAuthExecutor implements the OIDC authentication executor for Google.
type GoogleOIDCAuthExecutor struct {
	*oidcauth.OIDCAuthExecutor
	googleAuthService authngoogle.GoogleOIDCAuthnServiceInterface
}

// CheckInputData implements flow.ExecutorInterface.
// Subtle: this method shadows the method (*OIDCAuthExecutor).CheckInputData of GoogleOIDCAuthExecutor.OIDCAuthExecutor.
func (g *GoogleOIDCAuthExecutor) CheckInputData(ctx *flow.NodeContext, execResp *flow.ExecutorResponse) bool {
	panic("unimplemented")
}

// GetDefaultExecutorInputs implements flow.ExecutorInterface.
// Subtle: this method shadows the method (*OIDCAuthExecutor).GetDefaultExecutorInputs of GoogleOIDCAuthExecutor.OIDCAuthExecutor.
func (g *GoogleOIDCAuthExecutor) GetDefaultExecutorInputs() []flow.InputData {
	panic("unimplemented")
}

// GetID implements flow.ExecutorInterface.
// Subtle: this method shadows the method (*OIDCAuthExecutor).GetID of GoogleOIDCAuthExecutor.OIDCAuthExecutor.
func (g *GoogleOIDCAuthExecutor) GetID() string {
	panic("unimplemented")
}

// GetName implements flow.ExecutorInterface.
// Subtle: this method shadows the method (*OIDCAuthExecutor).GetName of GoogleOIDCAuthExecutor.OIDCAuthExecutor.
func (g *GoogleOIDCAuthExecutor) GetName() string {
	panic("unimplemented")
}

// GetPrerequisites implements flow.ExecutorInterface.
// Subtle: this method shadows the method (*OIDCAuthExecutor).GetPrerequisites of GoogleOIDCAuthExecutor.OIDCAuthExecutor.
func (g *GoogleOIDCAuthExecutor) GetPrerequisites() []flow.InputData {
	panic("unimplemented")
}

// GetProperties implements flow.ExecutorInterface.
// Subtle: this method shadows the method (*OIDCAuthExecutor).GetProperties of GoogleOIDCAuthExecutor.OIDCAuthExecutor.
func (g *GoogleOIDCAuthExecutor) GetProperties() flow.ExecutorProperties {
	panic("unimplemented")
}

// GetRequiredData implements flow.ExecutorInterface.
// Subtle: this method shadows the method (*OIDCAuthExecutor).GetRequiredData of GoogleOIDCAuthExecutor.OIDCAuthExecutor.
func (g *GoogleOIDCAuthExecutor) GetRequiredData(ctx *flow.NodeContext) []flow.InputData {
	panic("unimplemented")
}

// GetUserIDFromContext implements flow.ExecutorInterface.
// Subtle: this method shadows the method (*OIDCAuthExecutor).GetUserIDFromContext of GoogleOIDCAuthExecutor.OIDCAuthExecutor.
func (g *GoogleOIDCAuthExecutor) GetUserIDFromContext(ctx *flow.NodeContext) (string, error) {
	panic("unimplemented")
}

// ValidatePrerequisites implements flow.ExecutorInterface.
// Subtle: this method shadows the method (*OIDCAuthExecutor).ValidatePrerequisites of GoogleOIDCAuthExecutor.OIDCAuthExecutor.
func (g *GoogleOIDCAuthExecutor) ValidatePrerequisites(ctx *flow.NodeContext, execResp *flow.ExecutorResponse) bool {
	panic("unimplemented")
}

var _ flow.ExecutorInterface = (*GoogleOIDCAuthExecutor)(nil)

// NewGoogleOIDCAuthExecutorFromProps creates a new instance of GoogleOIDCAuthExecutor with the provided properties.
func NewGoogleOIDCAuthExecutorFromProps(execProps flow.ExecutorProperties,
	oAuthProps *model.BasicOAuthExecProperties) oidcauth.OIDCAuthExecutorInterface {
	// Prepare the complete OAuth properties for Google
	compOAuthProps := &model.OAuthExecProperties{
		AuthorizationEndpoint: authngoogle.AuthorizeEndpoint,
		TokenEndpoint:         authngoogle.TokenEndpoint,
		UserInfoEndpoint:      authngoogle.UserInfoEndpoint,
		JwksEndpoint:          authngoogle.JwksEndpoint,
		ClientID:              oAuthProps.ClientID,
		ClientSecret:          oAuthProps.ClientSecret,
		RedirectURI:           oAuthProps.RedirectURI,
		Scopes:                oAuthProps.Scopes,
		AdditionalParams:      oAuthProps.AdditionalParams,
	}
	endpoints := authnoauth.OAuthEndpoints{
		AuthorizationEndpoint: compOAuthProps.AuthorizationEndpoint,
		TokenEndpoint:         compOAuthProps.TokenEndpoint,
		UserInfoEndpoint:      compOAuthProps.UserInfoEndpoint,
	}

	defaultInputs := []flow.InputData{
		{
			Name:     "code",
			Type:     "string",
			Required: true,
		},
		{
			Name:     "nonce",
			Type:     "string",
			Required: false,
		},
	}

	oAuthSvc := authnoauth.NewOAuthAuthnService(
		httpservice.NewHTTPClientWithTimeout(flow.DefaultHTTPTimeout),
		idp.NewIDPService(),
		endpoints,
	)
	oidcAuthSvc := authnoidc.NewOIDCAuthnService(oAuthSvc, nil)
	authSvc := authngoogle.NewGoogleOIDCAuthnService(oidcAuthSvc)

	base := oidcauth.NewOIDCAuthExecutor("google_oidc_auth_executor", execProps.Name,
		defaultInputs, execProps.Properties, compOAuthProps)
	exec, ok := base.(*oidcauth.OIDCAuthExecutor)
	if !ok {
		panic("failed to cast GoogleOIDCAuthExecutor to OIDCAuthExecutor")
	}
	return &GoogleOIDCAuthExecutor{
		OIDCAuthExecutor:  exec,
		googleAuthService: authSvc,
	}
}

// NewGoogleOIDCAuthExecutor creates a new instance of GoogleOIDCAuthExecutor with the provided details.
func NewGoogleOIDCAuthExecutor(id, name string, properties map[string]string,
	clientID, clientSecret, redirectURI string, scopes []string,
	additionalParams map[string]string) oidcauth.OIDCAuthExecutorInterface {
	// Prepare the OAuth properties for Google
	oAuthProps := &model.OAuthExecProperties{
		AuthorizationEndpoint: authngoogle.AuthorizeEndpoint,
		TokenEndpoint:         authngoogle.TokenEndpoint,
		UserInfoEndpoint:      authngoogle.UserInfoEndpoint,
		JwksEndpoint:          authngoogle.JwksEndpoint,
		ClientID:              clientID,
		ClientSecret:          clientSecret,
		RedirectURI:           redirectURI,
		Scopes:                scopes,
		AdditionalParams:      additionalParams,
	}
	endpoints := authnoauth.OAuthEndpoints{
		AuthorizationEndpoint: oAuthProps.AuthorizationEndpoint,
		TokenEndpoint:         oAuthProps.TokenEndpoint,
		UserInfoEndpoint:      oAuthProps.UserInfoEndpoint,
	}

	defaultInputs := []flow.InputData{
		{
			Name:     "code",
			Type:     "string",
			Required: true,
		},
		{
			Name:     "nonce",
			Type:     "string",
			Required: false,
		},
	}

	oAuthSvc := authnoauth.NewOAuthAuthnService(
		httpservice.NewHTTPClientWithTimeout(flow.DefaultHTTPTimeout),
		idp.NewIDPService(),
		endpoints,
	)
	oidcAuthSvc := authnoidc.NewOIDCAuthnService(oAuthSvc, nil)
	authSvc := authngoogle.NewGoogleOIDCAuthnService(oidcAuthSvc)

	base := oidcauth.NewOIDCAuthExecutor(id, name, defaultInputs, properties, oAuthProps)
	exec, ok := base.(*oidcauth.OIDCAuthExecutor)
	if !ok {
		panic("failed to cast GoogleOIDCAuthExecutor to OIDCAuthExecutor")
	}
	return &GoogleOIDCAuthExecutor{
		OIDCAuthExecutor:  exec,
		googleAuthService: authSvc,
	}
}

// Execute executes the Google OIDC authentication flow.
func (g *GoogleOIDCAuthExecutor) Execute(ctx *flow.NodeContext) (*flow.ExecutorResponse, error) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, loggerComponentName))
	logger.Debug("Executing Google OIDC auth executor",
		log.String("executorID", g.GetID()), log.String("flowID", ctx.FlowID))

	execResp := &flow.ExecutorResponse{
		AdditionalData: make(map[string]string),
		RuntimeData:    make(map[string]string),
	}

	if g.CheckInputData(ctx, execResp) {
		logger.Debug("Required input data for Google OIDC auth executor is not provided")
		err := g.BuildAuthorizeFlow(ctx, execResp)
		if err != nil {
			return nil, err
		}

		logger.Debug("Google OIDC auth executor execution completed",
			log.String("status", string(execResp.Status)))
	} else {
		err := g.ProcessAuthFlowResponse(ctx, execResp)
		if err != nil {
			return nil, err
		}

		logger.Debug("Google OIDC auth executor execution completed",
			log.String("status", string(execResp.Status)),
			log.Bool("isAuthenticated", execResp.AuthenticatedUser.IsAuthenticated))
	}

	return execResp, nil
}

// ValidateIDToken validates the ID token received from Google.
func (g *GoogleOIDCAuthExecutor) ValidateIDToken(execResp *flow.ExecutorResponse, idToken string) error {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, loggerComponentName))
	logger.Debug("Validating ID token")

	svcErr := g.googleAuthService.ValidateIDToken(g.GetID(), idToken)
	if svcErr != nil {
		if svcErr.Type == serviceerror.ClientErrorType {
			execResp.Status = flow.ExecFailure
			execResp.FailureReason = svcErr.ErrorDescription
			return nil
		}

		logger.Error("Failed to validate ID token", log.String("errorCode", svcErr.Code),
			log.String("errorDescription", svcErr.ErrorDescription))
		return errors.New("failed to validate ID token")
	}

	return nil
}
