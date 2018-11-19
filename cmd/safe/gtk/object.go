// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// NewBuilder ...
func NewBuilder(xml string) (*gtk.Builder, error) {
	b, err := gtk.BuilderNew()
	if err != nil {
		return nil, err
	}
	if err = b.AddFromString(xml); err != nil {
		return nil, err
	}
	return b, nil
}

// Object ...
type Object struct {
	b    *gtk.Builder
	main glib.IObject
}

// NewObject ...
func NewObject(id, xml string) (*Object, error) {
	b, err := NewBuilder(xml)
	if err != nil {
		return nil, err
	}
	o := &Object{b: b}
	if o.main, err = o.ID(id); err != nil {
		return nil, err
	}
	return o, nil
}

// Builder ...
func (o *Object) Builder() *gtk.Builder {
	return o.b
}

// ButtonClicked implements the Container interface.
func (o *Object) ButtonClicked(id string, fn func()) (err error) {
	e, err := o.ID(id)
	if err != nil {
		return
	}
	switch b := e.(type) {
	case *gtk.Button:
		_, err = b.Connect("clicked", fn)
	case *gtk.RadioButton:
		_, err = b.Connect("clicked", fn)
	default:
		err = ErrContainer
	}
	return
}

// ReadEntry ...
func (o *Object) ReadEntry(id string) (string, error) {
	d, err := o.ID(id)
	if err != nil {
		return "", err
	}
	return d.(*gtk.Entry).GetText()
}

// WriteEntry ...
func (o *Object) WriteEntry(id, text string) error {
	d, err := o.ID(id)
	if err != nil {
		return err
	}
	d.(*gtk.Entry).SetText(text)
	return nil
}

// ReadSpinButton ...
func (o *Object) ReadSpinButton(id string) (int, error) {
	d, err := o.ID(id)
	if err != nil {
		return 0, err
	}
	return int(d.(*gtk.SpinButton).GetValue()), nil
}

// WriteSpinButton ...
func (o *Object) WriteSpinButton(id string, num int) error {
	d, err := o.ID(id)
	if err != nil {
		return err
	}
	d.(*gtk.SpinButton).SetValue(float64(num))
	return nil
}

// IsActivated ...
func (o *Object) IsActivated(id string) (bool, error) {
	d, err := o.ID(id)
	if err != nil {
		return false, err
	}
	return d.(*gtk.CheckButton).GetActive(), nil
}

// Activate ...
func (o *Object) Activate(id string, enable bool) error {
	d, err := o.ID(id)
	if err != nil {
		return err
	}
	d.(*gtk.CheckButton).SetActive(enable)
	return nil
}

// Error ...
func (o *Object) Error(id, msg string) (err error) {
	d, err := o.ID(id)
	if err != nil {
		return
	}
	l := d.(*gtk.Label)
	if msg == "" {
		l.Hide()
	} else {
		l.SetLabel(msg)
		l.Show()
	}
	return
}

// ID ...
func (o *Object) ID(name string) (glib.IObject, error) {
	if o.b == nil {
		return nil, ErrContainer
	}
	return o.b.GetObject(name)
}

// Object ...
func (o *Object) Box() glib.IObject {
	return o.main
}
