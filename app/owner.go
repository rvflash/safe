// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package app

import "github.com/rvflash/safe"

// SignUp ...
func (s *Safe) SignUp(pass string) (err error) {
	p := safe.NewPassPhrase(pass)
	if err = s.db.CreateOwner(p); err == nil {
		s.login(p)
	}
	return
}

// SignIn ...
func (s *Safe) SignIn(pass string) (ok bool) {
	p := safe.NewPassPhrase(pass)
	if ok = s.db.IsOwner(p); ok {
		s.login(p)
	}
	return
}

// Login ...
func (s *Safe) Login(pass string) error {
	if !s.db.HasOwner() {
		return s.SignUp(pass)
	}
	if !s.SignIn(pass) {
		return ErrNotLogged
	}
	return nil
}

// LogOut ..
func (s *Safe) LogOut() {
	s.mu.Lock()
	s.log = nil
	s.mu.Unlock()
}

// Logged ...
func (s *Safe) Logged() error {
	if !s.db.HasOwner() {
		return ErrNotFound
	}
	return s.logged()
}
