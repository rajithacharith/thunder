// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// defaultPackageWeight is the weight, in seconds, of a package that is not in packageWeights.
const defaultPackageWeight = 8

// packageWeights holds the approximate run time in seconds of the slower packages, keyed by their path
// below tests/integration. Weights only balance the shards: a package missing here still runs, in
// exactly one shard, so a new package never goes untested. Refresh them from a CI run when the
// shards drift apart.
var packageWeights = map[string]int{
	"governance":          465,
	"federated":           143,
	"gateway":             131,
	"flow/authentication": 110,
	"oauth/token":         97,
	"flow/execution":      90,
	"oauth/sso":           85,
	"agent":               67,
	"flow/registration":   57,
	"observability":       53,
	"oauth/idjag":         46,
	"oauth/introspect":    45,
	"flow/extension":      45,
	"application":         42,
	"user":                31,
	"connection":          29,
	"oauth/authz":         29,
	"oauth/ciba":          28,
	"authorizationdetail": 26,
	"consent":             25,
	"oauth/dcr":           24,
	"oauth/dpop":          24,
	"openid4vci":          24,
	"authn":               21,
	"role":                21,
	"openid4vp":           21,
	"oauth/claims":        20,
	"resource":            18,
	"flow/mgt":            16,
	"importexport":        15,
}

// shardPackages returns the packages that shard index (1-based) of total shards runs. Packages are
// placed heaviest first on the shard with the lightest load so far, so every shard computes the same
// split from the same package list and each package lands in exactly one shard.
func shardPackages(packages []string, index, total int) []string {
	ordered := append([]string(nil), packages...)
	sort.SliceStable(ordered, func(i, j int) bool {
		wi, wj := packageWeight(ordered[i]), packageWeight(ordered[j])
		if wi != wj {
			return wi > wj
		}
		return ordered[i] < ordered[j]
	})

	loads := make([]int, total)
	var selected []string
	for _, pkg := range ordered {
		lightest := 0
		for i := range loads {
			if loads[i] < loads[lightest] {
				lightest = i
			}
		}
		loads[lightest] += packageWeight(pkg)
		if lightest == index-1 {
			selected = append(selected, pkg)
		}
	}
	sort.Strings(selected)
	return selected
}

// packageWeight looks a package up by the part of its import path below tests/integration.
func packageWeight(pkg string) int {
	_, rel, found := strings.Cut(pkg, "/tests/integration/")
	if found {
		if weight, ok := packageWeights[rel]; ok {
			return weight
		}
	}
	return defaultPackageWeight
}

// parseShard parses a shard written as "index/total", for example "2/3".
func parseShard(spec string) (index, total int, err error) {
	indexText, totalText, found := strings.Cut(spec, "/")
	if found {
		index, err = strconv.Atoi(indexText)
		if err == nil {
			total, err = strconv.Atoi(totalText)
		}
	}
	if !found || err != nil || total < 1 || index < 1 || index > total {
		return 0, 0, fmt.Errorf("invalid INTEGRATION_SHARD %q, expected index/total such as 2/3", spec)
	}
	return index, total, nil
}
