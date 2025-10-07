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

// Package utils provides utility functions for flow processing.
package flow

import (
	"fmt"
)

// BuildGraphFromDefinition builds a graph from a graph definition json.
func BuildGraphFromDefinition(definition *GraphDefinition) (GraphInterface, error) {
	if definition == nil || len(definition.Nodes) == 0 {
		return nil, fmt.Errorf("graph definition is nil or has no nodes")
	}

	// Create a graph
	_type, err := getGraphType(definition.Type)
	if err != nil {
		return nil, fmt.Errorf("error while retrieving graph type: %w", err)
	}
	g := NewGraph(definition.ID, _type)

	// Add all nodes to the graph
	edges := make(map[string][]string)
	for _, nodeDef := range definition.Nodes {
		isFinalNode := len(nodeDef.Next) == 0

		// Construct a new node. Here we set isStartNode to false by default.
		node, err := NewNode(nodeDef.ID, nodeDef.Type, false, isFinalNode)
		if err != nil {
			return nil, fmt.Errorf("failed to create node %s: %w", nodeDef.ID, err)
		}

		// Set next nodes if defined
		if len(nodeDef.Next) > 0 {
			node.SetNextNodeList(nodeDef.Next)

			// Store edges based on the node definition
			_, exists := edges[nodeDef.ID]
			if !exists {
				edges[nodeDef.ID] = []string{}
			}
			edges[nodeDef.ID] = append(edges[nodeDef.ID], nodeDef.Next...)
		}

		// Convert and set input data from definition
		inputData := make([]InputData, len(nodeDef.InputData))
		for i, input := range nodeDef.InputData {
			inputData[i] = InputData{
				Name:     input.Name,
				Type:     input.Type,
				Required: input.Required,
			}
		}
		node.SetInputData(inputData)

		// Set the executor config if defined
		if nodeDef.Executor.Name != "" {
			executor, err := getExecutorConfigByName(nodeDef.Executor)
			if err != nil {
				return nil, fmt.Errorf("error while getting executor %s: %w", nodeDef.Executor, err)
			}
			node.SetExecutorConfig(executor)
		} else if nodeDef.Type == string(NodeTypeAuthSuccess) {
			executor, err := getExecutorConfigByName(ExecutorDefinition{
				Name: "AuthAssertExecutor",
			})
			if err != nil {
				return nil, fmt.Errorf("error while getting default AuthAssertExecutor: %w", err)
			}
			node.SetExecutorConfig(executor)
		}

		err = g.AddNode(node)
		if err != nil {
			return nil, fmt.Errorf("failed to add node %s to the graph: %w", nodeDef.ID, err)
		}
	}

	// Set edges in the graph
	for sourceID, targetIDs := range edges {
		for _, targetID := range targetIDs {
			err := g.AddEdge(sourceID, targetID)
			if err != nil {
				return nil, fmt.Errorf("failed to add edge from %s to %s: %w", sourceID, targetID, err)
			}
		}
	}

	// Determine the start node and set it in the graph
	startNodeID := ""
	for _, node := range g.GetNodes() {
		if len(node.GetPreviousNodeList()) == 0 {
			startNodeID = node.GetID()
			break
		}
	}
	if startNodeID == "" {
		return nil, fmt.Errorf("no start node found in the graph definition")
	}

	err = g.SetStartNode(startNodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to set start node ID: %w", err)
	}

	return g, nil
}

// getGraphType retrieves the graph type from a string representation.
func getGraphType(graphType string) (FlowType, error) {
	switch graphType {
	case string(FlowTypeAuthentication):
		return FlowTypeAuthentication, nil
	case string(FlowTypeRegistration):
		return FlowTypeRegistration, nil
	default:
		return "", fmt.Errorf("unsupported graph type: %s", graphType)
	}
}

// getExecutorConfigByName constructs an executor configuration by its definition if it exists.
func getExecutorConfigByName(execDef ExecutorDefinition) (*ExecutorConfig, error) {
	if execDef.Name == "" {
		return nil, fmt.Errorf("executor name cannot be empty")
	}

	// At this point, we assume executors and attached IDPs are already registered in the system.
	// Hence validations will not be done at this point.
	var executor ExecutorConfig
	switch execDef.Name {
	case "BasicAuthExecutor":
		executor = ExecutorConfig{
			Name:    "BasicAuthExecutor",
			IdpName: "Local",
		}
	case "SMSOTPAuthExecutor":
		executor = ExecutorConfig{
			Name:       "SMSOTPAuthExecutor",
			IdpName:    "Local",
			Properties: execDef.Properties,
		}
	case "GithubOAuthExecutor":
		executor = ExecutorConfig{
			Name:       "GithubOAuthExecutor",
			IdpName:    execDef.IdpName,
			Properties: execDef.Properties,
		}
	case "GoogleOIDCAuthExecutor":
		executor = ExecutorConfig{
			Name:       "GoogleOIDCAuthExecutor",
			IdpName:    execDef.IdpName,
			Properties: execDef.Properties,
		}
	case "AttributeCollector":
		executor = ExecutorConfig{
			Name:       "AttributeCollector",
			Properties: execDef.Properties,
		}
	case "ProvisioningExecutor":
		executor = ExecutorConfig{
			Name:       "ProvisioningExecutor",
			Properties: execDef.Properties,
		}
	case "AuthAssertExecutor":
		executor = ExecutorConfig{
			Name: "AuthAssertExecutor",
		}
	default:
		return nil, fmt.Errorf("executor with name %s not found", execDef.Name)
	}

	if executor.Name == "" {
		return nil, fmt.Errorf("executor with name %s could not be created", execDef.Name)
	}

	return &executor, nil
}
