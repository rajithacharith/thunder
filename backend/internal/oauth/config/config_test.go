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

package oauthconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thunder-id/thunderid/internal/system/config"
)

func TestFromServerRuntime(t *testing.T) {
	config.ResetServerRuntime()
	defer config.ResetServerRuntime()

	cfg := &config.Config{
		Server: config.ServerConfig{
			Identifier: "dep-1",
			Hostname:   "thunder.io",
			Port:       443,
			PublicURL:  "https://thunder.io",
		},
		Database: config.DatabaseConfig{
			Runtime: config.DataSource{Type: "sqlite"},
		},
		JWT: config.JWTConfig{
			Issuer:         "https://thunder.io",
			ValidityPeriod: 3600,
		},
		OAuth: config.OAuthConfig{
			PAR: config.PARConfig{ExpiresIn: 600},
		},
		GateClient: config.GateClientConfig{
			Scheme:   "https",
			Hostname: "localhost",
			Port:     3000,
		},
	}
	err := config.InitializeServerRuntime("/tmp/test-oauth-config", cfg)
	assert.NoError(t, err)

	result := FromServerRuntime()

	assert.Equal(t, "dep-1", result.DeploymentID)
	assert.Equal(t, "sqlite", result.RuntimeDBType)
	assert.Equal(t, "https://thunder.io", result.BaseURL)
	assert.Equal(t, "https://thunder.io", result.JWT.Issuer)
	assert.Equal(t, int64(600), result.OAuth.PAR.ExpiresIn)
	assert.Equal(t, "localhost", result.GateClient.Hostname)
}
