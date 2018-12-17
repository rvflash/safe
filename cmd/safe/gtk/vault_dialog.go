// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/rvflash/safe"
	"github.com/rvflash/safe/app"
)

// Password default options
const (
	pwdLength      = 32
	pwdDigits      = 10
	pwdSymbols     = 10
	pwdNoUppercase = false
	pwdAllowRepeat = true
)

// User interface
const (
	vaultCancel   = "vault_cancel"
	vaultError    = "vault_error"
	vaultName     = "vault_name"
	vaultSubmit   = "vault_submit"
	vaultUsername = "vault_username"
	vaultPwd      = "vault_password"
	vaultURL      = "vault_url"
	vaultNotes    = "vault_note"

	// Kinds of password'ls vault
	vaultSelfGenerated = "vault_self_generated"
	vaultCustomized    = "vault_customized"
	vaultCustomizedBox = "vault_password_customized"
	vaultHandwritten   = "vault_handwritten"

	// Customization options.
	vaultPwdLength      = "vault_password_length"
	vaultPwdDigits      = "vault_password_digits"
	vaultPwdSymbols     = "vault_password_symbols"
	vaultPwdUppercase   = "vault_password_uppercase"
	vaultPwdAllowRepeat = "vault_password_allow_repeat"
)

// VaultDialog ...
type VaultDialog struct {
	*Dialog
	v   *safe.Vault
	tag string
}

func (d *VaultDialog) Init() (err error) {
	// Cancellation.
	if err = d.Closed(d.Hide); err != nil {
		return
	}
	if err = d.ButtonClicked(vaultCancel, d.Hide); err != nil {
		return
	}
	
	// Options
	err = d.ButtonClicked(vaultSelfGenerated, func() {
		_ = d.showOptions("")
	})
	if err != nil {
		return
	}
	err = d.ButtonClicked(vaultCustomized, func() {
		_ = d.showOptions(vaultCustomizedBox)
	})
	if err != nil {
		return
	}
	err = d.ButtonClicked(vaultHandwritten, func() {
		_ = d.showOptions(vaultPwd)
	})
	if err != nil {
		return
	}

	// Submit
	return d.ButtonClicked(vaultSubmit, func() {
		var (
			n   string
			l   app.Login
			err error
			ok  bool
		)
		defer func() {
			if err != nil {
				d.Log("err=%q", err.Error())
				if err = d.Error(vaultError, err.Error()); err != nil {
					d.Log("fails to display the error: %s", err)
				}
			} else {
				d.Log("submitted")
				d.Hide()
			}
		}()

		n, err = d.ReadEntry(vaultName)
		if err != nil {
			return
		}
		l.Username, err = d.ReadEntry(vaultUsername)
		if err != nil {
			return
		}
		if ok, _ = d.IsActivated(vaultCustomized); ok {
			l.Password, err = d.customPassword()
			d.Log("generating customized password")
		} else if ok, _ = d.IsActivated(vaultHandwritten); ok {
			l.Password, err = d.ReadEntry(vaultPwd)
			d.Log("handwritten password")
		} else {
			l.Password, err = app.GeneratePassword(pwdLength, pwdDigits, pwdSymbols, pwdNoUppercase, pwdAllowRepeat)
			d.Log("self generated password")
		}
		if err != nil {
			return
		}
		l.URL, err = d.ReadEntry(vaultURL)
		if err != nil {
			return
		}
		l.Note, err = d.ReadEntry(vaultNotes)
		if err != nil {
			return
		}
		var (
			v   *safe.Vault
			add = d.v == nil
		)
		if add {
			d.Log("try to create a vault named %q in %q", n, d.tag)
			v, err = d.App().CreateVault(n, d.tag, l)
		} else {
			d.Log("try to update a vault named %q in %q", n, d.v.Tag().Name())
			v, err = d.App().UpdateVault(n, d.v.Tag().Name(), l)
		}
		if err != nil {
			return
		}
		err = d.Parent().UpsertVault(v, add)
	})
}

