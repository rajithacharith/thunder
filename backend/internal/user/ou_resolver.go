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

package user

import (
	"context"

	oupkg "github.com/asgardeo/thunder/internal/ou"
)

// ouUserResolverAdapter implements oupkg.OUUserResolver using the user store.
// This adapter allows the OU package to query user data without directly
// accessing the USER table, breaking the cross-DB access boundary.
type ouUserResolverAdapter struct {
	store userStoreInterface
}

// newOUUserResolver creates a new OUUserResolver backed by the given user store.
func newOUUserResolver(store userStoreInterface) oupkg.OUUserResolver {
	return &ouUserResolverAdapter{store: store}
}

// GetUserCountByOUID returns the count of users belonging to the given organization unit.
func (a *ouUserResolverAdapter) GetUserCountByOUID(ctx context.Context, ouID string) (int, error) {
	return a.store.GetUserListCountByOUIDs(ctx, []string{ouID}, nil)
}

// GetUserListByOUID returns a paginated list of users belonging to the given organization unit.
func (a *ouUserResolverAdapter) GetUserListByOUID(
	ctx context.Context, ouID string, limit, offset int,
) ([]oupkg.User, error) {
	users, err := a.store.GetUserListByOUIDs(ctx, []string{ouID}, limit, offset, nil)
	if err != nil {
		return nil, err
	}

	result := make([]oupkg.User, len(users))
	for i, u := range users {
		result[i] = oupkg.User{
			ID:         u.ID,
			Type:       u.Type,
			Attributes: u.Attributes,
		}
	}

	return result, nil
}
