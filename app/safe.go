// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package app

import (
	"errors"
	"path/filepath"
	"sync"
	"time"

	"github.com/rvflash/safe"
	"github.com/rvflash/safe/crypto"
	"github.com/sethvargo/go-password/password"
)

// Common errors
var (
	// ErrNotRegistered ...
	ErrNotFound = errors.New("no database")
	// ErrNotLogged ...
	ErrNotLogged = errors.New("not logged in")
)

// Dir defines the current workspace.
// An empty Dir is treated as ".".
type Dir string

// Join joins the given ptha to the current directory.
func (d Dir) Join(path string) string {
	return filepath.Clean(filepath.Join(d.String(), path))
}

// String implements the fmt.Stringer interface.
func (d Dir) String() string {
	dir := string(d)
	if dir == "" {
		return "."
	}
	return filepath.Clean(dir)
}

// Session ...
type Session interface {
	Lifetime() time.Duration
	Tick() time.Duration
}

type session struct {
	lifetime,
	tick time.Duration
}

// NewSession ...
func NewSession(lifetime, tick time.Duration) Session {
	return &session{lifetime: lifetime, tick: tick}
}

// Lifetime ...
func (s *session) Lifetime() time.Duration {
	return s.lifetime
}

// Tick ...
func (s *session) Tick() time.Duration {
	return s.tick
}

// Safe ...
type Safe struct {
	db   safe.Service
	log  *auth
	mu   sync.Mutex
	root Dir
	salt string
	tick *time.Ticker
}

// New ...
func New(db safe.Service, salt, root string, session Session) *Safe {
	s := &Safe{
		db:   db,
		root: Dir(root),
		salt: salt,
		tick: time.NewTicker(session.Tick()),
	}
	// Manages the session lifetime.
	go func() {
		for range s.tick.C {
			s.session(session.Lifetime())
		}
	}()
	return s
}

func (s *Safe) hash() (crypto.Hash, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.log == nil {
		return nil, ErrNotLogged
	}
	b, err := s.log.pass.NewCipher(s.salt)
	if err != nil {
		return nil, err
	}
	return crypto.New(b), nil
}

type auth struct {
	pass *safe.Passphrase
	at   time.Time
}

func (s *Safe) logged() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.log == nil {
		return ErrNotLogged
	}
	return nil
}

func (s *Safe) login(p *safe.Passphrase) {
	s.mu.Lock()
	s.log = &auth{pass: p, at: time.Now()}
	s.mu.Unlock()
}

func (s *Safe) session(d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.log == nil {
		// Not logged in
		return
	}
	if time.Now().After(s.log.at.Add(d)) {
		// Session has expired, close it.
		s.log = nil
	}
}

// Root ...
func (s *Safe) Root() Dir {
	return s.root
}

// Close ...
func (s *Safe) Close() (err error) {
	if s.tick != nil {
		s.tick.Stop()
	}
	if s.db != nil {
		err = s.db.Close()
	}
	return
}

// GeneratePassword ...
func GeneratePassword(length, numDigits, numSymbols int, noUpper, allowRepeat bool) (string, error) {
	return password.Generate(length, numDigits, numSymbols, noUpper, allowRepeat)
}
