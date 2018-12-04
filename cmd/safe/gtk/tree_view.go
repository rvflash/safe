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
	d   DataTable
	mf  *gtk.TreeModelFilter
	col int
	ls  *gtk.ListStore
	pos map[string]*gtk.TreeIter
	s   *gtk.TreeSelection
	v   *gtk.TreeView
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

	s, err := v.GetSelection()
	if err != nil {
		return nil, err
	}
	s.SetMode(gtk.SELECTION_SINGLE)

	t := &TreeView{
		v:   v,
		s:   s,
		d:   list,
		col: -1,
		pos: make(map[string]*gtk.TreeIter),
	}
	return t, t.init(changed)
}

func (t *TreeView) init(changed FuncTwo) (err error) {
	if err = t.buildColumns(); err != nil {
		return
	}
	if err = t.applyModel(); err != nil {
		return
	}
	if err = t.onChanged(changed); err != nil {
		return
	}
	for _, d := range t.d.Rows() {
		if err = t.AddRow(d); err != nil {
			return
		}
	}
	return
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
	// Columns header
	for i, s := range t.d.Cols() {
		c, err = gtk.TreeViewColumnNewWithAttribute(s, r, "text", i)
		if err != nil {
			return err
		}
		// Filtering
		if w[i] > 0 {
			c.SetMinWidth(w[i])
		} else {
			// Invisible column used to col it.
			t.col = i
			c.SetVisible(false)
		}
		// Default sort order
		if i == t.d.ColID() {
			c.SetSortIndicator(true)
			c.SetSortOrder(gtk.SORT_ASCENDING)
		}
		/* Changes the sort order.
		i := i
		c.SetClickable(true)
		c.Connect("clicked", func() {
			if t.v.GetColumn(i).GetSortIndicator() {
				t.ls.SetSortColumnId(i, gtk.SORT_ASCENDING)
			} else {
				t.ls.SetSortColumnId(i, gtk.SORT_DESCENDING)
			}
		})
		*/
		t.v.AppendColumn(c)
	}
	return err
}

func (t *TreeView) applyModel() (err error) {
	t.ls, err = gtk.ListStoreNew(t.d.Types()...)
	if err != nil {
		return
	}
	t.ls.SetSortColumnId(t.d.ColID(), gtk.SORT_ASCENDING)

	if t.col > -1 {
		p, _ := t.v.GetCursor()
		t.mf, err = t.ls.FilterNew(p)
		if err != nil {
			return
		}
		t.mf.SetVisibleColumn(t.col)
		t.v.SetModel(t.mf)
	} else {
		t.v.SetModel(t.ls)
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
	_, err = t.s.Connect("changed", func(c *gtk.TreeSelection) {
		m, i, ok := t.s.GetSelected()
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

// Len ...
func (t *TreeView) Len() int {
	return len(t.pos)
}

// Search ...
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
	i, ok := t.ls.GetIterFirst()
	for ok {
		v, err = t.ls.GetValue(i, t.d.ColID())
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
		if err = t.ls.SetValue(i, t.col, b); err != nil {
			return
		}
		ok = t.ls.IterNext(i)
	}
}

// AddRow ...
func (t *TreeView) AddRow(d []interface{}) error {
	if t.TreeView() == nil {
		return ErrUndObject
	}
	l := t.ls.Append()
	if err := t.setRow(l, d); err != nil {
		return err
	}
	t.pos[d[t.d.ColID()].(string)] = l

	return nil
}

// DelRow ...
func (t *TreeView) DelRow(name string) error {
	if t.TreeView() == nil {
		return ErrUndObject
	}
	_ = t.ls.Remove(t.pos[name])
	delete(t.pos, name)
	return nil
}

// UpdRow ...
func (t *TreeView) UpdRow(name string, d []interface{}) error {
	if t.TreeView() == nil {
		return ErrUndObject
	}
	l, ok := t.pos[name]
	if !ok || len(d) != len(t.d.Cols()) {
		return ErrUndColumn
	}
	return t.setRow(l, d)
}

func (t *TreeView) setRow(l *gtk.TreeIter, d []interface{}) (err error) {
	for i := range t.d.Cols() {
		if err = t.ls.SetValue(l, i, d[i]); err != nil {
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
