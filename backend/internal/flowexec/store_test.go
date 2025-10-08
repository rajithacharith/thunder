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
	"github.com/asgardeo/thunder/tests/mocks/database/clientmock"
	"github.com/asgardeo/thunder/tests/mocks/database/modelmock"
	"github.com/asgardeo/thunder/tests/mocks/database/providermock"
)

type FlowStoreTestSuite struct {
	suite.Suite
	mockDBProvider *providermock.DBProviderInterfaceMock
	mockDBClient   *clientmock.DBClientInterfaceMock
	mockTx         *modelmock.TxInterfaceMock
	flowStore      FlowStoreInterface
}

func TestFlowStoreTestSuite(t *testing.T) {
	suite.Run(t, new(FlowStoreTestSuite))
}

func (suite *FlowStoreTestSuite) SetupSuite() {
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

func (suite *FlowStoreTestSuite) TearDownSuite() {
	config.ResetThunderRuntime()
}

func (suite *FlowStoreTestSuite) SetupTest() {
	suite.mockDBProvider = providermock.NewDBProviderInterfaceMock(suite.T())
	suite.mockDBClient = clientmock.NewDBClientInterfaceMock(suite.T())
	suite.mockTx = modelmock.NewTxInterfaceMock(suite.T())

	// Create store with mocked dependencies
	suite.flowStore = &flowStore{
		DBProvider: suite.mockDBProvider,
	}
}

func (suite *FlowStoreTestSuite) TestNewFlowStore() {
	store := newFlowStore()
	assert.NotNil(suite.T(), store)
	assert.Implements(suite.T(), (*FlowStoreInterface)(nil), store)
}

// Test StoreFlowContext method with successful operation
func (suite *FlowStoreTestSuite) TestStoreFlowContext_Success() {
	// Create test context
	ctx := flow.EngineContext{
		FlowID:          "test-flow-id",
		AppID:           "test-app-id",
		CurrentActionID: "test-action-id",
		FlowType:        flow.FlowTypeAuthentication,
		UserInputData:   map[string]string{"username": "testuser"},
		RuntimeData:     map[string]string{"key": "value"},
		AuthenticatedUser: authncm.AuthenticatedUser{
			IsAuthenticated: true,
			UserID:          "test-user-id",
			Attributes:      map[string]interface{}{"email": "test@example.com"},
		},
	}

	// Set up mock expectations
	suite.mockDBProvider.On("GetDBClient", "runtime").Return(suite.mockDBClient, nil)
	suite.mockDBClient.On("BeginTx").Return(suite.mockTx, nil)

	// Mock execution of flow context query
	suite.mockTx.On("Exec", mock.MatchedBy(func(query string) bool {
		return query == QueryCreateFlowContext.Query
	}), mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	// Mock execution of user data query
	suite.mockTx.On("Exec", mock.MatchedBy(func(query string) bool {
		return query == QueryCreateFlowUserData.Query
	}), mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	suite.mockTx.On("Commit").Return(nil)

	// Execute
	err := suite.flowStore.StoreFlowContext(ctx)

	// Assertions
	assert.Nil(suite.T(), err)

	// Verify all expectations were met
	suite.mockDBProvider.AssertExpectations(suite.T())
	suite.mockDBClient.AssertExpectations(suite.T())
	suite.mockTx.AssertExpectations(suite.T())
}

// Test StoreFlowContext with database client error
func (suite *FlowStoreTestSuite) TestStoreFlowContext_DBClientError() {
	ctx := flow.EngineContext{
		FlowID: "test-flow-id",
		AppID:  "test-app-id",
	}

	// Set up mock expectations
	expectedError := errors.New("database client error")
	suite.mockDBProvider.On("GetDBClient", "runtime").Return(nil, expectedError)

	// Execute
	err := suite.flowStore.StoreFlowContext(ctx)

	// Assertions
	assert.NotNil(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to get database client")

	// Verify expectations
	suite.mockDBProvider.AssertExpectations(suite.T())
}

// Test StoreFlowContext with transaction begin error
func (suite *FlowStoreTestSuite) TestStoreFlowContext_TransactionBeginError() {
	ctx := flow.EngineContext{
		FlowID: "test-flow-id",
		AppID:  "test-app-id",
	}

	// Set up mock expectations
	expectedError := errors.New("begin transaction error")
	suite.mockDBProvider.On("GetDBClient", "runtime").Return(suite.mockDBClient, nil)
	suite.mockDBClient.On("BeginTx").Return(nil, expectedError)

	// Execute
	err := suite.flowStore.StoreFlowContext(ctx)

	// Assertions
	assert.NotNil(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to begin transaction")

	// Verify expectations
	suite.mockDBProvider.AssertExpectations(suite.T())
	suite.mockDBClient.AssertExpectations(suite.T())
}

// Test StoreFlowContext with query execution error
func (suite *FlowStoreTestSuite) TestStoreFlowContext_QueryExecutionError() {
	ctx := flow.EngineContext{
		FlowID: "test-flow-id",
		AppID:  "test-app-id",
	}

	// Set up mock expectations
	queryError := errors.New("query execution error")
	suite.mockDBProvider.On("GetDBClient", "runtime").Return(suite.mockDBClient, nil)
	suite.mockDBClient.On("BeginTx").Return(suite.mockTx, nil)

	// Mock first query failure
	suite.mockTx.On("Exec", mock.MatchedBy(func(query string) bool {
		return query == QueryCreateFlowContext.Query
	}), mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, queryError)

	suite.mockTx.On("Rollback").Return(nil)

	// Execute
	err := suite.flowStore.StoreFlowContext(ctx)

	// Assertions
	assert.NotNil(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "transaction failed")

	// Verify expectations
	suite.mockDBProvider.AssertExpectations(suite.T())
	suite.mockDBClient.AssertExpectations(suite.T())
	suite.mockTx.AssertExpectations(suite.T())
}

// Test GetFlowContext method with successful operation
func (suite *FlowStoreTestSuite) TestGetFlowContext_Success() {
	flowID := "test-flow-id"

	// Mock database result
	mockResult := []map[string]interface{}{
		{
			"flow_id":           "test-flow-id",
			"app_id":            "test-app-id",
			"current_node_id":   "test-node-id",
			"current_action_id": "test-action-id",
			"graph_id":          "test-graph-id",
			"runtime_data":      `{"key":"value"}`,
			"is_authenticated":  true,
			"user_id":           "test-user-id",
			"user_inputs":       `{"username":"testuser"}`,
			"user_attributes":   `{"email":"test@example.com"}`,
		},
	}

	// Set up mock expectations
	suite.mockDBProvider.On("GetDBClient", "runtime").Return(suite.mockDBClient, nil)
	suite.mockDBClient.On("Query", QueryGetFlowContextWithUserData, flowID).Return(mockResult, nil)

	// Execute
	result, err := suite.flowStore.GetFlowContext(flowID)

	// Assertions
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "test-flow-id", result.FlowID)
	assert.Equal(suite.T(), "test-app-id", result.AppID)
	assert.Equal(suite.T(), "test-node-id", *result.CurrentNodeID)
	assert.Equal(suite.T(), "test-action-id", *result.CurrentActionID)
	assert.Equal(suite.T(), "test-graph-id", result.GraphID)
	assert.True(suite.T(), result.IsAuthenticated)
	assert.Equal(suite.T(), "test-user-id", *result.UserID)

	// Verify expectations
	suite.mockDBProvider.AssertExpectations(suite.T())
	suite.mockDBClient.AssertExpectations(suite.T())
}

// Test GetFlowContext with no results (flow not found)
func (suite *FlowStoreTestSuite) TestGetFlowContext_NotFound() {
	flowID := "non-existent-flow-id"

	// Mock empty result
	mockResult := []map[string]interface{}{}

	// Set up mock expectations
	suite.mockDBProvider.On("GetDBClient", "runtime").Return(suite.mockDBClient, nil)
	suite.mockDBClient.On("Query", QueryGetFlowContextWithUserData, flowID).Return(mockResult, nil)

	// Execute
	result, err := suite.flowStore.GetFlowContext(flowID)

	// Assertions
	assert.Nil(suite.T(), err)
	assert.Nil(suite.T(), result)

	// Verify expectations
	suite.mockDBProvider.AssertExpectations(suite.T())
	suite.mockDBClient.AssertExpectations(suite.T())
}

// Test GetFlowContext with multiple results error
func (suite *FlowStoreTestSuite) TestGetFlowContext_MultipleResults() {
	flowID := "test-flow-id"

	// Mock multiple results
	mockResult := []map[string]interface{}{
		{"flow_id": "test-flow-id", "app_id": "test-app-id-1", "graph_id": "graph-1"},
		{"flow_id": "test-flow-id", "app_id": "test-app-id-2", "graph_id": "graph-2"},
	}

	// Set up mock expectations
	suite.mockDBProvider.On("GetDBClient", "runtime").Return(suite.mockDBClient, nil)
	suite.mockDBClient.On("Query", QueryGetFlowContextWithUserData, flowID).Return(mockResult, nil)

	// Execute
	result, err := suite.flowStore.GetFlowContext(flowID)

	// Assertions
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "unexpected number of results: 2")

	// Verify expectations
	suite.mockDBProvider.AssertExpectations(suite.T())
	suite.mockDBClient.AssertExpectations(suite.T())
}

// Test GetFlowContext with database error
func (suite *FlowStoreTestSuite) TestGetFlowContext_DatabaseError() {
	flowID := "test-flow-id"

	// Set up mock expectations
	expectedError := errors.New("database query error")
	suite.mockDBProvider.On("GetDBClient", "runtime").Return(suite.mockDBClient, nil)
	suite.mockDBClient.On("Query", QueryGetFlowContextWithUserData, flowID).Return(nil, expectedError)

	// Execute
	result, err := suite.flowStore.GetFlowContext(flowID)

	// Assertions
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "failed to execute query")

	// Verify expectations
	suite.mockDBProvider.AssertExpectations(suite.T())
	suite.mockDBClient.AssertExpectations(suite.T())
}

// Test UpdateFlowContext method with successful operation
func (suite *FlowStoreTestSuite) TestUpdateFlowContext_Success() {
	ctx := flow.EngineContext{
		FlowID:          "test-flow-id",
		AppID:           "test-app-id",
		CurrentActionID: "updated-action-id",
		RuntimeData:     map[string]string{"updated": "data"},
	}

	// Set up mock expectations
	suite.mockDBProvider.On("GetDBClient", "runtime").Return(suite.mockDBClient, nil)
	suite.mockDBClient.On("BeginTx").Return(suite.mockTx, nil)

	// Mock execution of update flow context query
	suite.mockTx.On("Exec", mock.MatchedBy(func(query string) bool {
		return query == QueryUpdateFlowContext.Query
	}), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	// Mock execution of update user data query
	suite.mockTx.On("Exec", mock.MatchedBy(func(query string) bool {
		return query == QueryUpdateFlowUserData.Query
	}), mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	suite.mockTx.On("Commit").Return(nil)

	// Execute
	err := suite.flowStore.UpdateFlowContext(ctx)

	// Assertions
	assert.Nil(suite.T(), err)

	// Verify expectations
	suite.mockDBProvider.AssertExpectations(suite.T())
	suite.mockDBClient.AssertExpectations(suite.T())
	suite.mockTx.AssertExpectations(suite.T())
}

// Test DeleteFlowContext method with successful operation
func (suite *FlowStoreTestSuite) TestDeleteFlowContext_Success() {
	flowID := "test-flow-id"

	// Set up mock expectations
	suite.mockDBProvider.On("GetDBClient", "runtime").Return(suite.mockDBClient, nil)
	suite.mockDBClient.On("BeginTx").Return(suite.mockTx, nil)

	// Mock execution of delete user data query (executed first)
	suite.mockTx.On("Exec", mock.MatchedBy(func(query string) bool {
		return query == QueryDeleteFlowUserData.Query
	}), flowID).Return(nil, nil)

	// Mock execution of delete flow context query
	suite.mockTx.On("Exec", mock.MatchedBy(func(query string) bool {
		return query == QueryDeleteFlowContext.Query
	}), flowID).Return(nil, nil)

	suite.mockTx.On("Commit").Return(nil)

	// Execute
	err := suite.flowStore.DeleteFlowContext(flowID)

	// Assertions
	assert.Nil(suite.T(), err)

	// Verify expectations
	suite.mockDBProvider.AssertExpectations(suite.T())
	suite.mockDBClient.AssertExpectations(suite.T())
	suite.mockTx.AssertExpectations(suite.T())
}

// Test buildFlowContextFromResultRow with valid data
func (suite *FlowStoreTestSuite) TestBuildFlowContextFromResultRow_Success() {
	row := map[string]interface{}{
		"flow_id":           "test-flow-id",
		"app_id":            "test-app-id",
		"current_node_id":   "test-node-id",
		"current_action_id": "test-action-id",
		"graph_id":          "test-graph-id",
		"runtime_data":      `{"key":"value"}`,
		"is_authenticated":  true,
		"user_id":           "test-user-id",
		"user_inputs":       `{"username":"testuser"}`,
		"user_attributes":   `{"email":"test@example.com"}`,
	}

	result, err := buildFlowContextFromResultRow(row)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "test-flow-id", result.FlowID)
	assert.Equal(suite.T(), "test-app-id", result.AppID)
	assert.Equal(suite.T(), "test-node-id", *result.CurrentNodeID)
	assert.Equal(suite.T(), "test-action-id", *result.CurrentActionID)
	assert.Equal(suite.T(), "test-graph-id", result.GraphID)
	assert.True(suite.T(), result.IsAuthenticated)
	assert.Equal(suite.T(), "test-user-id", *result.UserID)
}

// Test buildFlowContextFromResultRow with missing required fields
func (suite *FlowStoreTestSuite) TestBuildFlowContextFromResultRow_MissingRequiredField() {
	// Missing flow_id
	row := map[string]interface{}{
		"app_id":   "test-app-id",
		"graph_id": "test-graph-id",
	}

	result, err := buildFlowContextFromResultRow(row)

	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "failed to parse flow_id")
}

// Test buildFlowContextFromResultRow with nil optional fields
func (suite *FlowStoreTestSuite) TestBuildFlowContextFromResultRow_NilOptionalFields() {
	row := map[string]interface{}{
		"flow_id":           "test-flow-id",
		"app_id":            "test-app-id",
		"graph_id":          "test-graph-id",
		"current_node_id":   nil,
		"current_action_id": nil,
		"runtime_data":      nil,
		"is_authenticated":  false,
		"user_id":           nil,
		"user_inputs":       nil,
		"user_attributes":   nil,
	}

	result, err := buildFlowContextFromResultRow(row)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "test-flow-id", result.FlowID)
	assert.Equal(suite.T(), "test-app-id", result.AppID)
	assert.Equal(suite.T(), "test-graph-id", result.GraphID)
	assert.Nil(suite.T(), result.CurrentNodeID)
	assert.Nil(suite.T(), result.CurrentActionID)
	assert.Nil(suite.T(), result.RuntimeData)
	assert.False(suite.T(), result.IsAuthenticated)
	assert.Nil(suite.T(), result.UserID)
	assert.Nil(suite.T(), result.UserInputs)
	assert.Nil(suite.T(), result.UserAttributes)
}

// Test parseOptionalString function
func (suite *FlowStoreTestSuite) TestParseOptionalString() {
	// Test with valid string
	str := "test-string"
	result := parseOptionalString(str)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "test-string", *result)

	// Test with nil
	result = parseOptionalString(nil)
	assert.Nil(suite.T(), result)

	// Test with non-string type
	result = parseOptionalString(123)
	assert.Nil(suite.T(), result)
}