func (d *VaultDialog) customPassword() (string, error) {
	length, err := d.ReadSpinButton(vaultPwdLength)
	if err != nil {
		return "", err
	}
	digits, err := d.ReadSpinButton(vaultPwdDigits)
	if err != nil {
		return "", err
	}
	symbols, err := d.ReadSpinButton(vaultPwdSymbols)
	if err != nil {
		return "", err
	}
	uppercase, err := d.IsActivated(vaultPwdUppercase)
	if err != nil {
		return "", err
	}
	repeat, err := d.IsActivated(vaultPwdAllowRepeat)
	if err != nil {
		return "", err
	}
	return app.GeneratePassword(length, digits, symbols, uppercase, repeat)
}

func (d *VaultDialog) showOptions(kind string) (err error) {
	// Custom area
	o, err := d.ID(vaultCustomizedBox)
	if err != nil {
		return
	}
	if kind == vaultCustomizedBox {
		o.(*gtk.Box).Show()
	} else {
		o.(*gtk.Box).Hide()
	}

	// Password field
	o, err = d.ID(vaultPwd)
	if err != nil {
		return
	}
	if kind == vaultPwd {
		o.(*gtk.Entry).Show()
	} else {
		o.(*gtk.Entry).Hide()
	}
	return
}

// Reload ...
func (d *VaultDialog) Reload(tag string, vault ...string) {
	if len(vault) != 1 {
		// Creation
		d.Log("creation in tag=%q", tag)
		d.tag, d.v = tag, nil
	} else {
		// Update
		d.Log("updating %q in tag=%q", vault[0], tag)

		var err error
		d.v, err = d.App().Vault(vault[0], tag)
		if err != nil {
			d.Log("fails with %q", err.Error())
		}
	}
}

// Reset ...
func (d *VaultDialog) Reset() (err error) {
	defer func() {
		if err != nil {
			d.Log("reset with err=%q", err.Error())
		} else {
			d.Log("reset")
		}
	}()
	editable := d.v == nil
	if editable {
		err = d.Focus(vaultName)
	} else {
		err = d.Focus(vaultUsername)
	}
	if err != nil {
		return
	}
	// Vault name
	err = d.EditEntry(vaultName, d.name(), editable)
	if err != nil {
		return
	}
	// Username
	err = d.WriteEntry(vaultUsername, d.username())
	if err != nil {
		return
	}
	// Self-generated password by default
	err = d.showOptions("")
	if err != nil {
		return
	}
	err = d.Activate(vaultSelfGenerated, true)
	if err != nil {
		return
	}
	// Password length
	err = d.WriteSpinButton(vaultPwdLength, pwdLength)
	if err != nil {
		return
	}
	// Password digits
	err = d.WriteSpinButton(vaultPwdDigits, pwdDigits)
	if err != nil {
		return
	}
	// Password symbols
	err = d.WriteSpinButton(vaultPwdSymbols, pwdSymbols)
	if err != nil {
		return
	}
	// Password without uppercase
	err = d.Activate(vaultPwdUppercase, pwdNoUppercase)
	if err != nil {
		return
	}
	// Password allow repeat
	err = d.Activate(vaultPwdAllowRepeat, pwdAllowRepeat)
	if err != nil {
		return
	}
	// Password
	err = d.WriteEntry(vaultPwd, "")
	if err != nil {
		return
	}
	// URL
	err = d.WriteEntry(vaultURL, d.url())
	if err != nil {
		return
	}
	// Note
	err = d.WriteEntry(vaultNotes, d.note())
	if err != nil {
		return
	}
	// Error message
	return d.Error(vaultError, "")
}

func (d *VaultDialog) name() string {
	if d.v == nil {
		return ""
	}
	return d.v.Name()
}

func (d *VaultDialog) username() string {
	if d.v == nil {
		return ""
	}
	return d.v.Login().Name
}

func (d *VaultDialog) url() string {
	if d.v == nil || d.v.Login().URL == nil {
		return ""
	}
	return d.v.Login().URL.String()
}

func (d *VaultDialog) note() string {
	if d.v == nil {
		return ""
	}
	return d.v.Login().Note
}

// Log implements the Plugin interface.
func (d *VaultDialog) Log(format string, args ...interface{}) {
	d.Dialog.Log("dialog: vault: "+format, args...)
}
