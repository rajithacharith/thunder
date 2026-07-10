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

package session

import (
	"context"
	"fmt"

	"github.com/thunder-id/thunderid/internal/system/database/provider"
)

// CreateContext persists the session context. It rejects payloads exceeding MaxSessionContextBytes.
func (st *store) CreateContext(ctx context.Context, c SessionContext) error {
	payload, err := c.serializePayload()
	if err != nil {
		return err
	}
	if len(payload) > MaxSessionContextBytes {
		return errSessionContextTooLarge
	}

	return withOperationDBClient(st.dbProvider, func(dbClient provider.DBClientInterface) error {
		_, execErr := dbClient.ExecuteContext(ctx, queryCreateSessionContext,
			c.SessionID, st.deploymentID, c.CheckpointID, payload, c.ContextVersion)
		if execErr != nil {
			return fmt.Errorf("failed to create session context: %w", execErr)
		}
		return nil
	})
}

// GetByCheckpoint fetches one checkpoint's session context.
func (st *store) GetByCheckpoint(ctx context.Context, sessionID,
	checkpointID string) (*SessionContext, error) {
	var result *SessionContext

	err := withOperationDBClient(st.dbProvider, func(dbClient provider.DBClientInterface) error {
		results, queryErr := dbClient.QueryContext(ctx, queryGetSessionContextByCheckpoint,
			sessionID, st.deploymentID, checkpointID)
		if queryErr != nil {
			return fmt.Errorf("failed to execute query: %w", queryErr)
		}
		if len(results) == 0 {
			return nil
		}
		if len(results) != 1 {
			return fmt.Errorf("unexpected number of results: %d", len(results))
		}

		c, buildErr := st.buildSessionContextFromRow(results[0])
		if buildErr != nil {
			return buildErr
		}
		result = c
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ListCheckpointIDs returns the checkpoint ids a session has saved, without decrypting any payload.
func (st *store) ListCheckpointIDs(ctx context.Context, sessionID string) ([]string, error) {
	var ids []string

	err := withOperationDBClient(st.dbProvider, func(dbClient provider.DBClientInterface) error {
		results, queryErr := dbClient.QueryContext(ctx, queryListCheckpointsBySessionID, sessionID, st.deploymentID)
		if queryErr != nil {
			return fmt.Errorf("failed to execute query: %w", queryErr)
		}
		for _, row := range results {
			id, parseErr := parseString(row["checkpoint_id"], "checkpoint_id")
			if parseErr != nil {
				return parseErr
			}
			ids = append(ids, id)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// Delete removes a session's session context.
func (st *store) Delete(ctx context.Context, sessionID string) error {
	return withOperationDBClient(st.dbProvider, func(dbClient provider.DBClientInterface) error {
		_, err := dbClient.ExecuteContext(ctx, queryDeleteSessionContext, sessionID, st.deploymentID)
		if err != nil {
			return fmt.Errorf("failed to delete session context: %w", err)
		}
		return nil
	})
}

// buildSessionContextFromRow parses a result row into an SessionContext.
func (st *store) buildSessionContextFromRow(row map[string]interface{}) (*SessionContext, error) {
	sessionID, err := parseString(row["session_id"], "session_id")
	if err != nil {
		return nil, err
	}
	checkpointID, err := parseString(row["checkpoint_id"], "checkpoint_id")
	if err != nil {
		return nil, err
	}
	contextVersion, err := parseInt(row["context_version"], "context_version")
	if err != nil {
		return nil, err
	}
	payload, err := parseSessionContextPayload(parseNullableString(row["context"]))
	if err != nil {
		return nil, err
	}

	return &SessionContext{
		SessionID:      sessionID,
		CheckpointID:   checkpointID,
		RuntimeData:    payload.RuntimeData,
		AuthUser:       payload.AuthUser,
		CompletedSteps: payload.CompletedSteps,
		ContextVersion: contextVersion,
	}, nil
}
