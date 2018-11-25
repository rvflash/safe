// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

const (
	tagCancel = "tag_cancel"
	tagError  = "tag_error"
	tagName   = "tag_name"
	tagSubmit = "tag_submit"
)

// TagDialog ...
type TagDialog struct {
	*Dialog
}

// Init ...
func (d *TagDialog) Init() (err error) {
	// Cancellation.
	err = d.ButtonClicked(tagCancel, func() {
		d.Log("cancelled")
		d.Hide()
	})
	if err != nil {
		return
	}

	// Creation.
	return d.ButtonClicked(tagSubmit, func() {
		var err error
		defer func() {
			if err != nil {
				d.Log("err=%q", err.Error())
				d.Error(tagError, err.Error())
			} else {
				d.Log("created")
				d.Hide()
			}
		}()
		s, err := d.ReadEntry(tagName)
		if err != nil {
			return
		}
		_, err = d.App().CreateTag(s)
		if err != nil {
			return
		}
		if d.Parent().Window().GetChildren() == nil {
			// New be!
			if err = d.Parent().Build(); err != nil {
				return
			}
			d.Parent().Show()
		}
		// todo add new tag
	})
}

// Reset ...
func (d *TagDialog) Reset() (err error) {
	defer func() {
		if err != nil {
			d.Log("reset with err=%q", err.Error())
		} else {
			d.Log("reset")
		}
	}()
	// Input field
	if err = d.WriteEntry(tagName, ""); err != nil {
		return
	}
	// Error message
	return d.Error(tagError, "")
}

// Log implements the Plugin interface.
func (d *TagDialog) Log(format string, args ...interface{}) {
	d.Dialog.Log("dialog: tag: "+format, args...)
}
