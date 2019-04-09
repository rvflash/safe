// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package app

import (
	"net/url"

	"github.com/rvflash/safe"
)

// Login contains all login's data as plaintext.
type Login struct {
	Username string `json:"user"`
	Password string `json:"pass"`
	URL      string `json:"url,omitempty"`
	Note     string `json:"note,omitempty"`
}

// CreateVault creates a new vault.
func (s *Safe) CreateVault(name, tag string, data Login) (*safe.Vault, error) {
	h, err := s.hash()
	if err != nil {
		return nil, err
	}
	// Creates a new login
	l := safe.NewLogin(data.Username, data.Password)
	if data.URL != "" {
		l.URL, err = url.Parse(data.URL)
		if err != nil {
			return nil, err
		}
	}
	l.Note = data.Note

	// Adds it to the vault.
	v := safe.NewVault(h, name, safe.NewTag(tag), l)
	return v, s.db.CreateVault(v)
}

// DeleteVault deletes a vault by its name and tag.
func (s *Safe) DeleteVault(name, tag string) error {
	h, err := s.hash()
	if err != nil {
		return err
	}
	v := safe.NewVault(h, name, safe.NewTag(tag), nil)
	return s.db.DeleteVault(string(v.Key()))
}

// ListVaultByNames lists all tags prefixed by the given value and  using this tag.
func (s *Safe) ListVaultByNames(tag, prefix string) ([]*safe.Vault, error) {
	h, err := s.hash()
	if err != nil {
		return nil, err
	}
	d, err := s.db.Vaults(h, safe.NewTag(tag), prefix)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Vault returns the vault by its name and tag.
func (s *Safe) Vault(name, tag string) (*safe.Vault, error) {
	h, err := s.hash()
	if err != nil {
		return nil, err
	}
	v := safe.NewVault(h, name, safe.NewTag(tag), nil)
	return s.db.Vault(h, string(v.Key()))
}

// UpdateVault updates the vault.
func (s *Safe) UpdateVault(name, tag string, data Login) (*safe.Vault, error) {
	// Retrieves it
	v, err := s.Vault(name, tag)
	if err != nil {
		return nil, err
	}
	// Do some updates.
	l := v.Login()
	l.Name = data.Username
	l.Note = data.Note
	l.Password = data.Password
	if data.URL != "" {
		l.URL, err = url.Parse(data.URL)
		if err != nil {
			return nil, err
		}
	} else {
		l.URL = nil
	}
	// Ignores the error previously checked and signs data.
	h, _ := s.hash()
	if err = v.SignLogin(h, l); err != nil {
		return nil, err
	}
	return v, s.db.UpdateVault(v)
}
