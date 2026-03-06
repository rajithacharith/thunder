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
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/asgardeo/thunder/internal/system/cache"
	"github.com/asgardeo/thunder/internal/system/log"
	"github.com/asgardeo/thunder/tests/mocks/cachemock"
)

// CacheBackedStoreTestSuite tests the cachedBackedUserSchemaStore.
type CacheBackedStoreTestSuite struct {
	suite.Suite
	mockStore         *userSchemaStoreInterfaceMock
	schemaByIDCache   *cachemock.CacheInterfaceMock[*UserSchema]
	schemaByNameCache *cachemock.CacheInterfaceMock[*UserSchema]
	cachedStore       *cachedBackedUserSchemaStore
	// Helper maps to track cached values for verification.
	schemaByIDData   map[string]*UserSchema
	schemaByNameData map[string]*UserSchema
}

func TestCacheBackedStoreTestSuite(t *testing.T) {
	suite.Run(t, new(CacheBackedStoreTestSuite))
}

func (s *CacheBackedStoreTestSuite) SetupTest() {
	s.mockStore = newUserSchemaStoreInterfaceMock(s.T())
	s.schemaByIDData = make(map[string]*UserSchema)
	s.schemaByNameData = make(map[string]*UserSchema)

	s.schemaByIDCache = cachemock.NewCacheInterfaceMock[*UserSchema](s.T())
	s.schemaByNameCache = cachemock.NewCacheInterfaceMock[*UserSchema](s.T())

	setupCacheMock(s.schemaByIDCache, s.schemaByIDData)
	setupCacheMock(s.schemaByNameCache, s.schemaByNameData)

	s.schemaByIDCache.EXPECT().IsEnabled().Return(true).Maybe()
	s.schemaByNameCache.EXPECT().IsEnabled().Return(true).Maybe()

	s.cachedStore = &cachedBackedUserSchemaStore{
		schemaByIDCache:   s.schemaByIDCache,
		schemaByNameCache: s.schemaByNameCache,
		store:             s.mockStore,
		logger: log.GetLogger().With(
			log.String(log.LoggerKeyComponentName, "CacheBackedUserSchemaStore")),
	}
}

// setupCacheMock configures a cache mock to track Set/Get/Delete operations.
func setupCacheMock[T any](
	mockCache *cachemock.CacheInterfaceMock[T],
	data map[string]T,
) {
	mockCache.EXPECT().Set(mock.Anything, mock.Anything).
		RunAndReturn(func(key cache.CacheKey, value T) error {
			data[key.Key] = value
			return nil
		}).Maybe()

	mockCache.EXPECT().Get(mock.Anything).
		RunAndReturn(func(key cache.CacheKey) (T, bool) {
			if val, ok := data[key.Key]; ok {
				return val, true
			}
			var zero T
			return zero, false
		}).Maybe()

	mockCache.EXPECT().Delete(mock.Anything).
		RunAndReturn(func(key cache.CacheKey) error {
			delete(data, key.Key)
			return nil
		}).Maybe()

	mockCache.EXPECT().Clear().
		RunAndReturn(func() error {
			for k := range data {
				delete(data, k)
			}
			return nil
		}).Maybe()

	mockCache.EXPECT().GetName().Return("mockCache").Maybe()

	mockCache.EXPECT().CleanupExpired().Maybe()
}

// assertSchemaCachedByIDAndName verifies the schema is cached in both ID and Name caches.
func (s *CacheBackedStoreTestSuite) assertSchemaCachedByIDAndName(schema UserSchema) {
	cachedByID, ok := s.schemaByIDCache.Get(cache.CacheKey{Key: schema.ID})
	s.True(ok)
	s.Equal(schema.ID, cachedByID.ID)

	cachedByName, ok := s.schemaByNameCache.Get(cache.CacheKey{Key: schema.Name})
	s.True(ok)
	s.Equal(schema.Name, cachedByName.Name)
}

// createTestSchema returns a test user schema.
func (s *CacheBackedStoreTestSuite) createTestSchema() UserSchema {
	return UserSchema{
		ID:                    "schema-1",
		Name:                  "TestSchema",
		OrganizationUnitID:    "ou-1",
		AllowSelfRegistration: true,
		SystemAttributes:      &SystemAttributes{Display: "email"},
		Schema:                json.RawMessage(`{"email":{"type":"string"}}`),
	}
}

