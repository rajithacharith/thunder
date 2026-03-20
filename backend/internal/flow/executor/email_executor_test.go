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

package executor

import (
	"html"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/asgardeo/thunder/internal/flow/common"
	"github.com/asgardeo/thunder/internal/flow/core"
	"github.com/asgardeo/thunder/internal/system/email"
	"github.com/asgardeo/thunder/tests/mocks/emailmock"
	"github.com/asgardeo/thunder/tests/mocks/flow/coremock"
)

type EmailExecutorTestSuite struct {
	suite.Suite
	mockFlowFactory *coremock.FlowFactoryInterfaceMock
	mockEmailClient *emailmock.EmailClientInterfaceMock
	executor        *emailExecutor
}

func (suite *EmailExecutorTestSuite) SetupTest() {
	suite.mockFlowFactory = coremock.NewFlowFactoryInterfaceMock(suite.T())
	mockBaseExecutor := coremock.NewExecutorInterfaceMock(suite.T())
	suite.mockEmailClient = emailmock.NewEmailClientInterfaceMock(suite.T())

	suite.mockFlowFactory.On("CreateExecutor",
		ExecutorNameEmailExecutor,
		common.ExecutorTypeUtility,
		[]common.Input{},
		[]common.Input{
			{Identifier: userAttributeEmail, Type: common.InputTypeText, Required: true},
		},
	).Return(mockBaseExecutor)

	suite.executor = newEmailExecutor(suite.mockFlowFactory, suite.mockEmailClient)
}

func (suite *EmailExecutorTestSuite) TestExecute_SendMode_UserInviteTemplate_Success() {
	ctx := &core.NodeContext{
		FlowID:       "test-flow-id",
		ExecutorMode: ExecutorModeSend,
		UserInputs: map[string]string{
			"email": "user@example.com",
		},
		RuntimeData: map[string]string{
			common.RuntimeKeyInviteLink: "https://localhost:5190/gate/invite?flowId=test&inviteToken=abc",
		},
		NodeProperties: map[string]interface{}{
			"emailTemplate": "USER_INVITE",
		},
	}

	var sentEmail email.EmailData
	suite.mockEmailClient.On("Send", mock.Anything).Run(func(args mock.Arguments) {
		sentEmail = args.Get(0).(email.EmailData)
	}).Return(nil)

	resp, err := suite.executor.Execute(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), common.ExecComplete, resp.Status)
	assert.Equal(suite.T(), dataValueTrue, resp.AdditionalData[common.DataEmailSent])
	assert.Equal(suite.T(), []string{"user@example.com"}, sentEmail.To)
	assert.Equal(suite.T(), "You're Invited to Register", sentEmail.Subject)
	assert.True(suite.T(), sentEmail.IsHTML)
	assert.Contains(suite.T(), sentEmail.Body, "Complete Registration")
}

func (suite *EmailExecutorTestSuite) TestExecute_SendMode_EmailFromRuntimeData() {
	ctx := &core.NodeContext{
		FlowID:       "test-flow-id",
		ExecutorMode: ExecutorModeSend,
		UserInputs:   make(map[string]string),
		RuntimeData: map[string]string{
			"email":                     "runtime@example.com",
			common.RuntimeKeyInviteLink: "https://localhost:5190/gate/invite?flowId=test&inviteToken=abc",
		},
		NodeProperties: map[string]interface{}{
			"emailTemplate": "USER_INVITE",
		},
	}

	var sentEmail email.EmailData
	suite.mockEmailClient.On("Send", mock.Anything).Run(func(args mock.Arguments) {
		sentEmail = args.Get(0).(email.EmailData)
	}).Return(nil)

	resp, err := suite.executor.Execute(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), common.ExecComplete, resp.Status)
	assert.Equal(suite.T(), []string{"runtime@example.com"}, sentEmail.To)
}

func (suite *EmailExecutorTestSuite) TestExecute_SendMode_MissingRecipient() {
	ctx := &core.NodeContext{
		FlowID:       "test-flow-id",
		ExecutorMode: ExecutorModeSend,
		UserInputs:   make(map[string]string),
		RuntimeData: map[string]string{
			common.RuntimeKeyInviteLink: "https://localhost:5190/gate/invite?flowId=test&inviteToken=abc",
		},
		NodeProperties: map[string]interface{}{
			"emailTemplate": "USER_INVITE",
		},
	}

	resp, err := suite.executor.Execute(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), common.ExecFailure, resp.Status)
	assert.Equal(suite.T(), "Email recipient is required", resp.FailureReason)
	suite.mockEmailClient.AssertNotCalled(suite.T(), "Send", mock.Anything)
}

func (suite *EmailExecutorTestSuite) TestExecute_SendMode_MissingInviteLink() {
	ctx := &core.NodeContext{
		FlowID:       "test-flow-id",
		ExecutorMode: ExecutorModeSend,
		UserInputs: map[string]string{
			"email": "user@example.com",
		},
		RuntimeData: make(map[string]string),
		NodeProperties: map[string]interface{}{
			"emailTemplate": "USER_INVITE",
		},
	}

	resp, err := suite.executor.Execute(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), common.ExecFailure, resp.Status)
	assert.Contains(suite.T(), resp.FailureReason, "invite link not found")
	suite.mockEmailClient.AssertNotCalled(suite.T(), "Send", mock.Anything)
}

