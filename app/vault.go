// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package app

import (
	"net/url"

	"github.com/rvflash/safe"
)

// Names used for data properties in form
const (
	FormUser = "user"
	FormPass = "pwd"
	FormURL  = "url"
	FormNote = "note"
)

// CreateVault ...
func (s *Safe) CreateVault(name, tag string, data url.Values) (*safe.Vault, error) {
	h, err := s.hash()
	if err != nil {
		return nil, err
	}
	// Creates a new login
	l := safe.NewLogin(data.Get(FormUser), data.Get(FormPass))
	l.Note = data.Get(FormNote)
	l.URL, err = url.Parse(data.Get(FormURL))
	if err != nil {
		return nil, err
	}
	// Adds it to the vault.
	v := safe.NewVault(h, name, safe.NewTag(tag), l)
	return v, s.db.CreateVault(v)
}

// DeleteVault ...
func (s *Safe) DeleteVault(name, tag string) error {
	h, err := s.hash()
	if err != nil {
		return err
	}
	v := safe.NewVault(h, name, safe.NewTag(tag), nil)
	return s.db.DeleteVault(string(v.Key()))
}

// ListVaultByNames ...
func (s *Safe) ListVaultByNames(tag, prefix string) ([]*safe.Vault, error) {
	h, err := s.hash()
	if err != nil {
		return nil, err
	}
	return s.db.Vaults(h, safe.NewTag(tag), prefix)
}

// Vault ...
func (s *Safe) Vault(name, tag string) (*safe.Vault, error) {
	h, err := s.hash()
	if err != nil {
		return nil, err
	}
	v := safe.NewVault(h, name, safe.NewTag(tag), nil)
	return s.db.Vault(h, string(v.Key()))
}

// UpdateVault ...
func (s *Safe) UpdateVault(name, tag string, data url.Values) (*safe.Vault, error) {
	// Retrieves the vault to change.
	v, err := s.Vault(name, tag)
	if err != nil {
		return nil, err
	}
	// Do some updates.
	l := v.Login()
	l.Name = data.Get(FormUser)
	l.Note = data.Get(FormNote)
	l.Password = data.Get(FormPass)
	l.URL, err = url.Parse(data.Get(FormURL))
	if err != nil {
		return nil, err
	}
	// Ignores the error (previously retrieved)
	h, _ := s.hash()
	if err = v.SignLogin(h, l); err != nil {
		return nil, err
	}
	return v, s.db.UpdateVault(v)
}
