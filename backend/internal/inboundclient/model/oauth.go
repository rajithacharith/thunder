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

// Package model defines OAuth-related types for inbound client configuration.
//
//nolint:lll
package model

import (
	"github.com/thunder-id/thunderid/internal/system/jose/jwe"
	"github.com/thunder-id/thunderid/internal/system/jose/jws"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// InboundAuthType identifies the kind of inbound authentication configured for an entity.
type InboundAuthType string

const (
	// OAuthInboundAuthType is the OAuth 2.0 inbound authentication type.
	OAuthInboundAuthType InboundAuthType = "oauth2"
)

// Supported JOSE algorithms for userinfo responses.
var (
	SupportedUserInfoSigningAlgs = []string{
		string(jws.RS256), string(jws.RS512), string(jws.PS256),
		string(jws.ES256), string(jws.ES384), string(jws.ES512),
		string(jws.EdDSA),
	}
	SupportedUserInfoEncryptionAlgs = []string{string(jwe.RSAOAEP), string(jwe.RSAOAEP256)}
	SupportedUserInfoEncryptionEncs = []string{string(jwe.A128CBCHS256), string(jwe.A256GCM)}
)

// OAuthConfigWithSecret is the wire input shape and the create/update echo response shape.
// Carries ClientSecret (omitempty) so it appears only when freshly issued.
type OAuthConfigWithSecret struct {
	ClientID                           string                            `json:"clientId,omitempty"                 yaml:"clientId,omitempty"                 jsonschema:"OAuth client ID (auto-generated if not provided)"`
	ClientSecret                       string                            `json:"clientSecret,omitempty"             yaml:"clientSecret,omitempty"             jsonschema:"OAuth client secret (auto-generated if not provided)"`
	RedirectURIs                       []string                          `json:"redirectUris,omitempty"             yaml:"redirectUris,omitempty"             jsonschema:"Allowed redirect URIs. Required for Public (SPA/Mobile) and Confidential (Server) clients. Omit for M2M."`
	GrantTypes                         []providers.GrantType             `json:"grantTypes,omitempty"               yaml:"grantTypes,omitempty"               jsonschema:"OAuth grant types. Common: [authorization_code, refresh_token] for user apps, [client_credentials] for M2M."`
	ResponseTypes                      []providers.ResponseType          `json:"responseTypes,omitempty"            yaml:"responseTypes,omitempty"            jsonschema:"OAuth response types. Common: [code] for user apps. Omit for M2M."`
	TokenEndpointAuthMethod            providers.TokenEndpointAuthMethod `json:"tokenEndpointAuthMethod,omitempty"  yaml:"tokenEndpointAuthMethod,omitempty"  jsonschema:"Client authentication method. Use 'none' for Public clients, 'client_secret_basic' for Confidential/M2M."`
	PKCERequired                       bool                              `json:"pkceRequired"                       yaml:"pkceRequired"                       jsonschema:"Require PKCE for security. Recommended for all user-interactive flows."`
	PublicClient                       bool                              `json:"publicClient"                       yaml:"publicClient"                       jsonschema:"Identify if client is public (cannot store secrets). Set true for SPA/Mobile."`
	RequirePushedAuthorizationRequests bool                              `json:"requirePushedAuthorizationRequests" yaml:"requirePushedAuthorizationRequests" jsonschema:"Require Pushed Authorization Requests (PAR) per RFC 9126."`
	DPoPBoundAccessTokens              bool                              `json:"dpopBoundAccessTokens"              yaml:"dpopBoundAccessTokens"              jsonschema:"Require DPoP-bound access tokens (RFC 9449)."`
	IncludeActClaim                    bool                              `json:"includeActClaim"                    yaml:"includeActClaim"                    jsonschema:"Include an implicit on-behalf-of 'act' claim (identifying the application entity) in access tokens issued through this client's authorization code flow. Agents always include it regardless of this setting."`
	Token                              *providers.OAuthTokenConfig       `json:"token,omitempty"                    yaml:"token,omitempty"                    jsonschema:"Token configuration for access tokens and ID tokens"`
	Scopes                             []string                          `json:"scopes,omitempty"                   yaml:"scopes,omitempty"                   jsonschema:"Allowed OAuth scopes. Add custom scopes as needed for your application."`
	UserInfo                           *providers.UserInfoConfig         `json:"userInfo,omitempty"                 yaml:"userInfo,omitempty"                 jsonschema:"UserInfo endpoint configuration. Configure user attributes returned from the OIDC userinfo endpoint."`
	ScopeClaims                        map[string][]string               `json:"scopeClaims,omitempty"              yaml:"scopeClaims,omitempty"              jsonschema:"Scope-to-claims mapping. Maps OAuth scopes to user claims for both ID token and userinfo."`
	Certificate                        *providers.Certificate            `json:"certificate,omitempty"              yaml:"certificate,omitempty"              jsonschema:"Application certificate. Optional. For certificate-based authentication or JWT validation."`
	AcrValues                          []string                          `json:"acrValues,omitempty"                yaml:"acrValues,omitempty"                jsonschema:"Default ACR values applied when the request does not specify acr_values."`
}

// OAuthConfig is the wire output shape (GET responses). ClientSecret is structurally absent.
// Empty slice/map fields are omitted; booleans are always serialized in both JSON and YAML for
// explicit semantics.
type OAuthConfig struct {
	ClientID                           string                            `json:"clientId,omitempty"                 yaml:"clientId,omitempty"`
	RedirectURIs                       []string                          `json:"redirectUris,omitempty"             yaml:"redirectUris,omitempty"`
	GrantTypes                         []providers.GrantType             `json:"grantTypes,omitempty"               yaml:"grantTypes,omitempty"`
	ResponseTypes                      []providers.ResponseType          `json:"responseTypes,omitempty"            yaml:"responseTypes,omitempty"`
	TokenEndpointAuthMethod            providers.TokenEndpointAuthMethod `json:"tokenEndpointAuthMethod,omitempty"  yaml:"tokenEndpointAuthMethod,omitempty"`
	PKCERequired                       bool                              `json:"pkceRequired"                       yaml:"pkceRequired"`
	PublicClient                       bool                              `json:"publicClient"                       yaml:"publicClient"`
	RequirePushedAuthorizationRequests bool                              `json:"requirePushedAuthorizationRequests" yaml:"requirePushedAuthorizationRequests"`
	DPoPBoundAccessTokens              bool                              `json:"dpopBoundAccessTokens"              yaml:"dpopBoundAccessTokens"`
	IncludeActClaim                    bool                              `json:"includeActClaim"                    yaml:"includeActClaim"`
	Token                              *providers.OAuthTokenConfig       `json:"token,omitempty"                    yaml:"token,omitempty"`
	Scopes                             []string                          `json:"scopes,omitempty"                   yaml:"scopes,omitempty"`
	UserInfo                           *providers.UserInfoConfig         `json:"userInfo,omitempty"                 yaml:"userInfo,omitempty"`
	ScopeClaims                        map[string][]string               `json:"scopeClaims,omitempty"              yaml:"scopeClaims,omitempty"`
	Certificate                        *providers.Certificate            `json:"certificate,omitempty"              yaml:"certificate,omitempty"`
	AcrValues                          []string                          `json:"acrValues,omitempty"                yaml:"acrValues,omitempty"`
}

// SupportedIDTokenEncryptionAlgs lists JWE key-management algorithms supported for ID token encryption.
var SupportedIDTokenEncryptionAlgs = []string{string(jwe.RSAOAEP), string(jwe.RSAOAEP256)}

// SupportedIDTokenEncryptionEncs lists JWE content-encryption algorithms supported for ID token encryption.
var SupportedIDTokenEncryptionEncs = []string{string(jwe.A128CBCHS256), string(jwe.A256GCM)}

// InboundAuthConfigWithSecret is the wire input wrapper and create/update echo response wrapper.
type InboundAuthConfigWithSecret struct {
	Type        InboundAuthType        `json:"type"             yaml:"type"             jsonschema:"Inbound authentication type. Use 'oauth2' for OAuth/OIDC applications."`
	OAuthConfig *OAuthConfigWithSecret `json:"config,omitempty" yaml:"config,omitempty" jsonschema:"OAuth/OIDC configuration. Required when type is 'oauth2'. Defines OAuth grant types, redirect URIs, client authentication, and PKCE settings."`
}

// InboundAuthConfig is the wire output wrapper (GET responses).
type InboundAuthConfig struct {
	Type        InboundAuthType `json:"type"             yaml:"type"`
	OAuthConfig *OAuthConfig    `json:"config,omitempty" yaml:"config,omitempty"`
}

// InboundAuthConfigProcessed is the runtime wrapper.
type InboundAuthConfigProcessed struct {
	Type        InboundAuthType        `json:"type"             yaml:"type,omitempty"`
	OAuthConfig *providers.OAuthClient `json:"config,omitempty" yaml:"config,omitempty"`
}
