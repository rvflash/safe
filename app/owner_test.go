// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package app_test

import (
	"testing"
	"time"

	"github.com/rvflash/safe/app"
)

// test cases
func newService() *service {
	return &service{}
}

func hisService() *service {
	return &service{owned: true, owner: true}
}

func anotherService() *service {
	return &service{owned: true, owner: false}
}

func loggedApp() *app.Safe {
	a := app.New(hisService(), "salt", "", newSession())
	_ = a.SignIn("pass")
	return a
}

func newSession() app.Session {
	return app.NewSession(50*time.Millisecond, 5*time.Millisecond)
}

func TestSafe_Login(t *testing.T) {
	var dt = []struct {
		app *app.Safe
		err error
	}{
		{app: app.New(newService(), "", "", newSession())},
		{app: app.New(hisService(), "", "", newSession())},
		{app: app.New(anotherService(), "", "", newSession()), err: app.ErrNotLogged},
	}
	var err error
	for i, tt := range dt {
		if err = tt.app.Login("pass"); err != tt.err {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		}
		if err = tt.app.Logged(); err != tt.err {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		} else if tt.err == nil {
			time.Sleep(60 * time.Millisecond)
			if err = tt.app.Logged(); err != app.ErrNotLogged {
				t.Errorf("%d. unexpected error: got=%q", i, err)
			}
			if err = tt.app.Close(); err != nil {
				t.Errorf("%d. unexpected error: got=%q", i, err)
			}
		}
	}
}

func TestSafe_LogOut(t *testing.T) {
	a := app.New(hisService(), "", "", newSession())
	err := a.SignUp("pass")
	if err != nil {
		t.Fatalf("unexpected error: got=%q", err)
	}
	if err = a.Logged(); err != nil {
		t.Fatalf("unexpected error: got=%q", err)
	}
	a.LogOut()
	if err = a.Logged(); err != app.ErrNotLogged {
		t.Fatalf("unexpected error: got=%q", err)
	}
}

func TestSafe_SignUp(t *testing.T) {
	var dt = []struct {
		app *app.Safe
		err error
	}{
		{app: app.New(newService(), "", "", newSession())},
		{app: app.New(hisService(), "", "", newSession())},
		{app: app.New(anotherService(), "", "", newSession())},
	}
	var err error
	for i, tt := range dt {
		if err = tt.app.SignUp("pass"); err != nil {
			t.Errorf("%d. unexpected error: got=%q", i, err)
		}
		if err = tt.app.Logged(); err != tt.err {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		}
	}
}

func TestSafe_SignIn(t *testing.T) {
	var dt = []struct {
		app    *app.Safe
		logged bool
		err    error
	}{
		{app: app.New(newService(), "", "", newSession()), err: app.ErrNotFound},
		{app: app.New(hisService(), "", "", newSession()), logged: true},
		{app: app.New(anotherService(), "", "", newSession()), err: app.ErrNotLogged},
	}
	var (
		err    error
		logged bool
	)
	for i, tt := range dt {
		if logged = tt.app.SignIn("pass"); logged != tt.logged {
			t.Errorf("%d. mismatch result: got=%t exp=%t", i, logged, tt.logged)
		}
		if err = tt.app.Logged(); err != tt.err {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		}
	}
}