// Test parseBoolean function
func (suite *FlowStoreTestSuite) TestParseBoolean() {
	// Test with boolean true
	result := parseBoolean(true)
	assert.True(suite.T(), result)

	// Test with boolean false
	result = parseBoolean(false)
	assert.False(suite.T(), result)

	// Test with int64 non-zero
	result = parseBoolean(int64(1))
	assert.True(suite.T(), result)

	// Test with int64 zero
	result = parseBoolean(int64(0))
	assert.False(suite.T(), result)

	// Test with nil
	result = parseBoolean(nil)
	assert.False(suite.T(), result)

	// Test with unsupported type
	result = parseBoolean("true")
	assert.False(suite.T(), result)
}

// Test transaction rollback scenario
func (suite *FlowStoreTestSuite) TestExecuteTransaction_RollbackError() {
	ctx := flow.EngineContext{
		FlowID: "test-flow-id",
		AppID:  "test-app-id",
	}

	// Set up mock expectations
	queryError := errors.New("query execution error")
	rollbackError := errors.New("rollback error")

	suite.mockDBProvider.On("GetDBClient", "runtime").Return(suite.mockDBClient, nil)
	suite.mockDBClient.On("BeginTx").Return(suite.mockTx, nil)

	// Mock first query failure
	suite.mockTx.On("Exec", mock.MatchedBy(func(query string) bool {
		return query == QueryCreateFlowContext.Query
	}), mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, queryError)

	suite.mockTx.On("Rollback").Return(rollbackError)

	// Execute
	err := suite.flowStore.StoreFlowContext(ctx)

	// Assertions
	assert.NotNil(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to rollback transaction")

	// Verify expectations
	suite.mockDBProvider.AssertExpectations(suite.T())
	suite.mockDBClient.AssertExpectations(suite.T())
	suite.mockTx.AssertExpectations(suite.T())
}

// Test commit error scenario
func (suite *FlowStoreTestSuite) TestExecuteTransaction_CommitError() {
	ctx := flow.EngineContext{
		FlowID: "test-flow-id",
		AppID:  "test-app-id",
	}

	// Set up mock expectations
	commitError := errors.New("commit error")

	suite.mockDBProvider.On("GetDBClient", "runtime").Return(suite.mockDBClient, nil)
	suite.mockDBClient.On("BeginTx").Return(suite.mockTx, nil)

	// Mock successful query executions
	suite.mockTx.On("Exec", mock.MatchedBy(func(query string) bool {
		return query == QueryCreateFlowContext.Query
	}), mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	suite.mockTx.On("Exec", mock.MatchedBy(func(query string) bool {
		return query == QueryCreateFlowUserData.Query
	}), mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	suite.mockTx.On("Commit").Return(commitError)

	// Execute
	err := suite.flowStore.StoreFlowContext(ctx)

	// Assertions
	assert.NotNil(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to commit transaction")

	// Verify expectations
	suite.mockDBProvider.AssertExpectations(suite.T())
	suite.mockDBClient.AssertExpectations(suite.T())
	suite.mockTx.AssertExpectations(suite.T())
}
