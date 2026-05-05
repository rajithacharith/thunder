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

package usages

import (
	"context"

	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	"github.com/asgardeo/thunder/internal/system/log"
)

const serviceLogger = "DesignUsageService"

// DesignUsageServiceInterface defines the service for querying design resource usages.
type DesignUsageServiceInterface interface {
	GetResourceUsages(
		ctx context.Context,
		resourceType DesignUsageType,
		resourceID string,
	) (*ApplicationUsageResponse, *serviceerror.ServiceError)
}

// designUsageService is the default implementation of DesignUsageServiceInterface.
type designUsageService struct {
	resolver         ApplicationUsageResolver
	existenceChecker ResourceExistenceChecker
	logger           *log.Logger
}

// newDesignUsageService creates a new designUsageService with injected dependencies.
func newDesignUsageService(
	resolver ApplicationUsageResolver,
	existenceChecker ResourceExistenceChecker,
) DesignUsageServiceInterface {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, serviceLogger))
	return &designUsageService{
		resolver:         resolver,
		existenceChecker: existenceChecker,
		logger:           logger,
	}
}

// GetResourceUsages returns the applications that reference the given design resource.
func (s *designUsageService) GetResourceUsages(
	ctx context.Context,
	resourceType DesignUsageType,
	resourceID string,
) (*ApplicationUsageResponse, *serviceerror.ServiceError) {
	if resourceType == "" {
		return nil, &ErrorInvalidUsageType
	}

	if resourceID == "" {
		return nil, &ErrorMissingResourceID
	}

	switch resourceType {
	case DesignUsageTypeTheme, DesignUsageTypeLayout, DesignUsageTypeFlow:
		// supported
	default:
		return nil, &ErrorUnsupportedUsageType
	}

	if s.existenceChecker == nil {
		s.logger.Error("Resource existence checker is not configured")
		return nil, &serviceerror.InternalServerError
	}

	exists, svcErr := s.existenceChecker.ResourceExists(ctx, resourceType, resourceID)
	if svcErr != nil {
		s.logger.Error("Failed to check resource existence",
			log.String("type", string(resourceType)),
			log.String("id", resourceID))
		return nil, &serviceerror.InternalServerError
	}
	if !exists {
		return nil, &ErrorResourceNotFound
	}

	if s.resolver == nil {
		s.logger.Error("Application usage resolver is not configured")
		return nil, &serviceerror.InternalServerError
	}

	refs, err := s.resolver.GetApplicationRefsByResource(ctx, resourceType, resourceID)
	if err != nil {
		s.logger.Error("Failed to fetch application usages",
			log.String("type", string(resourceType)),
			log.String("id", resourceID),
			log.Error(err))
		return nil, &serviceerror.InternalServerError
	}

	if refs == nil {
		refs = []ApplicationRef{}
	}

	return &ApplicationUsageResponse{
		TotalResults: len(refs),
		Count:        len(refs),
		Applications: refs,
	}, nil
}
