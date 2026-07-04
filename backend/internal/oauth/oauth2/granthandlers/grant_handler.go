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

// Package granthandlers provides an interface and implementations for handling OAuth 2.0 grant types.
package granthandlers

import (
	"context"
	"errors"

	"github.com/thunder-id/thunderid/internal/oauth/oauth2/constants"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/model"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/revocation"
	"github.com/thunder-id/thunderid/internal/system/log"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// GrantHandlerInterface defines the interface for handling OAuth 2.0 grants.
type GrantHandlerInterface interface {
	ValidateGrant(
		ctx context.Context,
		tokenRequest *model.TokenRequest,
		oauthApp *providers.OAuthClient,
	) *model.ErrorResponse
	HandleGrant(ctx context.Context, tokenRequest *model.TokenRequest, oauthApp *providers.OAuthClient) (
		*model.TokenResponseDTO, *model.ErrorResponse)
}

// enforceRevocation runs the single-token deny-list check for jti and maps the outcome to an OAuth
// error response shared by the grant handlers. It returns nil when the token may proceed, revokedErr
// (the caller's grant-specific error) when the token is on the deny list, and a fail-closed
// server_error when the deny list cannot be consulted.
func enforceRevocation(
	ctx context.Context,
	enforcementService revocation.EnforcementServiceInterface,
	jti string,
	revokedErr *model.ErrorResponse,
	logger *log.Logger,
) *model.ErrorResponse {
	switch err := enforcementService.EnsureNotRevoked(ctx, jti); {
	case err == nil:
		return nil
	case errors.Is(err, revocation.ErrTokenRevoked):
		return revokedErr
	default:
		logger.Error(ctx, "Token revocation status could not be verified; failing closed",
			log.Error(err))
		return &model.ErrorResponse{
			Error:            constants.ErrorServerError,
			ErrorDescription: "Token revocation status could not be verified",
		}
	}
}

// RefreshTokenGrantHandlerInterface defines the interface for handling refresh token grants.
type RefreshTokenGrantHandlerInterface interface {
	GrantHandlerInterface
	IssueRefreshToken(
		ctx context.Context,
		tokenResponse *model.TokenResponseDTO,
		oauthApp *providers.OAuthClient,
		subject string, audiences []string, grantType string,
		scopes []string,
		claimsRequest *model.ClaimsRequest,
		claimsLocales string,
		attributeCacheID string,
	) *model.ErrorResponse
}
