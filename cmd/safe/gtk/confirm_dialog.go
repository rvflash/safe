// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

const (
	confirmCancel = "confirm_cancel"
	confirmSubmit = "confirm_ok"
)

// ConfirmDialog ...
type ConfirmDialog struct {
	*Dialog
	vault, tag string
}

func (d *ConfirmDialog) Init() (err error) {
	// Cancellation.
	err = d.ButtonClicked(confirmCancel, func() {
		d.Log("cancelled")
		d.Hide()
	})
	if err != nil {
		return
	}

	// Submit
	return d.ButtonClicked(confirmSubmit, func() {
		var err error
		defer func() {
			if err != nil {
				d.Log("err=%q", err.Error())
			}
			d.Log("confirmed")
			d.Hide()
		}()
		if err = d.App().DeleteVault(d.vault, d.tag); err != nil {
			d.Log("Rv: t:%q v:%q", d.vault, d.tag)
			return
		}
		err = d.Parent().DeleteVault(d.tag, d.vault)
	})
}

// Reload ...
func (d *ConfirmDialog) Reload(tag string, vault ...string) {
	if len(vault) != 1 {
		d.tag, d.vault = "", ""
	} else {
		d.tag, d.vault = tag, vault[0]
	}
	d.Log("reloaded (tag=%q, vault=%q)", d.tag, d.vault)
}

// Reset ...
func (d *ConfirmDialog) Reset() error {
	d.Log("reset")
	d.tag, d.vault = "", ""
	return nil
}

// Log implements the Plugin interface.
func (d *ConfirmDialog) Log(format string, args ...interface{}) {
	d.Dialog.Log("dialog: confirm: "+format, args...)
}
