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

// Package oauth provides the OAuth authentication executor for handling OAuth-based authentication flows.
package oauth

import (
	"errors"

	authncm "github.com/asgardeo/thunder/internal/authn/common"
	authnoauth "github.com/asgardeo/thunder/internal/authn/oauth"
	"github.com/asgardeo/thunder/internal/executor/oauth/model"
	"github.com/asgardeo/thunder/internal/flow"
	"github.com/asgardeo/thunder/internal/idp"
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	httpservice "github.com/asgardeo/thunder/internal/system/http"
	"github.com/asgardeo/thunder/internal/system/log"
	systemutils "github.com/asgardeo/thunder/internal/system/utils"
)

const loggerComponentName = "OAuthExecutor"

// OAuthExecutorInterface defines the interface for OAuth authentication executors.
type OAuthExecutorInterface interface {
	flow.ExecutorInterface
	BuildAuthorizeFlow(ctx *flow.NodeContext, execResp *flow.ExecutorResponse) error
	ProcessAuthFlowResponse(ctx *flow.NodeContext, execResp *flow.ExecutorResponse) error
	GetOAuthProperties() model.OAuthExecProperties
	GetCallBackURL() string
	GetAuthorizationEndpoint() string
	GetTokenEndpoint() string
	GetUserInfoEndpoint() string
	GetLogoutEndpoint() string
	GetJWKSEndpoint() string
	ExchangeCodeForToken(ctx *flow.NodeContext, execResp *flow.ExecutorResponse,
		code string) (*model.TokenResponse, error)
	GetUserInfo(ctx *flow.NodeContext, execResp *flow.ExecutorResponse,
		accessToken string) (map[string]string, error)
}

// OAuthExecutor implements the OAuthExecutorInterface for handling generic OAuth authentication flows.
type OAuthExecutor struct {
	internal        flow.Executor
	oAuthProperties model.OAuthExecProperties
	authService     authnoauth.OAuthAuthnServiceInterface
}

var _ flow.ExecutorInterface = (*OAuthExecutor)(nil)

// NewOAuthExecutor creates a new instance of OAuthExecutor.
func NewOAuthExecutor(id, name string, defaultInputs []flow.InputData, properties map[string]string,
	oAuthProps *model.OAuthExecProperties) OAuthExecutorInterface {
	endpoints := authnoauth.OAuthEndpoints{
		AuthorizationEndpoint: oAuthProps.AuthorizationEndpoint,
		TokenEndpoint:         oAuthProps.TokenEndpoint,
		UserInfoEndpoint:      oAuthProps.UserInfoEndpoint,
		LogoutEndpoint:        oAuthProps.LogoutEndpoint,
		JwksEndpoint:          oAuthProps.JwksEndpoint,
	}
	authService := authnoauth.NewOAuthAuthnService(
		httpservice.NewHTTPClientWithTimeout(flow.DefaultHTTPTimeout),
		idp.NewIDPService(),
		endpoints,
	)

	return NewOAuthExecutorWithAuthService(id, name, defaultInputs, properties, oAuthProps, authService)
}

// NewOAuthExecutorWithAuthService creates a new instance of OAuthExecutor with a provided
// OAuth authentication service.
// Use this function instead of NewOAuthExecutor when you need to supply a custom OAuth authentication service,
// such as for testing, dependency injection, or when using a specialized implementation.
func NewOAuthExecutorWithAuthService(id, name string, defaultInputs []flow.InputData,
	properties map[string]string, oAuthProps *model.OAuthExecProperties,
	authService authnoauth.OAuthAuthnServiceInterface) OAuthExecutorInterface {
	if len(defaultInputs) == 0 {
		defaultInputs = []flow.InputData{
			{
				Name:     "code",
				Type:     "string",
				Required: true,
			},
		}
	}

	return &OAuthExecutor{
		internal:        *flow.NewExecutor(id, name, defaultInputs, []flow.InputData{}, properties),
		oAuthProperties: *oAuthProps,
		authService:     authService,
	}
}

// GetID returns the ID of the OAuthExecutor.
func (o *OAuthExecutor) GetID() string {
	return o.internal.GetID()
}

// GetName returns the name of the OAuthExecutor.
func (o *OAuthExecutor) GetName() string {
	return o.internal.GetName()
}

// GetProperties returns the properties of the OAuthExecutor.
func (o *OAuthExecutor) GetProperties() flow.ExecutorProperties {
	return o.internal.Properties
}

