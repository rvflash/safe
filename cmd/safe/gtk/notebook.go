// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import (
	"fmt"
	"sort"

	"github.com/gotk3/gotk3/gtk"
)

// Note ...
type Note struct {
	title  string
	widget *gtk.Paned
}

// NewNote ...
func NewNote(list DataTable, add FuncOne, upd, del FuncTwo, cp FuncOne) (*Note, error) {
	l, err := newNoteSideBar(list, add, func(a, b string) {
		fmt.Printf("show %q in %q\n", a, b)
	})
	if err != nil {
		return nil, err
	}
	r, err := newNoteInfo(upd, del, cp)
	if err != nil {
		return nil, err
	}
	p, err := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	p.Pack1(l, true, false)
	p.Pack2(r, true, false)

	return &Note{title: list.Title(), widget: p}, nil
}

func newNoteInfo(upd, del FuncTwo, cp FuncOne) (gtk.IWidget, error) {
	// Vault actions
	hb, err := NewHBox()
	if err != nil {
		return nil, err
	}
	b, err := NewButton("Edit", func() { fmt.Println("edit") }, 5, 0, 5)
	if err != nil {
		return nil, err
	}
	hb.Add(b)

	b, err = NewButton("Delete", func() { fmt.Println("delete") }, 5, 10, 5)
	if err != nil {
		return nil, err
	}
	hb.Add(b)

	// Header (name and update, delete buttons)
	h, err := gtk.InfoBarNew()
	if err != nil {
		return nil, err
	}
	h.SetMessageType(gtk.MESSAGE_INFO)
	h.PackEnd(hb, false, true, 0)
	//h.SetSizeRequest(480, )

	// Vaults properties
	vb, err := NewVBox()
	if err != nil {
		return nil, err
	}
	vb.SetMarginBottom(defaultMargin)
	vb.SetSizeRequest(200, 400)

	sw, err := NewScrolledWindow(false, true)
	if err != nil {
		return nil, err
	}
	sw.Add(vb)

	mb, err := NewVBox()
	if err != nil {
		return nil, err
	}
	mb.SetMarginTop(defaultMargin / 2)
	mb.PackStart(h, false, true, 0)
	mb.Add(sw)

	l, err := NewLabel("Vault name", font(titleFontSize), defaultMargin*2, defaultMargin, defaultMargin, defaultMargin)
	if err != nil {
		return nil, err
	}
	l.SetHAlign(gtk.ALIGN_START)
	vb.Add(l)

	l, err = NewLabel("Username:", font(defaultFontSize), defaultMargin, defaultMargin, 0, defaultMargin)
	if err != nil {
		return nil, err
	}
	l.SetHAlign(gtk.ALIGN_START)
	//l.SetJustify(gtk.JUSTIFY_LEFT)
	//l.SetHExpand(true)
	vb.Add(l)

	hb, err = NewHBox()
	if err != nil {
		return nil, err
	}
	l, err = NewLabel("hgouchet@gmail.com", font(defaultFontSize+1), defaultMargin/2, defaultMargin, defaultMargin, defaultMargin)
	if err != nil {
		return nil, err
	}
	l.SetHAlign(gtk.ALIGN_START)
	l.SetSelectable(true)
	hb.PackStart(l, true, true, 0)

	b, err = NewButton("Copy", func() { fmt.Println("copy") }, 0, defaultMargin)
	if err != nil {
		return nil, err
	}
	//b.SetHAlign(gtk.ALIGN_START)
	hb.PackEnd(b, false, true, 0)
	vb.Add(hb)

	l, err = NewLabel("Password:", font(defaultFontSize), defaultMargin, defaultMargin, 0, defaultMargin)
	if err != nil {
		return nil, err
	}
	l.SetHAlign(gtk.ALIGN_START)
	vb.Add(l)

	hb, err = NewHBox()
	if err != nil {
		return nil, err
	}
	l, err = NewLabel("********", font(defaultFontSize+1), defaultMargin/2, defaultMargin, defaultMargin, defaultMargin)
	if err != nil {
		return nil, err
	}
	l.SetHAlign(gtk.ALIGN_START)
	hb.PackStart(l, true, true, 0)

	b, err = NewButton("Copy", func() { fmt.Println("copy") }, 0, defaultMargin)
	if err != nil {
		return nil, err
	}
	//b.SetHAlign(gtk.ALIGN_START)
	hb.PackEnd(b, false, true, 0)
	vb.Add(hb)

	l, err = NewLabel("Note:", font(defaultFontSize), defaultMargin, defaultMargin, 0, defaultMargin)
	if err != nil {
		return nil, err
	}
	l.SetHAlign(gtk.ALIGN_START)
	vb.Add(l)

	l, err = NewLabel(`a
b
c
d
e
f
g
h
i
j
-`, font(defaultFontSize+1), defaultMargin/2, defaultMargin, defaultMargin, defaultMargin)
	if err != nil {
		return nil, err
	}
	l.SetHAlign(gtk.ALIGN_START)
	l.SetSelectable(true)
	vb.Add(l)

	lk, err := gtk.LinkButtonNewWithLabel("https://www.google.fr", "www.google.fr")
	if err != nil {
		return nil, err
	}
	lk.SetHAlign(gtk.ALIGN_START)
	vb.Add(lk)

	lb, err := NewLevelBar(2, 0, 4)
	if err != nil {
		return nil, err
	}
	vb.Add(lb)

	l, err = NewLabel("Last updated: 3 days ago", font(defaultFontSize), defaultMargin, defaultMargin, 0, defaultMargin)
	if err != nil {
		return nil, err
	}
	l.SetHAlign(gtk.ALIGN_START)
	//l.SetJustify(gtk.JUSTIFY_LEFT)
	//l.SetHExpand(true)
	vb.Add(l)

	l, err = NewLabel("Creation: 5 days ago", font(defaultFontSize), defaultMargin, defaultMargin, 0, defaultMargin)
	if err != nil {
		return nil, err
	}
	l.SetHAlign(gtk.ALIGN_START)
	//l.SetJustify(gtk.JUSTIFY_LEFT)
	//l.SetHExpand(true)
	vb.Add(l)

	return mb, nil
}