// TestNewCachedBackedUserSchemaStore verifies suite setup.
func (s *CacheBackedStoreTestSuite) TestNewCachedBackedUserSchemaStore() {
	s.NotNil(s.cachedStore)
	s.IsType(&cachedBackedUserSchemaStore{}, s.cachedStore)
	s.NotNil(s.cachedStore.schemaByIDCache)
	s.NotNil(s.cachedStore.schemaByNameCache)
	s.NotNil(s.cachedStore.store)
}

// GetUserSchemaByID tests

func (s *CacheBackedStoreTestSuite) TestGetUserSchemaByID_CacheHit() {
	schema := s.createTestSchema()
	s.schemaByIDData[schema.ID] = &schema

	result, err := s.cachedStore.GetUserSchemaByID(context.Background(), schema.ID)
	s.Nil(err)
	s.Equal(schema.ID, result.ID)
	s.Equal(schema.Name, result.Name)
	// Store should NOT be called on cache hit.
	s.mockStore.AssertNotCalled(s.T(), "GetUserSchemaByID")
}

func (s *CacheBackedStoreTestSuite) TestGetUserSchemaByID_CacheMiss() {
	schema := s.createTestSchema()
	s.mockStore.On("GetUserSchemaByID", mock.Anything, schema.ID).Return(schema, nil).Once()

	result, err := s.cachedStore.GetUserSchemaByID(context.Background(), schema.ID)
	s.Nil(err)
	s.Equal(schema.ID, result.ID)
	s.mockStore.AssertExpectations(s.T())
	s.assertSchemaCachedByIDAndName(schema)
}

func (s *CacheBackedStoreTestSuite) TestGetUserSchemaByID_StoreError() {
	storeErr := errors.New("db error")
	s.mockStore.On("GetUserSchemaByID", mock.Anything, "bad-id").Return(UserSchema{}, storeErr).Once()

	_, err := s.cachedStore.GetUserSchemaByID(context.Background(), "bad-id")
	s.Equal(storeErr, err)

	// Verify nothing cached.
	_, ok := s.schemaByIDCache.Get(cache.CacheKey{Key: "bad-id"})
	s.False(ok)
}

// GetUserSchemaByName tests

func (s *CacheBackedStoreTestSuite) TestGetUserSchemaByName_CacheHit() {
	schema := s.createTestSchema()
	s.schemaByNameData[schema.Name] = &schema

	result, err := s.cachedStore.GetUserSchemaByName(context.Background(), schema.Name)
	s.Nil(err)
	s.Equal(schema.Name, result.Name)
	s.mockStore.AssertNotCalled(s.T(), "GetUserSchemaByName")
}

func (s *CacheBackedStoreTestSuite) TestGetUserSchemaByName_CacheMiss() {
	schema := s.createTestSchema()
	s.mockStore.On("GetUserSchemaByName", mock.Anything, schema.Name).Return(schema, nil).Once()

	result, err := s.cachedStore.GetUserSchemaByName(context.Background(), schema.Name)
	s.Nil(err)
	s.Equal(schema.Name, result.Name)
	s.mockStore.AssertExpectations(s.T())
	s.assertSchemaCachedByIDAndName(schema)
}

func (s *CacheBackedStoreTestSuite) TestGetUserSchemaByName_StoreError() {
	storeErr := errors.New("db error")
	s.mockStore.On("GetUserSchemaByName", mock.Anything, "bad-name").Return(UserSchema{}, storeErr).Once()

	_, err := s.cachedStore.GetUserSchemaByName(context.Background(), "bad-name")
	s.Equal(storeErr, err)

	_, ok := s.schemaByNameCache.Get(cache.CacheKey{Key: "bad-name"})
	s.False(ok)
}

// CreateUserSchema tests

func (s *CacheBackedStoreTestSuite) TestCreateUserSchema_Success() {
	schema := s.createTestSchema()
	s.mockStore.On("CreateUserSchema", mock.Anything, schema).Return(nil).Once()

	err := s.cachedStore.CreateUserSchema(context.Background(), schema)
	s.Nil(err)
	s.mockStore.AssertExpectations(s.T())
	s.assertSchemaCachedByIDAndName(schema)
}

