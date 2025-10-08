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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	appmodel "github.com/asgardeo/thunder/internal/application/model"
	"github.com/asgardeo/thunder/internal/flow"
	"github.com/asgardeo/thunder/internal/system/config"
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
)

type FlowExecServiceTestSuite struct {
	suite.Suite
	flowExecService FlowExecServiceInterface
}

func TestFlowExecServiceSuite(t *testing.T) {
	suite.Run(t, new(FlowExecServiceTestSuite))
}

func (suite *FlowExecServiceTestSuite) SetupSuite() {
	// Initialize ThunderRuntime for tests
	mockConfig := &config.Config{}
	config.ResetThunderRuntime()
	err := config.InitializeThunderRuntime("/test/thunder/home", mockConfig)
	if err != nil {
		suite.T().Fatal("Failed to initialize ThunderRuntime:", err)
	}
}

func (suite *FlowExecServiceTestSuite) TearDownSuite() {
	config.ResetThunderRuntime()
}

func (suite *FlowExecServiceTestSuite) SetupTest() {
	// Create minimal mocks
	mockAppService := &MockApplicationService{}
	mockFlowMgtService := &MockFlowMgtService{}
	mockFlowEngine := &MockFlowEngine{}
	mockflowStore := &MockFlowStore{}

	// Create service with mocked dependencies
	suite.flowExecService = &FlowExecService{
		appService:     mockAppService,
		flowMgtService: mockFlowMgtService,
		flowEngine:     mockFlowEngine,
		flowStore:      mockflowStore,
	}
}

// Test validateFlowType function
func (suite *FlowExecServiceTestSuite) TestValidateFlowType() {
	// Test valid authentication flow type
	flowType, err := validateFlowType("AUTHENTICATION")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), flow.FlowTypeAuthentication, flowType)

	// Test valid registration flow type
	flowType, err = validateFlowType("REGISTRATION")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), flow.FlowTypeRegistration, flowType)

	// Test invalid flow type
	flowType, err = validateFlowType("INVALID_FLOW_TYPE")
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), flow.ErrorInvalidFlowType.Code, err.Code)
	assert.Equal(suite.T(), flow.FlowType(""), flowType)
}

// Test isNewFlow function
func (suite *FlowExecServiceTestSuite) TestIsNewFlow() {
	// Test with empty flow ID (new flow)
	assert.True(suite.T(), isNewFlow(""))

	// Test with non-empty flow ID (existing flow)
	assert.False(suite.T(), isNewFlow("existing-flow-id"))
}

// Test isComplete function
func (suite *FlowExecServiceTestSuite) TestIsComplete() {
	// Test with complete status
	completeStep := flow.FlowStep{Status: flow.FlowStatusComplete}
	assert.True(suite.T(), isComplete(completeStep))

	// Test with incomplete status
	incompleteStep := flow.FlowStep{Status: flow.FlowStatusIncomplete}
	assert.False(suite.T(), isComplete(incompleteStep))

	// Test with error status
	errorStep := flow.FlowStep{Status: flow.FlowStatusError}
	assert.False(suite.T(), isComplete(errorStep))
}

// Test prepareContext function
func (suite *FlowExecServiceTestSuite) TestPrepareContext() {
	// Create test context
	ctx := &flow.EngineContext{}

	// Test with input data and action ID
	inputData := map[string]string{"username": "testuser", "password": "testpass"}
	actionID := "test-action"

	prepareContext(ctx, actionID, inputData)

	// Assert that context was properly prepared
	assert.Equal(suite.T(), actionID, ctx.CurrentActionID)
	assert.Equal(suite.T(), inputData, ctx.UserInputData)
	assert.NotNil(suite.T(), ctx.RuntimeData)

	// Test merging additional input data
	additionalData := map[string]string{"email": "test@example.com"}
	prepareContext(ctx, "", additionalData)

	// Assert that data was merged
	assert.Equal(suite.T(), "testuser", ctx.UserInputData["username"])
	assert.Equal(suite.T(), "testpass", ctx.UserInputData["password"])
	assert.Equal(suite.T(), "test@example.com", ctx.UserInputData["email"])
}

// Test Execute method with invalid flow type
func (suite *FlowExecServiceTestSuite) TestExecute_InvalidFlowType() {
	// Test data
	appID := "test-app-id"
	flowID := ""
	actionID := "test-action"
	flowType := "invalid-flow-type"
	inputData := map[string]string{}

	// Execute
	result, err := suite.flowExecService.Execute(appID, flowID, actionID, flowType, inputData)

	// Assertions
	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), flow.ErrorInvalidFlowType.Code, err.Code)
}

// Mock implementations for testing without import cycles

// MockApplicationService implements a minimal ApplicationServiceInterface for testing
type MockApplicationService struct {
	mock.Mock
}

func (m *MockApplicationService) CreateApplication(app *appmodel.ApplicationDTO) (*appmodel.ApplicationDTO, *serviceerror.ServiceError) {
	args := m.Called(app)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*serviceerror.ServiceError)
	}
	return args.Get(0).(*appmodel.ApplicationDTO), nil
}

