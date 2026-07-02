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

package core

import (
	"github.com/thunder-id/thunderid/internal/flow/common"
	"github.com/thunder-id/thunderid/internal/system/log"
	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// CallNodeInterface extends NodeInterface for CALL nodes, which transfer execution to
// a referenced flow and return control to the caller when the callee's END node is reached.
type CallNodeInterface interface {
	NodeInterface
	GetReferencedFlow() string
	SetReferencedFlow(flowID string)
	GetOnSuccess() string
	SetOnSuccess(nodeID string)
	GetOnFailure() string
	SetOnFailure(nodeID string)
}

// callNode implements CallNodeInterface and represents a CALL node in the flow graph.
type callNode struct {
	*node
	referencedFlow string
	onSuccess      string
	onFailure      string
	logger         *log.Logger
}

var _ CallNodeInterface = (*callNode)(nil)

// newCallNode creates a new instance of callNode with the given parameters.
func newCallNode(id string, properties map[string]interface{}, isStartNode, isFinalNode bool) NodeInterface {
	if properties == nil {
		properties = make(map[string]interface{})
	}
	return &callNode{
		node: &node{
			id:               id,
			_type:            common.NodeTypeCall,
			properties:       properties,
			isStartNode:      isStartNode,
			isFinalNode:      isFinalNode,
			nextNodeList:     []string{},
			previousNodeList: []string{},
		},
		logger: log.GetLogger().With(log.String(log.LoggerKeyComponentName, "CallNode"),
			log.String(log.LoggerKeyNodeID, id)),
	}
}

// Execute executes the CALL node logic, transferring control to the referenced flow.
func (n *callNode) Execute(ctx *providers.NodeContext) (*common.NodeResponse, *tidcommon.ServiceError) {
	if n.referencedFlow == "" {
		n.logger.Error(ctx.Context, "Referenced flow ID is not set for CALL node")
		return nil, &tidcommon.InternalServerError
	}
	return &common.NodeResponse{
		Status:           common.NodeStatusCall,
		CallTargetFlowID: n.referencedFlow,
	}, nil
}

// GetReferencedFlow returns the ID of the flow referenced by this CALL node.
func (n *callNode) GetReferencedFlow() string {
	return n.referencedFlow
}

// SetReferencedFlow sets the ID of the flow to be referenced by this CALL node.
func (n *callNode) SetReferencedFlow(flowID string) {
	n.referencedFlow = flowID
}

// GetOnSuccess returns the ID of the node to transition to upon successful completion of the referenced flow.
func (n *callNode) GetOnSuccess() string {
	return n.onSuccess
}

// SetOnSuccess sets the ID of the node to transition to upon successful completion of the referenced flow.
func (n *callNode) SetOnSuccess(nodeID string) {
	n.onSuccess = nodeID
}

// GetOnFailure returns the ID of the node to transition to upon failure of the referenced flow.
func (n *callNode) GetOnFailure() string {
	return n.onFailure
}

// SetOnFailure sets the ID of the node to transition to upon failure of the referenced flow.
func (n *callNode) SetOnFailure(nodeID string) {
	n.onFailure = nodeID
}
