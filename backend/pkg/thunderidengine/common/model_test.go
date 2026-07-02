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

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ModelTestSuite struct {
	suite.Suite
}

func TestModelSuite(t *testing.T) {
	suite.Run(t, new(ModelTestSuite))
}

func (suite *ModelTestSuite) TestI18nMessage_String() {
	msg := I18nMessage{
		Key:          "key",
		DefaultValue: "default value",
	}
	assert.Equal(suite.T(), "default value", msg.String())

	msgEmpty := I18nMessage{}
	assert.Equal(suite.T(), "", msgEmpty.String())
}

func (suite *ModelTestSuite) TestI18nMessage_IsEmpty() {
	suite.T().Run("empty message", func(t *testing.T) {
		assert.True(t, I18nMessage{}.IsEmpty())
	})

	suite.T().Run("non-empty message", func(t *testing.T) {
		assert.False(t, I18nMessage{Key: "key"}.IsEmpty())
	})

	suite.T().Run("value only without key", func(t *testing.T) {
		assert.True(t, I18nMessage{DefaultValue: "val"}.IsEmpty())
	})
}

func (suite *ModelTestSuite) TestCustomServiceError() {
	base := ServiceError{
		Type: ClientErrorType,
		Code: "SSE-4030",
		Error: I18nMessage{
			Key:          "error.unauthorized",
			DefaultValue: "Unauthorized",
		},
		ErrorDescription: I18nMessage{
			Key:          "error.unauthorized_description",
			DefaultValue: "The caller is not authorized",
		},
	}

	suite.T().Run("empty custom description preserves base description", func(t *testing.T) {
		err := CustomServiceError(base, I18nMessage{})
		assert.Equal(t, base.Type, err.Type)
		assert.Equal(t, base.Code, err.Code)
		assert.Equal(t, base.Error, err.Error)
		assert.Equal(t, base.ErrorDescription, err.ErrorDescription)
	})

	suite.T().Run("non-empty custom description overrides base description", func(t *testing.T) {
		customDesc := I18nMessage{
			Key:          "error.custom",
			DefaultValue: "Custom description",
		}
		err := CustomServiceError(base, customDesc)
		assert.Equal(t, customDesc, err.ErrorDescription)
		assert.Equal(t, base.Error, err.Error)
	})
}
