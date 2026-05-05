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

package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	usages "github.com/asgardeo/thunder/internal/design/usages"
	"github.com/asgardeo/thunder/internal/entityprovider"
	inboundclient "github.com/asgardeo/thunder/internal/inboundclient"
	inboundmodel "github.com/asgardeo/thunder/internal/inboundclient/model"
)

// designUsageResolverAdapter implements usages.ApplicationUsageResolver using the inbound client
// service and entity provider.
type designUsageResolverAdapter struct {
	inboundClientService inboundclient.InboundClientServiceInterface
	entityProvider       entityprovider.EntityProviderInterface
}

// NewDesignUsageResolver returns an ApplicationUsageResolver backed by the inbound client
// service and entity provider.
func NewDesignUsageResolver(
	inboundSvc inboundclient.InboundClientServiceInterface,
	ep entityprovider.EntityProviderInterface,
) usages.ApplicationUsageResolver {
	return &designUsageResolverAdapter{
		inboundClientService: inboundSvc,
		entityProvider:       ep,
	}
}

// GetApplicationRefsByResource returns the applications that reference the given design resource.
func (a *designUsageResolverAdapter) GetApplicationRefsByResource(
	ctx context.Context,
	resourceType usages.DesignUsageType,
	resourceID string,
) ([]usages.ApplicationRef, error) {
	var clients []inboundmodel.InboundClient
	var err error

	switch resourceType {
	case usages.DesignUsageTypeTheme:
		clients, err = a.inboundClientService.GetInboundClientsByThemeID(ctx, resourceID)
	case usages.DesignUsageTypeLayout:
		clients, err = a.inboundClientService.GetInboundClientsByLayoutID(ctx, resourceID)
	case usages.DesignUsageTypeFlow:
		clients, err = a.inboundClientService.GetInboundClientsByFlowID(ctx, resourceID)
	default:
		return nil, errors.New("unsupported resource type")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query inbound clients: %w", err)
	}
	if len(clients) == 0 {
		return []usages.ApplicationRef{}, nil
	}

	entityIDs := make([]string, len(clients))
	for i, c := range clients {
		entityIDs[i] = c.ID
	}

	entities, epErr := a.entityProvider.GetEntitiesByIDs(entityIDs)
	if epErr != nil {
		return nil, fmt.Errorf("failed to get entities: %v", epErr)
	}

	entityMap := make(map[string]*entityprovider.Entity, len(entities))
	for i := range entities {
		entityMap[entities[i].ID] = &entities[i]
	}

	refs := make([]usages.ApplicationRef, 0, len(clients))
	for _, c := range clients {
		ref := usages.ApplicationRef{ID: c.ID}
		if e := entityMap[c.ID]; e != nil {
			var sysAttrs map[string]interface{}
			if len(e.SystemAttributes) > 0 {
				_ = json.Unmarshal(e.SystemAttributes, &sysAttrs)
			}
			if sysAttrs != nil {
				if name, ok := sysAttrs[fieldName].(string); ok {
					ref.Name = name
				}
				if clientID, ok := sysAttrs[fieldClientID].(string); ok {
					ref.ClientID = clientID
				}
			}
		}
		refs = append(refs, ref)
	}
	return refs, nil
}
