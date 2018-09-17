// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package app_test

import (
	"net/url"
	"testing"

	"github.com/rvflash/safe"
	"github.com/rvflash/safe/app"
)

func vault(name, pass, rawURL, note string) url.Values {
	m := make(url.Values)
	if name != "" {
		m[app.FormUser] = []string{name}
	}
	if pass != "" {
		m[app.FormPass] = []string{pass}
	}
	if rawURL != "" {
		m[app.FormURL] = []string{rawURL}
	}
	if note != "" {
		m[app.FormNote] = []string{note}
	}
	return m
}

func TestSafe_CreateVault(t *testing.T) {
	var dt = []struct {
		app       *app.Safe
		tag, name string
		data      url.Values
		err       error
	}{
		{app: app.New(newService(), "", "", newSession()), err: app.ErrNotLogged},
		{
			app: loggedApp(),
			tag: "tag", name: "vault",
			data: vault("username", "password", "", ""),
		},
		{
			app: loggedApp(),
			tag: "tag", name: "vault",
			data: vault("username", "password", "http://localhost", "note"),
		},
	}
	var (
		s     string
		u     *url.URL
		vault *safe.Vault
		err   error
	)
	for i, tt := range dt {
		vault, err = tt.app.CreateVault(tt.name, tt.tag, tt.data)
		if err != tt.err {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		}
		if err == nil {
			if s = vault.Name(); s != tt.name {
				t.Errorf("%d. mismatch name: got=%q exp=%q", i, s, tt.name)
			}
			if s = vault.Tag().Name(); s != tt.tag {
				t.Errorf("%d. mismatch tag name: got=%q exp=%q", i, s, tt.tag)
			}
			if s = vault.Login().Name; s != tt.data.Get(app.FormUser) {
				t.Errorf("%d. mismatch username: got=%q exp=%q", i, s, tt.data.Get(app.FormUser))
			}
			if s = vault.Login().Password; s != tt.data.Get(app.FormPass) {
				t.Errorf("%d. mismatch username: got=%q exp=%q", i, s, tt.data.Get(app.FormPass))
			}
			if u = vault.Login().URL; u != nil && u.String() != tt.data.Get(app.FormURL) {
				t.Errorf("%d. mismatch username: got=%q exp=%q", i, s, tt.data.Get(app.FormURL))
			}
			if s = vault.Login().Note; s != tt.data.Get(app.FormNote) {
				t.Errorf("%d. mismatch username: got=%q exp=%q", i, s, tt.data.Get(app.FormNote))
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
		data      url.Values
		err       error
	}{
		{app: app.New(newService(), "", "", newSession()), err: app.ErrNotLogged},
		{
			app: loggedApp(),
			tag: "tag", name: "vault",
			data: vault("username", "password", "", ""),
		},
		{
			app: loggedApp(),
			tag: "tag", name: "vault",
			data: vault("username", "password", "http://localhost", "note"),
		},
	}
	var (
		s     string
		u     *url.URL
		vault *safe.Vault
		err   error
	)
	for i, tt := range dt {
		vault, err = tt.app.UpdateVault(tt.name, tt.tag, tt.data)
		if err != tt.err {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		}
		if err == nil {
			if s = vault.Name(); s != tt.name {
				t.Errorf("%d. mismatch name: got=%q exp=%q", i, s, tt.name)
			}
			if s = vault.Tag().Name(); s != tt.tag {
				t.Errorf("%d. mismatch tag name: got=%q exp=%q", i, s, tt.tag)
			}
			if s = vault.Login().Name; s != tt.data.Get(app.FormUser) {
				t.Errorf("%d. mismatch username: got=%q exp=%q", i, s, tt.data.Get(app.FormUser))
			}
			if s = vault.Login().Password; s != tt.data.Get(app.FormPass) {
				t.Errorf("%d. mismatch username: got=%q exp=%q", i, s, tt.data.Get(app.FormPass))
			}
			if u = vault.Login().URL; u != nil && u.String() != tt.data.Get(app.FormURL) {
				t.Errorf("%d. mismatch username: got=%q exp=%q", i, s, tt.data.Get(app.FormURL))
			}
			if s = vault.Login().Note; s != tt.data.Get(app.FormNote) {
				t.Errorf("%d. mismatch username: got=%q exp=%q", i, s, tt.data.Get(app.FormNote))
			}
		}
	}
}
