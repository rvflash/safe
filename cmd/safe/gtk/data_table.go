// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/rvflash/safe"
)

// DataTable ...
type DataTable interface {
	Cols() []string
	ColSizes() []int
	ColID() int
	Rows() [][]interface{}
	Types() []glib.Type
	Title() string
}

type vaultTable struct {
	parent string
	data   []*safe.Vault
}

func newVaultTable(tag string, list []*safe.Vault) *vaultTable {
	return &vaultTable{parent: tag, data: list}
}

// Cols ...
func (d vaultTable) Cols() []string {
	return []string{"Name", "Username", "Filter"}
}

// ColID ...
// Use this column as identifier (name) of the line.
func (d vaultTable) ColID() int {
	return 0
}

// ColSizes ...
func (d vaultTable) ColSizes() []int {
	return []int{130, 170, 0}
}

// Rows ...
func (d *vaultTable) Rows() [][]interface{} {
	const (
		name = iota
		username
		filter
	)
	rs := make([][]interface{}, len(d.data))
	for p, v := range d.data {
		rs[p] = make([]interface{}, len(d.Cols()))
		rs[p][name] = v.Name()
		rs[p][username] = v.Login().Name
		rs[p][filter] = true
	}
	return rs
}

// Title ...
func (d vaultTable) Title() string {
	return d.parent
}

// Types ...
func (d vaultTable) Types() []glib.Type {
	return []glib.Type{glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_BOOLEAN}
}
