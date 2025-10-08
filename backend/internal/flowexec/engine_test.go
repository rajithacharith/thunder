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

package flowexec

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	authncm "github.com/asgardeo/thunder/internal/authn/common"
	"github.com/asgardeo/thunder/internal/flow"
	"github.com/asgardeo/thunder/internal/system/config"
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	sysutils "github.com/asgardeo/thunder/internal/system/utils"
)

type FlowEngineTestSuite struct {
	suite.Suite
	engine FlowEngineInterface
}

func TestFlowEngineTestSuite(t *testing.T) {
	suite.Run(t, new(FlowEngineTestSuite))
}

func (suite *FlowEngineTestSuite) SetupSuite() {
	// Initialize ThunderRuntime for tests
	mockConfig := &config.Config{
		Database: config.DatabaseConfig{
			Runtime: config.DataSource{
				Type: "sqlite",
				Name: ":memory:",
			},
		},
		Cache: config.CacheConfig{
			Disabled: true,
		},
	}
	config.ResetThunderRuntime()
	err := config.InitializeThunderRuntime("/test/thunder/home", mockConfig)
	if err != nil {
		suite.T().Fatal("Failed to initialize ThunderRuntime:", err)
	}
}

func (suite *FlowEngineTestSuite) TearDownSuite() {
	config.ResetThunderRuntime()
}

func (suite *FlowEngineTestSuite) SetupTest() {
	suite.engine = NewFlowEngine()
}

// Test NewFlowEngine
func (suite *FlowEngineTestSuite) TestNewFlowEngine() {
	engine := NewFlowEngine()
	assert.NotNil(suite.T(), engine)
	assert.Implements(suite.T(), (*FlowEngineInterface)(nil), engine)
}

// Test Execute with nil graph
func (suite *FlowEngineTestSuite) TestExecute_NilGraph() {
	ctx := &flow.EngineContext{
		FlowID: "test-flow-id",
		Graph:  nil,
	}

	flowStep, err := suite.engine.Execute(ctx)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), flow.ErrorFlowGraphNotInitialized.Code, err.Code)
	assert.Equal(suite.T(), "test-flow-id", flowStep.FlowID)
}

// Test Execute with graph but no start node
func (suite *FlowEngineTestSuite) TestExecute_NoStartNode() {
	mockGraph := &MockGraphForEngine{}
	mockGraph.On("GetStartNode").Return(nil, errors.New("start node not found"))

	ctx := &flow.EngineContext{
		FlowID:      "test-flow-id",
		Graph:       mockGraph,
		CurrentNode: nil,
	}

	flowStep, err := suite.engine.Execute(ctx)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), flow.ErrorStartNodeNotFoundInGraph.Code, err.Code)
	assert.Equal(suite.T(), "test-flow-id", flowStep.FlowID)
	mockGraph.AssertExpectations(suite.T())
}

// Test Execute with successful single node execution
func (suite *FlowEngineTestSuite) TestExecute_SingleNodeSuccess() {
	mockNode := &MockNode{}
	mockGraph := &MockGraphForEngine{}
	mockExecutor := &MockExecutor{}

	// Set up node behavior
	mockNode.On("GetID").Return("node-1")
	mockNode.On("GetType").Return(flow.NodeTypeTaskExecution)
	mockNode.On("GetExecutor").Return(mockExecutor)
	mockNode.On("GetInputData").Return([]flow.InputData{})
	mockNode.On("Execute", mock.AnythingOfType("*flow.NodeContext")).Return(&flow.NodeResponse{
		Status:    flow.NodeStatusComplete,
		Assertion: "test-assertion",
	}, nil)
	mockNode.On("GetNextNodeList").Return([]string{})

	// Set up graph behavior
	mockGraph.On("GetStartNode").Return(mockNode, nil)

	ctx := &flow.EngineContext{
		FlowID:      "test-flow-id",
		FlowType:    flow.FlowTypeAuthentication,
		AppID:       "test-app-id",
		Graph:       mockGraph,
		CurrentNode: nil,
	}

	flowStep, err := suite.engine.Execute(ctx)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "test-flow-id", flowStep.FlowID)
	assert.Equal(suite.T(), flow.FlowStatusComplete, flowStep.Status)
	assert.Equal(suite.T(), "test-assertion", flowStep.Assertion)
	assert.Equal(suite.T(), mockNode, ctx.CurrentNode)

	mockNode.AssertExpectations(suite.T())
	mockGraph.AssertExpectations(suite.T())
}

