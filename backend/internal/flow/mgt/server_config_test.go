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

package flowmgt

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	flowconfig "github.com/thunder-id/thunderid/internal/flow/config"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

type FlowConfigHandlerTestSuite struct {
	suite.Suite
	handler *FlowConfigHandler
}

func TestFlowConfigHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(FlowConfigHandlerTestSuite))
}

func (s *FlowConfigHandlerTestSuite) SetupTest() {
	s.handler = NewFlowConfigHandler()
}

func (s *FlowConfigHandlerTestSuite) TestDecode_NilInput() {
	result, err := s.handler.Decode(nil)
	s.NoError(err)
	s.Equal(flowconfig.FlowSectionConfig{}, result)
}

func (s *FlowConfigHandlerTestSuite) TestDecode_EmptyBytes() {
	result, err := s.handler.Decode(json.RawMessage{})
	s.NoError(err)
	s.Equal(flowconfig.FlowSectionConfig{}, result)
}

func (s *FlowConfigHandlerTestSuite) TestDecode_ValidJSON() {
	raw := json.RawMessage(`{"authFlow":{"defaultHandle":"my-auth","expirySeconds":900}}`)

	result, err := s.handler.Decode(raw)
	s.Require().NoError(err)

	cfg, ok := result.(flowconfig.FlowSectionConfig)
	s.Require().True(ok)
	s.Equal("my-auth", cfg.AuthFlow.DefaultHandle)
	s.Equal(int64(900), cfg.AuthFlow.ExpirySeconds)
}

func (s *FlowConfigHandlerTestSuite) TestDecode_InvalidJSON() {
	_, err := s.handler.Decode(json.RawMessage(`{invalid`))
	s.Error(err)
}

func (s *FlowConfigHandlerTestSuite) TestValidate_WrongType() {
	err := s.handler.Validate("not-a-config", nil, nil)
	s.Error(err)
}

func (s *FlowConfigHandlerTestSuite) TestValidate_NegativeExpiry() {
	cfg := flowconfig.FlowSectionConfig{
		AuthFlow: flowconfig.FlowTypeConfig{ExpirySeconds: -1},
	}
	err := s.handler.Validate(cfg, nil, nil)
	s.Error(err)
}

func (s *FlowConfigHandlerTestSuite) TestValidate_ValidConfig() {
	cfg := flowconfig.FlowSectionConfig{
		AuthFlow:         flowconfig.FlowTypeConfig{DefaultHandle: "", ExpirySeconds: 1800},
		RegistrationFlow: flowconfig.FlowTypeConfig{ExpirySeconds: 3600},
		RecoveryFlow:     flowconfig.FlowTypeConfig{ExpirySeconds: 1800},
		SignOutFlow:      flowconfig.FlowTypeConfig{ExpirySeconds: 1800},
	}
	err := s.handler.Validate(cfg, nil, nil)
	s.NoError(err)
}

func (s *FlowConfigHandlerTestSuite) TestValidate_HandleValidatorCalled() {
	called := false
	s.handler.SetHandleValidator(func(_ context.Context, handle string, _ providers.FlowType) bool {
		called = true
		return handle == "valid-handle"
	})

	cfg := flowconfig.FlowSectionConfig{
		AuthFlow: flowconfig.FlowTypeConfig{DefaultHandle: "valid-handle"},
	}
	err := s.handler.Validate(cfg, nil, nil)
	s.NoError(err)
	s.True(called)
}

func (s *FlowConfigHandlerTestSuite) TestValidate_HandleValidatorRejectsUnknown() {
	s.handler.SetHandleValidator(func(_ context.Context, _ string, _ providers.FlowType) bool {
		return false
	})

	cfg := flowconfig.FlowSectionConfig{
		SignOutFlow: flowconfig.FlowTypeConfig{DefaultHandle: "nonexistent"},
	}
	err := s.handler.Validate(cfg, nil, nil)
	s.Error(err)
}

func (s *FlowConfigHandlerTestSuite) TestValidate_NoValidatorSkipsHandleCheck() {
	cfg := flowconfig.FlowSectionConfig{
		AuthFlow: flowconfig.FlowTypeConfig{DefaultHandle: "any-handle"},
	}
	err := s.handler.Validate(cfg, nil, nil)
	s.NoError(err)
}

func (s *FlowConfigHandlerTestSuite) TestMerge_WritableWins() {
	ro := flowconfig.FlowSectionConfig{
		AuthFlow: flowconfig.FlowTypeConfig{DefaultHandle: "ro-handle", ExpirySeconds: 1800},
	}
	wr := flowconfig.FlowSectionConfig{
		AuthFlow: flowconfig.FlowTypeConfig{DefaultHandle: "wr-handle", ExpirySeconds: 900},
	}

	result := s.handler.Merge(ro, wr).(flowconfig.FlowSectionConfig)
	s.Equal("wr-handle", result.AuthFlow.DefaultHandle)
	s.Equal(int64(900), result.AuthFlow.ExpirySeconds)
}

func (s *FlowConfigHandlerTestSuite) TestMerge_ReadOnlyFallsBackWhenWritableZero() {
	ro := flowconfig.FlowSectionConfig{
		AuthFlow: flowconfig.FlowTypeConfig{DefaultHandle: "ro-handle", ExpirySeconds: 1800},
	}
	wr := flowconfig.FlowSectionConfig{
		AuthFlow: flowconfig.FlowTypeConfig{DefaultHandle: "", ExpirySeconds: 0},
	}

	result := s.handler.Merge(ro, wr).(flowconfig.FlowSectionConfig)
	s.Equal("ro-handle", result.AuthFlow.DefaultHandle)
	s.Equal(int64(1800), result.AuthFlow.ExpirySeconds)
}

func (s *FlowConfigHandlerTestSuite) TestMerge_NilInputsReturnZero() {
	result := s.handler.Merge(nil, nil).(flowconfig.FlowSectionConfig)
	s.Equal(flowconfig.FlowSectionConfig{}, result)
}

func (s *FlowConfigHandlerTestSuite) TestMergeFlowTypeConfig_WritableHandleWins() {
	ro := flowconfig.FlowTypeConfig{DefaultHandle: "ro", ExpirySeconds: 500}
	wr := flowconfig.FlowTypeConfig{DefaultHandle: "wr", ExpirySeconds: 0}

	merged := mergeFlowTypeConfig(ro, wr)
	s.Equal("wr", merged.DefaultHandle)
	s.Equal(int64(500), merged.ExpirySeconds)
}

func (s *FlowConfigHandlerTestSuite) TestMergeFlowTypeConfig_WritableExpiryWins() {
	ro := flowconfig.FlowTypeConfig{DefaultHandle: "ro", ExpirySeconds: 500}
	wr := flowconfig.FlowTypeConfig{DefaultHandle: "", ExpirySeconds: 900}

	merged := mergeFlowTypeConfig(ro, wr)
	s.Equal("ro", merged.DefaultHandle)
	s.Equal(int64(900), merged.ExpirySeconds)
}
