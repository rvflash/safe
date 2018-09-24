// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package app_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/rvflash/safe"
	"github.com/rvflash/safe/app"
)

func checkVault(v *safe.Vault, tag, name string, data app.Login) error {
	var (
		s string
		u *url.URL
	)
	if s = v.Name(); s != name {
		return fmt.Errorf("mismatch name: got=%q exp=%q", s, name)
	}
	if s = v.Tag().Name(); s != tag {
		return fmt.Errorf("mismatch tag name: got=%q exp=%q", s, tag)
	}
	if s = v.Login().Name; s != data.Username {
		return fmt.Errorf("mismatch username: got=%q exp=%q", s, data.Username)
	}
	if s = v.Login().Password; s != data.Password {
		return fmt.Errorf("mismatch username: got=%q exp=%q", s, data.Password)
	}
	if u = v.Login().URL; u != nil && u.String() != data.URL {
		return fmt.Errorf("mismatch username: got=%q exp=%q", s, data.URL)
	}
	if s = v.Login().Note; s != data.Note {
		return fmt.Errorf("mismatch username: got=%q exp=%q", s, data.Note)
	}
	return nil
}

func TestSafe_CreateVault(t *testing.T) {
	var dt = []struct {
		app       *app.Safe
		tag, name string
		data      app.Login
		err       error
	}{
		{app: app.New(newService(), "", "", newSession()), err: app.ErrNotLogged},
		{
			app: loggedApp(),
			tag: "tag", name: "vault",
			data: app.Login{Username: "username", Password: "password"},
		},
		{
			app: loggedApp(),
			tag: "tag", name: "vault",
			data: app.Login{Username: "username", Password: "password", URL: "http://localhost", Note: "note"},
		},
	}
	var (
		err error
		v   *safe.Vault
	)
	for i, tt := range dt {
		v, err = tt.app.CreateVault(tt.name, tt.tag, tt.data)
		if err != tt.err {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		} else if err == nil {
			if err = checkVault(v, tt.tag, tt.name, tt.data); err != nil {
				t.Errorf("%d. %s", i, err)
			}
		}
	}
}

func TestSafe_DeleteVault(t *testing.T) {
	var dt = []struct {
		app       *app.Safe
		tag, name string
		err       error
	}{
		{app: app.New(newService(), "", "", newSession()), err: app.ErrNotLogged},
		{app: loggedApp(), tag: "tag", name: "vault"},
	}
	var err error
	for i, tt := range dt {
		if err = tt.app.DeleteVault(tt.name, tt.tag); err != tt.err {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		}
	}
}

func TestSafe_ListVaultByNames(t *testing.T) {
	var dt = []struct {
		app *app.Safe
		tag string
		err error
	}{
		{app: app.New(newService(), "", "", newSession()), err: app.ErrNotLogged},
		{app: loggedApp(), tag: "tag"},
	}
	var err error
	for i, tt := range dt {
		if _, err = tt.app.ListVaultByNames(tt.tag, ""); err != tt.err {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		}
	}
}

func TestSafe_Vault(t *testing.T) {
	var dt = []struct {
		app       *app.Safe
		tag, name string
		err       error
	}{
		{app: app.New(newService(), "", "", newSession()), err: app.ErrNotLogged},
		{app: loggedApp(), tag: "tag"},
	}
	var err error
	for i, tt := range dt {
		if _, err = tt.app.Vault(tt.name, tt.tag); err != tt.err {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		}
	}
}

func TestSafe_UpdateVault(t *testing.T) {
	var dt = []struct {
		app       *app.Safe
		tag, name string
		data      app.Login
		err       error
	}{
		{app: app.New(newService(), "", "", newSession()), err: app.ErrNotLogged},
		{
			app: loggedApp(),
			tag: "tag", name: "vault",
			data: app.Login{Username: "username", Password: "password"},
		},
		{
			app: loggedApp(),
			tag: "tag", name: "vault",
			data: app.Login{Username: "username", Password: "password", URL: "http://localhost", Note: "note"},
		},
	}
	var (
		err error
		v   *safe.Vault
	)
	for i, tt := range dt {
		v, err = tt.app.UpdateVault(tt.name, tt.tag, tt.data)
		if err != tt.err {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		} else if err == nil {
			if err = checkVault(v, tt.tag, tt.name, tt.data); err != nil {
				t.Errorf("%d. %s", i, err)
			}
		}
	}
}