// Test Execute with node that requires executor setup
func (suite *FlowEngineTestSuite) TestExecute_NodeRequiresExecutorSetup() {
	mockNode := &MockNode{}
	mockGraph := &MockGraphForEngine{}
	mockExecutor := &MockExecutor{}

	// Set up node behavior - executor initially nil
	mockNode.On("GetID").Return("node-1")
	mockNode.On("GetType").Return(flow.NodeTypeTaskExecution)
	mockNode.On("GetExecutor").Return(nil).Once()
	mockNode.On("GetExecutorConfig").Return(&flow.ExecutorConfig{
		Name: "BasicAuthExecutor",
	})
	mockNode.On("SetExecutor", mock.AnythingOfType("flow.ExecutorInterface")).Return()
	mockNode.On("GetExecutor").Return(mockExecutor).Once()
	mockNode.On("GetInputData").Return([]flow.InputData{})
	mockNode.On("Execute", mock.AnythingOfType("*flow.NodeContext")).Return(&flow.NodeResponse{
		Status: flow.NodeStatusComplete,
	}, nil)
	mockNode.On("GetNextNodeList").Return([]string{})

	ctx := &flow.EngineContext{
		FlowID:      "test-flow-id",
		FlowType:    flow.FlowTypeAuthentication,
		AppID:       "test-app-id",
		Graph:       mockGraph,
		CurrentNode: mockNode,
	}

	flowStep, err := suite.engine.Execute(ctx)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), flow.FlowStatusComplete, flowStep.Status)

	mockNode.AssertExpectations(suite.T())
}

// Test Execute with executor construction failure
func (suite *FlowEngineTestSuite) TestExecute_ExecutorConstructionFailure() {
	mockNode := &MockNode{}

	// Set up node behavior - executor construction will fail
	mockNode.On("GetID").Return("node-1")
	mockNode.On("GetType").Return(flow.NodeTypeTaskExecution)
	mockNode.On("GetExecutor").Return(nil)
	mockNode.On("GetExecutorConfig").Return(&flow.ExecutorConfig{
		Name: "UnknownExecutor", // This will cause construction to fail
	})

	ctx := &flow.EngineContext{
		FlowID:      "test-flow-id",
		CurrentNode: mockNode,
	}

	flowStep, err := suite.engine.Execute(ctx)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), flow.ErrorConstructingNodeExecutor.Code, err.Code)
	assert.Equal(suite.T(), "test-flow-id", flowStep.FlowID)

	mockNode.AssertExpectations(suite.T())
}

// Test Execute with node execution failure
func (suite *FlowEngineTestSuite) TestExecute_NodeExecutionFailure() {
	mockNode := &MockNode{}
	mockExecutor := &MockExecutor{}

	mockNode.On("GetID").Return("node-1")
	mockNode.On("GetType").Return(flow.NodeTypeTaskExecution)
	mockNode.On("GetExecutor").Return(mockExecutor)
	mockNode.On("GetInputData").Return([]flow.InputData{})
	mockNode.On("Execute", mock.AnythingOfType("*flow.NodeContext")).Return(nil, &serviceerror.ServiceError{
		Code:             "NODE_EXECUTION_FAILED",
		ErrorDescription: "Node execution failed",
	})

	ctx := &flow.EngineContext{
		FlowID:      "test-flow-id",
		CurrentNode: mockNode,
	}

	flowStep, err := suite.engine.Execute(ctx)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "NODE_EXECUTION_FAILED", err.Code)
	assert.Equal(suite.T(), "test-flow-id", flowStep.FlowID)

	mockNode.AssertExpectations(suite.T())
}

