// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package app

import "github.com/rvflash/safe"

// ListTagByNames ...
func (s *Safe) ListTagByNames() ([]*safe.Tag, error) {
	if err := s.logged(); err != nil {
		return nil, err
	}
	return s.db.Tags()
}

// CreateTag ...
func (s *Safe) CreateTag(name string) (*safe.Tag, error) {
	if err := s.logged(); err != nil {
		return nil, err
	}
	tag := safe.NewTag(name)
	return tag, s.db.CreateTag(tag)
}

// DeleteTag ...
func (s *Safe) DeleteTag(name string) error {
	if err := s.logged(); err != nil {
		return err
	}
	tag := safe.NewTag(name)
	return s.db.DeleteTag(string(tag.Key()))
}
