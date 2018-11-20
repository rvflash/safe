// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import (
	"errors"

	"github.com/gotk3/gotk3/gtk"
	"github.com/rvflash/safe"
)

const (
	signCancel            = "sign_cancel"
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
	err = d.ButtonClicked(signCancel, func() {
		gtk.MainQuit()
	})
	if err != nil {
		return
	}
	// Registration.
	return d.ButtonClicked(signSubmit, func() {
		var err error
		defer func() {
			if err != nil {
				d.Error(signError, err.Error())
			} else {
				d.Hide()
				d.Parent().ShowAll()
			}
		}()
		p, err := d.ReadEntry(signPassphrase)
		if err != nil {
			return
		}
		if d.App().Logged() == safe.ErrNotFound {
			// Sign up behavior.
			pc, err := d.ReadEntry(signConfirmPassphrase)
			if err != nil {
				return
			}
			if p != pc {
				err = errors.New("the confirm passphrase does not match the expected one")
				return
			}
		}
		err = d.App().Login(p)
	})
}

// Reset ...
func (d *SignDialog) Reset() (err error) {
	// Input fields
	if err = d.WriteEntry(signPassphrase, ""); err != nil {
		return nil
	}
	if err = d.WriteEntry(signConfirmPassphrase, ""); err != nil {
		return nil
	}

	// Sign in or sign up layout?
	o, err := d.ID(signConfirmBox)
	if err != nil {
		return
	}
	cb := o.(*gtk.Box)
	if d.App().Logged() == safe.ErrNotFound {
		cb.Show()
	} else {
		cb.Hide()
	}

	// Error message
	return d.Error(signError, "")
}
