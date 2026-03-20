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

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildPaginationLinks_MiddlePage(t *testing.T) {
	links := BuildPaginationLinks("/items", 5, 5, 20, "")
	require.Len(t, links, 4)
	assert.Equal(t, "first", links[0].Rel)
	assert.Equal(t, "/items?offset=0&limit=5", links[0].Href)
	assert.Equal(t, "prev", links[1].Rel)
	assert.Equal(t, "/items?offset=0&limit=5", links[1].Href)
	assert.Equal(t, "next", links[2].Rel)
	assert.Equal(t, "/items?offset=10&limit=5", links[2].Href)
	assert.Equal(t, "last", links[3].Rel)
	assert.Equal(t, "/items?offset=15&limit=5", links[3].Href)
}

func TestBuildPaginationLinks_FirstPage(t *testing.T) {
	links := BuildPaginationLinks("/items", 10, 0, 25, "")
	require.Len(t, links, 2)
	assert.Equal(t, "next", links[0].Rel)
	assert.Equal(t, "last", links[1].Rel)
}

func TestBuildPaginationLinks_LastPage(t *testing.T) {
	links := BuildPaginationLinks("/items", 10, 20, 25, "")
	require.Len(t, links, 2)
	assert.Equal(t, "first", links[0].Rel)
	assert.Equal(t, "prev", links[1].Rel)
}

func TestBuildPaginationLinks_SinglePage(t *testing.T) {
	links := BuildPaginationLinks("/items", 10, 0, 5, "")
	require.Len(t, links, 0)
}

func TestBuildPaginationLinks_ZeroLimit(t *testing.T) {
	links := BuildPaginationLinks("/items", 0, 0, 10, "")
	require.Len(t, links, 0)
}

func TestBuildPaginationLinks_NegativeLimit(t *testing.T) {
	links := BuildPaginationLinks("/items", -1, 0, 10, "")
	require.Len(t, links, 0)
}

func TestBuildPaginationLinks_WithExtraQuery(t *testing.T) {
	links := BuildPaginationLinks("/items", 5, 5, 20, "&include=display")
	require.Len(t, links, 4)
	assert.Equal(t, "/items?offset=0&limit=5&include=display", links[0].Href)
	assert.Equal(t, "/items?offset=0&limit=5&include=display", links[1].Href)
	assert.Equal(t, "/items?offset=10&limit=5&include=display", links[2].Href)
	assert.Equal(t, "/items?offset=15&limit=5&include=display", links[3].Href)
}
