// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import (
	"errors"

	"github.com/rvflash/safe"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rvflash/safe/app"
)

// ErrContainer ...
var ErrContainer = errors.New("invalid window")

// Container ...
type Container interface {
	ID(name string) (glib.IObject, error)
	Init() error
}

// Listener ...
// todo
type Listener interface {
	ButtonClicked(id string, fn func()) error
}

/// Visibility ...
type Visibility interface {
	Hide()
	Show()
	Reset() error
}

// VisibleContainer ...
type VisibleContainer interface {
	Container
	Visibility
}

// Window ...
type Window struct {
	*Plug
	sign, tag, vault VisibleContainer
}

// NewApplicationWindow ...
func NewWindow(id, xml string, app *app.Safe, log Logger) (*Window, error) {
	p, err := NewPlug(id, xml, app, log)
	if err != nil {
		return nil, err
	}
	return &Window{Plug: p}, nil
}

// AddDialog ...
func (w *Window) AddDialog(d VisibleContainer) error {
	switch d.(type) {
	case *SignDialog:
		w.sign = d
	case *TagDialog:
		w.tag = d
	case *VaultDialog:
		w.vault = d
	default:
		return ErrUnkDialog
	}
	return nil
}

// AddTag ...
func (w *Window) AddTag(t *safe.Tag) error {
	return nil
}

// AddTag ...
func (w *Window) AddVault(v *safe.Vault) error {
	return nil
}

// Hide implements the Visibility interface.
func (w *Window) Hide() {
	if c := w.Window(); c != nil {
		c.Hide()
	}
}

// Init ...
func (w *Window) Init() (err error) {
	if err = w.connectDestroy(); err != nil {
		return
	}
	if err = w.connectLogOut(); err != nil {
		return
	}
	if err = w.connectNewTag(); err != nil {
		return
	}
	return w.login()
}

func (w *Window) connectDestroy() (err error) {
	_, err = w.Window().Connect("destroy", func() {
		gtk.MainQuit()
	})
	return
}

func (w *Window) connectLogOut() error {
	return w.ButtonClicked("logout", func() {
		gtk.MainQuit()
	})
}

func (w *Window) connectNewTag() error {
	if w.tag == nil {
		return ErrContainer
	}
	err := w.tag.Init()
	if err != nil {
		return err
	}
	return w.ButtonClicked("add_tag", func() {
		w.tag.Reset()
		w.tag.Show()
	})
}

func (w *Window) login() (err error) {
	o, err := w.ID("app_box")
	if err != nil {
		return
	}
	o.(*gtk.Box).Hide()
	w.sign.Init()
	w.sign.Show()
	return
}

// Show implements the Visibility interface.
func (w *Window) Show() {
	if c := w.Window(); c != nil {
		c.Show()
	}
}

// Show implements the Visibility interface.
func (w *Window) ShowAll() {
	if c := w.Window(); c != nil {
		c.ShowAll()
	}
}

// Window ...
func (w *Window) Window() *gtk.Window {
	if w.Box() == nil {
		return nil
	}
	return w.Box().(*gtk.Window)
}
