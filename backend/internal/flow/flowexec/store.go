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
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/asgardeo/thunder/internal/system/config"
	"github.com/asgardeo/thunder/internal/system/database/provider"
)

// flowStoreInterface defines the methods for flow context storage operations.
type flowStoreInterface interface {
	StoreFlowContext(ctx context.Context, engineCtx EngineContext, expirySeconds int64) error
	GetFlowContext(ctx context.Context, flowID string) (*FlowContextWithUserDataDB, error)
	UpdateFlowContext(ctx context.Context, engineCtx EngineContext) error
	DeleteFlowContext(ctx context.Context, flowID string) error
}

// flowStore implements the FlowStoreInterface for managing flow contexts.
type flowStore struct {
	dbProvider   provider.DBProviderInterface
	deploymentID string
}

// newFlowStore creates a new instance of FlowStore.
func newFlowStore(dbProvider provider.DBProviderInterface) flowStoreInterface {
	return &flowStore{
		dbProvider:   dbProvider,
		deploymentID: config.GetThunderRuntime().Config.Server.Identifier,
	}
}

// StoreFlowContext stores the complete flow context in the database.
func (s *flowStore) StoreFlowContext(ctx context.Context, engineCtx EngineContext, expirySeconds int64) error {
	// Convert engine context to database model
	dbModel, err := FromEngineContext(engineCtx)
	if err != nil {
		return fmt.Errorf("failed to convert engine context to database model: %w", err)
	}

	expiryTime := time.Now().UTC().Add(time.Duration(expirySeconds) * time.Second)

	return withRuntimeDBClientContext(ctx, s.dbProvider, func(dbClient provider.DBClientInterface) error {
		if _, err := dbClient.ExecuteContext(ctx, QueryCreateFlowContext,
			dbModel.FlowID, dbModel.AppID, dbModel.Verbose,
			dbModel.CurrentNodeID, dbModel.CurrentAction, dbModel.GraphID,
			dbModel.RuntimeData, dbModel.ExecutionHistory, expiryTime, s.deploymentID); err != nil {
			return err
		}
		_, err := dbClient.ExecuteContext(ctx, QueryCreateFlowUserData, dbModel.FlowID,
			dbModel.IsAuthenticated, dbModel.UserID, dbModel.OrganizationUnitID,
			dbModel.UserType, dbModel.UserInputs, dbModel.UserAttributes, dbModel.Token,
			dbModel.AvailableAttributes, s.deploymentID)
		return err
	})
}

