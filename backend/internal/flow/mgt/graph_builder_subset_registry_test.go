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
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/internal/flow/common"
	"github.com/thunder-id/thunderid/internal/flow/executor"
	"github.com/thunder-id/thunderid/internal/system/config"
	"github.com/thunder-id/thunderid/internal/system/log"
	"github.com/thunder-id/thunderid/tests/mocks/authn/githubmock"
	"github.com/thunder-id/thunderid/tests/mocks/authn/googlemock"
	"github.com/thunder-id/thunderid/tests/mocks/authn/oauthmock"
	"github.com/thunder-id/thunderid/tests/mocks/authn/oidcmock"
	"github.com/thunder-id/thunderid/tests/mocks/entityprovidermock"
	"github.com/thunder-id/thunderid/tests/mocks/flow/coremock"
)

type GraphBuilderSubsetRegistryTestSuite struct {
	suite.Suite
}

func TestGraphBuilderSubsetRegistrySuite(t *testing.T) {
	suite.Run(t, new(GraphBuilderSubsetRegistryTestSuite))
}

func (s *GraphBuilderSubsetRegistryTestSuite) SetupSuite() {
	_ = config.InitializeServerRuntime("test", &config.Config{
		Server: config.ServerConfig{Identifier: "test-deployment"},
	})
}

func (s *GraphBuilderSubsetRegistryTestSuite) subsetExecutorRegistry() executor.ExecutorRegistryInterface {
	mockFactory := coremock.NewFlowFactoryInterfaceMock(s.T())
	mockBase := coremock.NewExecutorInterfaceMock(s.T())
	mockBase.On("GetName").Return("").Maybe()
	mockBase.On("GetType").Return(common.ExecutorTypeUtility).Maybe()
	mockFactory.On("CreateExecutor", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(mockBase).Maybe()

	reg, err := executor.Initialize(executor.ExecutorDependencies{
		FlowFactory:    mockFactory,
		EntityProvider: entityprovidermock.NewEntityProviderInterfaceMock(s.T()),
		OAuthSvc:       oauthmock.NewOAuthAuthnServiceInterfaceMock(s.T()),
		OIDCSvc:        oidcmock.NewOIDCAuthnServiceInterfaceMock(s.T()),
		GithubSvc:      githubmock.NewGithubOAuthAuthnServiceInterfaceMock(s.T()),
		GoogleSvc:      googlemock.NewGoogleOIDCAuthnServiceInterfaceMock(s.T()),
	}, config.FlowConfig{
		Executors: []string{executor.ExecutorNameInviteExecutor},
	})
	require.NoError(s.T(), err)
	return reg
}

func (s *GraphBuilderSubsetRegistryTestSuite) TestBuildGraph_SubsetRegistryRejectsUnregisteredExecutor() {
	execRegistry := s.subsetExecutorRegistry()
	require.False(s.T(), execRegistry.IsRegistered(executor.ExecutorNameBasicAuth))

	mockFlowFactory := coremock.NewFlowFactoryInterfaceMock(s.T())
	builder := &graphBuilder{
		flowFactory:      mockFlowFactory,
		executorRegistry: execRegistry,
		logger:           log.GetLogger().With(log.String(log.LoggerKeyComponentName, "FlowGraphBuilder")),
	}

	flow := &CompleteFlowDefinition{
		ID:       "flow-1",
		Handle:   "test-handle",
		Name:     "Test Flow",
		FlowType: common.FlowTypeAuthentication,
		Nodes: []NodeDefinition{
			{
				ID:       "task",
				Type:     "TASK_EXECUTION",
				Executor: &ExecutorDefinition{Name: executor.ExecutorNameBasicAuth},
			},
		},
	}

	mockGraph := coremock.NewGraphInterfaceMock(s.T())
	mockTaskNode := coremock.NewExecutorBackedNodeInterfaceMock(s.T())

	mockFlowFactory.EXPECT().CreateGraph("flow-1", common.FlowTypeAuthentication).Return(mockGraph)
	mockFlowFactory.EXPECT().CreateNode(
		"task", "TASK_EXECUTION", map[string]interface{}(nil), false, true).Return(mockTaskNode, nil)
	mockTaskNode.EXPECT().SetInputs([]common.Input{})

	graph, err := builder.buildGraph(context.Background(), flow)

	s.Nil(graph)
	s.Require().Error(err)
	s.Contains(err.Error(), "executor with name "+executor.ExecutorNameBasicAuth+" not registered")
}
