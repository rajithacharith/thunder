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

package userschema

import (
	"context"
	"errors"

	"github.com/asgardeo/thunder/internal/system/cache"
	"github.com/asgardeo/thunder/internal/system/log"
)

// cachedBackedUserSchemaStore wraps a userSchemaStoreInterface with in-memory caching
// for individual schema lookups by ID and Name.
type cachedBackedUserSchemaStore struct {
	schemaByIDCache   cache.CacheInterface[*UserSchema]
	schemaByNameCache cache.CacheInterface[*UserSchema]
	store             userSchemaStoreInterface
	logger            *log.Logger
}

// newCachedBackedUserSchemaStore creates a cache-backed wrapper around the given store.
func newCachedBackedUserSchemaStore(store userSchemaStoreInterface) userSchemaStoreInterface {
	return &cachedBackedUserSchemaStore{
		schemaByIDCache:   cache.GetCache[*UserSchema]("UserSchemaByIDCache"),
		schemaByNameCache: cache.GetCache[*UserSchema]("UserSchemaByNameCache"),
		store:             store,
		logger: log.GetLogger().With(
			log.String(log.LoggerKeyComponentName, "CacheBackedUserSchemaStore")),
	}
}

// GetUserSchemaByID retrieves a user schema by ID, checking cache first.
func (s *cachedBackedUserSchemaStore) GetUserSchemaByID(ctx context.Context, schemaID string) (UserSchema, error) {
	cacheKey := cache.CacheKey{Key: schemaID}
	if cached, ok := s.schemaByIDCache.Get(cacheKey); ok {
		return *cached, nil
	}

	schema, err := s.store.GetUserSchemaByID(ctx, schemaID)
	if err != nil {
		return schema, err
	}

	s.cacheUserSchema(&schema)

	return schema, nil
}

// GetUserSchemaByName retrieves a user schema by name, checking cache first.
func (s *cachedBackedUserSchemaStore) GetUserSchemaByName(ctx context.Context, name string) (UserSchema, error) {
	cacheKey := cache.CacheKey{Key: name}
	if cached, ok := s.schemaByNameCache.Get(cacheKey); ok {
		return *cached, nil
	}

	schema, err := s.store.GetUserSchemaByName(ctx, name)
	if err != nil {
		return schema, err
	}

	s.cacheUserSchema(&schema)

	return schema, nil
}

// CreateUserSchema creates a user schema and populates the cache.
func (s *cachedBackedUserSchemaStore) CreateUserSchema(ctx context.Context, userSchema UserSchema) error {
	if err := s.store.CreateUserSchema(ctx, userSchema); err != nil {
		return err
	}

	s.cacheUserSchema(&userSchema)

	return nil
}

// UpdateUserSchemaByID updates a user schema, invalidates old cache entries, and caches the new state.
func (s *cachedBackedUserSchemaStore) UpdateUserSchemaByID(
	ctx context.Context, schemaID string, userSchema UserSchema,
) error {
	// Fetch existing schema for cache invalidation (check cache first, then store).
	existingCacheKey := cache.CacheKey{Key: schemaID}
	existing, ok := s.schemaByIDCache.Get(existingCacheKey)
	if !ok {
		existingSchema, err := s.store.GetUserSchemaByID(ctx, schemaID)
		if err == nil {
			existing = &existingSchema
		}
	}

	if err := s.store.UpdateUserSchemaByID(ctx, schemaID, userSchema); err != nil {
		return err
	}

	// Invalidate old cache entries.
	if existing != nil {
		s.invalidateUserSchemaCache(existing.ID, existing.Name)
	}

	// Cache the updated schema.
	s.cacheUserSchema(&userSchema)

	return nil
}

