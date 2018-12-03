// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import (
	"fmt"
	"sort"

	"github.com/gotk3/gotk3/gtk"
)

// Notebook ...
type Notebook struct {
	c     *gtk.Notebook
	index []string
}

// NewNotebook ...
func NewNotebook(name string, clicked Func, pages map[string]*tagNote) (*Notebook, error) {
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
	for i, name := range book.index {
		if err = book.add(pages[name], i); err != nil {
			return nil, err
		}
	}
	c.SetCurrentPage(0)

	return book, nil
}

func (n *Notebook) init(pages map[string]*tagNote) {
	for name := range pages {
		n.index = append(n.index, name)
	}
	sort.Strings(n.index)
}

func (n *Notebook) add(p *tagNote, pos int) error {
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
func (n *Notebook) AddPage(p *tagNote) error {
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
