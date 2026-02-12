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

	"github.com/stretchr/testify/suite"
)

type ModelTestSuite struct {
	suite.Suite
}

func TestModelTestSuite(t *testing.T) {
	suite.Run(t, new(ModelTestSuite))
}

func (s *ModelTestSuite) TestNodeExecutionRecord_GetDuration() {
	tests := []struct {
		name      string
		startTime int64
		endTime   int64
		expected  int64
	}{
		{"Valid duration calculation", 100, 150, 50000},
		{"Zero start time", 0, 150, 0},
		{"Zero end time", 100, 0, 0},
		{"Both times zero", 0, 0, 0},
		{"Large time values", 1000000, 1000100, 100000},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			record := NodeExecutionRecord{StartTime: tt.startTime, EndTime: tt.endTime}
			duration := record.GetDuration()
			s.Equal(tt.expected, duration)
		})
	}
}

func (s *ModelTestSuite) TestInput_IsSensitive() {
	tests := []struct {
		name     string
		input    Input
		expected bool
	}{
		{"Password input is sensitive", Input{Identifier: "password", Type: InputTypePassword}, true},
		{"OTP input is sensitive", Input{Identifier: "otp", Type: InputTypeOTP}, true},
		{"Text input is not sensitive", Input{Identifier: "username", Type: InputTypeText}, false},
		{"Empty type is not sensitive", Input{Identifier: "field", Type: ""}, false},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.Equal(tt.expected, tt.input.IsSensitive())
		})
	}
}

func (s *ModelTestSuite) TestExecutionAttempt_GetDuration() {
	tests := []struct {
		name      string
		startTime int64
		endTime   int64
		expected  int64
	}{
		{"Valid duration calculation", 200, 250, 50000},
		{"Zero start time", 0, 250, 0},
		{"Zero end time", 200, 0, 0},
		{"Both times zero", 0, 0, 0},
		{"Same start and end time", 500, 500, 0},
		{"One millisecond duration", 1, 2, 1000},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			attempt := ExecutionAttempt{StartTime: tt.startTime, EndTime: tt.endTime}
			duration := attempt.GetDuration()
			s.Equal(tt.expected, duration)
		})
	}
}
