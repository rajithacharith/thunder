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
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	oupkg "github.com/asgardeo/thunder/internal/ou"
)

func TestOUUserResolver_GetUserCountByOUID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		store := newUserStoreInterfaceMock(t)
		store.On("GetUserListCountByOUIDs", context.Background(), []string{"ou-1"}, (map[string]interface{})(nil)).
			Return(5, nil).Once()

		resolver := newOUUserResolver(store)
		count, err := resolver.GetUserCountByOUID(context.Background(), "ou-1")

		require.NoError(t, err)
		require.Equal(t, 5, count)
	})

	t.Run("store error", func(t *testing.T) {
		store := newUserStoreInterfaceMock(t)
		store.On("GetUserListCountByOUIDs", context.Background(), []string{"ou-1"}, (map[string]interface{})(nil)).
			Return(0, errors.New("db error")).Once()

		resolver := newOUUserResolver(store)
		count, err := resolver.GetUserCountByOUID(context.Background(), "ou-1")

		require.Error(t, err)
		require.Equal(t, 0, count)
	})
}

func TestOUUserResolver_GetUserListByOUID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		store := newUserStoreInterfaceMock(t)
		store.On("GetUserListByOUIDs", context.Background(), []string{"ou-1"}, 10, 0, (map[string]interface{})(nil)).
			Return([]User{{ID: "user-1"}, {ID: "user-2"}}, nil).Once()

		resolver := newOUUserResolver(store)
		users, err := resolver.GetUserListByOUID(context.Background(), "ou-1", 10, 0)

		require.NoError(t, err)
		require.Len(t, users, 2)
		require.Equal(t, oupkg.User{ID: "user-1"}, users[0])
		require.Equal(t, oupkg.User{ID: "user-2"}, users[1])
	})

	t.Run("store error", func(t *testing.T) {
		store := newUserStoreInterfaceMock(t)
		store.On("GetUserListByOUIDs", context.Background(), []string{"ou-1"}, 10, 0, (map[string]interface{})(nil)).
			Return([]User(nil), errors.New("db error")).Once()

		resolver := newOUUserResolver(store)
		users, err := resolver.GetUserListByOUID(context.Background(), "ou-1", 10, 0)

		require.Error(t, err)
		require.Nil(t, users)
	})

	t.Run("empty results", func(t *testing.T) {
		store := newUserStoreInterfaceMock(t)
		store.On("GetUserListByOUIDs", context.Background(), []string{"ou-1"}, 10, 0, (map[string]interface{})(nil)).
			Return([]User{}, nil).Once()

		resolver := newOUUserResolver(store)
		users, err := resolver.GetUserListByOUID(context.Background(), "ou-1", 10, 0)

		require.NoError(t, err)
		require.Empty(t, users)
	})
}
