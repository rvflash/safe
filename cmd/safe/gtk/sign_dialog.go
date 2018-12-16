// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import (
	"errors"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rvflash/safe"
	"github.com/rvflash/safe/app"
)

const (
	signConfirmBox        = "sign_confirm_passphrase"
	signConfirmPassphrase = "confirm_passphrase"
	signError             = "sign_error"
	signPassphrase        = "passphrase"
	signSubmit            = "sign_submit"
)

// SignDialog ...
type SignDialog struct {
	*Dialog
}

// Init ...
func (d *SignDialog) Init() (err error) {
	// Cancellation.
	if err = d.Closed(gtk.MainQuit); err != nil {
		return
	}
	// Registration.
	/*
	if err = d.Focus(signPassphrase); err != nil {
		return
	}
	*/
	if err = d.EnterPressed(signPassphrase, d.login); err != nil {
		return
	}
	return d.ButtonClicked(signSubmit, d.login)
}

func (d *SignDialog) login() {
	var (
		err   error
		p, cp string
	)
	defer func() {
		if err != nil {
			d.Log("err=%q", err.Error())
			if err = d.Error(signError, err.Error()); err != nil {
				d.Log("fails to display the error: %s", err)
			}
		} else {
			d.Log("logged")
			d.Hide()
		}
	}()
	if p, err = d.ReadEntry(signPassphrase); err != nil {
		return
	}
	if d.App().Logged() == app.ErrNotFound {
		// Sign up behavior.
		if cp, err = d.ReadEntry(signConfirmPassphrase); err != nil {
			return
		}
		if p != cp {
			err = errors.New("the confirm passphrase does not match the expected one")
			return
		}
	}
	if err = d.App().Login(p); err != nil {
		return
	}

	// What's next?
	var l []*safe.Tag
	l, err = d.App().ListTagByNames()
	if err != nil {
		return
	}
	if len(l) == 0 {
		// New be!
		d.Parent().ShowTagDialog()
		return
	}
	if err = d.Parent().Build(); err != nil {
		return
	}
	d.Parent().Show()
}

// Reset ...
func (d *SignDialog) Reset() (err error) {
	defer func() {
		if err != nil {
			d.Log("reset with err=%q", err.Error())
		} else {
			d.Log("reset")
		}
	}()
	// Input fields
	if err = d.WriteEntry(signPassphrase, ""); err != nil {
		return
	}
	if err = d.WriteEntry(signConfirmPassphrase, ""); err != nil {
		return
	}

	// Sign in or sign up layout?
	o, err := d.ID(signConfirmBox)
	if err != nil {
		return
	}
	if d.App().Logged() == app.ErrNotFound {
		o.(*gtk.Box).Show()
	} else {
		o.(*gtk.Box).Hide()
	}
	return d.Error(signError, "")
}

// Log implements the Plugin interface.
func (d *SignDialog) Log(format string, args ...interface{}) {
	d.Dialog.Log("dialog: sign: "+format, args...)
}
