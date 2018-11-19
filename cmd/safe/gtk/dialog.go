// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/rvflash/safe/app"
)

// Dialog ...
type Dialog struct {
	*Object
	p *Window
}

// NewDialog ...
func NewDialog(parent *Window, id, xml string) (*Dialog, error) {
	if parent == nil {
		return nil, ErrContainer
	}
	o, err := NewObject(id, xml)
	if err != nil {
		return nil, err
	}
	// Checks if the interface is valid.
	d := &Dialog{Object: o, p: parent}
	// Attaches this dialog on the main window.
	d.Dialog().SetTransientFor(d.Parent())

	return d, nil
}

// App ...
func (d *Dialog) App() *app.Safe {
	return d.p.App()
}

// Log ...
func (d *Dialog) Log(format string, args ...interface{}) {
	d.p.Log(format, args)
}

// Dialog ...
func (d *Dialog) Dialog() *gtk.Dialog {
	if d.Box() == nil {
		return nil
	}
	return d.Box().(*gtk.Dialog)
}

// Hide implements the Visibility interface.
func (d *Dialog) Hide() {
	if c := d.Dialog(); c != nil {
		c.Hide()
	}
}

// Parent ...
func (d *Dialog) Parent() gtk.IWindow {
	if d.p == nil {
		return nil
	}
	return d.p.Box().(gtk.IWindow)
}

// Show implements the Visibility interface.
func (d *Dialog) Show() {
	if c := d.Dialog(); c != nil {
		c.Show()
	}
}