func (m *MockApplicationService) GetApplicationList() (*appmodel.ApplicationListResponse, *serviceerror.ServiceError) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Get(1).(*serviceerror.ServiceError)
	}
	return args.Get(0).(*appmodel.ApplicationListResponse), nil
}

func (m *MockApplicationService) GetOAuthApplication(clientID string) (*appmodel.OAuthAppConfigProcessedDTO, *serviceerror.ServiceError) {
	args := m.Called(clientID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*serviceerror.ServiceError)
	}
	return args.Get(0).(*appmodel.OAuthAppConfigProcessedDTO), nil
}

func (m *MockApplicationService) GetApplication(appID string) (*appmodel.ApplicationProcessedDTO, *serviceerror.ServiceError) {
	args := m.Called(appID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*serviceerror.ServiceError)
	}
	return args.Get(0).(*appmodel.ApplicationProcessedDTO), nil
}

func (m *MockApplicationService) UpdateApplication(appID string, app *appmodel.ApplicationDTO) (*appmodel.ApplicationDTO, *serviceerror.ServiceError) {
	args := m.Called(appID, app)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*serviceerror.ServiceError)
	}
	return args.Get(0).(*appmodel.ApplicationDTO), nil
}

func (m *MockApplicationService) DeleteApplication(appID string) *serviceerror.ServiceError {
	args := m.Called(appID)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*serviceerror.ServiceError)
}

// MockFlowMgtService implements a minimal FlowMgtServiceInterface for testing
type MockFlowMgtService struct {
	mock.Mock
}

func (m *MockFlowMgtService) RegisterGraph(graphID string, g flow.GraphInterface) {
	m.Called(graphID, g)
}

func (m *MockFlowMgtService) GetGraph(graphID string) (flow.GraphInterface, bool) {
	args := m.Called(graphID)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(flow.GraphInterface), args.Bool(1)
}

func (m *MockFlowMgtService) IsValidGraphID(graphID string) bool {
	args := m.Called(graphID)
	return args.Bool(0)
}

// MockFlowEngine implements FlowEngineInterface for testing
type MockFlowEngine struct {
	mock.Mock
}

func (m *MockFlowEngine) Execute(ctx *flow.EngineContext) (flow.FlowStep, *serviceerror.ServiceError) {
	args := m.Called(ctx)
	if args.Get(1) == nil {
		return args.Get(0).(flow.FlowStep), nil
	}
	return flow.FlowStep{}, args.Get(1).(*serviceerror.ServiceError)
}

// MockGraph implements flow.GraphInterface for testing
type MockGraph struct {
	id        string
	flowType  flow.FlowType
	nodes     map[string]flow.NodeInterface
	edges     map[string][]string
	startNode string
}

func (m *MockGraph) GetID() string {
	return m.id
}

func (m *MockGraph) GetType() flow.FlowType {
	return m.flowType
}

func (m *MockGraph) GetNodes() map[string]flow.NodeInterface {
	if m.nodes == nil {
		m.nodes = make(map[string]flow.NodeInterface)
	}
	return m.nodes
}

func (m *MockGraph) GetEdges() map[string][]string {
	if m.edges == nil {
		m.edges = make(map[string][]string)
	}
	return m.edges
}

func (m *MockGraph) GetStartNodeID() string {
	return m.startNode
}

func (m *MockGraph) SetStartNode(nodeID string) error {
	m.startNode = nodeID
	return nil
}

func (m *MockGraph) SetNodes(nodes map[string]flow.NodeInterface) {
	m.nodes = nodes
}

func (m *MockGraph) SetEdges(edges map[string][]string) {
	m.edges = edges
}

func (m *MockGraph) AddNode(node flow.NodeInterface) error {
	if m.nodes == nil {
		m.nodes = make(map[string]flow.NodeInterface)
	}
	m.nodes[node.GetID()] = node
	return nil
}

func (m *MockGraph) AddEdge(from, to string) error {
	if m.edges == nil {
		m.edges = make(map[string][]string)
	}
	m.edges[from] = append(m.edges[from], to)
	return nil
}

func (m *MockGraph) RemoveEdge(from, to string) error {
	if m.edges == nil {
		return nil
	}
	edges := m.edges[from]
	for i, edge := range edges {
		if edge == to {
			m.edges[from] = append(edges[:i], edges[i+1:]...)
			break
		}
	}
	return nil
}

func (m *MockGraph) ToJSON() (string, error) {
	return `{"id":"` + m.id + `","type":"` + string(m.flowType) + `"}`, nil
}

// MockFlowStore implements FlowStoreInterface for testing
type MockFlowStore struct {
	mock.Mock
}

func (m *MockFlowStore) StoreFlowContext(ctx flow.EngineContext) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockFlowStore) GetFlowContext(flowID string) (*FlowContextWithUserDataDB, error) {
	args := m.Called(flowID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*FlowContextWithUserDataDB), nil
}

func (m *MockFlowStore) UpdateFlowContext(ctx flow.EngineContext) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockFlowStore) DeleteFlowContext(flowID string) error {
	args := m.Called(flowID)
	return args.Error(0)
}
