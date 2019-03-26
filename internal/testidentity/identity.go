// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package testidentity

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	"storj.io/storj/pkg/identity"
	"storj.io/storj/pkg/storj"
)

type IdentityTest func(*testing.T, *identity.FullIdentity)

// NewTestIdentity is a helper function to generate new node identities with
// correct difficulty and concurrency
func NewTestIdentity(ctx context.Context) (*identity.FullIdentity, error) {
	ca, err := identity.NewCA(ctx, identity.NewCAOptions{
		//VersionNumber: storj.V1,
		Difficulty:  10,
		Concurrency: 4,
	})
	if err != nil {
		return nil, err
	}
	identity, err := ca.NewIdentity()
	if err != nil {
		return nil, err
	}
	return identity, err
}

// NewTestCA returns a ca with a default difficulty and concurrency for use in tests
func NewTestCA(ctx context.Context) (*identity.FullCertificateAuthority, error) {
	return identity.NewCA(ctx, identity.NewCAOptions{
		Difficulty:  8,
		Concurrency: 4,
	})
}

func IdentityVersionsTest(t *testing.T, test IdentityTest) {
	for versionNumber := range storj.IDVersions {
		t.Run(fmt.Sprintf("identity version %d", versionNumber), func(t *testing.T) {
			ident, err := IdentityVersions[versionNumber].NewIdentity()
			require.NoError(t, err)

			test(t, ident)
		})
	}
}

func SignedIdentityVersionsTest(t *testing.T, test IdentityTest) {
	for versionNumber := range storj.IDVersions {
		t.Run(fmt.Sprintf("identity version %d", versionNumber), func(t *testing.T) {
			fmt.Printf("t.Run version %d\n", versionNumber)
			ident, err := SignedIdentityVersions[versionNumber].NewIdentity()
			require.NoError(t, err)

			fmt.Printf("actual version %d\n", ident.ID.Version().Number)
			test(t, ident)
		})
	}
}

func CompleteIdentityVersionsTest(t *testing.T, test IdentityTest) {
	t.Run("unsigned identity", func(t *testing.T) {
		IdentityVersionsTest(t, test)
	})

	t.Run("signed identity", func(t *testing.T) {
		SignedIdentityVersionsTest(t, test)
	})
}

// NewTestManageablePeerIdentity returns a new manageable peer identity for use in tests.
func NewTestManageablePeerIdentity(ctx context.Context) (*identity.ManageablePeerIdentity, error) {
	ca, err := NewTestCA(ctx)
	if err != nil {
		return nil, err
	}

	ident, err := ca.NewIdentity()
	if err != nil {
		return nil, err
	}
	return identity.NewManageablePeerIdentity(ident.PeerIdentity(), ca), nil
}
