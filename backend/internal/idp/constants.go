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

package idp

// IDPType represents the type of an identity provider.
type IDPType string

const (
	// IDPTypeOAuth represents an OAuth2 identity provider.
	IDPTypeOAuth IDPType = "OAUTH"
	// IDPTypeOIDC represents an OIDC identity provider.
	IDPTypeOIDC IDPType = "OIDC"
	// IDPTypeGoogle represents a Google identity provider.
	IDPTypeGoogle IDPType = "GOOGLE"
	// IDPTypeGitHub represents a GitHub identity provider.
	IDPTypeGitHub IDPType = "GITHUB"
	// IDPTypeLDAP represents an LDAP identity provider.
	IDPTypeLDAP IDPType = "LDAP"
	// IDPTypeSAML represents a SAML identity provider.
	IDPTypeSAML IDPType = "SAML"
)

// supportedIDPTypes lists all the supported identity provider types.
var supportedIDPTypes = []IDPType{
	IDPTypeOAuth,
	IDPTypeOIDC,
	IDPTypeGoogle,
	IDPTypeGitHub,
	IDPTypeLDAP,
	IDPTypeSAML,
}

// supportedIDPProperties lists all the supported identity provider properties.
var supportedIDPProperties = []string{
	"client_id",
	"client_secret",
	"redirect_uri",
	"scopes",
	"authorization_endpoint",
	"token_endpoint",
	"userinfo_endpoint",
	"logout_endpoint",
	"jwks_endpoint",
	"prompt",
	// LDAP properties
	"ldap_url",
	"bind_dn",
	"bind_password",
	"user_base_dn",
	"user_filter",
	// SAML properties
	"idp_entity_id",
	"sso_url",
	"x509_cert",
	"acs_url",
}