func newNoteSideBar(list DataTable, add FuncOne, show FuncTwo) (gtk.IWidget, error) {
	// Main content
	tv, err := NewTreeView(list, show)
	if err != nil {
		return nil, err
	}
	// Header (search bar + add vault)
	hb, err := NewHBox()
	if err != nil {
		return nil, err
	}
	// Search bar
	s, err := NewSearchEntry("Search ...", tv.Search, defaultMargin, defaultMargin, 0, defaultMargin)
	if err != nil {
		return nil, err
	}
	hb.Add(s)
	// Adds a button to create a new vault.
	b, err := NewButton("New Vault", func() {
		add(list.Title())
	}, defaultMargin, defaultMargin, 0, defaultMargin)
	if err != nil {
		return nil, err
	}
	hb.Add(b)

	// Side bar (header, vaults list)
	vb, err := NewVBox()
	if err != nil {
		return nil, err
	}
	vb.Add(hb)

	// Lists all tag's vaults.
	sp, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	vb.Add(sp)

	sw, err := tv.ScrollTreeView(true, true)
	if err != nil {
		return nil, err
	}
	vb.PackStart(sw, true, true, 0)

	return vb, nil
}

// Title ...
func (n *Note) Title() string {
	return n.title
}

// Widget ...
func (n *Note) Widget() *gtk.Paned {
	return n.widget
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
	c.SetScrollable(true)
	c.SetShowBorder(true)

	// Title
	l, err := NewLabel(name, font(titleFontSize), 0, defaultMargin, 0, defaultMargin)
	if err != nil {
		return nil, err
	}
	l.Show()
	c.SetActionWidget(l, gtk.PACK_START)

	// New Page
	b, err := NewButton("Add", clicked, defaultMargin, defaultMargin, defaultMargin, defaultMargin)
	if err != nil {
		return nil, err
	}
	b.Show()
	c.SetActionWidget(b, gtk.PACK_END)

	// Adds pages
	book := &Notebook{c: c}
	book.init(pages)
	for i, note := range pages {
		if err = book.add(note, i); err != nil {
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

func (n *Notebook) add(p *Note, pos int) error {
	// Creates the tab.
	l, err := NewLabel(p.Title(), font(defaultFontSize))
	if err != nil {
		return err
	}
	i := n.c.InsertPage(p.Widget(), l, pos)
	n.c.SetCurrentPage(i)
	n.c.ShowAll()

	return nil
}

// AddPage ...
func (n *Notebook) AddPage(p *Note) error {
	// Adds a page to the book.
	if err := n.add(p, n.numPage(p.Title())); err != nil {
		return err
	}
	// Updates the summary.
	n.index = append(n.index, p.Title())
	sort.Strings(n.index)

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
