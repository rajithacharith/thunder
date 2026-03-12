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

// OU resolve from strategy values.
const (
	// ouResolveFromCaller indicates that the caller's OU should be used when creating the user.
	ouResolveFromCaller = "caller"
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
// It reads the "resolveFrom" node property to determine the OU resolution strategy.
// When set to "caller", it overrides the default OU with the caller's OU from the security context.
func (e *ouResolverExecutor) Execute(ctx *core.NodeContext) (*common.ExecutorResponse, error) {
	logger := e.logger.With(log.String(log.LoggerKeyFlowID, ctx.FlowID))

	execResp := &common.ExecutorResponse{
		Status:      common.ExecComplete,
		RuntimeData: make(map[string]string),
	}

	resolveFrom := e.getResolveFrom(ctx)
	if resolveFrom == "" {
		logger.Debug("resolveFrom not configured, skipping OU override")
		return execResp, nil
	}

	switch resolveFrom {
	case ouResolveFromCaller:
		return e.resolveFromCaller(ctx, execResp, logger)
	default:
		logger.Error("Unsupported resolveFrom value", log.String("resolveFrom", resolveFrom))
		execResp.Status = common.ExecFailure
		execResp.FailureReason = "Unsupported OU resolution strategy: " + resolveFrom
		return execResp, nil
	}
}

// resolveFromCaller resolves the OU from the caller's security context.
func (e *ouResolverExecutor) resolveFromCaller(ctx *core.NodeContext,
	execResp *common.ExecutorResponse, logger *log.Logger) (*common.ExecutorResponse, error) {
	callerOUID := security.GetOUID(ctx.Context)
	if callerOUID == "" {
		logger.Error("Caller OU not found in security context")
		execResp.Status = common.ExecFailure
		execResp.FailureReason = "Unable to determine caller organization unit"
		return execResp, nil
	}

	logger.Debug("Overriding user OU with caller's OU", log.String("callerOUID", callerOUID))
	execResp.RuntimeData[ouIDKey] = callerOUID

	return execResp, nil
}

// getResolveFrom retrieves the resolveFrom strategy from the node properties.
func (e *ouResolverExecutor) getResolveFrom(ctx *core.NodeContext) string {
	if ctx.NodeProperties == nil {
		return ""
	}
	val, ok := ctx.NodeProperties[common.NodePropertyOUResolveFrom]
	if !ok {
		return ""
	}
	strVal, ok := val.(string)
	if !ok {
		return ""
	}
	return strVal
}
