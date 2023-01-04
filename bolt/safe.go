// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package bolt

import (
	"encoding/json"
	"hash/adler32"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rvflash/safe"

	bolt "go.etcd.io/bbolt"
)

// Lists of data tables.
const (
	// OwnerTable is the name of the table where store the information's owner.
	OwnerTable = "owner"
	// VaultTable is the name of the table to store the vaults.
	VaultTable = "vault"
	// TagTable is the table's name to store tags.
	TagTable = "tag"
)

// Safe is the BoldDB database.
type Safe struct {
	db *bolt.DB
}

// Open opens the database in the given path.
// An error occurred if it failed to do it.
func Open(path, name string) (safe.Service, error) {
	name, err := sum32(name)
	if err != nil {
		return nil, err
	}
	path = filepath.Join(path, name+".db")
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	// Initializes the database by creating the default buckets (if not exists).
	err = db.Update(func(tx *bolt.Tx) (err error) {
		for _, s := range []string{VaultTable, OwnerTable, TagTable} {
			if _, err = tx.CreateBucketIfNotExists([]byte(s)); err != nil {
				return
			}
		}
		return
	})
	return &Safe{db: db}, err
}

// Close closes the connection to database.
// It implements the io.Closer interface.
func (s *Safe) Close() error {
	return s.db.Close()
}

func (s *Safe) create(in string, data safe.Data, orReplace bool) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		if data == nil || !data.Valid() {
			return safe.ErrInvalid
		}
		k := data.Key()
		b := tx.Bucket([]byte(in))
		if len(b.Get(k)) > 0 && !orReplace {
			// Already exists
			return safe.ErrInvalid
		}
		buf, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return b.Put(k, buf)
	})
}

func (s *Safe) delete(in, key string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		if in == "" {
			return safe.ErrMissing
		}
		if key == "" {
			return safe.ErrNotFound
		}
		return tx.Bucket([]byte(in)).Delete([]byte(key))
	})
}

func sum32(s string) (string, error) {
	h := adler32.New()
	if _, err := h.Write([]byte(s)); err != nil {
		return "", err
	}
	return strconv.Itoa(int(h.Sum32())), nil
}