func (suite *EmailExecutorTestSuite) TestExecute_SendMode_MissingTemplateProperty_Success() {
	ctx := &core.NodeContext{
		FlowID:       "test-flow-id",
		ExecutorMode: ExecutorModeSend,
		UserInputs: map[string]string{
			"email": "user@example.com",
		},
		RuntimeData: map[string]string{
			common.RuntimeKeyInviteLink: "https://localhost:5190/gate/invite?flowId=test&inviteToken=abc",
		},
		NodeProperties: map[string]interface{}{},
	}

	var sentEmail email.EmailData
	suite.mockEmailClient.On("Send", mock.Anything).Run(func(args mock.Arguments) {
		sentEmail = args.Get(0).(email.EmailData)
	}).Return(nil)

	resp, err := suite.executor.Execute(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), common.ExecComplete, resp.Status)
	assert.Equal(suite.T(), []string{"user@example.com"}, sentEmail.To)
}

func (suite *EmailExecutorTestSuite) TestExecute_SendMode_ClientError() {
	ctx := &core.NodeContext{
		FlowID:       "test-flow-id",
		ExecutorMode: ExecutorModeSend,
		UserInputs: map[string]string{
			"email": "user@example.com",
		},
		RuntimeData: map[string]string{
			common.RuntimeKeyInviteLink: "https://localhost:5190/gate/invite?flowId=test&inviteToken=abc",
		},
		NodeProperties: map[string]interface{}{
			"emailTemplate": "USER_INVITE",
		},
	}

	suite.mockEmailClient.On("Send", mock.Anything).Return(email.ErrorInvalidRecipient)

	resp, err := suite.executor.Execute(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), common.ExecFailure, resp.Status)
	assert.Equal(suite.T(), "Failed to send email", resp.FailureReason)
}

func (suite *EmailExecutorTestSuite) TestExecute_SendMode_ServerError() {
	ctx := &core.NodeContext{
		FlowID:       "test-flow-id",
		ExecutorMode: ExecutorModeSend,
		UserInputs: map[string]string{
			"email": "user@example.com",
		},
		RuntimeData: map[string]string{
			common.RuntimeKeyInviteLink: "https://localhost:5190/gate/invite?flowId=test&inviteToken=abc",
		},
		NodeProperties: map[string]interface{}{
			"emailTemplate": "USER_INVITE",
		},
	}

	suite.mockEmailClient.On("Send", mock.Anything).Return(email.ErrorSMTPConnection)

	resp, err := suite.executor.Execute(ctx)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
}

func (suite *EmailExecutorTestSuite) TestExecute_SendMode_NilEmailClient_NoOp() {
	// Create executor with nil email client (SMTP not configured)
	mockBaseExecutor := coremock.NewExecutorInterfaceMock(suite.T())
	mockFactory := coremock.NewFlowFactoryInterfaceMock(suite.T())
	mockFactory.On("CreateExecutor",
		ExecutorNameEmailExecutor,
		common.ExecutorTypeUtility,
		[]common.Input{},
		[]common.Input{
			{Identifier: userAttributeEmail, Type: common.InputTypeText, Required: true},
		},
	).Return(mockBaseExecutor)

	noEmailExecutor := newEmailExecutor(mockFactory, nil)

	ctx := &core.NodeContext{
		FlowID:       "test-flow-id",
		ExecutorMode: ExecutorModeSend,
		UserInputs: map[string]string{
			"email": "user@example.com",
		},
		RuntimeData: map[string]string{
			common.RuntimeKeyInviteLink: "https://localhost:5190/gate/invite?flowId=test&inviteToken=abc",
		},
		NodeProperties: map[string]interface{}{
			"emailTemplate": "USER_INVITE",
		},
	}

	resp, err := noEmailExecutor.Execute(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), common.ExecComplete, resp.Status)
	// emailSent should NOT be set when email client is not configured
	assert.Empty(suite.T(), resp.AdditionalData[common.DataEmailSent])
}

func (suite *EmailExecutorTestSuite) TestExecute_InvalidMode() {
	ctx := &core.NodeContext{
		FlowID:       "test-flow-id",
		ExecutorMode: "invalid",
		UserInputs:   make(map[string]string),
		RuntimeData:  make(map[string]string),
	}

	resp, err := suite.executor.Execute(ctx)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "invalid executor mode for EmailExecutor")
}

func (suite *EmailExecutorTestSuite) TestBuildInviteEmailBody_EscapesLink() {
	// Craft a link that would break naive interpolation if unescaped
	malicious := `https://example.com/?q="'><script>alert(1)</script>`
	body := buildInviteEmailBody(malicious)

	// Original string should not appear unescaped
	assert.NotContains(suite.T(), body, malicious)
	escaped := html.EscapeString(malicious)
	assert.Contains(suite.T(), body, escaped)
	// Ensure the escaped version appears in both href and visible text sections
	assert.Equal(suite.T(), 3, strings.Count(body, escaped))
}

func TestEmailExecutorSuite(t *testing.T) {
	suite.Run(t, new(EmailExecutorTestSuite))
}
