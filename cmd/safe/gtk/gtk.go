// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const (
	defaultMargin   = 10
	defaultFont     = "Helvetica"
	defaultFontSize = 12
	titleFontSize   = 16
)

// NewSearchEntry ...
func NewSearchEntry(placeholder string, search FuncOne, margin ...int) (*gtk.SearchEntry, error) {
	top, right, bottom, left := spaces(margin)
	w, err := gtk.SearchEntryNew()
	if err != nil {
		return nil, err
	}
	w.SetPlaceholderText(placeholder)
	w.SetMarginBottom(bottom)
	w.SetMarginEnd(right)
	w.SetMarginStart(left)
	w.SetMarginTop(top)
	if _, err = w.Connect("search-changed", func(e *gtk.SearchEntry) {
		s, _ := e.GetText()
		search(s)
	}); err != nil {
		return nil, err
	}
	if _, err = w.Connect("stop-search", func() {
		search("")
	}); err != nil {
		return nil, err
	}
	return w, nil
}

// NewButton ...
func NewButton(label string, clicked Func, margin ...int) (*gtk.Button, error) {
	top, right, bottom, left := spaces(margin)
	w, err := gtk.ButtonNewWithLabel(label)
	if err != nil {
		return nil, err
	}
	w.SetMarginBottom(bottom)
	w.SetMarginEnd(right)
	w.SetMarginStart(left)
	w.SetMarginTop(top)
	if _, err = w.Connect("clicked", clicked); err != nil {
		return nil, err
	}
	return w, nil
}

// NewLabel ...
func NewLabel(label, font string, margin ...int) (*gtk.Label, error) {
	top, right, bottom, left := spaces(margin)
	w, err := gtk.LabelNew(label)
	if err != nil {
		return nil, err
	}
	w.SetMarginBottom(bottom)
	w.SetMarginEnd(right)
	w.SetMarginStart(left)
	w.SetMarginTop(top)
	w.SetFont(font)

	return w, nil
}

func spaces(side []int) (top, right, bottom, left int) {
	for i, v := range side {
		switch i {
		case 0:
			top = v
		case 1:
			right = v
		case 2:
			bottom = v
		case 3:
			left = v
		}
	}
	return
}

// NewLevelBar ...
func NewLevelBar(val, min, max float64) (*gtk.LevelBar, error) {
	l, err := gtk.LevelBarNewForInterval(min, max)
	if err != nil {
		return nil, err
	}
	l.SetValue(val)

	return l, nil
}

// Clipboard ...
type Clipboard struct {
	c *gtk.Clipboard
}

// NewClipboard ...
func NewClipboard() (*Clipboard, error) {
	c, err := gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
	if err != nil {
		return nil, err
	}
	return &Clipboard{c: c}, nil
}

// Copy ...
func (c *Clipboard) Copy(text string) {
	c.c.SetText(text)
	if gtk.MainIterationDo(true) {
		c.c.Store()
	}
	_ = gtk.MainIterationDo(true)
}

// MenuButton ...
type MenuButton struct {
	b *gtk.MenuButton
	m *gtk.Menu
}

func NewMenuButton() (*MenuButton, error) {
	b, err := gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	b.SetDirection(gtk.ARROW_DOWN)
	b.SetHAlign(gtk.ALIGN_END)

	// Creates the menu to show on click.
	c, err := gtk.MenuNew()
	if err != nil {
		return nil, err
	}
	b.SetPopup(c)

	return &MenuButton{b: b, m: c}, nil
}

// Add ...
func (m *MenuButton) Add(item gtk.IWidget) {
	if m.m != nil {
		m.m.Add(item)
	}
}

// MenuButton ...
func (m *MenuButton) MenuButton() *gtk.MenuButton {
	return m.b
}

// NewHBox ...
func NewHBox() (*gtk.Box, error) {
	return gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, defaultMargin)
}

// NewVBox ...
func NewVBox() (*gtk.Box, error) {
	return gtk.BoxNew(gtk.ORIENTATION_VERTICAL, defaultMargin)
}
