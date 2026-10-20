// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package gateway

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// autoApplyFixture is a version fixture whose only gateway is the default one, with auto apply set.
func autoApplyFixture(t *testing.T, enabled, isDefault bool) *versionFixture {
	t.Helper()
	f := newVersionFixture(t)
	svc := f.svc.(*versionService)
	svc.autoApply = enabled
	svc.gateways.(*fakeStore).gateways[0].IsDefault = isDefault
	f.capture(t, "ou-a")
	return f
}

func TestAutoApplyAppliesTheLatestVersionToTheDefaultGateway(t *testing.T) {
	f := autoApplyFixture(t, true, true)

	f.svc.(*versionService).applyToDefaultGateway(context.Background())

	assert.Len(t, f.client.sent, 1)
}

func TestAutoApplyAppliesAgainWhenTheConfigurationIsTheSame(t *testing.T) {
	f := autoApplyFixture(t, true, true)
	svc := f.svc.(*versionService)

	svc.applyToDefaultGateway(context.Background())
	svc.applyToDefaultGateway(context.Background())

	assert.Len(t, f.client.sent, 2, "a changed secret leaves the version as it was, so it is sent again")
}

func TestAutoApplyDoesNothingWhenItIsOff(t *testing.T) {
	f := autoApplyFixture(t, false, true)

	f.svc.(*versionService).applyToDefaultGateway(context.Background())

	assert.Empty(t, f.client.sent)
}

func TestAutoApplyDoesNothingWithoutADefaultGateway(t *testing.T) {
	f := autoApplyFixture(t, true, false)

	f.svc.(*versionService).applyToDefaultGateway(context.Background())

	assert.Empty(t, f.client.sent)
}
