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
			return
		}
		err = d.Parent().DeleteVault(d.vault, d.tag)
	})
}

// Reload ...
func (d *ConfirmDialog) Reload(tag string, vault ...string) error {
	if len(vault) != 1 {
		return d.Reload("", "")
	}
	d.tag, d.vault = tag, vault[0]
	d.Log("reloaded (tag=%q, vault=%q)", tag, vault[0])
	return nil
}

// Reset ...
func (d *ConfirmDialog) Reset() error {
	d.Log("reset")
	return d.Reload("", "")
}

// Log implements the Plugin interface.
func (d *ConfirmDialog) Log(format string, args ...interface{}) {
	d.Dialog.Log("dialog: confirm: "+format, args...)
}
