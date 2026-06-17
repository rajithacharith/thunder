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
	"github.com/thunder-id/thunderid/internal/system/log"
)

// actorProvider delegates actor resolution to inbound-client and entity-provider services.
type actorProvider struct {
	inboundClient  inboundclient.InboundClientServiceInterface
	entityProvider entityprovider.EntityProviderInterface
	logger         *log.Logger
}

// newActorProvider creates a new actorProvider backed by the given inbound-client and entity-provider.
func newActorProvider(
	inboundClient inboundclient.InboundClientServiceInterface,
	entityProvider entityprovider.EntityProviderInterface,
) ActorProviderInterface {
	return &actorProvider{
		inboundClient:  inboundClient,
		entityProvider: entityProvider,
		logger:         log.GetLogger().With(log.String(log.LoggerKeyComponentName, "ActorProvider")),
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
		p.logger.Error(ctx, "Failed to fetch OAuth client", log.String("id", id), log.Error(err))
		return nil, &serviceerror.InternalServerError
	}
	return client, nil
}

// GetOAuthProfileByID returns the stored OAuth profile for the given entity UUID.
func (p *actorProvider) GetOAuthProfileByID(
	ctx context.Context, id string,
) (*inboundmodel.OAuthProfile, *serviceerror.ServiceError) {
	profile, err := p.inboundClient.GetOAuthProfileByEntityID(ctx, id)
	if err != nil {
		if errors.Is(err, inboundclient.ErrInboundClientNotFound) {
			return nil, &ErrorActorNotFound
		}
		p.logger.Error(ctx, "Failed to fetch OAuth profile by entity ID",
			log.String("id", id), log.Error(err))
		return nil, &serviceerror.InternalServerError
	}
	return profile, nil
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
		p.logger.Error(ctx, "Failed to fetch inbound client", log.String("id", id), log.Error(err))
		return nil, &serviceerror.InternalServerError
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
