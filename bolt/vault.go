// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package bolt

import (
	"bytes"
	"encoding/json"

	"github.com/rvflash/safe"
	"github.com/rvflash/safe/crypto"

	bolt "go.etcd.io/bbolt"
)

// CreateVault implements the VaultService interface.
func (s *Safe) CreateVault(e *safe.Vault) error {
	if e == nil || e.Tag() == nil || !s.exists([]byte(TagTable), e.Tag().Key()) {
		// Tag not exists
		return safe.ErrInvalid
	}
	return s.create(VaultTable, e, false)
}

func (s *Safe) exists(in, key []byte) bool {
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(in).Get(key)
		if len(b) == 0 {
			return safe.ErrNotFound
		}
		return nil
	})
	return err == nil
}

// DeleteVault implements the VaultService interface.
func (s *Safe) DeleteVault(key string) error {
	return s.delete(VaultTable, key)
}

// Vaults implements the VaultService interface.
func (s *Safe) Vaults(hash crypto.Hash, tag *safe.Tag, prefix string) ([]*safe.Vault, error) {
	res := make([]*safe.Vault, 0)
	err := s.db.View(func(tx *bolt.Tx) (err error) {
		c := tx.Bucket([]byte(VaultTable)).Cursor()
		p := bytes.Join([][]byte{tag.Key(), []byte(prefix)}, []byte(""))
		for k, v := c.Seek(p); k != nil && bytes.HasPrefix(k, p); k, v = c.Next() {
			d := safe.EmptyVault(hash)
			if err = json.Unmarshal(v, &d); err != nil {
				return
			}
			res = append(res, d)
		}
		return
	})
	return res, err
}

// Vault implements the VaultService interface.
func (s *Safe) Vault(hash crypto.Hash, key string) (*safe.Vault, error) {
	res := safe.EmptyVault(hash)
	err := s.db.View(func(tx *bolt.Tx) error {
		buf := tx.Bucket([]byte(VaultTable)).Get([]byte(key))
		if len(buf) == 0 {
			return safe.ErrNotFound
		}
		return json.Unmarshal(buf, &res)
	})
	return res, err
}

// UpdateVault implements the VaultService interface.
func (s *Safe) UpdateVault(e *safe.Vault) error {
	return s.create(VaultTable, e, true)
}
