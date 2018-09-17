// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package bolt_test

import (
	"io/ioutil"
	"testing"

	"github.com/rvflash/safe"
	"github.com/rvflash/safe/bolt"
)

func newDB() (safe.Service, error) {
	path, err := ioutil.TempDir("", "safe")
	if err != nil {
		return nil, err
	}
	return bolt.Open(path, "test")

}

func TestSafe_CreateOwner(t *testing.T) {
	// Opens temporary DB
	db, err := newDB()
	if err != nil {
		t.Fatalf("workspace: %s", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			t.Fatalf("db closing: %s", err)
		}
	}()
	// Sort order matters in this tests suite.
	var dt = []struct {
		in                *safe.Passphrase
		hasOwner, isOwner bool
		err               error
	}{
		{in: safe.NewPassPhrase("tooShort"), err: safe.ErrInvalid},
		{in: safe.NewPassPhrase("minimumLength:16"), hasOwner: true, isOwner: true},
		{in: safe.NewPassPhrase("minimumLength-16"), err: safe.ErrInvalid, hasOwner: true, isOwner: false},
	}
	var ok bool
	for i, tt := range dt {
		err = db.CreateOwner(tt.in)
		if err != tt.err {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		}
		if ok = db.HasOwner(); ok != tt.hasOwner {
			t.Errorf("%d. unexpected result: got=%t exp=%t", i, ok, tt.hasOwner)
		}
		if ok = db.IsOwner(tt.in); ok != tt.isOwner {
			t.Errorf("%d. unexpected owner: got=%t exp=%t", i, ok, tt.isOwner)
		}
	}
}