// Test Execute with incomplete node response (view type)
func (suite *FlowEngineTestSuite) TestExecute_IncompleteNodeResponse_View() {
	mockNode := &MockNode{}
	mockExecutor := &MockExecutor{}

	mockNode.On("GetID").Return("node-1")
	mockNode.On("GetType").Return(flow.NodeTypeTaskExecution)
	mockNode.On("GetExecutor").Return(mockExecutor)
	mockNode.On("GetInputData").Return([]flow.InputData{})
	mockNode.On("Execute", mock.AnythingOfType("*flow.NodeContext")).Return(&flow.NodeResponse{
		Status: flow.NodeStatusIncomplete,
		Type:   flow.NodeResponseTypeView,
		RequiredData: []flow.InputData{
			{Name: "username", Type: "string", Required: true},
		},
		Actions: []flow.Action{
			{Type: "submit", ID: "submit-action"},
		},
	}, nil)

	ctx := &flow.EngineContext{
		FlowID:      "test-flow-id",
		CurrentNode: mockNode,
	}

	flowStep, err := suite.engine.Execute(ctx)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), flow.FlowStatusIncomplete, flowStep.Status)
	assert.Equal(suite.T(), flow.StepTypeView, flowStep.Type)
	assert.Len(suite.T(), flowStep.Data.Inputs, 1)
	assert.Equal(suite.T(), "username", flowStep.Data.Inputs[0].Name)
	assert.Len(suite.T(), flowStep.Data.Actions, 1)

	mockNode.AssertExpectations(suite.T())
}

// Test Execute with incomplete node response (redirection type)
func (suite *FlowEngineTestSuite) TestExecute_IncompleteNodeResponse_Redirection() {
	mockNode := &MockNode{}
	mockExecutor := &MockExecutor{}

	mockNode.On("GetID").Return("node-1")
	mockNode.On("GetType").Return(flow.NodeTypeTaskExecution)
	mockNode.On("GetExecutor").Return(mockExecutor)
	mockNode.On("GetInputData").Return([]flow.InputData{})
	mockNode.On("Execute", mock.AnythingOfType("*flow.NodeContext")).Return(&flow.NodeResponse{
		Status:      flow.NodeStatusIncomplete,
		Type:        flow.NodeResponseTypeRedirection,
		RedirectURL: "https://example.com/auth",
		RequiredData: []flow.InputData{
			{Name: "code", Type: "string", Required: true},
		},
		AdditionalData: map[string]string{
			"state": "random-state",
		},
	}, nil)

	ctx := &flow.EngineContext{
		FlowID:      "test-flow-id",
		CurrentNode: mockNode,
	}

	flowStep, err := suite.engine.Execute(ctx)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), flow.FlowStatusIncomplete, flowStep.Status)
	assert.Equal(suite.T(), flow.StepTypeRedirection, flowStep.Type)
	assert.Equal(suite.T(), "https://example.com/auth", flowStep.Data.RedirectURL)
	assert.Len(suite.T(), flowStep.Data.Inputs, 1)
	assert.Equal(suite.T(), "random-state", flowStep.Data.AdditionalData["state"])

	mockNode.AssertExpectations(suite.T())
}

// Test Execute with failure node response
func (suite *FlowEngineTestSuite) TestExecute_FailureNodeResponse() {
	mockNode := &MockNode{}
	mockExecutor := &MockExecutor{}

	mockNode.On("GetID").Return("node-1")
	mockNode.On("GetType").Return(flow.NodeTypeTaskExecution)
	mockNode.On("GetExecutor").Return(mockExecutor)
	mockNode.On("GetInputData").Return([]flow.InputData{})
	mockNode.On("Execute", mock.AnythingOfType("*flow.NodeContext")).Return(&flow.NodeResponse{
		Status:        flow.NodeStatusFailure,
		FailureReason: "Authentication failed",
	}, nil)

	ctx := &flow.EngineContext{
		FlowID:      "test-flow-id",
		CurrentNode: mockNode,
	}

	flowStep, err := suite.engine.Execute(ctx)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), flow.FlowStatusError, flowStep.Status)
	assert.Equal(suite.T(), "Authentication failed", flowStep.FailureReason)

	mockNode.AssertExpectations(suite.T())
}

