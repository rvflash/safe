// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/rvflash/safe"
)

type tagNote struct {
	name   string
	tree   *TreeView
	widget *gtk.Paned
}

func newTagNote(list DataTable, add FuncOne, show FuncTwo) (*tagNote, error) {
	// Right
	w, tv, err := vaultList(list, add, show)
	if err != nil {
		return nil, err
	}
	p, err := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	p.Pack1(w, true, false)

	return &tagNote{widget: p, tree: tv, name: list.Title()}, nil
}

func vaultList(list DataTable, add FuncOne, show FuncTwo) (*gtk.Box, *TreeView, error) {
	// Vaults list
	tv, err := NewTreeView(list, show)
	if err != nil {
		return nil, nil, err
	}

	// Search bar
	hb, err := vaultSearch(list, add, tv.Search)
	if err != nil {
		return nil, nil, err
	}

	// Lists all tag'ls vaults.
	vb, err := NewVBox()
	if err != nil {
		return nil, nil, err
	}
	vb.Add(hb)

	sp, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, nil, err
	}
	vb.Add(sp)

	sw, err := tv.ScrollTreeView(true, true)
	if err != nil {
		return nil, nil, err
	}
	vb.PackStart(sw, true, true, 0)

	return vb, tv, nil
}

func vaultSearch(list DataTable, add, search FuncOne) (*gtk.Box, error) {
	// Header (search bar + add vault)
	hb, err := NewHBox()
	if err != nil {
		return nil, err
	}

	// Search bar
	s, err := NewSearchEntry("Search ...", search, defaultMargin, defaultMargin, 0, defaultMargin)
	if err != nil {
		return nil, err
	}
	hb.Add(s)

	// Adds a button to create a new vault.
	b, err := NewButton("New Vault", func() { add(list.Title()) }, defaultMargin, defaultMargin, 0, defaultMargin)
	if err != nil {
		return nil, err
	}
	hb.Add(b)

	return hb, nil
}

// Add ...
func (n tagNote) Add(v *safe.Vault) error {
	if v == nil || n.tree == nil {
		return safe.ErrMissing
	}
	return n.tree.AddRow(row(v))
}

// Delete ...
func (n tagNote) Delete(name string) error {
	if n.tree == nil {
		return safe.ErrMissing
	}
	if err := n.tree.DelRow(name); err != nil {
		return err
	}
	if n.tree.Len() == 0 {
		// No more data to display on the right sidebar.
		if i, err := n.widget.GetChild2(); err == nil {
			n.widget.Remove(i)
		}
	}
	return nil
}

// Title ...
func (n *tagNote) Title() string {
	return n.name
}

// Update ...
func (n tagNote) Update(v *safe.Vault, upd, del FuncTwo, cp FuncOne) error {
	if v == nil || n.tree == nil {
		return safe.ErrMissing
	}
	if err := n.tree.UpdRow(v.Name(), row(v)); err != nil {
		return err
	}
	return n.View(v, upd, del, cp)
}

// View ...
func (n tagNote) View(v *safe.Vault, upd, del FuncTwo, cp FuncOne) error {
	vi, err := newVaultInfo(v, upd, del, cp)
	if err != nil {
		return err
	}
	if i, err := n.widget.GetChild2(); err == nil {
		// Removes the previous one.
		n.widget.Remove(i)
	}
	n.widget.Pack2(vi, true, false)
	n.widget.ShowAll()

	return nil
}

// Widget ...
func (n *tagNote) Widget() *gtk.Paned {
	return n.widget
}
