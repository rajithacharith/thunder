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

package model

import (
	"errors"

	"github.com/asgardeo/thunder/internal/system/jose/jwe"
	"github.com/asgardeo/thunder/internal/system/jose/jws"
)

// InboundAuthType represents the type of inbound authentication.
type InboundAuthType string

const (
	// OAuthInboundAuthType represents the OAuth 2.0 inbound authentication type.
	OAuthInboundAuthType InboundAuthType = "oauth2"
)

// UserInfoResponseType represents the response format of the UserInfo endpoint.
type UserInfoResponseType string

const (
	// UserInfoResponseTypeJSON represents the JSON userinfo response type.
	UserInfoResponseTypeJSON UserInfoResponseType = "JSON"

	// UserInfoResponseTypeJWS represents the signed JWT (JWS) userinfo response type.
	UserInfoResponseTypeJWS UserInfoResponseType = "JWS"

	// UserInfoResponseTypeJWE represents the encrypted (JWE) userinfo response type.
	UserInfoResponseTypeJWE UserInfoResponseType = "JWE"

	// UserInfoResponseTypeNESTEDJWT represents the signed-then-encrypted (Nested JWT) userinfo response type.
	UserInfoResponseTypeNESTEDJWT UserInfoResponseType = "NESTED_JWT"
)

// SupportedUserInfoSigningAlgs lists JWS algorithms supported for userinfo signing.
var SupportedUserInfoSigningAlgs = []string{
	string(jws.RS256), string(jws.RS512), string(jws.PS256),
	string(jws.ES256), string(jws.ES384), string(jws.ES512),
	string(jws.EdDSA),
}

// SupportedUserInfoEncryptionAlgs lists JWE key-management algorithms supported for userinfo encryption.
var SupportedUserInfoEncryptionAlgs = []string{string(jwe.RSAOAEP), string(jwe.RSAOAEP256)}

// SupportedUserInfoEncryptionEncs lists JWE content-encryption algorithms supported for userinfo encryption.
var SupportedUserInfoEncryptionEncs = []string{string(jwe.A128CBCHS256), string(jwe.A256GCM)}

// ApplicationNotFoundError is the error returned when an application is not found.
var ApplicationNotFoundError error = errors.New("application not found")

// ApplicationDataCorruptedError is the error returned when application data is corrupted.
var ApplicationDataCorruptedError error = errors.New("application data is corrupted")

// Constants for MCP tool defaults
var (
	// DefaultUserAttributes are the standard user attributes for application templates.
	DefaultUserAttributes = []string{
		"email", "name", "given_name", "family_name",
		"profile", "picture", "phone_number", "address", "created_at",
	}
	// DefaultScopes are the standard OAuth scopes for application templates.
	DefaultScopes = []string{"openid", "profile", "email"}
)
