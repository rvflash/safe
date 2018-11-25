// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import (
	"fmt"
	"sort"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rvflash/elapsed"
	"github.com/rvflash/safe"
)

type vaults struct {
	tag  string
	data []*safe.Vault
}

func newDataTable(tag string, data []*safe.Vault) *vaults {
	return &vaults{tag: tag, data: data}
}

// Cols ...
func (d vaults) Cols() []string {
	return []string{"Name", "Username", "Password", "URL", "Last updated", "Note", "*"}
}

// ColSizes ...
func (d vaults) ColSizes() []int {
	return []int{100, 170, 130, 150, 100, 120, 10}
}

// Rows ...
func (d *vaults) Rows(upd, del FuncTwo, cp FuncOne) (rs [][]interface{}, err error) {
	const (
		name = iota
		username
		password
		url
		lastUpdated
		note
		actions
	)
	rs = make([][]interface{}, len(d.data))
	for p, v := range d.data {
		rs[p] = make([]interface{}, len(d.Cols()))
		rs[p][name] = v.Name()
		rs[p][username] = v.Login().Name
		rs[p][password], err = NewLevelBar(float64(v.Login().Strength()), 0, 4)
		if err != nil {
			return nil, err
		}
		if v.Login().URL == nil {
			rs[p][url], err = gtk.LabelNew("-")
		} else {
			rs[p][url], err = gtk.LinkButtonNewWithLabel(v.Login().URL.String(), v.Login().URL.Host)
		}
		if err != nil {
			return nil, err
		}
		rs[p][lastUpdated] = elapsed.Time(v.Login().LastUpdate)
		rs[p][note] = v.Login().Note
		rs[p][actions], err = d.newMenuButton(upd, del, cp, p)
		if err != nil {
			return nil, err
		}
	}
	return rs, nil
}

func (d vaults) newMenuButton(upd, del FuncTwo, cp FuncOne, line int) (*gtk.MenuButton, error) {
	m, err := NewMenuButton()
	if err != nil {
		return nil, err
	}
	b, err := NewButton("Edit Vault", func() {
		upd(d.tag, d.data[line].Name())
		m.MenuButton().Hide()
	})
	if err != nil {
		return nil, err
	}
	m.Add(b)

	b, err = NewButton("Copy Username", func() {
		cp(d.data[line].Login().Name)
		m.MenuButton().Hide()
	})
	if err != nil {
		return nil, err
	}
	m.Add(b)

	b, err = NewButton("Copy Password", func() {
		cp(d.data[line].Login().Password)
		m.MenuButton().Hide()
	})
	if err != nil {
		return nil, err
	}
	m.Add(b)

	sep, err := gtk.SeparatorMenuItemNew()
	if err != nil {
		return nil, err
	}
	m.Add(sep)

	b, err = NewButton("Delete Vault", func() {
		del(d.tag, d.data[line].Name())
		m.MenuButton().Hide()
	})
	if err != nil {
		return nil, err
	}
	return m.MenuButton(), nil
}

// Title ...
func (d vaults) Title() string {
	return d.tag
}

// Types ...
func (d vaults) Types() []glib.Type {
	return []glib.Type{
		glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_POINTER, glib.TYPE_POINTER,
		glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_POINTER,
	}
}

// Note ...
type Note struct {
	box   *gtk.Box
	title string
}

// NewNote ...
func NewNote(v DataTable, add FuncOne, upd, del FuncTwo, cp FuncOne) (*Note, error) {
	hb, err := NewVBox()
	if err != nil {
		return nil, err
	}

	// Adds a button to create a new vault.
	b, err := NewButton("New Vault", func() {
		add(v.Title())
	})
	if err != nil {
		return nil, err
	}
	b.SetHExpand(false)
	hb.Add(b)

	// Lists all tag's vaults.
	tv, err := NewTreeView(v.Cols(), v.ColSizes(), v.Types())
	if err != nil {
		return nil, err
	}
	rs, err := v.Rows(upd, del, cp)
	if err != nil {
		return nil, err
	}
	for _, d := range rs {
		if err = tv.AddRow(d...); err != nil {
			return nil, err
		}
	}
	sw, err := tv.ScrollTreeView(true, true)
	if err != nil {
		return nil, err
	}
	hb.PackStart(sw, true, true, 0)

	return &Note{title: v.Title(), box: hb}, nil
}

// Title ...
func (n *Note) Title() string {
	return n.title
}

// Widget ...
func (n *Note) Widget() *gtk.Box {
	return n.box
}

// Notebook ...
type Notebook struct {
	c     *gtk.Notebook
	index []string
}

// NewNotebook ...
func NewNotebook(name string, clicked Func, pages ...*Note) (*Notebook, error) {
	c, err := gtk.NotebookNew()
	if err != nil {
		return nil, err
	}

	// Title
	l, err := NewLabel(name, font(titleFontSize), 0, defaultMargin, 0, defaultMargin)
	if err != nil {
		return nil, err
	}
	l.Show()
	c.SetActionWidget(l, gtk.PACK_START)

	// New Page
	b, err := NewButton("Add", clicked)
	if err != nil {
		return nil, err
	}
	b.Show()
	c.SetActionWidget(b, gtk.PACK_END)

	// Adds pages
	book := &Notebook{c: c}
	book.init(pages)
	for _, note := range pages {
		if err = book.AddPage(note); err != nil {
			return nil, err
		}
	}
	c.SetCurrentPage(0)

	return book, nil
}

func (n *Notebook) init(pages []*Note) {
	n.index = make([]string, len(pages))
	for i, note := range pages {
		n.index[i] = note.Title()
	}
	sort.Strings(n.index)
}

// AddPage ...
func (n *Notebook) AddPage(note *Note) error {
	// Creates the tab
	l, err := NewLabel(note.Title(), font(defaultFontSize))
	if err != nil {
		return err
	}
	// Adds the page content
	i := n.c.InsertPage(note.Widget(), l, n.numPage(note.Title()))
	// Updates the summary.
	if len(n.index) == i {
		n.index = append(n.index, "")
		copy(n.index[i+1:], n.index[i:])
		n.index[i] = note.Title()
	}
	return nil
}

func (n *Notebook) numPage(name string) int {
	if n.index != nil {
		return sort.SearchStrings(n.index, name)
	}
	return 0
}

// Notebook ...
func (n *Notebook) Notebook() *gtk.Notebook {
	return n.c
}

func font(size int) string {
	return fmt.Sprintf("%s %d", defaultFont, size)
}
