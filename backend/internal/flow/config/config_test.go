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

package flowconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thunder-id/thunderid/internal/system/config"
)

func TestFromServerRuntime(t *testing.T) {
	config.ResetServerRuntime()
	defer config.ResetServerRuntime()

	cfg := &config.Config{
		Flow: config.FlowConfig{UserOnboardingFlowHandle: "onboarding-handle"},
		Server: config.ServerConfig{
			Identifier: "dep-1",
		},
		Database: config.DatabaseConfig{
			Runtime: config.DataSource{Type: "postgres"},
		},
	}
	err := config.InitializeServerRuntime("/tmp/test-flow-config", cfg)
	assert.NoError(t, err)

	result := FromServerRuntime()

	assert.Equal(t, "onboarding-handle", result.Flow.UserOnboardingFlowHandle)
	assert.Equal(t, "dep-1", result.DeploymentID)
	assert.Equal(t, "postgres", result.RuntimeDBType)
}