// GetOAuthProperties returns the OAuth properties of the executor.
func (o *OAuthExecutor) GetOAuthProperties() model.OAuthExecProperties {
	return o.oAuthProperties
}

// GetCallBackURL returns the callback URL for the OAuth authentication.
func (o *OAuthExecutor) GetCallBackURL() string {
	return o.oAuthProperties.RedirectURI
}

// GetAuthorizationEndpoint returns the authorization endpoint of the OAuth authentication.
func (o *OAuthExecutor) GetAuthorizationEndpoint() string {
	return o.oAuthProperties.AuthorizationEndpoint
}

// GetTokenEndpoint returns the token endpoint of the OAuth authentication.
func (o *OAuthExecutor) GetTokenEndpoint() string {
	return o.oAuthProperties.TokenEndpoint
}

// GetUserInfoEndpoint returns the user info endpoint of the OAuth authentication.
func (o *OAuthExecutor) GetUserInfoEndpoint() string {
	return o.oAuthProperties.UserInfoEndpoint
}

// GetLogoutEndpoint returns the logout endpoint of the OAuth authentication.
func (o *OAuthExecutor) GetLogoutEndpoint() string {
	return o.oAuthProperties.LogoutEndpoint
}

// GetJWKSEndpoint returns the JWKs endpoint of the OAuth authentication.
func (o *OAuthExecutor) GetJWKSEndpoint() string {
	return o.oAuthProperties.JwksEndpoint
}

// Execute executes the OAuth authentication flow.
func (o *OAuthExecutor) Execute(ctx *flow.NodeContext) (*flow.ExecutorResponse, error) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, loggerComponentName),
		log.String(log.LoggerKeyExecutorID, o.GetID()),
		log.String(log.LoggerKeyFlowID, ctx.FlowID))
	logger.Debug("Executing OAuth authentication executor")

	execResp := &flow.ExecutorResponse{
		AdditionalData: make(map[string]string),
		RuntimeData:    make(map[string]string),
	}

	// Check if the required input data is provided
	if o.CheckInputData(ctx, execResp) {
		// If required input data is not provided, return incomplete status with redirection.
		logger.Debug("Required input data for OAuth authentication executor is not provided")
		err := o.BuildAuthorizeFlow(ctx, execResp)
		if err != nil {
			return nil, err
		}
	} else {
		err := o.ProcessAuthFlowResponse(ctx, execResp)
		if err != nil {
			return nil, err
		}
	}

	logger.Debug("OAuth authentication executor execution completed",
		log.String("status", string(execResp.Status)),
		log.Bool("isAuthenticated", execResp.AuthenticatedUser.IsAuthenticated))

	return execResp, nil
}

// BuildAuthorizeFlow constructs the redirection to the external OAuth provider for user authentication.
func (o *OAuthExecutor) BuildAuthorizeFlow(ctx *flow.NodeContext, execResp *flow.ExecutorResponse) error {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, loggerComponentName),
		log.String(log.LoggerKeyExecutorID, o.GetID()),
		log.String(log.LoggerKeyFlowID, ctx.FlowID))
	logger.Debug("Initiating OAuth authentication flow")

	authorizeURL, svcErr := o.authService.BuildAuthorizeURL(o.GetID())
	if svcErr != nil {
		if svcErr.Type == serviceerror.ClientErrorType {
			execResp.Status = flow.ExecFailure
			execResp.FailureReason = svcErr.ErrorDescription
			return nil
		}

		logger.Error("Failed to build authorize URL", log.String("errorCode", svcErr.Code),
			log.String("errorDescription", svcErr.ErrorDescription))
		return errors.New("failed to build authorize URL")
	}

	// Set the response to redirect the user to the authorization URL.
	execResp.Status = flow.ExecExternalRedirection
	execResp.RedirectURL = authorizeURL
	execResp.AdditionalData = map[string]string{
		flow.DataIDPName: o.GetName(),
	}

	return nil
}

