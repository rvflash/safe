// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/rvflash/safe/app"
	"github.com/rvflash/safe/cmd/safe/static"
)

// Application ...
type Application struct {
	fs  http.FileSystem
	win *Window
}

// Init ...
func Init(args *[]string, db *app.Safe, log Logger, debug bool) (*Application, error) {
	// First at all, initializes GTK.
	gtk.Init(args)

	// Defines the theme to use as default.
	if err := applyDarkTheme(); err != nil {
		return nil, err
	}
	// Debug mode.
	if !debug {
		log = nil
	}

	// Creates the main window.
	w, err := NewWindow(db, log, "Safe", 800, 600)
	if err != nil {
		return nil, err
	}
	a := &Application{
		fs:  static.FS(!debug),
		win: w,
	}

	// Adds dialogs.
	for _, path := range []string{
		"/static/ui/confirm.ui",
		"/static/ui/sign.ui",
		"/static/ui/tag.ui",
		"/static/ui/vault.ui",
	} {
		d, err := a.parseXMLFile(path)
		if err != nil {
			return nil, err
		}
		if err = w.AttachDialog(d); err != nil {
			return nil, err
		}
		a.win.Log("app: %s attached", path)
	}
	return a, nil
}

func applyDarkTheme() error {
	s, err := gtk.SettingsGetDefault()
	if err != nil {
		return err
	}
	return s.Set("gtk-application-prefer-dark-theme", true)
}

func (a *Application) readFile(path string) (string, error) {
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

// /static/ui/vault.ui > vault
func (a *Application) fileID(path string) (string, error) {
	s := filepath.Base(path)
	if s = strings.SplitN(s, ".", 2)[0]; s == "" {
		return "", ErrUndObject
	}
	return s, nil
}

func (a *Application) parseXMLFile(path string) (VisibleWidgetContainer, error) {
	xml, err := a.readFile(path)
	if err != nil {
		return nil, err
	}
	id, err := a.fileID(path)
	if err != nil {
		return nil, err
	}
	d, err := NewDialog(a.win, id, xml)
	if err != nil {
		return nil, err
	}
	switch id {
	case "confirm":
		return &ConfirmDialog{Dialog: d}, nil
	case "sign":
		return &SignDialog{Dialog: d}, nil
	case "tag":
		return &TagDialog{Dialog: d}, nil
	case "vault":
		return &VaultDialog{Dialog: d}, nil
	default:
		return nil, ErrUnkDialog
	}
}

// Run ...
func (a *Application) Run() {
	if a.win == nil {
		return
	}
	// Launches the application.
	a.win.Run()
	// Waiting for exit's demand.
	gtk.Main()
	// Graceful closing
	if err := a.win.Close(); err != nil {
		a.win.Log("app: closing on error: %s", err)
	}
	a.win.Log("app: closed")
}
