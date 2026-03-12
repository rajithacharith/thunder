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

package executor

import (
	"github.com/asgardeo/thunder/internal/flow/common"
	"github.com/asgardeo/thunder/internal/flow/core"
	"github.com/asgardeo/thunder/internal/system/log"
	"github.com/asgardeo/thunder/internal/system/security"
)

// ouResolverExecutor resolves the organization unit for a user being onboarded.
type ouResolverExecutor struct {
	core.ExecutorInterface
	logger *log.Logger
}

// newOUResolverExecutor creates a new OU resolver executor.
func newOUResolverExecutor(flowFactory core.FlowFactoryInterface) *ouResolverExecutor {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "OUResolverExecutor"))
	base := flowFactory.CreateExecutor(
		ExecutorNameOUResolver,
		common.ExecutorTypeUtility,
		[]common.Input{},
		[]common.Input{},
	)
	return &ouResolverExecutor{
		ExecutorInterface: base,
		logger:            logger,
	}
}

// Execute resolves the organization unit for the user being onboarded.
// If the createInAdminOU node property is true, it sets the ouId runtime data
// to the admin's organization unit, overriding the default OU from the user type.
func (e *ouResolverExecutor) Execute(ctx *core.NodeContext) (*common.ExecutorResponse, error) {
	logger := e.logger.With(log.String(log.LoggerKeyFlowID, ctx.FlowID))

	execResp := &common.ExecutorResponse{
		Status:      common.ExecComplete,
		RuntimeData: make(map[string]string),
	}

	// Check if the node property requests creating the user in the admin's OU.
	if !e.shouldCreateInAdminOU(ctx) {
		logger.Debug("createInAdminOU not enabled, skipping OU override")
		return execResp, nil
	}

	// Extract the admin's OU from the security context.
	adminOUID := security.GetOUID(ctx.Context)
	if adminOUID == "" {
		logger.Error("Admin OU not found in security context")
		execResp.Status = common.ExecFailure
		execResp.FailureReason = "Unable to determine admin organization unit"
		return execResp, nil
	}

	logger.Debug("Overriding user OU with admin's OU", log.String("adminOUID", adminOUID))
	execResp.RuntimeData[ouIDKey] = adminOUID

	return execResp, nil
}

// shouldCreateInAdminOU checks if the node property "createInAdminOU" is set to true.
func (e *ouResolverExecutor) shouldCreateInAdminOU(ctx *core.NodeContext) bool {
	if ctx.NodeProperties == nil {
		return false
	}
	val, ok := ctx.NodeProperties[common.NodePropertyCreateInAdminOU]
	if !ok {
		return false
	}
	boolVal, ok := val.(bool)
	return ok && boolVal
}
