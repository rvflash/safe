// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package safe_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/rvflash/safe"
	"github.com/rvflash/safe/crypto"
)

func TestNewTag(t *testing.T) {
	var dt = []struct {
		key  []byte
		name string
		ok   bool
	}{
		{},
		{name: " "},
		{name: "#rv"},
		{name: "rv ", key: []byte("rv#"), ok: true},
	}
	var (
		key []byte
		tg  *safe.Tag
		ok  bool
	)
	for i, tt := range dt {
		tg = safe.NewTag(tt.name)
		if key = tg.Key(); !bytes.Equal(key, tt.key) {
			t.Fatalf("%d. mismatch key: got=%q exp=%q", i, key, tt.key)
		}
		if ok = tg.Valid(); ok != tt.ok {
			t.Fatalf("%d. mismatch result: got=%t exp=%t", i, ok, tt.ok)
		}
	}
}

func TestTag_MarshalJSON(t *testing.T) {
	t1 := safe.NewTag("rv")
	b, err := t1.MarshalJSON()
	if err != nil {
		t.Fatalf("unexpected error: %q", err)
	}
	var t2 safe.Tag
	if err = t2.UnmarshalJSON(b); err != nil {
		t.Fatalf("unexpected error: %q", err)
	}
	if t2.Name() != t1.Name() {
		t.Fatalf("mismatch content: got=%q exp=%q", t2.Name(), t1.Name())
	}
}

func TestNewVault(t *testing.T) {
	// Basic workflow
	h := hash()
	// Create the vault
	v1 := safe.NewVault(h, "vault", safe.NewTag("tag"), login())
	b, err := json.Marshal(v1)
	if err != nil {
		t.Fatalf("unexpected error: %q", err)
	}
	// Get the vault from a JSON string with The hash key.
	v2 := safe.EmptyVault(h)
	if err = json.Unmarshal(b, &v2); err != nil {
		t.Fatalf("unexpected error: %q", err)
	}
	if v2.Login() == nil || v2.Login().Name != "name" {
		t.Fatalf("mismatch content: got=%q exp=%q", v2.Login(), v1.Login())
	}
	// Get the vault from a JSON string with a wrong hash key.
	v3 := safe.EmptyVault(crypto.New([]byte("what-a-beautiful-public-fail-key")))
	if err = json.Unmarshal(b, &v3); err == nil {
		t.Fatal("expected error")
	}
}

func hash() crypto.Hash {
	return crypto.New([]byte("what-a-beautiful-public-hash-key")) // 32 bytes
}

func login() *safe.Login {
	return safe.NewLogin("name", "pass")
}

func TestVault_Valid(t *testing.T) {
	// Environments
	var (
		h  = hash()
		l  = login()
		tg = safe.NewTag("tag")
	)
	// Tests
	var (
		dt = []struct {
			login *safe.Login
			name  string
			ok    bool
			tag   *safe.Tag
		}{
			{name: "rv"},
			{name: "rv", tag: tg},
			{name: "rv", login: l},
			{name: "rv", tag: tg, login: l, ok: true},
		}
		v  *safe.Vault
		ok bool
	)
	for i, tt := range dt {
		v = safe.NewVault(h, tt.name, tt.tag, tt.login)
		if ok = v.Valid(); ok != tt.ok {
			t.Errorf("%d. mismatch result: got=%t exp=%t", i, ok, tt.ok)
		} else if ok {
			if n := strings.TrimSpace(tt.name); v.Name() != n {
				t.Errorf("%d. unexpected name: got=%q exp=%q", i, v.Name(), n)
			}
			if now := time.Now(); v.AddDate().After(now) {
				t.Errorf("%d. unexpected creation date: got=%q exp=%q", i, v.AddDate(), now)
			}
		}
		if v.Tag() != nil {
			if k := bytes.Join([][]byte{v.Tag().Key(), []byte(v.Name())}, []byte("")); !bytes.Equal(v.Key(), k) {
				t.Errorf("%d. unexpected key: got=%q exp=%q", i, v.Key(), k)
			}
		}
	}
}