func (s *CacheBackedStoreTestSuite) TestCreateUserSchema_StoreError() {
	schema := s.createTestSchema()
	storeErr := errors.New("store error")
	s.mockStore.On("CreateUserSchema", mock.Anything, schema).Return(storeErr).Once()

	err := s.cachedStore.CreateUserSchema(context.Background(), schema)
	s.Equal(storeErr, err)

	// Verify nothing cached on error.
	_, ok := s.schemaByIDCache.Get(cache.CacheKey{Key: schema.ID})
	s.False(ok)
}

// UpdateUserSchemaByID tests

func (s *CacheBackedStoreTestSuite) TestUpdateUserSchemaByID_Success() {
	oldSchema := s.createTestSchema()
	s.schemaByIDData[oldSchema.ID] = &oldSchema
	s.schemaByNameData[oldSchema.Name] = &oldSchema

	updatedSchema := oldSchema
	updatedSchema.Name = "UpdatedSchema"
	updatedSchema.SystemAttributes = &SystemAttributes{Display: "firstName"}

	s.mockStore.On("UpdateUserSchemaByID", mock.Anything, oldSchema.ID, updatedSchema).Return(nil).Once()

	err := s.cachedStore.UpdateUserSchemaByID(context.Background(), oldSchema.ID, updatedSchema)
	s.Nil(err)
	s.mockStore.AssertExpectations(s.T())

	// Old name key should be invalidated.
	_, ok := s.schemaByNameCache.Get(cache.CacheKey{Key: "TestSchema"})
	s.False(ok)

	// New name key should be cached.
	cachedByNewName, ok := s.schemaByNameCache.Get(cache.CacheKey{Key: "UpdatedSchema"})
	s.True(ok)
	s.Equal("UpdatedSchema", cachedByNewName.Name)

	// ID cache should now point at the updated schema.
	cachedByID, ok := s.schemaByIDCache.Get(cache.CacheKey{Key: oldSchema.ID})
	s.True(ok)
	s.Equal("UpdatedSchema", cachedByID.Name)
	s.Equal("firstName", cachedByID.SystemAttributes.Display)
}

func (s *CacheBackedStoreTestSuite) TestUpdateUserSchemaByID_StoreError() {
	oldSchema := s.createTestSchema()
	s.schemaByIDData[oldSchema.ID] = &oldSchema
	s.schemaByNameData[oldSchema.Name] = &oldSchema

	updatedSchema := oldSchema
	updatedSchema.Name = "UpdatedSchema"

	storeErr := errors.New("update error")
	s.mockStore.On("UpdateUserSchemaByID", mock.Anything, oldSchema.ID, updatedSchema).Return(storeErr).Once()

	err := s.cachedStore.UpdateUserSchemaByID(context.Background(), oldSchema.ID, updatedSchema)
	s.Equal(storeErr, err)

	// Original cache entries should still exist (not invalidated on error).
	cachedByID, ok := s.schemaByIDCache.Get(cache.CacheKey{Key: oldSchema.ID})
	s.True(ok)
	s.Equal("TestSchema", cachedByID.Name)

	cachedByName, ok := s.schemaByNameCache.Get(cache.CacheKey{Key: oldSchema.Name})
	s.True(ok)
	s.Equal("TestSchema", cachedByName.Name)
}

// DeleteUserSchemaByID tests

func (s *CacheBackedStoreTestSuite) TestDeleteUserSchemaByID_ExistsInCache() {
	schema := s.createTestSchema()
	s.schemaByIDData[schema.ID] = &schema
	s.schemaByNameData[schema.Name] = &schema

	s.mockStore.On("DeleteUserSchemaByID", mock.Anything, schema.ID).Return(nil).Once()

	err := s.cachedStore.DeleteUserSchemaByID(context.Background(), schema.ID)
	s.Nil(err)
	s.mockStore.AssertExpectations(s.T())

	// Both caches should be invalidated.
	_, ok := s.schemaByIDCache.Get(cache.CacheKey{Key: schema.ID})
	s.False(ok)

	_, ok = s.schemaByNameCache.Get(cache.CacheKey{Key: schema.Name})
	s.False(ok)
}