// ProcessAuthFlowResponse processes the response from the OAuth authentication flow and authenticates the user.
func (o *OAuthExecutor) ProcessAuthFlowResponse(ctx *flow.NodeContext,
	execResp *flow.ExecutorResponse) error {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, loggerComponentName),
		log.String(log.LoggerKeyExecutorID, o.GetID()),
		log.String(log.LoggerKeyFlowID, ctx.FlowID))
	logger.Debug("Processing OAuth authentication response")

	code, ok := ctx.UserInputData["code"]
	if ok && code != "" {
		tokenResp, err := o.ExchangeCodeForToken(ctx, execResp, code)
		if err != nil {
			return err
		}
		if execResp.Status == flow.ExecFailure {
			return nil
		}

		if tokenResp.Scope == "" {
			logger.Error("Scopes are empty in the token response")
			execResp.AuthenticatedUser = authncm.AuthenticatedUser{
				IsAuthenticated: false,
			}
		} else {
			authenticatedUser, err := o.getAuthenticatedUserWithAttributes(ctx, execResp, tokenResp.AccessToken)
			if err != nil {
				return err
			}
			if authenticatedUser == nil {
				return nil
			}
			execResp.AuthenticatedUser = *authenticatedUser
		}
	} else {
		execResp.AuthenticatedUser = authncm.AuthenticatedUser{
			IsAuthenticated: false,
		}
	}

	if execResp.AuthenticatedUser.IsAuthenticated {
		execResp.Status = flow.ExecComplete
	} else if ctx.FlowType != flow.FlowTypeRegistration {
		execResp.Status = flow.ExecFailure
		execResp.FailureReason = "Authentication failed. Authorization code not provided or invalid."
	}

	return nil
}

// GetDefaultExecutorInputs returns the default required input data for the OAuthExecutor.
func (o *OAuthExecutor) GetDefaultExecutorInputs() []flow.InputData {
	return o.internal.GetDefaultExecutorInputs()
}

// CheckInputData checks if the required input data is provided in the context.
func (o *OAuthExecutor) CheckInputData(ctx *flow.NodeContext, execResp *flow.ExecutorResponse) bool {
	if code, ok := ctx.UserInputData["code"]; ok && code != "" {
		return false
	}
	return o.internal.CheckInputData(ctx, execResp)
}

// GetPrerequisites returns the prerequisites for the OAuthExecutor.
func (o *OAuthExecutor) GetPrerequisites() []flow.InputData {
	return o.internal.GetPrerequisites()
}

// ValidatePrerequisites validates whether the prerequisites for the OAuthExecutor are met.
func (o *OAuthExecutor) ValidatePrerequisites(ctx *flow.NodeContext, execResp *flow.ExecutorResponse) bool {
	return o.internal.ValidatePrerequisites(ctx, execResp)
}

// GetUserIDFromContext retrieves the user ID from the context.
func (o *OAuthExecutor) GetUserIDFromContext(ctx *flow.NodeContext) (string, error) {
	return o.internal.GetUserIDFromContext(ctx)
}

// GetRequiredData returns the required input data for the OAuthExecutor.
func (o *OAuthExecutor) GetRequiredData(ctx *flow.NodeContext) []flow.InputData {
	return o.internal.GetRequiredData(ctx)
}

// ExchangeCodeForToken exchanges the authorization code for an access token.
func (o *OAuthExecutor) ExchangeCodeForToken(ctx *flow.NodeContext, execResp *flow.ExecutorResponse,
	code string) (*model.TokenResponse, error) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, loggerComponentName),
		log.String(log.LoggerKeyExecutorID, o.GetID()),
		log.String(log.LoggerKeyFlowID, ctx.FlowID))
	logger.Debug("Exchanging authorization code for a token", log.String("tokenEndpoint", o.GetTokenEndpoint()))

	tokenResp, svcErr := o.authService.ExchangeCodeForToken(o.GetID(), code, true)
	if svcErr != nil {
		if svcErr.Type == serviceerror.ClientErrorType {
			execResp.Status = flow.ExecFailure
			execResp.FailureReason = svcErr.ErrorDescription
			return nil, nil
		}

		logger.Error("Failed to exchange code for a token", log.String("errorCode", svcErr.Code),
			log.String("errorDescription", svcErr.ErrorDescription))
		return nil, errors.New("failed to exchange code for token")
	}

	return &model.TokenResponse{
		AccessToken:  tokenResp.AccessToken,
		TokenType:    tokenResp.TokenType,
		Scope:        tokenResp.Scope,
		RefreshToken: tokenResp.RefreshToken,
		IDToken:      tokenResp.IDToken,
		ExpiresIn:    tokenResp.ExpiresIn,
	}, nil
}