// Test Execute with multi-node flow
func (suite *FlowEngineTestSuite) TestExecute_MultiNodeFlow() {
	mockNode1 := &MockNode{}
	mockNode2 := &MockNode{}
	mockGraph := &MockGraphForEngine{}
	mockExecutor1 := &MockExecutor{}
	mockExecutor2 := &MockExecutor{}

	// Set up first node
	mockNode1.On("GetID").Return("node-1")
	mockNode1.On("GetType").Return(flow.NodeTypeTaskExecution)
	mockNode1.On("GetExecutor").Return(mockExecutor1)
	mockNode1.On("GetInputData").Return([]flow.InputData{})
	mockNode1.On("Execute", mock.AnythingOfType("*flow.NodeContext")).Return(&flow.NodeResponse{
		Status: flow.NodeStatusComplete,
		RuntimeData: map[string]string{
			"step1": "completed",
		},
	}, nil)
	mockNode1.On("GetNextNodeList").Return([]string{"node-2"})

	// Set up second node
	mockNode2.On("GetID").Return("node-2")
	mockNode2.On("GetType").Return(flow.NodeTypeTaskExecution)
	mockNode2.On("GetExecutor").Return(mockExecutor2)
	mockNode2.On("GetInputData").Return([]flow.InputData{})
	mockNode2.On("Execute", mock.AnythingOfType("*flow.NodeContext")).Return(&flow.NodeResponse{
		Status:    flow.NodeStatusComplete,
		Assertion: "flow-completed",
	}, nil)
	mockNode2.On("GetNextNodeList").Return([]string{})

	// Set up graph
	mockGraph.On("GetNode", "node-2").Return(mockNode2, true)

	ctx := &flow.EngineContext{
		FlowID:      "test-flow-id",
		FlowType:    flow.FlowTypeAuthentication,
		AppID:       "test-app-id",
		Graph:       mockGraph,
		CurrentNode: mockNode1,
		RuntimeData: make(map[string]string),
	}

	flowStep, err := suite.engine.Execute(ctx)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), flow.FlowStatusComplete, flowStep.Status)
	assert.Equal(suite.T(), "flow-completed", flowStep.Assertion)
	assert.Equal(suite.T(), mockNode2, ctx.CurrentNode)
	assert.Equal(suite.T(), "completed", ctx.RuntimeData["step1"])

	mockNode1.AssertExpectations(suite.T())
	mockNode2.AssertExpectations(suite.T())
	mockGraph.AssertExpectations(suite.T())
}

// Test Execute with decision node
func (suite *FlowEngineTestSuite) TestExecute_DecisionNode() {
	mockDecisionNode := &MockNode{}
	mockTargetNode := &MockNode{}
	mockGraph := &MockGraphForEngine{}
	mockExecutor := &MockExecutor{}

	// Set up decision node
	mockDecisionNode.On("GetID").Return("decision-node")
	mockDecisionNode.On("GetType").Return(flow.NodeTypeDecision)
	mockDecisionNode.On("GetInputData").Return([]flow.InputData{})
	mockDecisionNode.On("Execute", mock.AnythingOfType("*flow.NodeContext")).Return(&flow.NodeResponse{
		Status:     flow.NodeStatusComplete,
		NextNodeID: "target-node",
	}, nil)

	// Set up target node
	mockTargetNode.On("GetID").Return("target-node")
	mockTargetNode.On("GetType").Return(flow.NodeTypeTaskExecution)
	mockTargetNode.On("GetExecutor").Return(mockExecutor)
	mockTargetNode.On("GetInputData").Return([]flow.InputData{})
	mockTargetNode.On("Execute", mock.AnythingOfType("*flow.NodeContext")).Return(&flow.NodeResponse{
		Status: flow.NodeStatusComplete,
	}, nil)
	mockTargetNode.On("GetNextNodeList").Return([]string{})

	// Set up graph
	mockGraph.On("GetNode", "target-node").Return(mockTargetNode, true)

	ctx := &flow.EngineContext{
		FlowID:      "test-flow-id",
		Graph:       mockGraph,
		CurrentNode: mockDecisionNode,
	}

	flowStep, err := suite.engine.Execute(ctx)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), flow.FlowStatusComplete, flowStep.Status)
	assert.Equal(suite.T(), mockTargetNode, ctx.CurrentNode)

	mockDecisionNode.AssertExpectations(suite.T())
	mockTargetNode.AssertExpectations(suite.T())
	mockGraph.AssertExpectations(suite.T())
}

