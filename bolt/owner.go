// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package bolt

import (
	"bytes"

	"github.com/coreos/bbolt"
	"github.com/rvflash/safe"
)

// todo adds UpdateOwner with modification of all vaults.

// CreateOwner implements the safe.OwnerService.
func (s *Safe) CreateOwner(p *safe.Passphrase) error {
	return s.create(OwnerTable, p, false)
}

// HasOwner implements the safe.OwnerService.
func (s *Safe) HasOwner() bool {
	return s.privateKey(nil) != safe.ErrNotFound
}

// IsOwner implements the safe.OwnerService.
func (s *Safe) IsOwner(p *safe.Passphrase) bool {
	return s.privateKey(p) == nil
}

func (s *Safe) privateKey(p *safe.Passphrase) error {
	return s.db.View(func(tx *bolt.Tx) error {
		buf := bytes.Trim(tx.Bucket([]byte(OwnerTable)).Get(p.Key()), `"`)
		if len(buf) == 0 {
			return safe.ErrNotFound
		}
		if p == nil || p.Compare(buf) != nil {
			return safe.ErrInvalid
		}
		return nil
	})
}
