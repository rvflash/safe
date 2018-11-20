// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gotk3/gotk3/gtk"
	"github.com/rvflash/safe/app"
	"github.com/rvflash/safe/cmd/safe/static"
)

const (
	mainUIID   = "app"
	mainUIPath = "/static/ui/app.ui"

	signUIID   = "sign"
	signUIPath = "/static/ui/sign.ui"

	tagUIID   = "tag"
	tagUIPath = "/static/ui/tag.ui"

	vaultUIID   = "vault"
	vaultUIPath = "/static/ui/vault.ui"
)

// ErrUnkDialog ...
var ErrUnkDialog = errors.New("unknown dialog")

// App ...
type App struct {
	fs http.FileSystem
	w  *Window
}

// Init ...
func Init(db *app.Safe, log Logger, debug bool) (*App, error) {
	a := &App{
		fs: static.FS(!debug),
	}
	if !debug {
		log = nil
	}
	gtk.Init(&os.Args)

	// Default theme
	s, err := gtk.SettingsGetDefault()
	if err != nil {
		return nil, err
	}
	s.Set("gtk-application-prefer-dark-theme", true)

	// Containers (window and dialogs)
	if err = a.window(mainUIID, mainUIPath, db, log); err != nil {
		return nil, err
	}
	if err = a.dialog(signUIID, signUIPath); err != nil {
		return nil, err
	}
	if err = a.dialog(tagUIID, tagUIPath); err != nil {
		return nil, err
	}
	if err = a.dialog(vaultUIID, vaultUIPath); err != nil {
		return nil, err
	}
	if err = a.w.Init(); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) dialog(id, path string) (err error) {
	xml, err := a.readFile(path)
	if err != nil {
		return
	}
	d, err := NewDialog(a.w, id, xml)
	if err != nil {
		return
	}
	var c VisibleContainer
	switch id {
	case signUIID:
		c = &SignDialog{Dialog: d}
	case tagUIID:
		c = &TagDialog{Dialog: d}
	case vaultUIID:
		c = &VaultDialog{Dialog: d}
	default:
		return ErrUnkDialog
	}
	return a.w.AddDialog(c)
}

func (a *App) window(id, path string, db *app.Safe, log Logger) (err error) {
	xml, err := a.readFile(path)
	if err != nil {
		return
	}
	a.w, err = NewWindow(id, xml, db, log)
	return
}

func (a *App) readFile(path string) (string, error) {
	f, err := a.fs.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Run ...
func (a *App) Run() {
	if a.w != nil {
		a.w.Show()
		gtk.Main()
	}
}