// Test Execute with next node not found in graph
func (suite *FlowEngineTestSuite) TestExecute_NextNodeNotFoundInGraph() {
	mockNode := &MockNode{}
	mockGraph := &MockGraphForEngine{}
	mockExecutor := &MockExecutor{}

	mockNode.On("GetID").Return("node-1")
	mockNode.On("GetType").Return(flow.NodeTypeTaskExecution)
	mockNode.On("GetExecutor").Return(mockExecutor)
	mockNode.On("GetInputData").Return([]flow.InputData{})
	mockNode.On("Execute", mock.AnythingOfType("*flow.NodeContext")).Return(&flow.NodeResponse{
		Status: flow.NodeStatusComplete,
	}, nil)
	mockNode.On("GetNextNodeList").Return([]string{"non-existent-node"})

	mockGraph.On("GetNode", "non-existent-node").Return(nil, false)

	ctx := &flow.EngineContext{
		FlowID:      "test-flow-id",
		Graph:       mockGraph,
		CurrentNode: mockNode,
	}

	flowStep, err := suite.engine.Execute(ctx)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), flow.ErrorMovingToNextNode.Code, err.Code)
	assert.Contains(suite.T(), err.ErrorDescription, "next node not found in the graph")
	assert.Equal(suite.T(), "test-flow-id", flowStep.FlowID)

	mockNode.AssertExpectations(suite.T())
	mockGraph.AssertExpectations(suite.T())
}

// Test Execute with unsupported node response status
func (suite *FlowEngineTestSuite) TestExecute_UnsupportedNodeResponseStatus() {
	mockNode := &MockNode{}
	mockExecutor := &MockExecutor{}

	mockNode.On("GetID").Return("node-1")
	mockNode.On("GetType").Return(flow.NodeTypeTaskExecution)
	mockNode.On("GetExecutor").Return(mockExecutor)
	mockNode.On("GetInputData").Return([]flow.InputData{})
	mockNode.On("Execute", mock.AnythingOfType("*flow.NodeContext")).Return(&flow.NodeResponse{
		Status: "UNKNOWN_STATUS",
	}, nil)

	ctx := &flow.EngineContext{
		FlowID:      "test-flow-id",
		CurrentNode: mockNode,
	}

	_, err := suite.engine.Execute(ctx)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), flow.ErrorUnsupportedNodeResponseStatus.Code, err.Code)
	assert.Contains(suite.T(), err.ErrorDescription, "unsupported status returned from the node: UNKNOWN_STATUS")

	mockNode.AssertExpectations(suite.T())
}

// Test Execute with missing node response status
func (suite *FlowEngineTestSuite) TestExecute_MissingNodeResponseStatus() {
	mockNode := &MockNode{}
	mockExecutor := &MockExecutor{}

	mockNode.On("GetID").Return("node-1")
	mockNode.On("GetType").Return(flow.NodeTypeTaskExecution)
	mockNode.On("GetExecutor").Return(mockExecutor)
	mockNode.On("GetInputData").Return([]flow.InputData{})
	mockNode.On("Execute", mock.AnythingOfType("*flow.NodeContext")).Return(&flow.NodeResponse{
		// Status is empty
	}, nil)

	ctx := &flow.EngineContext{
		FlowID:      "test-flow-id",
		CurrentNode: mockNode,
	}

	flowStep, err := suite.engine.Execute(ctx)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), flow.ErrorNodeResponseStatusNotFound.Code, err.Code)
	assert.Equal(suite.T(), "test-flow-id", flowStep.FlowID)

	mockNode.AssertExpectations(suite.T())
}

