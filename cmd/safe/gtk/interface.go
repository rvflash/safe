// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import (
	"errors"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rvflash/safe/app"
)

var (
	// ErrUnkDialog ...
	ErrUnkDialog = errors.New("unknown dialog")
	// ErrFileSystem ...
	ErrFileSystem = errors.New("invalid file system")
	// ErrUndColumn ...
	ErrUndColumn = errors.New("missing columns")
	// ErrUndObject ...
	ErrUndObject = errors.New("missing object")
	// ErrContainer ...
	ErrContainer = errors.New("invalid window")
)

// Func ...
type Func func()

// FuncOne ...
type FuncOne func(string)

// FuncTwo ....
type FuncTwo func(string, string)

// Logger must be implemented by any logger.
type Logger interface {
	// Printf logs a message at level Info on the standard logger.
	Printf(format string, args ...interface{})
}

// Container ...
type Container interface {
	ID(name string) (glib.IObject, error)
	Init() error
	Reset() error
}

// Plugin ...
type Plugin interface {
	App() *app.Safe
	Log(format string, args ...interface{})
}

// WidgetContainer
type WidgetContainer interface {
	Container
	Plugin
	IWidget() gtk.IWidget
}

/// Visibility ...
type Visibility interface {
	Hide()
	Show()
}

// VisibleContainer ...
type VisibleWidgetContainer interface {
	WidgetContainer
	Visibility
}

// LoadVisibleWidgetContainer ...
type LoadVisibleWidgetContainer interface {
	VisibleWidgetContainer
	Reload(tag string, vault ...string) error
}

// Listener ...
// todo
type Listener interface {
	//ButtonClicked(id string, fn Func) error
}

// Builder ...
// todo
type Builder interface {
}

// DataTable ...
type DataTable interface {
	Cols() []string
	ColSizes() []int
	Rows(upd, del FuncTwo, cp FuncOne) ([][]interface{}, error)
	Types() []glib.Type
	Title() string
}
