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
	sysutils "github.com/thunder-id/thunderid/internal/system/utils"
)

// Record inserts or refreshes a participant under the upsert query.
func (st *store) Record(ctx context.Context, p Participant) error {
	return withOperationDBClient(st.dbProvider, func(dbClient provider.DBClientInterface) error {
		_, err := dbClient.ExecuteContext(ctx, queryUpsertParticipant,
			p.SessionID, st.deploymentID, p.AppID, p.FirstJoinedAt, p.LastActiveAt)
		if err != nil {
			return fmt.Errorf("failed to record session participant: %w", err)
		}
		return nil
	})
}

// ListBySessionID returns the participants of a session, oldest first.
func (st *store) ListBySessionID(ctx context.Context, sessionID string) ([]Participant, error) {
	var result []Participant

	err := withOperationDBClient(st.dbProvider, func(dbClient provider.DBClientInterface) error {
		results, queryErr := dbClient.QueryContext(ctx, queryListParticipantsBySessionID, sessionID, st.deploymentID)
		if queryErr != nil {
			return fmt.Errorf("failed to execute query: %w", queryErr)
		}
		for _, row := range results {
			p, buildErr := buildParticipantFromRow(row)
			if buildErr != nil {
				return buildErr
			}
			result = append(result, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteBySessionID removes all participants of a session.
func (st *store) DeleteBySessionID(ctx context.Context, sessionID string) error {
	return withOperationDBClient(st.dbProvider, func(dbClient provider.DBClientInterface) error {
		_, err := dbClient.ExecuteContext(ctx, queryDeleteParticipantsBySessionID, sessionID, st.deploymentID)
		if err != nil {
			return fmt.Errorf("failed to delete session participants: %w", err)
		}
		return nil
	})
}

// buildParticipantFromRow maps a database result row into a Participant.
func buildParticipantFromRow(row map[string]interface{}) (Participant, error) {
	sessionID, err := parseString(row["session_id"], "session_id")
	if err != nil {
		return Participant{}, err
	}
	appID, err := parseString(row["app_id"], "app_id")
	if err != nil {
		return Participant{}, err
	}
	firstJoinedAt, err := sysutils.ParseDBTimeField(row["first_joined_at"], "first_joined_at")
	if err != nil {
		return Participant{}, err
	}
	lastActiveAt, err := sysutils.ParseDBTimeField(row["last_active_at"], "last_active_at")
	if err != nil {
		return Participant{}, err
	}
	return Participant{
		SessionID:     sessionID,
		AppID:         appID,
		FirstJoinedAt: firstJoinedAt,
		LastActiveAt:  lastActiveAt,
	}, nil
}