// Test Execute with unsupported node response type for incomplete status
func (suite *FlowEngineTestSuite) TestExecute_UnsupportedNodeResponseType() {
	mockNode := &MockNode{}
	mockExecutor := &MockExecutor{}

	mockNode.On("GetID").Return("node-1")
	mockNode.On("GetType").Return(flow.NodeTypeTaskExecution)
	mockNode.On("GetExecutor").Return(mockExecutor)
	mockNode.On("GetInputData").Return([]flow.InputData{})
	mockNode.On("Execute", mock.AnythingOfType("*flow.NodeContext")).Return(&flow.NodeResponse{
		Status: flow.NodeStatusIncomplete,
		Type:   "UNKNOWN_TYPE",
	}, nil)

	ctx := &flow.EngineContext{
		FlowID:      "test-flow-id",
		CurrentNode: mockNode,
	}

	flowStep, err := suite.engine.Execute(ctx)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), flow.ErrorUnsupportedNodeResponseType.Code, err.Code)
	assert.Contains(suite.T(), err.ErrorDescription, "unsupported node response type: UNKNOWN_TYPE")
	assert.Equal(suite.T(), "test-flow-id", flowStep.FlowID)

	mockNode.AssertExpectations(suite.T())
}

// Test Execute with redirection but missing redirect URL
func (suite *FlowEngineTestSuite) TestExecute_RedirectionMissingURL() {
	mockNode := &MockNode{}
	mockExecutor := &MockExecutor{}

	mockNode.On("GetID").Return("node-1")
	mockNode.On("GetType").Return(flow.NodeTypeTaskExecution)
	mockNode.On("GetExecutor").Return(mockExecutor)
	mockNode.On("GetInputData").Return([]flow.InputData{})
	mockNode.On("Execute", mock.AnythingOfType("*flow.NodeContext")).Return(&flow.NodeResponse{
		Status: flow.NodeStatusIncomplete,
		Type:   flow.NodeResponseTypeRedirection,
		// RedirectURL is missing
	}, nil)

	ctx := &flow.EngineContext{
		FlowID:      "test-flow-id",
		CurrentNode: mockNode,
	}

	flowStep, err := suite.engine.Execute(ctx)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), flow.ErrorResolvingStepForRedirection.Code, err.Code)
	assert.Contains(suite.T(), err.ErrorDescription, "redirect URL not found in the node response")
	assert.Equal(suite.T(), "test-flow-id", flowStep.FlowID)

	mockNode.AssertExpectations(suite.T())
}

// Test Execute with view response but no required data or actions
func (suite *FlowEngineTestSuite) TestExecute_ViewResponseMissingData() {
	mockNode := &MockNode{}
	mockExecutor := &MockExecutor{}

	mockNode.On("GetID").Return("node-1")
	mockNode.On("GetType").Return(flow.NodeTypeTaskExecution)
	mockNode.On("GetExecutor").Return(mockExecutor)
	mockNode.On("GetInputData").Return([]flow.InputData{})
	mockNode.On("Execute", mock.AnythingOfType("*flow.NodeContext")).Return(&flow.NodeResponse{
		Status: flow.NodeStatusIncomplete,
		Type:   flow.NodeResponseTypeView,
		// No RequiredData or Actions
	}, nil)

	ctx := &flow.EngineContext{
		FlowID:      "test-flow-id",
		CurrentNode: mockNode,
	}

	flowStep, err := suite.engine.Execute(ctx)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), flow.ErrorResolvingStepForPrompt.Code, err.Code)
	assert.Contains(suite.T(), err.ErrorDescription, "no required data or actions found in the node response")
	assert.Equal(suite.T(), "test-flow-id", flowStep.FlowID)

	mockNode.AssertExpectations(suite.T())
}

