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

// Package event provides event models and types for the analytics system.
package event

import (
	"time"

	"github.com/thunder-id/thunderid/internal/system/utils"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// NewEvent creates a new Event with required fields.
// Additional data should be added using WithData().
func NewEvent(traceID string, eventType string, component string) *providers.Event {
	eventID, err := utils.GenerateUUIDv7()
	if err != nil {
		return &providers.Event{}
	}

	return &providers.Event{
		TraceID:   traceID,
		EventID:   eventID,
		Type:      eventType,
		Timestamp: time.Now(),
		Component: component,
		Status:    providers.StatusInProgress,
		Data:      make(map[string]interface{}),
	}
}
