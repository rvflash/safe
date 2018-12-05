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
	d *gtk.Dialog
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
	w := o.IObject().(*gtk.Dialog)
	w.SetTransientFor(parent.Window())
	//w.SetPosition(gtk.WIN_POS_CENTER_ON_PARENT)
	//w.SetModal(true)
	//w.SetFocusVisible(true)

	return &Dialog{Object: o, p: parent, d: w}, nil
}

// App implements the Plugin interface.
func (d *Dialog) App() *app.Safe {
	return d.p.App()
}

// Log implements the Plugin interface.
func (d *Dialog) Log(format string, args ...interface{}) {
	d.p.Log(format, args...)
}

// Hide implements the Visibility interface.
func (d *Dialog) Hide() {
	if d.d != nil {
		d.d.Hide()
	}
}

// Show implements the Visibility interface.
func (d *Dialog) Show() {
	if d.d != nil {
		d.d.Show()
	}
}

// Parent ...
func (d *Dialog) Parent() *Window {
	return d.p
}

// IWidget ...
func (d *Dialog) IWidget() gtk.IWidget {
	return d.d
}
