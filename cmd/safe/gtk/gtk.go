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
func NewLevelBar(val, min, max int) (*gtk.LevelBar, error) {
	l, err := gtk.LevelBarNewForInterval(float64(min), float64(max))
	if err != nil {
		return nil, err
	}
	l.SetValue(float64(val))

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

// NewHBox ...
func NewHBox() (*gtk.Box, error) {
	return gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, defaultMargin)
}

// NewVBox ...
func NewVBox() (*gtk.Box, error) {
	return gtk.BoxNew(gtk.ORIENTATION_VERTICAL, defaultMargin)
}

// NewScrolledWindow ...
func NewScrolledWindow(h, v bool) (*gtk.ScrolledWindow, error) {
	s, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	policy := func(need bool) gtk.PolicyType {
		if need {
			return gtk.POLICY_AUTOMATIC
		}
		return gtk.POLICY_NEVER
	}
	s.SetHExpand(h)
	s.SetVExpand(v)
	s.SetPolicy(policy(h), policy(v))
	return s, nil
}
