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

// Package usages provides functionality for resolving the applications that reference a
// given design resource (theme, layout, or flow).
package usages

import (
	"context"

	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
)

// DesignUsageType identifies the kind of design resource whose usages are being queried.
type DesignUsageType string

const (
	DesignUsageTypeTheme  DesignUsageType = "THEME"
	DesignUsageTypeLayout DesignUsageType = "LAYOUT"
	DesignUsageTypeFlow   DesignUsageType = "FLOW"
)

// ApplicationRef is a slim representation of an application that references a design resource.
type ApplicationRef struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ClientID string `json:"clientId,omitempty"`
}

// ApplicationUsageResponse is the response envelope returned by the usages endpoint.
type ApplicationUsageResponse struct {
	TotalResults int              `json:"totalResults"`
	Count        int              `json:"count"`
	Applications []ApplicationRef `json:"applications"`
}

// ApplicationUsageResolver fetches applications referencing a given design resource.
// The interface is owned by this package (consumer) and implemented by the application package
// (provider), following the Scenario A boundary rule.
type ApplicationUsageResolver interface {
	GetApplicationRefsByResource(
		ctx context.Context,
		resourceType DesignUsageType,
		resourceID string,
	) ([]ApplicationRef, error)
}

// ResourceExistenceChecker verifies that a design resource with the given ID exists.
// Used to distinguish "resource not found" (404) from "resource has no usages" (200 empty list).
type ResourceExistenceChecker interface {
	ResourceExists(
		ctx context.Context,
		resourceType DesignUsageType,
		resourceID string,
	) (bool, *serviceerror.ServiceError)
}
