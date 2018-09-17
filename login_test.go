// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package safe_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/rvflash/safe"
)

const (
	name = "rv"
	pass = "secret"
)

func TestLogin_Valid(t *testing.T) {
	var dt = []struct {
		name,
		pass string
		ok bool
	}{
		{},
		{name: name, pass: " "},
		{name: " ", pass: pass},
		{name: name},
		{pass: name},
		{name: name, pass: pass, ok: true},
	}
	var ok bool
	for i, tt := range dt {
		if ok = safe.NewLogin(tt.name, tt.pass).Valid(); ok != tt.ok {
			t.Fatalf("%d. mismatch result: got=%t exp=%t", i, ok, tt.ok)
		}
	}
}

func TestLogin_Safe(t *testing.T) {
	var dt = []struct {
		login *safe.Login
		ok    bool
		err   error
	}{
		{login: &safe.Login{}, err: safe.ErrInvalid},
		{login: &safe.Login{Name: name}, err: safe.ErrInvalid},
		{login: &safe.Login{Name: name, Password: pass}, err: safe.ErrTooShort},
		{
			login: &safe.Login{
				Name:       name,
				Password:   "best-secret-ever",
				LastUpdate: time.Now().Add(-time.Hour * 24 * 365),
			},
			err: safe.ErrOutdated,
		},
		{login: &safe.Login{Name: name, Password: "best-secret-ever", LastUpdate: time.Now()}, ok: true},
	}
	var (
		ok  bool
		err error
	)
	for i, tt := range dt {
		ok, err = tt.login.Safe()
		if !reflect.DeepEqual(err, tt.err) {
			t.Fatalf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		}
		if ok != tt.ok {
			t.Errorf("%d. mismatch result: got=%t exp=%t", i, ok, tt.ok)
		}
	}
}
