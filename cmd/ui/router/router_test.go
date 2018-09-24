// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package router_test

import (
	"time"

	"github.com/rvflash/safe"
	"github.com/rvflash/safe/app"
	"github.com/rvflash/safe/crypto"
)

type service struct {
	logged bool
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
	if !s.logged {
		return nil, nil
	}
	return []*safe.Vault{
		safe.NewVault(hash, "vault #1", safe.NewTag("Job"), safe.NewLogin("name", "pass")),
		safe.NewVault(hash, "vault #2", safe.NewTag("Job"), safe.NewLogin("username", "pwd")),
	}, nil
}

// Vault implements the safe.Service interface.
func (s *service) Vault(hash crypto.Hash, key string) (*safe.Vault, error) {
	if !s.logged {
		return nil, nil
	}
	return safe.NewVault(hash, "vault #1", safe.NewTag("Job"), safe.NewLogin("name", "pass")), nil
}

// UpdateVault implements the safe.Service interface.
func (s *service) UpdateVault(v *safe.Vault) error {
	return nil
}

// CreateOwner implements the safe.Service interface.
func (s *service) CreateOwner(p *safe.Passphrase) error {
	return nil
}

// HasOwner implements the safe.Service interface.
func (s *service) HasOwner() bool {
	return s.logged
}

// IsOwner implements the safe.Service interface.
func (s *service) IsOwner(p *safe.Passphrase) bool {
	s.logged = true
	return s.logged
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
	if !s.logged {
		return nil, nil
	}
	return []*safe.Tag{safe.NewTag("Job"), safe.NewTag("Private")}, nil
}

// Close implements the safe.Service interface.
func (s *service) Close() error {
	return nil
}

var (
	connected    = newApp(true)
	notConnected = newApp(false)
)

func newApp(logged bool) *app.Safe {
	db := &service{logged: logged}
	a := app.New(db, "salt", "./..", app.NewSession(time.Minute, time.Second))
	if db.logged {
		_ = a.Login("pass")
	}
	return a
}
