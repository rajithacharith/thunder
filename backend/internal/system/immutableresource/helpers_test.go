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

package immutableresource

import (
	"testing"

	"github.com/asgardeo/thunder/internal/system/config"
	"github.com/stretchr/testify/assert"
)

func TestIsImmutableModeEnabled(t *testing.T) {
	t.Run("Returns true when immutable resources are enabled", func(t *testing.T) {
		config.ResetThunderRuntime()
		testConfig := &config.Config{
			ImmutableResources: config.ImmutableResources{
				Enabled: true,
			},
		}
		err := config.InitializeThunderRuntime("", testConfig)
		assert.NoError(t, err)

		result := IsImmutableModeEnabled()
		assert.True(t, result)
	})

	t.Run("Returns false when immutable resources are disabled", func(t *testing.T) {
		config.ResetThunderRuntime()
		testConfig := &config.Config{
			ImmutableResources: config.ImmutableResources{
				Enabled: false,
			},
		}
		err := config.InitializeThunderRuntime("", testConfig)
		assert.NoError(t, err)

		result := IsImmutableModeEnabled()
		assert.False(t, result)
	})
}

func TestCheckImmutableCreate(t *testing.T) {
	t.Run("Returns error when immutable mode is enabled", func(t *testing.T) {
		config.ResetThunderRuntime()
		testConfig := &config.Config{
			ImmutableResources: config.ImmutableResources{
				Enabled: true,
			},
		}
		err := config.InitializeThunderRuntime("", testConfig)
		assert.NoError(t, err)

		result := CheckImmutableCreate()
		assert.NotNil(t, result)
		assert.Contains(t, result.Error, "Immutable resource create operation")
	})

	t.Run("Returns nil when immutable mode is disabled", func(t *testing.T) {
		config.ResetThunderRuntime()
		testConfig := &config.Config{
			ImmutableResources: config.ImmutableResources{
				Enabled: false,
			},
		}
		err := config.InitializeThunderRuntime("", testConfig)
		assert.NoError(t, err)

		result := CheckImmutableCreate()
		assert.Nil(t, result)
	})
}

func TestCheckImmutableUpdate(t *testing.T) {
	t.Run("Returns error when immutable mode is enabled", func(t *testing.T) {
		config.ResetThunderRuntime()
		testConfig := &config.Config{
			ImmutableResources: config.ImmutableResources{
				Enabled: true,
			},
		}
		err := config.InitializeThunderRuntime("", testConfig)
		assert.NoError(t, err)

		result := CheckImmutableUpdate()
		assert.NotNil(t, result)
		assert.Contains(t, result.Error, "Immutable resource update operation")
	})

	t.Run("Returns nil when immutable mode is disabled", func(t *testing.T) {
		config.ResetThunderRuntime()
		testConfig := &config.Config{
			ImmutableResources: config.ImmutableResources{
				Enabled: false,
			},
		}
		err := config.InitializeThunderRuntime("", testConfig)
		assert.NoError(t, err)

		result := CheckImmutableUpdate()
		assert.Nil(t, result)
	})
}

func TestCheckImmutableDelete(t *testing.T) {
	t.Run("Returns error when immutable mode is enabled", func(t *testing.T) {
		config.ResetThunderRuntime()
		testConfig := &config.Config{
			ImmutableResources: config.ImmutableResources{
				Enabled: true,
			},
		}
		err := config.InitializeThunderRuntime("", testConfig)
		assert.NoError(t, err)

		result := CheckImmutableDelete()
		assert.NotNil(t, result)
		assert.Contains(t, result.Error, "Immutable resource delete operation")
	})

	t.Run("Returns nil when immutable mode is disabled", func(t *testing.T) {
		config.ResetThunderRuntime()
		testConfig := &config.Config{
			ImmutableResources: config.ImmutableResources{
				Enabled: false,
			},
		}
		err := config.InitializeThunderRuntime("", testConfig)
		assert.NoError(t, err)

		result := CheckImmutableDelete()
		assert.Nil(t, result)
	})
}
