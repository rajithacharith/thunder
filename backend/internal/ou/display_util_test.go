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

package ou

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	"github.com/asgardeo/thunder/internal/system/log"
)

// --- resolveDisplayAttributePaths tests ---

func TestResolveDisplayAttributePaths_Success(t *testing.T) {
	schemaMock := NewDisplayAttributeResolverMock(t)
	schemaMock.On("GetDisplayAttributesByNames", mock.Anything,
		mock.MatchedBy(func(names []string) bool {
			if len(names) != 2 {
				return false
			}
			has := map[string]bool{names[0]: true, names[1]: true}
			return has["employee"] && has["contractor"]
		})).Return(map[string]string{"employee": "email", "contractor": "name"},
		(*serviceerror.ServiceError)(nil))

	result := resolveDisplayAttributePaths(context.Background(),
		[]string{"employee", "contractor", "employee"}, schemaMock, nil)
	assert.Equal(t, "email", result["employee"])
	assert.Equal(t, "name", result["contractor"])
}

func TestResolveDisplayAttributePaths_NilSchemaService(t *testing.T) {
	result := resolveDisplayAttributePaths(context.Background(), []string{"employee"}, nil, nil)
	assert.Nil(t, result)
}

func TestResolveDisplayAttributePaths_EmptyTypes(t *testing.T) {
	schemaMock := NewDisplayAttributeResolverMock(t)
	result := resolveDisplayAttributePaths(context.Background(), []string{}, schemaMock, nil)
	assert.Nil(t, result)
}

func TestResolveDisplayAttributePaths_AllEmptyStrings(t *testing.T) {
	schemaMock := NewDisplayAttributeResolverMock(t)
	result := resolveDisplayAttributePaths(context.Background(), []string{"", ""}, schemaMock, nil)
	assert.Nil(t, result)
}

func TestResolveDisplayAttributePaths_SchemaServiceError(t *testing.T) {
	schemaMock := NewDisplayAttributeResolverMock(t)
	schemaMock.On("GetDisplayAttributesByNames", mock.Anything, []string{"employee"}).
		Return(nil, &serviceerror.ServiceError{Code: "500", Error: "schema unavailable"})

	logger := log.GetLogger()
	result := resolveDisplayAttributePaths(context.Background(), []string{"employee"}, schemaMock, logger)
	assert.Nil(t, result)
}
