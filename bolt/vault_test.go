// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package bolt_test

import (
	"fmt"
	"testing"

	"github.com/rvflash/safe"
	"github.com/rvflash/safe/crypto"
)

func newHash() (crypto.Hash, error) {
	b, err := newPassPhrase().NewCipher("salt")
	if err != nil {
		return nil, err
	}
	return crypto.New(b), nil
}

func newDBWithTags() (safe.Service, error) {
	db, err := newDBWithOwner()
	if err != nil {
		return nil, err
	}
	for _, s := range []string{"Jobs", "Social"} {
		if err = db.CreateTag(safe.NewTag(s)); err != nil {
			return nil, err
		}
	}
	return db, nil
}

func newJobTag() *safe.Tag {
	return safe.NewTag("Jobs")
}
func newUnknownTag() *safe.Tag {
	return safe.NewTag("Oops")
}

func TestSafe_Vaults(t *testing.T) {
	// Opens temporary DB
	db, err := newDBWithTags()
	if err != nil {
		t.Fatalf("workspace: %s", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			t.Fatalf("db closing: %s", err)
		}
	}()
	// Hash
	h, err := newHash()
	if err != nil {
		t.Fatalf("workspace: %s", err)
	}
	// Sort order matters in this tests suite.
	var dt = []struct {
		key, prefix string
		list        []*safe.Vault
		tag         *safe.Tag
		v1, v2      *safe.Vault
		err,
		errDel,
		errList,
		errNew,
		errUpd error
	}{
		{
			v1:     safe.NewVault(h, "Vault 1", newJobTag(), safe.NewLogin("user", "pass1")),
			v2:     safe.NewVault(h, "Vault 1", newJobTag(), safe.NewLogin("user", "pass2")),
			tag:    newJobTag(),
			errDel: safe.ErrNotFound,
			list: []*safe.Vault{
				safe.NewVault(h, "Vault 1", newJobTag(), safe.NewLogin("user", "pass2")),
			},
		},
		{
			v1:     safe.NewVault(h, "Vault 2", newUnknownTag(), safe.NewLogin("user", "pass1")),
			v2:     safe.NewVault(h, "Vault 2", newUnknownTag(), safe.NewLogin("user", "pass2")),
			tag:    newJobTag(),
			errDel: safe.ErrNotFound,
			list: []*safe.Vault{
				safe.NewVault(h, "Vault 1", newJobTag(), safe.NewLogin("user", "pass2")),
			},
			errNew: safe.ErrInvalid,
		},
	}
	var (
		ok   bool
		v    *safe.Vault
		list []*safe.Vault
	)
	for i, tt := range dt {
		if err = db.CreateVault(tt.v1); err != tt.errNew {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.errNew)
		}
		if err = db.DeleteVault(tt.key); err != tt.errDel {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.errDel)
		}
		if list, err = db.Vaults(h, tt.tag, tt.prefix); err != tt.errList {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.errList)
		}
		if ok, err = equalVaults(list, tt.list); !ok {
			t.Errorf("%d. mismatch vaults: %q", i, err)
		}
		if err = db.UpdateVault(tt.v2); err != tt.errUpd {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.errNew)
		}
		if v, err = db.Vault(h, string(tt.v2.Key())); err != tt.err {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		}
		if ok, err = equalVaults([]*safe.Vault{v}, []*safe.Vault{tt.v2}); !ok {
			t.Errorf("%d. mismatch vaults: got=%v exp=%v", i, v, tt.v2)
		}
	}
}

func equalVaults(l1, l2 []*safe.Vault) (bool, error) {
	if a, b := len(l1), len(l2); a != b {
		return false, fmt.Errorf("len: got=%d exp=%d", a, b)
	}
	if l1 == nil {
		return true, nil
	}
	for i, t := range l1 {
		if a, b := t.Name(), l2[i].Name(); a != b {
			return false, fmt.Errorf("name: got=%q exp=%q", a, b)
		}
	}
	return true, nil
}
