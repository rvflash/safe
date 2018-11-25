// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import (
	"log"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	defaultMargin   = 10
	defaultFont     = "Helvetica"
	defaultFontSize = 12
	titleFontSize   = 16
)

// NewButton ...
func NewButton(label string, clicked Func) (*gtk.Button, error) {
	b, err := gtk.ButtonNewWithLabel(label)
	if err != nil {
		return nil, err
	}
	b.SetMarginBottom(defaultMargin)
	b.SetMarginEnd(defaultMargin)
	b.SetMarginStart(defaultMargin)
	b.SetMarginTop(defaultMargin)
	if _, err = b.Connect("clicked", clicked); err != nil {
		return nil, err
	}
	return b, nil
}

// NewLabel ...
func NewLabel(label, font string, margin ...int) (*gtk.Label, error) {
	top, right, bottom, left := spaces(margin)
	l, err := gtk.LabelNew(label)
	if err != nil {
		return nil, err
	}
	l.SetMarginBottom(bottom)
	l.SetMarginEnd(right)
	l.SetMarginStart(left)
	l.SetMarginTop(top)
	l.SetFont(font)

	return l, nil
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

// TreeView ...
type TreeView struct {
	v    *gtk.TreeView
	s    *gtk.ListStore
	cols []int
}

// NewTreeView ...
func NewTreeView(cols []string, sizes []int, types []glib.Type) (*TreeView, error) {
	v, err := gtk.TreeViewNew()
	if err != nil {
		return nil, err
	}
	// connect + select: tmp (debug)
	c, err := v.GetSelection()
	if err != nil {
		return nil, err
	}
	c.SetMode(gtk.SELECTION_SINGLE)
	c.Connect("changed", func(s *gtk.TreeSelection) {
		var iter *gtk.TreeIter
		var model gtk.ITreeModel
		var ok bool
		model, iter, ok = s.GetSelected()
		if ok {
			tpath, err := model.(*gtk.TreeModel).GetPath(iter)
			if err != nil {
				log.Printf("treeSelectionChangedCB: Could not get path from model: %s\n", err)
				return
			}
			log.Printf("treeSelectionChangedCB: selected path: %s\n", tpath)
		}
	})

	t := &TreeView{v: v}
	if err = t.withColumns(cols, sizes); err != nil {
		return nil, err
	}
	if err = t.withStoreTypes(types); err != nil {
		return nil, err
	}
	return t, nil
}

// AddRow ...
func (t *TreeView) AddRow(d ...interface{}) error {
	if t.s == nil {
		return ErrUndObject
	}
	l := t.s.Append()
	for i := range t.cols {
		if i > 1 {
			return nil
		}
		t.s.SetValue(l, i, d[i])
	}
	return nil
}

// withColumns ...
func (t *TreeView) withColumns(d []string, w []int) error {
	if len(d) == 0 || len(d) != len(w) {
		return ErrUndColumn
	}
	r, err := gtk.CellRendererTextNew()
	if err != nil {
		return err
	}
	for i, s := range d {
		c, err := gtk.TreeViewColumnNewWithAttribute(s, r, "text", i)
		if err != nil {
			return err
		}
		c.SetMinWidth(w[i])
		t.v.AppendColumn(c)
		t.cols = append(t.cols, i)
	}
	return nil
}

// withStoreTypes ...
func (t *TreeView) withStoreTypes(d []glib.Type) (err error) {
	t.s, err = gtk.ListStoreNew(d...)
	if err != nil {
		return
	}
	t.v.SetModel(t.s)
	return
}

// ScrollTreeView ...
func (t *TreeView) ScrollTreeView(v, h bool) (*gtk.ScrolledWindow, error) {
	if t.s == nil {
		return nil, ErrUndObject
	}
	s, err := NewScrolledWindow(v, h)
	if err != nil {
		return nil, err
	}
	s.Add(t.TreeView())
	return s, nil
}

// TreeView ...
func (t *TreeView) TreeView() *gtk.TreeView {
	return t.v
}

// NewScrolledWindow
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

// NewHBox ...
func NewHBox() (*gtk.Box, error) {
	return gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, defaultMargin)
}

// NewVBox ...
func NewVBox() (*gtk.Box, error) {
	return gtk.BoxNew(gtk.ORIENTATION_VERTICAL, defaultMargin)
}
