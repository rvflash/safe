// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package bolt

import (
	"encoding/json"

	"github.com/coreos/bbolt"
	"github.com/rvflash/safe"
)

// CreateTag implements the TagService.
func (s *Safe) CreateTag(t *safe.Tag) error {
	return s.create(TagTable, t, false)
}

// DeleteTag implements the TagService.
func (s *Safe) DeleteTag(key string) error {
	// todo Check if it is still used?
	return s.delete(TagTable, key)
}

// Tags implements the TagService.
func (s *Safe) Tags() ([]*safe.Tag, error) {
	res := make([]*safe.Tag, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(TagTable)).ForEach(func(k, v []byte) (err error) {
			var d safe.Tag
			if err = json.Unmarshal(v, &d); err != nil {
				return
			}
			res = append(res, &d)
			return
		})
	})
	return res, err
}
