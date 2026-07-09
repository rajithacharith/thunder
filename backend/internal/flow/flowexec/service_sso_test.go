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

package flowexec

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/internal/flow/common"
	"github.com/thunder-id/thunderid/internal/flow/core"
	"github.com/thunder-id/thunderid/internal/flow/executor"
	"github.com/thunder-id/thunderid/internal/flow/session"
	"github.com/thunder-id/thunderid/internal/system/cache"
	"github.com/thunder-id/thunderid/internal/system/config"
	"github.com/thunder-id/thunderid/internal/system/log"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

const testFlowID = "auth-graph-1"

type ServiceSSOTestSuite struct {
	suite.Suite
}

func TestServiceSSOTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceSSOTestSuite))
}

func (s *ServiceSSOTestSuite) SetupTest() {
	s.Require().NoError(config.InitializeServerRuntime(s.T().TempDir(), &config.Config{}))
}

func (s *ServiceSSOTestSuite) TearDownTest() {
	config.ResetServerRuntime()
}

func (s *ServiceSSOTestSuite) newTestGraph() core.GraphInterface {
	flowFactory, _ := core.Initialize(cache.Initialize(config.GetServerRuntime().Config.Cache, "test-deployment"))
	return flowFactory.CreateGraph(testFlowID, providers.FlowTypeAuthentication, 1)
}

// newGraphWithExecutor builds a single-node graph whose node is backed by the given executor.
func (s *ServiceSSOTestSuite) newGraphWithExecutor(executorName string) core.GraphInterface {
	flowFactory, _ := core.Initialize(cache.Initialize(config.GetServerRuntime().Config.Cache, "test-deployment"))
	graph := flowFactory.CreateGraph(testFlowID, providers.FlowTypeAuthentication, 1)
	node, err := flowFactory.CreateNode("n1", string(common.NodeTypeTaskExecution), nil, false, false)
	s.Require().NoError(err)
	node.(core.ExecutorBackedNodeInterface).SetExecutorName(executorName)
	s.Require().NoError(graph.AddNode(node))
	return graph
}

func (s *ServiceSSOTestSuite) TestApplyInboundSSO_SelectsHandleForFlow() {
	engineCtx := &EngineContext{Graph: s.newTestGraph()}

	ih := session.InboundHandle{
		Cookies: map[string]string{session.CookieName(testFlowID): "handle-1"},
	}
	ctx := session.WithInbound(context.Background(), ih)

	applyInboundSSO(engineCtx, ctx)

	s.Equal("handle-1", engineCtx.SSOHandleIn)
}

func (s *ServiceSSOTestSuite) TestApplyInboundSSO_NoInbound() {
	engineCtx := &EngineContext{Graph: s.newTestGraph()}

	applyInboundSSO(engineCtx, context.Background())

	s.Empty(engineCtx.SSOHandleIn)
}

func (s *ServiceSSOTestSuite) TestApplyInboundSSO_NilGraph() {
	engineCtx := &EngineContext{}
	ctx := session.WithInbound(context.Background(),
		session.InboundHandle{Cookies: map[string]string{}})

	applyInboundSSO(engineCtx, ctx)

	s.Empty(engineCtx.SSOHandleIn)
}

func (s *ServiceSSOTestSuite) TestResolveActiveFlowVersion_FromProvider() {
	provider := NewFlowProviderMock(s.T())
	provider.EXPECT().GetFlow(mock.Anything, testFlowID).
		Return(&providers.CompleteFlowDefinition{ID: testFlowID, ActiveVersion: 5}, nil)
	svc := &flowExecService{flowProvider: provider}
	engineCtx := &EngineContext{Graph: s.newTestGraph()}

	version := svc.resolveActiveFlowVersion(context.Background(), engineCtx, log.GetLogger())

	s.Equal(5, version)
}

func (s *ServiceSSOTestSuite) TestResolveActiveFlowVersion_NilGraph() {
	// A nil graph short-circuits before the provider is consulted, so no GetFlow call is expected.
	svc := &flowExecService{flowProvider: NewFlowProviderMock(s.T())}

	version := svc.resolveActiveFlowVersion(context.Background(), &EngineContext{}, log.GetLogger())

	s.Equal(0, version)
}

// TestFlowUsesSSOSession covers the version-lookup gate: a flow that establishes (Session) or consults
// (SSO-Check) a session must have its active version resolved on every path — including the
// fresh-login save path, which carries no inbound handle. Gating the version lookup on an inbound
// handle would persist sessions at version 0 and fail the version check on the next login.
func (s *ServiceSSOTestSuite) TestFlowUsesSSOSession() {
	s.True(flowUsesSSOSession(s.newGraphWithExecutor(executor.ExecutorNameSession)))
	s.True(flowUsesSSOSession(s.newGraphWithExecutor(executor.ExecutorNameSSOCheck)))
	s.False(flowUsesSSOSession(s.newGraphWithExecutor(executor.ExecutorNameCredentialsAuth)))
	s.False(flowUsesSSOSession(nil))
}
