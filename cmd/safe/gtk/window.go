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
	w.SetSizeRequest(width, height)
	w.SetIconName("dialog-password")

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

type notes map[string]*tagNote

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
	tabs    notes
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
	w.tabs, err = w.toNotes(t)
	if err != nil {
		return err
	}
	w.book, err = NewNotebook("Safe", w.ShowTagDialog, w.tabs)
	if err != nil {
		return err
	}

	// Logout button
	b, err := NewButton("Logout", gtk.MainQuit, 0, defaultMargin, defaultMargin, defaultMargin)
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

	return nil
}

func (w *Window) toNotes(tags []*safe.Tag) (notes, error) {
	notes := make(map[string]*tagNote, len(tags))
	for _, t := range tags {
		d, err := w.App().ListVaultByNames(t.Name(), "")
		if err != nil {
			return nil, err
		}
		notes[t.Name()], err = newTagNote(newVaultTable(t.Name(), d), w.ShowNewVaultDialog, w.ShowVaultInfo)
		if err != nil {
			return nil, err
		}
		w.Log("win: tag: %s initialized with %d vault(s)", t.Name(), len(d))
	}
	return notes, nil
}

// ShowVaultInfo ...
func (w *Window) ShowVaultInfo(tag, vault string) {
	n, ok := w.tabs[tag]
	if !ok {
		w.Log("win: vault: unknown tag named %q", tag)
		return
	}
	v, err := w.App().Vault(vault, tag)
	if err != nil {
		w.Log("win: vault: %q is unknown in %q", vault, tag)
		return
	}
	if err = n.View(v, w.ShowUpdVaultDialog, w.ShowDelVaultConfirm, w.cb.Copy); err != nil {
		w.Log("win: vault: fails to show %q in %q", vault, tag)
	}
	w.Log("win: vault: %q in %q displayed", vault, tag)
}

// AddTag ...
func (w *Window) AddTag(name string) error {
	n, err := newTagNote(newVaultTable(name, nil), w.ShowNewVaultDialog, w.ShowVaultInfo)
	if err != nil {
		return err
	}
	if err = w.book.AddPage(n); err != nil {
		return err
	}
	w.tabs[name] = n

	return nil
}

// DeleteVault ...
func (w *Window) DeleteVault(tag, vault string) error {
	n, ok := w.tabs[tag]
	if !ok {
		return safe.ErrMissing
	}
	return n.Delete(vault)
}

// AddVault ...
func (w *Window) UpsertVault(v *safe.Vault, add bool) error {
	if v == nil {
		return safe.ErrMissing
	}
	n, ok := w.tabs[v.Tag().Name()]
	if !ok {
		return safe.ErrMissing
	}
	if add {
		return n.Add(v)
	}
	return n.Update(v, w.ShowUpdVaultDialog, w.ShowDelVaultConfirm, w.cb.Copy)
}

// Close ...
func (w *Window) Close() error {
	w.App().LogOut()
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
	w.reloadDialog(w.confirm, tag, vault)
}

// ShowVaultDialog ...
func (w *Window) ShowNewVaultDialog(tag string) {
	w.reloadDialog(w.vault, tag)
}

// ShowTagDialog ...
func (w *Window) ShowTagDialog() {
	w.showDialog(w.tag)
}

// ShowVaultDialog ...
func (w *Window) ShowUpdVaultDialog(tag, vault string) {
	w.reloadDialog(w.vault, tag, vault)
}

func (w *Window) reloadDialog(d LoadVisibleWidgetContainer, args ...string) {
	switch len(args) {
	case 1:
		// tag name
		d.Reload(args[0])
	case 2:
		//tag name + vault name
		d.Reload(args[0], args[1])
	}
	w.showDialog(d)
}

func (w *Window) showDialog(d VisibleWidgetContainer) {
	if err := d.Reset(); err != nil {
		d.Log("fails to display, err=%q", err.Error())
		return
	}
	d.Show()
	d.Log("displayed")
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
