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

package main

import (
	"context"

	layoutmgt "github.com/asgardeo/thunder/internal/design/layout/mgt"
	thememgt "github.com/asgardeo/thunder/internal/design/theme/mgt"
	usages "github.com/asgardeo/thunder/internal/design/usages"
	flowmgt "github.com/asgardeo/thunder/internal/flow/mgt"
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
)

// designResourceExistenceAdapter implements usages.ResourceExistenceChecker by delegating to
// the individual design resource services.
type designResourceExistenceAdapter struct {
	themeSvc  thememgt.ThemeMgtServiceInterface
	layoutSvc layoutmgt.LayoutMgtServiceInterface
	flowSvc   flowmgt.FlowMgtServiceInterface
}

// ResourceExists checks whether the given design resource ID exists.
func (a *designResourceExistenceAdapter) ResourceExists(
	ctx context.Context,
	resourceType usages.DesignUsageType,
	resourceID string,
) (bool, *serviceerror.ServiceError) {
	switch resourceType {
	case usages.DesignUsageTypeTheme:
		exists, svcErr := a.themeSvc.IsThemeExist(resourceID)
		if svcErr != nil {
			return false, svcErr
		}
		return exists, nil
	case usages.DesignUsageTypeLayout:
		exists, svcErr := a.layoutSvc.IsLayoutExist(resourceID)
		if svcErr != nil {
			return false, svcErr
		}
		return exists, nil
	case usages.DesignUsageTypeFlow:
		_, svcErr := a.flowSvc.GetFlow(ctx, resourceID)
		if svcErr == nil {
			return true, nil
		}
		if svcErr.Code == flowmgt.ErrorFlowNotFound.Code {
			return false, nil
		}
		return false, svcErr
	default:
		return false, &serviceerror.InternalServerError
	}
}