// Test Execute with authenticated user context update
func (suite *FlowEngineTestSuite) TestExecute_AuthenticatedUserContextUpdate() {
	mockNode := &MockNode{}
	mockExecutor := &MockExecutor{}

	mockNode.On("GetID").Return("node-1")
	mockNode.On("GetType").Return(flow.NodeTypeTaskExecution)
	mockNode.On("GetExecutor").Return(mockExecutor)
	mockNode.On("GetInputData").Return([]flow.InputData{})
	mockNode.On("Execute", mock.AnythingOfType("*flow.NodeContext")).Return(&flow.NodeResponse{
		Status: flow.NodeStatusComplete,
		AuthenticatedUser: authncm.AuthenticatedUser{
			IsAuthenticated: true,
			UserID:          "user-123",
			Attributes: map[string]interface{}{
				"email": "test@example.com",
				"name":  "Test User",
			},
		},
		RuntimeData: map[string]string{
			"authTime": "2024-01-01T10:00:00Z",
		},
	}, nil)
	mockNode.On("GetNextNodeList").Return([]string{})

	ctx := &flow.EngineContext{
		FlowID:      "test-flow-id",
		CurrentNode: mockNode,
		AuthenticatedUser: authncm.AuthenticatedUser{
			Attributes: map[string]interface{}{
				"existingAttr": "existingValue",
			},
		},
		RuntimeData: map[string]string{
			"existingData": "existingValue",
		},
	}

	flowStep, err := suite.engine.Execute(ctx)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), flow.FlowStatusComplete, flowStep.Status)
	assert.True(suite.T(), ctx.AuthenticatedUser.IsAuthenticated)
	assert.Equal(suite.T(), "user-123", ctx.AuthenticatedUser.UserID)
	assert.Equal(suite.T(), "test@example.com", ctx.AuthenticatedUser.Attributes["email"])
	assert.Equal(suite.T(), "existingValue", ctx.AuthenticatedUser.Attributes["existingAttr"])
	assert.Equal(suite.T(), "user-123", ctx.RuntimeData["userID"])
	assert.Equal(suite.T(), "2024-01-01T10:00:00Z", ctx.RuntimeData["authTime"])
	assert.Equal(suite.T(), "existingValue", ctx.RuntimeData["existingData"])

	mockNode.AssertExpectations(suite.T())
}

// Mock implementations

type MockGraphForEngine struct {
	mock.Mock
}

func (m *MockGraphForEngine) GetID() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockGraphForEngine) GetType() flow.FlowType {
	args := m.Called()
	return args.Get(0).(flow.FlowType)
}

func (m *MockGraphForEngine) AddNode(node flow.NodeInterface) error {
	args := m.Called(node)
	return args.Error(0)
}

func (m *MockGraphForEngine) GetNode(nodeID string) (flow.NodeInterface, bool) {
	args := m.Called(nodeID)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(flow.NodeInterface), args.Bool(1)
}

func (m *MockGraphForEngine) AddEdge(fromNodeID, toNodeID string) error {
	args := m.Called(fromNodeID, toNodeID)
	return args.Error(0)
}

func (m *MockGraphForEngine) RemoveEdge(fromNodeID, toNodeID string) error {
	args := m.Called(fromNodeID, toNodeID)
	return args.Error(0)
}

func (m *MockGraphForEngine) GetNodes() map[string]flow.NodeInterface {
	args := m.Called()
	return args.Get(0).(map[string]flow.NodeInterface)
}

func (m *MockGraphForEngine) SetNodes(nodes map[string]flow.NodeInterface) {
	m.Called(nodes)
}

func (m *MockGraphForEngine) GetEdges() map[string][]string {
	args := m.Called()
	return args.Get(0).(map[string][]string)
}

func (m *MockGraphForEngine) SetEdges(edges map[string][]string) {
	m.Called(edges)
}

func (m *MockGraphForEngine) GetStartNodeID() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockGraphForEngine) GetStartNode() (flow.NodeInterface, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(flow.NodeInterface), args.Error(1)
}

