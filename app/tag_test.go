// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package app_test

import (
	"testing"

	"github.com/rvflash/safe"
	"github.com/rvflash/safe/app"
)

func TestSafe_ListTagByNames(t *testing.T) {
	var dt = []struct {
		app *app.Safe
		err error
	}{
		{app: app.New(newService(), "", "", newSession()), err: app.ErrNotLogged},
		{app: loggedApp()},
	}
	var err error
	for i, tt := range dt {
		if _, err = tt.app.ListTagByNames(); err != tt.err {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		}
	}
}

func TestSafe_CreateTag(t *testing.T) {
	var dt = []struct {
		app  *app.Safe
		name string
		err  error
	}{
		{app: app.New(newService(), "", "", newSession()), name: "tag", err: app.ErrNotLogged},
		{app: loggedApp(), name: "tag"},
	}
	var (
		name string
		tag  *safe.Tag
		err  error
	)
	for i, tt := range dt {
		tag, err = tt.app.CreateTag(tt.name)
		if err != tt.err {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		}
		if err == nil {
			if name = tag.Name(); name != tt.name {
				t.Errorf("%d. mismatch result: got=%q exp=%q", i, name, tt.name)
			}
		}
	}
}

func TestSafe_DeleteTag(t *testing.T) {
	var dt = []struct {
		app  *app.Safe
		name string
		err  error
	}{
		{app: app.New(newService(), "", "", newSession()), name: "tag", err: app.ErrNotLogged},
		{app: loggedApp(), name: "tag"},
	}
	var err error
	for i, tt := range dt {
		if err = tt.app.DeleteTag(tt.name); err != tt.err {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		}
	}
}
