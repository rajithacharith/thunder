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

package actorprovider

import (
	"context"
	"errors"

	"github.com/thunder-id/thunderid/internal/entityprovider"
	"github.com/thunder-id/thunderid/internal/inboundclient"
	inboundmodel "github.com/thunder-id/thunderid/internal/inboundclient/model"
	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
)

// actorProvider delegates actor resolution to inbound-client and entity-provider services.
type actorProvider struct {
	inboundClient  inboundclient.InboundClientServiceInterface
	entityProvider entityprovider.EntityProviderInterface
}

func newActorProvider(
	inboundClient inboundclient.InboundClientServiceInterface,
	entityProvider entityprovider.EntityProviderInterface,
) ActorProviderInterface {
	return &actorProvider{
		inboundClient:  inboundClient,
		entityProvider: entityProvider,
	}
}

// GetOAuthClientByID returns the OAuth client registered for the given ID.
func (p *actorProvider) GetOAuthClientByID(
	ctx context.Context, id string,
) (*inboundmodel.OAuthClient, *serviceerror.ServiceError) {
	client, err := p.inboundClient.GetOAuthClientByClientID(ctx, id)
	if err != nil {
		if errors.Is(err, inboundclient.ErrInboundClientNotFound) {
			return nil, &ErrorActorNotFound
		}
		return nil, &ErrorActorFetchFailed
	}
	return client, nil
}

// GetInboundClientByID returns the inbound-client row for the given ID.
func (p *actorProvider) GetInboundClientByID(
	ctx context.Context, id string,
) (*inboundmodel.InboundClient, *serviceerror.ServiceError) {
	client, err := p.inboundClient.GetInboundClientByEntityID(ctx, id)
	if err != nil {
		if errors.Is(err, inboundclient.ErrInboundClientNotFound) {
			return nil, &ErrorActorNotFound
		}
		return nil, &ErrorActorFetchFailed
	}
	return client, nil
}

// GetActor returns the backing entity record for the given actor ID.
func (p *actorProvider) GetActor(actorID string) (*entityprovider.Entity, *entityprovider.EntityProviderError) {
	return p.entityProvider.GetEntity(actorID)
}

// GetActorGroups returns transitive group memberships for the given actor ID.
func (p *actorProvider) GetActorGroups(
	actorID string,
) ([]entityprovider.EntityGroup, *entityprovider.EntityProviderError) {
	return p.entityProvider.GetTransitiveEntityGroups(actorID)
}
