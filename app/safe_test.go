// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package app_test

import (
	"testing"
	"time"

	"github.com/rvflash/safe"
	"github.com/rvflash/safe/app"
	"github.com/rvflash/safe/crypto"
	"github.com/sethvargo/go-password/password"
)

func TestDir_Join(t *testing.T) {
	var dt = []struct {
		in, join, out string
	}{
		{out: "."},
		{in: ".", out: "."},
		{in: "../", out: ".."},
		{in: "//", out: "/"},
		{in: "/rv/", out: "/rv"},
	}
	var d app.Dir
	for i, tt := range dt {
		d = app.Dir(tt.in)
		if out := d.Join(tt.join); out != tt.out {
			t.Errorf("%d. mismatch result: got=%q exp=%q", i, out, tt.out)
		}
	}
}

func TestNewSession(t *testing.T) {
	lifetime := time.Minute
	tick := time.Second
	s := app.NewSession(lifetime, tick)
	if s.Lifetime() != lifetime {
		t.Fatal("unexpected result")
	}
	if s.Tick() != tick {
		t.Fatal("unexpected result")
	}
}

func TestGeneratePassword(t *testing.T) {
	var dt = []struct {
		length,
		numDigits,
		numSymbols int
		noUpper,
		allowRepeat bool
		err error
	}{
		{},
		{numDigits: 1, err: password.ErrExceedsTotalLength},
	}
	for i, tt := range dt {
		out, err := app.GeneratePassword(tt.length, tt.numDigits, tt.numSymbols, tt.noUpper, tt.allowRepeat)
		if err != tt.err {
			t.Fatalf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		}
		if len(out) != tt.length {
			t.Fatalf("%d. unexpected result: got=%q exp=%q", i, err, tt.err)
		}
	}
}

func TestNew(t *testing.T) {
	s := app.New(nil, "", "", app.NewSession(time.Minute, time.Second))
	if root := s.Root(); root.String() != "." {
		t.Errorf("unexpected root: got=%q", root)
	}
	if err := s.Close(); err != nil {
		t.Errorf("unexpected close error: got=%q", err)
	}
}

type service struct {
	owned bool
	owner bool
}

// CreateVault implements the safe.Service interface.
func (s *service) CreateVault(v *safe.Vault) error {
	return nil
}

// DeleteVault implements the safe.Service interface.
func (s *service) DeleteVault(key string) error {
	return nil
}

// Vaults implements the safe.Service interface.
func (s *service) Vaults(hash crypto.Hash, tag *safe.Tag, prefix string) ([]*safe.Vault, error) {
	return nil, nil
}

// Vault implements the safe.Service interface.
func (s *service) Vault(hash crypto.Hash, key string) (*safe.Vault, error) {
	return safe.NewVault(hash, "vault", safe.NewTag("tag"), safe.NewLogin("name", "pass")), nil
}

// UpdateVault implements the safe.Service interface.
func (s *service) UpdateVault(v *safe.Vault) error {
	return nil
}

// CreateOwner implements the safe.Service interface.
func (s *service) CreateOwner(p *safe.Passphrase) error {
	s.owned, s.owner = true, true
	return nil
}

// HasOwner implements the safe.Service interface.
func (s *service) HasOwner() bool {
	return s.owned
}

// IsOwner implements the safe.Service interface.
func (s *service) IsOwner(p *safe.Passphrase) bool {
	return s.owned && s.owner
}

// CreateTag implements the safe.Service interface.
func (s *service) CreateTag(t *safe.Tag) error {
	return nil
}

// DeleteTag implements the safe.Service interface.
func (s *service) DeleteTag(key string) error {
	return nil
}

// Tags implements the safe.Service interface.
func (s *service) Tags() ([]*safe.Tag, error) {
	return nil, nil
}

// Close implements the safe.Service interface.
func (s *service) Close() error {
	return nil
}