// GetUserInfo fetches user information from the OAuth provider using the access token.
func (o *OAuthExecutor) GetUserInfo(ctx *flow.NodeContext, execResp *flow.ExecutorResponse,
	accessToken string) (map[string]string, error) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, loggerComponentName),
		log.String(log.LoggerKeyExecutorID, o.GetID()),
		log.String(log.LoggerKeyFlowID, ctx.FlowID))
	logger.Debug("Fetching user info from OAuth provider", log.String("userInfoEndpoint", o.GetUserInfoEndpoint()))

	userInfo, svcErr := o.authService.FetchUserInfo(o.GetID(), accessToken)
	if svcErr != nil {
		if svcErr.Type == serviceerror.ClientErrorType {
			execResp.Status = flow.ExecFailure
			execResp.FailureReason = svcErr.ErrorDescription
			return nil, nil
		}

		logger.Error("Failed to fetch user info", log.String("errorCode", svcErr.Code),
			log.String("errorDescription", svcErr.ErrorDescription))
		return nil, errors.New("failed to fetch user information")
	}

	return systemutils.ConvertInterfaceMapToStringMap(userInfo), nil
}

// getAuthenticatedUserWithAttributes retrieves the authenticated user information with additional attributes
// from the OAuth provider using the access token.
func (o *OAuthExecutor) getAuthenticatedUserWithAttributes(ctx *flow.NodeContext,
	execResp *flow.ExecutorResponse, accessToken string) (*authncm.AuthenticatedUser, error) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, loggerComponentName),
		log.String(log.LoggerKeyExecutorID, o.GetID()),
		log.String(log.LoggerKeyFlowID, ctx.FlowID))

	// Get user info using the access token
	userInfo, err := o.GetUserInfo(ctx, execResp, accessToken)
	if err != nil {
		return nil, err
	}
	if execResp.Status == flow.ExecFailure {
		return nil, nil
	}

	// Resolve user with the sub claim.
	sub, ok := userInfo["sub"]
	if !ok || sub == "" {
		execResp.Status = flow.ExecFailure
		execResp.FailureReason = "sub claim not found in the response."
		return nil, nil
	}

	user, svcErr := o.authService.GetInternalUser(sub)
	if svcErr != nil {
		if svcErr.Code == authnoauth.ErrorUserNotFound.Code {
			if ctx.FlowType == flow.FlowTypeRegistration {
				logger.Debug("User not found for the provided sub claim. Proceeding with registration flow.")
				execResp.Status = flow.ExecComplete
				execResp.FailureReason = ""

				if execResp.RuntimeData == nil {
					execResp.RuntimeData = make(map[string]string)
				}
				execResp.RuntimeData["sub"] = sub

				return &authncm.AuthenticatedUser{
					IsAuthenticated: false,
					Attributes:      getUserAttributes(userInfo, ""),
				}, nil
			} else {
				execResp.Status = flow.ExecFailure
				execResp.FailureReason = "User not found"
				return nil, nil
			}
		} else {
			if svcErr.Type == serviceerror.ClientErrorType {
				execResp.Status = flow.ExecFailure
				execResp.FailureReason = svcErr.ErrorDescription
				return nil, nil
			}
			logger.Error("Error while retrieving internal user", log.String("errorCode", svcErr.Code),
				log.String("description", svcErr.ErrorDescription))
			return nil, errors.New("error while retrieving internal user")
		}
	}

	if ctx.FlowType == flow.FlowTypeRegistration {
		// At this point, a unique user is found in the system. Hence fail the execution.
		execResp.Status = flow.ExecFailure
		execResp.FailureReason = "User already exists with the provided sub claim."
		return nil, nil
	}

	if user == nil || user.ID == "" {
		return nil, errors.New("retrieved user is nil or has an empty ID")
	}
	userID := user.ID

	if execResp.Status == flow.ExecFailure {
		return nil, nil
	}

	authenticatedUser := authncm.AuthenticatedUser{
		IsAuthenticated: true,
		UserID:          userID,
		Attributes:      getUserAttributes(userInfo, userID),
	}

	return &authenticatedUser, nil
}

// getUserAttributes extracts user attributes from the user info map, excluding certain keys.
// TODO: Need to convert attributes as per the IDP to local attribute mapping when the support is implemented.
func getUserAttributes(userInfo map[string]string, userID string) map[string]interface{} {
	attributes := make(map[string]interface{})
	for key, value := range userInfo {
		if key != "username" && key != "sub" {
			attributes[key] = value
		}
	}
	if userID != "" {
		attributes["user_id"] = userID
	}

	return attributes
}
