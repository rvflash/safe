// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.
package gtk

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/rvflash/safe"
	"github.com/rvflash/safe/app"
)

// NewWindow ...
func NewWindow(app *app.Safe, log Logger, title string, width, height int) (*Window, error) {
	w, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return nil, err
	}
	w.SetTitle(title)
	w.SetResizable(true)
	w.SetPosition(gtk.WIN_POS_CENTER)
	//w.SetDefaultSize(width, height)
	//w.SetFocusVisible(true)
	w.SetSizeRequest(width, height)
	//w.SetTypeHint(gdk.WINDOW_TYPE_HINT_NORMAL)

	// Main windows, destroys on closing.
	if _, err = w.Connect("destroy", gtk.MainQuit); err != nil {
		return nil, err
	}
	// Deals with clipboard.
	c, err := NewClipboard()
	if err != nil {
		return nil, err
	}
	return &Window{db: app, log: log, win: w, cb: c}, nil
}

// Window ...
type Window struct {
	book    *Notebook
	cb      *Clipboard
	db      *app.Safe
	log     Logger
	confirm *ConfirmDialog
	sign    *SignDialog
	tag     *TagDialog
	vault   *VaultDialog
	win     *gtk.Window
}

// App ...
func (w *Window) App() *app.Safe {
	return w.db
}

// AttachDialog ...
func (w *Window) AttachDialog(c VisibleWidgetContainer) error {
	if err := c.Init(); err != nil {
		return err
	}
	switch d := c.(type) {
	case *ConfirmDialog:
		w.confirm = d
	case *SignDialog:
		w.sign = d
	case *TagDialog:
		w.tag = d
	case *VaultDialog:
		w.vault = d
	default:
		return ErrUnkDialog
	}
	return nil
}

// Build ...
func (w *Window) Build() error {
	if w.win == nil {
		return ErrContainer
	}
	// Lists of available tags.
	t, err := w.App().ListTagByNames()
	if err != nil {
		return err
	}
	w.Log("win: expected %d tag(s)", len(t))

	n, err := w.toNotes(t)
	if err != nil {
		return err
	}
	w.Log("win: %d notes built", len(n))

	w.book, err = NewNotebook("Safe", w.ShowTagDialog, n...)
	if err != nil {
		return err
	}
	b, err := NewButton("Logout", func() {
		w.App().LogOut()
		gtk.MainQuit()
	})
	if err != nil {
		return err
	}
	hb, err := NewVBox()
	if err != nil {
		return err
	}
	// Adds a notebook to list the tags.
	hb.Add(w.book.Notebook())
	// Adds a button to logout.
	hb.Add(b)
	// Attaches the container to the window.
	w.win.Add(hb)
	w.Log("win: built")

	return nil
}

func (w *Window) toNotes(tags []*safe.Tag) ([]*Note, error) {
	notes := make([]*Note, len(tags))
	for i, t := range tags {
		d, err := w.App().ListVaultByNames(t.Name(), "")
		if err != nil {
			return nil, err
		}
		v := newDataTable(t.Name(), d)
		n, err := NewNote(v, w.ShowNewVaultDialog, w.ShowUpdVaultDialog, w.ShowDelVaultConfirm, w.cb.Copy)
		if err != nil {
			return nil, err
		}
		notes[i] = n
		w.Log("win: tag: %s initialized with %d vault(s)", t.Name(), len(d))
	}
	return notes, nil
}

// Close ...
func (w *Window) Close() error {
	return w.App().Close()
}

// Log ...
func (w *Window) Log(format string, args ...interface{}) {
	if w.log != nil {
		w.log.Printf(format, args...)
	}
}

// Hide implements the Visibility interface.
func (w *Window) Hide() {
	if c := w.Window(); c != nil {
		c.Hide()
		w.Log("win: hidden")
	}
}

// Show implements the Visibility interface.
func (w *Window) Show() {
	if c := w.Window(); c != nil {
		c.ShowAll()
		w.Log("win: displayed")
	}
}

// ShowDelVaultConfirm ...
func (w *Window) ShowDelVaultConfirm(tag, vault string) {
	if w.confirm != nil {
		w.confirm.Reload(tag, vault)
		w.showDialog(w.confirm)
	}
}

// ShowVaultDialog ...
func (w *Window) ShowNewVaultDialog(tag string) {
	if w.vault != nil {
		w.vault.Reload(tag)
		w.showDialog(w.vault)
	}
}

// ShowTagDialog ...
func (w *Window) ShowTagDialog() {
	if w.tag != nil {
		w.showDialog(w.tag)
	}
}

// ShowVaultDialog ...
func (w *Window) ShowUpdVaultDialog(tag, vault string) {
	if w.vault != nil {
		w.vault.Reload(tag, vault)
		w.showDialog(w.vault)
	}
}

func (w *Window) showDialog(d VisibleWidgetContainer) (err error) {
	if err = d.Reset(); err == nil {
		d.Show()
	}
	d.Log("displayed")

	return
}

// ShowAll implements the Visibility interface.
func (w *Window) Run() {
	if c := w.Window(); c != nil {
		c.ShowAll()
		w.showDialog(w.sign)
	}
	w.Log("win: running")
}

// Window ...
func (w *Window) Window() *gtk.Window {
	return w.win
}
