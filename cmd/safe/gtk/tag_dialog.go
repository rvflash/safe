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
				d.Error(tagError, err.Error())
			} else {
				d.Hide()
			}
		}()
		t, err := d.ReadEntry(tagName)
		if err != nil {
			return
		}
		// todo add the new tag
		_, err = d.App().CreateTag(t)
	})
}

// Reset ...
func (d *TagDialog) Reset() (err error) {
	// Input field
	if err = d.WriteEntry(tagName, ""); err != nil {
		return
	}
	// Error message
	return d.Error(tagError, "")
}
