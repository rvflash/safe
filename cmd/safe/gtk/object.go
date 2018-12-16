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

// Object ...
type Object struct {
	b    *gtk.Builder
	main glib.IObject
}

// Container ...
func (o *Object) IObject() glib.IObject {
	return o.main
}

// Closed ...
func (o *Object) Closed(fn func()) (err error) {
	d := o.main.(*gtk.Dialog)
	_, err = d.Connect("response", fn)
	if err != nil {
		return
	}
	_, err = d.Connect("delete-event", func() bool {
		return true
	})
	return
}

// Error ...
func (o *Object) Error(id, msg string) (err error) {
	d, err := o.ID(id + "_bar")
	if err != nil {
		return
	}
	m := d.(*gtk.InfoBar)
	if msg == "" {
		m.Hide()
		return
	}
	l, err := o.ID(id)
	if err != nil {
		return
	}
	l.(*gtk.Label).SetLabel(msg)
	m.Show()

	return
}

// ID ...
func (o *Object) ID(name string) (glib.IObject, error) {
	if o.b == nil {
		return nil, ErrContainer
	}
	return o.b.GetObject(name)
}

// ButtonClicked ...
func (o *Object) ButtonClicked(id string, fn Func) (err error) {
	return o.connect("clicked", id, fn)
}

// EnterPressed ...
func (o *Object) EnterPressed(id string, fn Func) error {
	return o.connect("activate", id, fn)
}

func (o *Object) connect(signal string, id string, fn Func) (err error) {
	e, err := o.ID(id)
	if err != nil {
		return
	}
	switch b := e.(type) {
	case *gtk.Button:
		_, err = b.Connect(signal, fn)
	case *gtk.RadioButton:
		_, err = b.Connect(signal, fn)
	case *gtk.Entry:
		_, err = b.Connect(signal, fn)
	default:
		err = ErrContainer
	}
	return
}

// Focus ...
func (o *Object) Focus(id string)  (err error) {
	e, err := o.ID(id)
	if err != nil {
		return
	}
	switch b := e.(type) {
	case *gtk.Entry:
		b.GrabFocus()
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
	return d.(*gtk.SpinButton).GetValueAsInt(), nil
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
func (o *Object) IsActivated(id string) (enable bool, err error) {
	d, err := o.ID(id)
	if err != nil {
		return false, err
	}
	switch b := d.(type) {
	case *gtk.CheckButton:
		enable = b.GetActive()
	case *gtk.RadioButton:
		enable = b.GetActive()
	}
	return
}

// Activate ...
func (o *Object) Activate(id string, enable bool) error {
	d, err := o.ID(id)
	if err != nil {
		return err
	}
	switch b := d.(type) {
	case *gtk.CheckButton:
		b.SetActive(enable)
	case *gtk.RadioButton:
		b.SetActive(enable)
	}
	return nil
}