// GetFlowContext retrieves the flow context from the database.
func (s *flowStore) GetFlowContext(ctx context.Context, flowID string) (*FlowContextWithUserDataDB, error) {
	var result *FlowContextWithUserDataDB

	err := withRuntimeDBClientContext(ctx, s.dbProvider, func(dbClient provider.DBClientInterface) error {
		results, err := dbClient.QueryContext(ctx, QueryGetFlowContextWithUserData,
			flowID, s.deploymentID, time.Now().UTC())
		if err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}

		if len(results) == 0 {
			return nil
		}

		if len(results) != 1 {
			return fmt.Errorf("unexpected number of results: %d", len(results))
		}

		row := results[0]
		var buildErr error
		result, buildErr = s.buildFlowContextFromResultRow(row)
		return buildErr
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateFlowContext updates the flow context in the database.
func (s *flowStore) UpdateFlowContext(ctx context.Context, engineCtx EngineContext) error {
	// Convert engine context to database model
	dbModel, err := FromEngineContext(engineCtx)
	if err != nil {
		return fmt.Errorf("failed to convert engine context to database model: %w", err)
	}

	return withRuntimeDBClientContext(ctx, s.dbProvider, func(dbClient provider.DBClientInterface) error {
		if _, err := dbClient.ExecuteContext(ctx, QueryUpdateFlowContext, dbModel.FlowID,
			dbModel.CurrentNodeID, dbModel.CurrentAction, dbModel.RuntimeData, dbModel.ExecutionHistory,
			s.deploymentID); err != nil {
			return err
		}
		_, err := dbClient.ExecuteContext(ctx, QueryUpdateFlowUserData, dbModel.FlowID,
			dbModel.IsAuthenticated, dbModel.UserID, dbModel.OrganizationUnitID, dbModel.UserType,
			dbModel.UserInputs, dbModel.UserAttributes, dbModel.Token,
			dbModel.AvailableAttributes, s.deploymentID)
		return err
	})
}

// DeleteFlowContext removes the flow context from the database.
func (s *flowStore) DeleteFlowContext(ctx context.Context, flowID string) error {
	return withRuntimeDBClientContext(ctx, s.dbProvider, func(dbClient provider.DBClientInterface) error {
		if _, err := dbClient.ExecuteContext(ctx, QueryDeleteFlowUserData, flowID, s.deploymentID); err != nil {
			return err
		}
		_, err := dbClient.ExecuteContext(ctx, QueryDeleteFlowContext, flowID, s.deploymentID)
		return err
	})
}

// withRuntimeDBClientContext is a helper to execute a function with a runtime database client.
func withRuntimeDBClientContext(_ context.Context, dbProvider provider.DBProviderInterface,
	fn func(provider.DBClientInterface) error) error {
	dbClient, err := dbProvider.GetRuntimeDBClient()
	if err != nil {
		return fmt.Errorf("failed to get database client: %w", err)
	}
	return fn(dbClient)
}

// buildFlowContextFromResultRow builds a FlowContextWithUserDataDB from a database result row.
func (s *flowStore) buildFlowContextFromResultRow(row map[string]interface{}) (*FlowContextWithUserDataDB, error) {
	// Parse required fields
	flowID, ok := row["flow_id"].(string)
	if !ok {
		return nil, errors.New("failed to parse flow_id as string")
	}

	appID, ok := row["app_id"].(string)
	if !ok {
		return nil, errors.New("failed to parse app_id as string")
	}

	graphID, ok := row["graph_id"].(string)
	if !ok {
		return nil, errors.New("failed to parse graph_id as string")
	}

	expiryTime, err := s.parseTimeField(row["expiry_time"], "expiry_time")
	if err != nil {
		return nil, err
	}

	// Parse optional fields
	currentNodeID := s.parseOptionalString(row["current_node_id"])
	currentAction := s.parseOptionalString(row["current_action"])
	userID := s.parseOptionalString(row["user_id"])
	organizationUnitID := s.parseOptionalString(row["ou_id"])
	userType := s.parseOptionalString(row["user_type"])
	userInputs := s.parseOptionalString(row["user_inputs"])
	runtimeData := s.parseOptionalString(row["runtime_data"])
	userAttributes := s.parseOptionalString(row["user_attributes"])
	token := s.parseOptionalString(row["token"])
	availableAttributes := s.parseOptionalString(row["available_attributes"])
	executionHistory := s.parseOptionalString(row["execution_history"])

	// Parse boolean fields with type conversion support
	isAuthenticated := s.parseBoolean(row["is_authenticated"])
	verbose := s.parseBoolean(row["verbose"])

	return &FlowContextWithUserDataDB{
		FlowID:              flowID,
		AppID:               appID,
		CurrentNodeID:       currentNodeID,
		CurrentAction:       currentAction,
		GraphID:             graphID,
		RuntimeData:         runtimeData,
		Verbose:             verbose,
		IsAuthenticated:     isAuthenticated,
		UserID:              userID,
		OrganizationUnitID:  organizationUnitID,
		UserType:            userType,
		UserInputs:          userInputs,
		UserAttributes:      userAttributes,
		Token:               token,
		AvailableAttributes: availableAttributes,
		ExecutionHistory:    executionHistory,
		ExpiryTime:          expiryTime,
	}, nil
}

// parseOptionalString safely parses an optional string field from the database row
func (s *flowStore) parseOptionalString(value interface{}) *string {
	if value == nil {
		return nil
	}
	if str, ok := value.(string); ok {
		return &str
	}
	// Handle []byte type (PostgreSQL may return TEXT/JSON as []byte)
	if bytes, ok := value.([]byte); ok {
		str := string(bytes)
		return &str
	}
	return nil
}

// parseBoolean safely parses a boolean field from the database row with type conversion support
func (s *flowStore) parseBoolean(value interface{}) bool {
	if value == nil {
		return false
	}

	if boolVal, ok := value.(bool); ok {
		return boolVal
	}

	if intVal, ok := value.(int64); ok {
		return intVal != 0
	}

	return false
}

// parseTimeField safely parses a time field from the database row handling multiple formats.
// This follows the pattern used in other stores for consistency.
func (s *flowStore) parseTimeField(field interface{}, fieldName string) (time.Time, error) {
	const customTimeFormat = "2006-01-02 15:04:05.999999999"

	switch v := field.(type) {
	case string:
		// Handle SQLite datetime strings
		trimmedTime := s.trimTimeString(v)
		parsedTime, err := time.Parse(customTimeFormat, trimmedTime)
		if err != nil {
			// Try alternative ISO 8601 format as fallback
			parsedTime, err = time.Parse(time.RFC3339, v)
			if err != nil {
				return time.Time{}, fmt.Errorf("error parsing %s: %w", fieldName, err)
			}
		}
		return parsedTime, nil
	case time.Time:
		return v, nil
	case nil:
		return time.Time{}, fmt.Errorf("%s is nil", fieldName)
	default:
		return time.Time{}, fmt.Errorf("unexpected type for %s: %T", fieldName, field)
	}
}

// trimTimeString trims extra information from a time string to match the expected format.
func (s *flowStore) trimTimeString(timeStr string) string {
	parts := strings.SplitN(timeStr, " ", 3)
	if len(parts) >= 2 {
		return parts[0] + " " + parts[1]
	}
	return timeStr
}