func (s *CacheBackedStoreTestSuite) TestDeleteUserSchemaByID_NotInCache() {
	schema := s.createTestSchema()
	s.mockStore.On("GetUserSchemaByID", mock.Anything, schema.ID).Return(schema, nil).Once()
	s.mockStore.On("DeleteUserSchemaByID", mock.Anything, schema.ID).Return(nil).Once()

	err := s.cachedStore.DeleteUserSchemaByID(context.Background(), schema.ID)
	s.Nil(err)
	s.mockStore.AssertExpectations(s.T())

	// Both caches should be invalidated (even though fetched from store).
	_, ok := s.schemaByIDCache.Get(cache.CacheKey{Key: schema.ID})
	s.False(ok)

	_, ok = s.schemaByNameCache.Get(cache.CacheKey{Key: schema.Name})
	s.False(ok)
}

func (s *CacheBackedStoreTestSuite) TestDeleteUserSchemaByID_NotFound() {
	s.mockStore.On("GetUserSchemaByID", mock.Anything, "nonexistent").
		Return(UserSchema{}, ErrUserSchemaNotFound).Once()

	err := s.cachedStore.DeleteUserSchemaByID(context.Background(), "nonexistent")
	s.Nil(err)
	s.mockStore.AssertNotCalled(s.T(), "DeleteUserSchemaByID")
}

// Pass-through method tests

func (s *CacheBackedStoreTestSuite) TestGetUserSchemaListCount_Delegated() {
	s.mockStore.On("GetUserSchemaListCount", mock.Anything).Return(5, nil).Once()

	count, err := s.cachedStore.GetUserSchemaListCount(context.Background())
	s.Nil(err)
	s.Equal(5, count)
	s.mockStore.AssertExpectations(s.T())
}

func (s *CacheBackedStoreTestSuite) TestGetUserSchemaList_Delegated() {
	expected := []UserSchemaListItem{{ID: "s1", Name: "Schema1"}}
	s.mockStore.On("GetUserSchemaList", mock.Anything, 10, 0).Return(expected, nil).Once()

	result, err := s.cachedStore.GetUserSchemaList(context.Background(), 10, 0)
	s.Nil(err)
	s.Equal(expected, result)
	s.mockStore.AssertExpectations(s.T())
}

func (s *CacheBackedStoreTestSuite) TestIsUserSchemaDeclarative_Delegated() {
	s.mockStore.On("IsUserSchemaDeclarative", "schema-1").Return(false).Once()

	result := s.cachedStore.IsUserSchemaDeclarative("schema-1")
	s.False(result)
	s.mockStore.AssertExpectations(s.T())
}

func (s *CacheBackedStoreTestSuite) TestGetUserSchemaListByOUIDs_Delegated() {
	expected := []UserSchemaListItem{{ID: "s1", Name: "Schema1"}}
	s.mockStore.On("GetUserSchemaListByOUIDs", mock.Anything,
		[]string{"ou-1"}, 10, 0).Return(expected, nil).Once()

	result, err := s.cachedStore.GetUserSchemaListByOUIDs(
		context.Background(), []string{"ou-1"}, 10, 0)
	s.Nil(err)
	s.Equal(expected, result)
	s.mockStore.AssertExpectations(s.T())
}

func (s *CacheBackedStoreTestSuite) TestGetUserSchemaListCountByOUIDs_Delegated() {
	s.mockStore.On("GetUserSchemaListCountByOUIDs", mock.Anything,
		[]string{"ou-1", "ou-2"}).Return(3, nil).Once()

	count, err := s.cachedStore.GetUserSchemaListCountByOUIDs(
		context.Background(), []string{"ou-1", "ou-2"})
	s.Nil(err)
	s.Equal(3, count)
	s.mockStore.AssertExpectations(s.T())
}

func (s *CacheBackedStoreTestSuite) TestGetDisplayAttributesByNames_Delegated() {
	expected := map[string]string{"Schema1": "email", "Schema2": "firstName"}
	s.mockStore.On("GetDisplayAttributesByNames", mock.Anything,
		[]string{"Schema1", "Schema2"}).Return(expected, nil).Once()

	result, err := s.cachedStore.GetDisplayAttributesByNames(
		context.Background(), []string{"Schema1", "Schema2"})
	s.Nil(err)
	s.Equal(expected, result)
	s.mockStore.AssertExpectations(s.T())
}
