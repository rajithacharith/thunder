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

package authz

// Subject identifies the principal for an access evaluation.
type Subject struct {
	Type       string                 `json:"type,omitempty"`
	ID         string                 `json:"id"`
	GroupIDs   []string               `json:"groupIds,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// ResourceServer identifies the resource server for an access evaluation.
type ResourceServer struct {
	Handle     string                 `json:"handle"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// Permission identifies the permission string being evaluated.
type Permission struct {
	Name       string                 `json:"name"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// AccessEvaluationRequest represents a single fine-grained access evaluation request.
type AccessEvaluationRequest struct {
	Subject        Subject                `json:"subject"`
	ResourceServer ResourceServer         `json:"resourceServer"`
	Permission     Permission             `json:"permission"`
	Context        map[string]interface{} `json:"context,omitempty"`
}

// AccessEvaluationResponse represents a single fine-grained access evaluation response.
type AccessEvaluationResponse struct {
	Decision bool                   `json:"decision"`
	Context  map[string]interface{} `json:"context,omitempty"`
}

// AccessEvaluationsRequest represents a batched fine-grained access evaluation request.
type AccessEvaluationsRequest struct {
	Evaluations []AccessEvaluationRequest `json:"evaluations"`
}

// AccessEvaluationsResponse represents a batched fine-grained access evaluation response.
type AccessEvaluationsResponse struct {
	Evaluations []AccessEvaluationResponse `json:"evaluations"`
}
