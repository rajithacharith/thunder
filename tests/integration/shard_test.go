// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func testPackages(names ...string) []string {
	packages := make([]string, len(names))
	for i, name := range names {
		packages[i] = "github.com/thunder-id/thunderid/tests/integration/" + name
	}
	return packages
}

func TestShardPackagesCoversEveryPackageOnce(t *testing.T) {
	packages := testPackages("a", "b", "c", "d", "e", "f", "g", "h")
	for total := 1; total <= 5; total++ {
		var all []string
		for index := 1; index <= total; index++ {
			all = append(all, shardPackages(packages, index, total)...)
		}
		sort.Strings(all)
		want := append([]string(nil), packages...)
		sort.Strings(want)
		require.Equal(t, want, all, "total=%d", total)
	}
}

func TestShardPackagesBalancesByWeight(t *testing.T) {
	packageWeights = map[string]int{"heavy": 100, "mid1": 50, "mid2": 50}
	t.Cleanup(func() { packageWeights = map[string]int{} })

	packages := testPackages("heavy", "mid1", "mid2")
	require.Equal(t, testPackages("heavy"), shardPackages(packages, 1, 2))
	require.Equal(t, testPackages("mid1", "mid2"), shardPackages(packages, 2, 2))
}

func TestShardPackagesIsDeterministic(t *testing.T) {
	packages := testPackages("a", "b", "c", "d")
	reversed := testPackages("d", "c", "b", "a")
	for index := 1; index <= 3; index++ {
		require.Equal(t, shardPackages(packages, index, 3), shardPackages(reversed, index, 3))
	}
}

func TestParseShard(t *testing.T) {
	index, total, err := parseShard("2/3")
	require.NoError(t, err)
	require.Equal(t, 2, index)
	require.Equal(t, 3, total)

	for _, spec := range []string{"", "3", "0/3", "4/3", "a/b", "1/0", "1/2/3"} {
		_, _, err := parseShard(spec)
		require.Error(t, err, fmt.Sprintf("spec %q", spec))
	}
}
