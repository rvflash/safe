// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import (
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// TreeView ...
type TreeView struct {
	v         *gtk.TreeView
	s         *gtk.ListStore
	f         *gtk.TreeModelFilter
	d         DataTable
	filtering int
}

// NewTreeView ...
func NewTreeView(list DataTable, changed FuncTwo) (*TreeView, error) {
	v, err := gtk.TreeViewNew()
	if err != nil {
		return nil, err
	}
	v.SetMarginBottom(defaultMargin)
	v.SetMarginEnd(defaultMargin)
	v.SetMarginStart(defaultMargin)
	v.SetMarginTop(0)
	v.SetName(list.Title())

	t := &TreeView{v: v, d: list, filtering: -1}
	if err := t.buildColumns(); err != nil {
		return nil, err
	}
	if err = t.applyModel(); err != nil {
		return nil, err
	}
	if err = t.onChanged(changed); err != nil {
		return nil, err
	}
	for _, d := range list.Rows() {
		if err = t.AddRow(d...); err != nil {
			return nil, err
		}
	}
	return t, nil
}
func (t *TreeView) buildColumns() error {
	// Integrity controls.
	if len(t.d.Cols()) == 0 || len(t.d.Cols()) != len(t.d.ColSizes()) {
		return ErrUndColumn
	}
	r, err := gtk.CellRendererTextNew()
	if err != nil {
		return err
	}
	var (
		c *gtk.TreeViewColumn
		w = t.d.ColSizes()
	)
	for i, s := range t.d.Cols() {
		c, err = gtk.TreeViewColumnNewWithAttribute(s, r, "text", i)
		if err != nil {
			return err
		}
		if w[i] > 0 {
			c.SetMinWidth(w[i])
		} else {
			t.filtering = i
			c.SetVisible(false)
		}
		t.v.AppendColumn(c)
	}
	return err
}

func (t *TreeView) applyModel() (err error) {
	t.s, err = gtk.ListStoreNew(t.d.Types()...)
	if err != nil {
		return
	}
	if t.filtering > -1 {
		p, _ := t.v.GetCursor()
		t.f, err = t.s.FilterNew(p)
		if err != nil {
			return
		}
		t.f.SetVisibleColumn(t.filtering)
		t.v.SetModel(t.f)
	} else {
		t.v.SetModel(t.s)
	}
	return
}

func (t *TreeView) onChanged(changed FuncTwo) error {
	if changed == nil {
		// Nothing to do.
		return nil
	}
	n, err := t.v.GetName()
	if err != nil {
		return err
	}
	s, err := t.v.GetSelection()
	if err != nil {
		return err
	}
	s.SetMode(gtk.SELECTION_SINGLE)
	_, err = s.Connect("changed", func(c *gtk.TreeSelection) {
		m, i, ok := s.GetSelected()
		if !ok {
			return
		}
		v, err := m.(*gtk.TreeModel).GetValue(i, t.d.ColID())
		if err != nil {
			return
		}
		s, err := v.GetString()
		if err != nil {
			return
		}
		// The treeView's name and the value of the reference's column.
		// ex: tag and vault names (at colID #0)
		changed(n, s)
	})
	return err
}

func (t *TreeView) Search(q string) {
	if t.TreeView() == nil {
		return
	}
	var (
		b   bool
		err error
		s   string
		v   *glib.Value
	)
	i, ok := t.s.GetIterFirst()
	for ok {
		v, err = t.s.GetValue(i, t.d.ColID())
		if err != nil {
			return
		}
		s, err = v.GetString()
		if err != nil {
			return
		}
		if q == "" {
			// No query, so show all.
			b = true
		} else {
			// The reference's column must contain the query string.
			b = strings.Contains(strings.ToLower(s), strings.ToLower(q))
		}
		if err = t.s.SetValue(i, t.filtering, b); err != nil {
			return
		}
		ok = t.s.IterNext(i)
	}
}

// AddRow ...
func (t *TreeView) AddRow(d ...interface{}) (err error) {
	if t.TreeView() == nil {
		return ErrUndObject
	}
	l := t.s.Append()
	for i := range t.d.Cols() {
		if err = t.s.SetValue(l, i, d[i]); err != nil {
			return
		}
	}
	return
}

// ScrollTreeView ...
func (t *TreeView) ScrollTreeView(v, h bool) (*gtk.ScrolledWindow, error) {
	if t.TreeView() == nil {
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