func (m *MockGraphForEngine) SetStartNode(startNodeID string) error {
	args := m.Called(startNodeID)
	return args.Error(0)
}

func (m *MockGraphForEngine) ToJSON() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

type MockNode struct {
	mock.Mock
}

func (m *MockNode) Clone() (sysutils.ClonableInterface, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(sysutils.ClonableInterface), args.Error(1)
}

func (m *MockNode) Execute(ctx *flow.NodeContext) (*flow.NodeResponse, *serviceerror.ServiceError) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*serviceerror.ServiceError)
	}
	return args.Get(0).(*flow.NodeResponse), args.Get(1).(*serviceerror.ServiceError)
}

func (m *MockNode) GetID() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockNode) GetType() flow.NodeType {
	args := m.Called()
	return args.Get(0).(flow.NodeType)
}

func (m *MockNode) IsStartNode() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockNode) SetAsStartNode() {
	m.Called()
}

func (m *MockNode) IsFinalNode() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockNode) SetAsFinalNode() {
	m.Called()
}

func (m *MockNode) GetNextNodeList() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockNode) SetNextNodeList(nextNodeIDList []string) {
	m.Called(nextNodeIDList)
}

func (m *MockNode) AddNextNodeID(nextNodeID string) {
	m.Called(nextNodeID)
}

func (m *MockNode) RemoveNextNodeID(nextNodeID string) {
	m.Called(nextNodeID)
}

func (m *MockNode) GetPreviousNodeList() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockNode) SetPreviousNodeList(previousNodeIDList []string) {
	m.Called(previousNodeIDList)
}

func (m *MockNode) AddPreviousNodeID(previousNodeID string) {
	m.Called(previousNodeID)
}

func (m *MockNode) RemovePreviousNodeID(previousNodeID string) {
	m.Called(previousNodeID)
}

func (m *MockNode) GetInputData() []flow.InputData {
	args := m.Called()
	return args.Get(0).([]flow.InputData)
}

func (m *MockNode) SetInputData(inputData []flow.InputData) {
	m.Called(inputData)
}

func (m *MockNode) GetExecutorConfig() *flow.ExecutorConfig {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*flow.ExecutorConfig)
}

func (m *MockNode) SetExecutorConfig(executorConfig *flow.ExecutorConfig) {
	m.Called(executorConfig)
}

func (m *MockNode) GetExecutor() flow.ExecutorInterface {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(flow.ExecutorInterface)
}

func (m *MockNode) SetExecutor(executor flow.ExecutorInterface) {
	m.Called(executor)
}

type MockExecutor struct {
	mock.Mock
}

func (m *MockExecutor) Execute(ctx *flow.NodeContext) (*flow.ExecutorResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*flow.ExecutorResponse), args.Error(1)
}

func (m *MockExecutor) GetID() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockExecutor) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockExecutor) GetProperties() flow.ExecutorProperties {
	args := m.Called()
	return args.Get(0).(flow.ExecutorProperties)
}

func (m *MockExecutor) GetDefaultExecutorInputs() []flow.InputData {
	args := m.Called()
	return args.Get(0).([]flow.InputData)
}

func (m *MockExecutor) GetPrerequisites() []flow.InputData {
	args := m.Called()
	return args.Get(0).([]flow.InputData)
}

func (m *MockExecutor) CheckInputData(ctx *flow.NodeContext, execResp *flow.ExecutorResponse) bool {
	args := m.Called(ctx, execResp)
	return args.Bool(0)
}

func (m *MockExecutor) ValidatePrerequisites(ctx *flow.NodeContext, execResp *flow.ExecutorResponse) bool {
	args := m.Called(ctx, execResp)
	return args.Bool(0)
}

func (m *MockExecutor) GetUserIDFromContext(ctx *flow.NodeContext) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockExecutor) GetRequiredData(ctx *flow.NodeContext) []flow.InputData {
	args := m.Called(ctx)
	return args.Get(0).([]flow.InputData)
}
