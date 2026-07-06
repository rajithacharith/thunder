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

package dbstore

import (
	"context"
	"errors"

	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// ErrNotImplemented is returned by all DBStore methods until the implementation is complete.
var ErrNotImplemented = errors.New("dbstore: not implemented")

// dbStore implements the RuntimeStoreProvider interface using the database as the backend.
type dbStore struct {
	deploymentID string
}

func newDBStore(deploymentID string) providers.RuntimeStoreProvider {
	return &dbStore{
		deploymentID: deploymentID,
	}
}

// Put stores a value in the database runtime store with the specified TTL.
func (d *dbStore) Put(ctx context.Context, namespace providers.RuntimeStoreNamespace,
	key string, value []byte, ttlSeconds int64) error {
	// TODO: Implementation for putting data into the database
	return ErrNotImplemented
}

// Get retrieves a value from the database runtime store by its key.
func (d *dbStore) Get(ctx context.Context, namespace providers.RuntimeStoreNamespace,
	key string) ([]byte, error) {
	// TODO: Implementation for getting data from the database
	return nil, ErrNotImplemented
}

// Update updates the value associated with a key in the database runtime store.
func (d *dbStore) Update(ctx context.Context, namespace providers.RuntimeStoreNamespace,
	key string, value []byte) error {
	// TODO: Implementation for updating data in the database
	return ErrNotImplemented
}

// Delete removes a value from the database runtime store by its key.
func (d *dbStore) Delete(ctx context.Context, namespace providers.RuntimeStoreNamespace,
	key string) error {
	// TODO: Implementation for deleting data from the database
	return ErrNotImplemented
}

// Take retrieves and removes a value from the database runtime store by its key.
func (d *dbStore) Take(ctx context.Context, namespace providers.RuntimeStoreNamespace,
	key string) ([]byte, error) {
	// TODO: Implementation for taking data from the database
	return nil, ErrNotImplemented
}
