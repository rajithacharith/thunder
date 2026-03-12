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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/asgardeo/thunder/internal/flow/common"
	"github.com/asgardeo/thunder/internal/flow/core"
	"github.com/asgardeo/thunder/internal/system/security"
	"github.com/asgardeo/thunder/tests/mocks/flow/coremock"
)

type OUResolverExecutorTestSuite struct {
	suite.Suite
	mockFlowFactory *coremock.FlowFactoryInterfaceMock
	executor        *ouResolverExecutor
}

func (suite *OUResolverExecutorTestSuite) SetupTest() {
	suite.mockFlowFactory = coremock.NewFlowFactoryInterfaceMock(suite.T())
	mockBaseExecutor := coremock.NewExecutorInterfaceMock(suite.T())

	suite.mockFlowFactory.On("CreateExecutor",
		ExecutorNameOUResolver,
		common.ExecutorTypeUtility,
		[]common.Input{},
		[]common.Input{}).Return(mockBaseExecutor)

	suite.executor = newOUResolverExecutor(suite.mockFlowFactory)
}

func (suite *OUResolverExecutorTestSuite) TestExecute_CreateInAdminOU_Success() {
	adminOUID := "admin-ou-123"
	httpCtx := context.Background()
	authCtx := security.NewSecurityContextForTest(
		"admin-user", adminOUID, "token",
		[]string{"system"}, nil,
	)
	httpCtx = security.WithSecurityContextTest(httpCtx, authCtx)

	ctx := &core.NodeContext{
		FlowID:  "test-flow",
		Context: httpCtx,
		NodeProperties: map[string]interface{}{
			common.NodePropertyCreateInAdminOU: true,
		},
		RuntimeData: map[string]string{
			defaultOUIDKey: "default-ou-456",
		},
	}

	resp, err := suite.executor.Execute(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), common.ExecComplete, resp.Status)
	assert.Equal(suite.T(), adminOUID, resp.RuntimeData[ouIDKey])
}

func (suite *OUResolverExecutorTestSuite) TestExecute_CreateInAdminOU_AdminOUMissing() {
	httpCtx := context.Background()
	// Security context without OU.
	authCtx := security.NewSecurityContextForTest(
		"admin-user", "", "token",
		[]string{"system"}, nil,
	)
	httpCtx = security.WithSecurityContextTest(httpCtx, authCtx)

	ctx := &core.NodeContext{
		FlowID:  "test-flow",
		Context: httpCtx,
		NodeProperties: map[string]interface{}{
			common.NodePropertyCreateInAdminOU: true,
		},
	}

	resp, err := suite.executor.Execute(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), common.ExecFailure, resp.Status)
	assert.Equal(suite.T(), "Unable to determine admin organization unit", resp.FailureReason)
}

func (suite *OUResolverExecutorTestSuite) TestExecute_CreateInAdminOU_False() {
	httpCtx := context.Background()
	authCtx := security.NewSecurityContextForTest(
		"admin-user", "admin-ou-123", "token",
		[]string{"system"}, nil,
	)
	httpCtx = security.WithSecurityContextTest(httpCtx, authCtx)

	ctx := &core.NodeContext{
		FlowID:  "test-flow",
		Context: httpCtx,
		NodeProperties: map[string]interface{}{
			common.NodePropertyCreateInAdminOU: false,
		},
		RuntimeData: map[string]string{
			defaultOUIDKey: "default-ou-456",
		},
	}

	resp, err := suite.executor.Execute(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), common.ExecComplete, resp.Status)
	assert.Empty(suite.T(), resp.RuntimeData[ouIDKey])
}

func (suite *OUResolverExecutorTestSuite) TestExecute_PropertyMissing() {
	httpCtx := context.Background()

	ctx := &core.NodeContext{
		FlowID:         "test-flow",
		Context:        httpCtx,
		NodeProperties: map[string]interface{}{},
		RuntimeData: map[string]string{
			defaultOUIDKey: "default-ou-456",
		},
	}

	resp, err := suite.executor.Execute(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), common.ExecComplete, resp.Status)
	assert.Empty(suite.T(), resp.RuntimeData[ouIDKey])
}

func (suite *OUResolverExecutorTestSuite) TestExecute_NilNodeProperties() {
	httpCtx := context.Background()

	ctx := &core.NodeContext{
		FlowID:         "test-flow",
		Context:        httpCtx,
		NodeProperties: nil,
		RuntimeData: map[string]string{
			defaultOUIDKey: "default-ou-456",
		},
	}

	resp, err := suite.executor.Execute(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), common.ExecComplete, resp.Status)
	assert.Empty(suite.T(), resp.RuntimeData[ouIDKey])
}

func (suite *OUResolverExecutorTestSuite) TestExecute_PropertyWrongType_String() {
	httpCtx := context.Background()
	authCtx := security.NewSecurityContextForTest(
		"admin-user", "admin-ou-123", "token",
		[]string{"system"}, nil,
	)
	httpCtx = security.WithSecurityContextTest(httpCtx, authCtx)

	ctx := &core.NodeContext{
		FlowID:  "test-flow",
		Context: httpCtx,
		NodeProperties: map[string]interface{}{
			common.NodePropertyCreateInAdminOU: "true", // String instead of bool.
		},
	}

	resp, err := suite.executor.Execute(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), common.ExecComplete, resp.Status)
	assert.Empty(suite.T(), resp.RuntimeData[ouIDKey])
}

func (suite *OUResolverExecutorTestSuite) TestExecute_NilContext() {
	ctx := &core.NodeContext{
		FlowID:  "test-flow",
		Context: nil,
		NodeProperties: map[string]interface{}{
			common.NodePropertyCreateInAdminOU: true,
		},
	}

	resp, err := suite.executor.Execute(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), common.ExecFailure, resp.Status)
	assert.Equal(suite.T(), "Unable to determine admin organization unit", resp.FailureReason)
}

func TestOUResolverExecutorSuite(t *testing.T) {
	suite.Run(t, new(OUResolverExecutorTestSuite))
}