// DeleteUserSchemaByID deletes a user schema and invalidates its cache entries.
func (s *cachedBackedUserSchemaStore) DeleteUserSchemaByID(ctx context.Context, schemaID string) error {
	cacheKey := cache.CacheKey{Key: schemaID}
	existing, ok := s.schemaByIDCache.Get(cacheKey)
	if !ok {
		existingSchema, err := s.store.GetUserSchemaByID(ctx, schemaID)
		if err != nil {
			if errors.Is(err, ErrUserSchemaNotFound) {
				return nil
			}
			return err
		}
		existing = &existingSchema
	}

	if err := s.store.DeleteUserSchemaByID(ctx, schemaID); err != nil {
		return err
	}

	if existing != nil {
		s.invalidateUserSchemaCache(existing.ID, existing.Name)
	}

	return nil
}

// Pass-through methods: list/count operations are not cached.

// GetUserSchemaListCount delegates to the underlying store.
func (s *cachedBackedUserSchemaStore) GetUserSchemaListCount(ctx context.Context) (int, error) {
	return s.store.GetUserSchemaListCount(ctx)
}

// GetUserSchemaList delegates to the underlying store.
func (s *cachedBackedUserSchemaStore) GetUserSchemaList(
	ctx context.Context, limit, offset int,
) ([]UserSchemaListItem, error) {
	return s.store.GetUserSchemaList(ctx, limit, offset)
}

// GetUserSchemaListByOUIDs delegates to the underlying store.
func (s *cachedBackedUserSchemaStore) GetUserSchemaListByOUIDs(
	ctx context.Context, ouIDs []string, limit, offset int,
) ([]UserSchemaListItem, error) {
	return s.store.GetUserSchemaListByOUIDs(ctx, ouIDs, limit, offset)
}

// GetUserSchemaListCountByOUIDs delegates to the underlying store.
func (s *cachedBackedUserSchemaStore) GetUserSchemaListCountByOUIDs(
	ctx context.Context, ouIDs []string,
) (int, error) {
	return s.store.GetUserSchemaListCountByOUIDs(ctx, ouIDs)
}

// IsUserSchemaDeclarative delegates to the underlying store.
func (s *cachedBackedUserSchemaStore) IsUserSchemaDeclarative(schemaID string) bool {
	return s.store.IsUserSchemaDeclarative(schemaID)
}

// GetDisplayAttributesByNames delegates to the underlying store.
func (s *cachedBackedUserSchemaStore) GetDisplayAttributesByNames(
	ctx context.Context, names []string,
) (map[string]string, error) {
	return s.store.GetDisplayAttributesByNames(ctx, names)
}

// cacheUserSchema populates both ID and Name caches for the given schema.
func (s *cachedBackedUserSchemaStore) cacheUserSchema(schema *UserSchema) {
	if schema == nil {
		return
	}

	if schema.ID != "" {
		key := cache.CacheKey{Key: schema.ID}
		if err := s.schemaByIDCache.Set(key, schema); err != nil {
			s.logger.Error("Failed to cache user schema by ID",
				log.String("schemaID", schema.ID), log.Error(err))
		}
	}

	if schema.Name != "" {
		key := cache.CacheKey{Key: schema.Name}
		if err := s.schemaByNameCache.Set(key, schema); err != nil {
			s.logger.Error("Failed to cache user schema by name",
				log.String("schemaName", schema.Name), log.Error(err))
		}
	}
}

// invalidateUserSchemaCache removes entries from both ID and Name caches.
func (s *cachedBackedUserSchemaStore) invalidateUserSchemaCache(schemaID, schemaName string) {
	if schemaID != "" {
		key := cache.CacheKey{Key: schemaID}
		if err := s.schemaByIDCache.Delete(key); err != nil {
			s.logger.Error("Failed to invalidate user schema cache by ID",
				log.String("schemaID", schemaID), log.Error(err))
		}
	}

	if schemaName != "" {
		key := cache.CacheKey{Key: schemaName}
		if err := s.schemaByNameCache.Delete(key); err != nil {
			s.logger.Error("Failed to invalidate user schema cache by name",
				log.String("schemaName", schemaName), log.Error(err))
		}
	}
}
