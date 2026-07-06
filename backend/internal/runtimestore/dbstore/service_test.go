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
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

const (
	testDeploymentID = "test-deployment"
	testNamespace    = providers.RuntimeStoreNamespace("test-ns")
	testKey          = "key1"
)

type DBStoreTestSuite struct {
	suite.Suite
	store *dbStore
	ctx   context.Context
}

func TestDBStoreTestSuite(t *testing.T) {
	suite.Run(t, new(DBStoreTestSuite))
}

func (s *DBStoreTestSuite) SetupTest() {
	s.store = &dbStore{deploymentID: testDeploymentID}
	s.ctx = context.Background()
}

func (s *DBStoreTestSuite) TestNewDBStore() {
	store := newDBStore(testDeploymentID)
	s.Equal(&dbStore{deploymentID: testDeploymentID}, store)
}

func (s *DBStoreTestSuite) TestPut_ReturnsNotImplemented() {
	err := s.store.Put(s.ctx, testNamespace, testKey, []byte("value"), 60)
	s.True(errors.Is(err, ErrNotImplemented))
}

func (s *DBStoreTestSuite) TestGet_ReturnsNotImplemented() {
	got, err := s.store.Get(s.ctx, testNamespace, testKey)
	s.Nil(got)
	s.True(errors.Is(err, ErrNotImplemented))
}

func (s *DBStoreTestSuite) TestUpdate_ReturnsNotImplemented() {
	err := s.store.Update(s.ctx, testNamespace, testKey, []byte("value"))
	s.True(errors.Is(err, ErrNotImplemented))
}

func (s *DBStoreTestSuite) TestDelete_ReturnsNotImplemented() {
	err := s.store.Delete(s.ctx, testNamespace, testKey)
	s.True(errors.Is(err, ErrNotImplemented))
}

func (s *DBStoreTestSuite) TestTake_ReturnsNotImplemented() {
	got, err := s.store.Take(s.ctx, testNamespace, testKey)
	s.Nil(got)
	s.True(errors.Is(err, ErrNotImplemented))
}
