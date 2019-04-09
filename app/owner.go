// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package app

import "github.com/rvflash/safe"

// SignUp creates the owner if not exists and signs in.
func (s *Safe) SignUp(pass string) (err error) {
	p := safe.NewPassPhrase(pass)
	if err = s.db.CreateOwner(p); err == nil {
		s.login(p)
	}
	return
}

// SignIn signs in the owner if the password is valid.
func (s *Safe) SignIn(pass string) (ok bool) {
	p := safe.NewPassPhrase(pass)
	if ok = s.db.IsOwner(p); ok {
		s.login(p)
	}
	return
}

// Login deals with new or not owner to login it .
func (s *Safe) Login(pass string) error {
	if !s.db.HasOwner() {
		return s.SignUp(pass)
	}
	if !s.SignIn(pass) {
		return ErrNotLogged
	}
	return nil
}

// LogOut disconnects the owner.
func (s *Safe) LogOut() {
	s.mu.Lock()
	s.log = nil
	s.mu.Unlock()
}

// Logged returns an error:
// > if the owner is unknown
// > if the owner is not logged
func (s *Safe) Logged() error {
	if !s.db.HasOwner() {
		return ErrNotFound
	}
	return s.logged()
}
