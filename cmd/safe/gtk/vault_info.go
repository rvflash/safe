// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import (
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/rvflash/elapsed"
	"github.com/rvflash/safe"
)

type vaultInfo struct {
	d        *safe.Vault
	upd, del FuncTwo
	copy     FuncOne
}

func newVaultInfo(v *safe.Vault, upd, del FuncTwo, copy FuncOne) (*gtk.Box, error) {
	if v == nil {
		return nil, safe.ErrMissing
	}
	// Vault properties
	vi := &vaultInfo{d: v, upd: upd, del: del, copy: copy}
	p, err := vi.properties()
	if err != nil {
		return nil, err
	}
	sw, err := NewScrolledWindow(false, true)
	if err != nil {
		return nil, err
	}
	sw.Add(p)

	// Menu actions
	vb, err := NewVBox()
	if err != nil {
		return nil, err
	}
	m, err := vi.menu()
	if err != nil {
		return nil, err
	}
	vb.SetMarginTop(defaultMargin / 2)
	vb.PackStart(m, false, true, 0)
	vb.Add(sw)

	return vb, nil
}

func (v *vaultInfo) menu() (*gtk.InfoBar, error) {
	hb, err := NewHBox()
	if err != nil {
		return nil, err
	}
	var (
		hMargin = defaultMargin / 2
		tag     = v.d.Tag().Name()
		vault   = v.d.Name()
	)
	b, err := NewButton("Edit", func() { v.upd(tag, vault) }, hMargin, 0, hMargin)
	if err != nil {
		return nil, err
	}
	hb.Add(b)

	b, err = NewButton("Delete", func() { v.del(tag, vault) }, hMargin, defaultMargin, hMargin)
	if err != nil {
		return nil, err
	}
	hb.Add(b)

	h, err := gtk.InfoBarNew()
	if err != nil {
		return nil, err
	}
	h.SetSizeRequest(infoWidth, infoBarHeight)
	h.SetMessageType(gtk.MESSAGE_INFO)
	h.PackEnd(hb, false, false, 0)

	return h, nil
}

const (
	// Sidebar sizes
	infoWidth     = 320
	infoHeight    = 410
	infoBarHeight = 42

	// Icons
	// see https://commons.wikimedia.org/wiki/GNOME_Desktop_icons
	secLowIco = "security-low"
	secMedIco = "security-medium"
	//secHighIco = "security-high"
)

func (v *vaultInfo) properties() (*gtk.Box, error) {
	vb, err := NewVBox()
	if err != nil {
		return nil, err
	}
	vb.SetMarginBottom(defaultMargin)
	vb.SetSizeRequest(infoWidth, infoHeight)

	var dMargin = defaultMargin * 2

	// Title with security icon
	hb, err := NewHBox()
	if err != nil {
		return nil, err
	}
	// > Icon
	img, err := v.securityIcon()
	if err != nil {
		return nil, err
	}
	img.SetMarginEnd(defaultMargin)
	hb.PackEnd(img, false, false, 0)

	// > Title
	l, err := NewLabel(v.d.Name(), font(titleFontSize), dMargin, defaultMargin, defaultMargin, defaultMargin)
	if err != nil {
		return nil, err
	}
	l.SetHAlign(gtk.ALIGN_START)
	hb.Add(l)

	vb.Add(hb)

	// Username
	usr := v.d.Login().Name
	if err = v.packField(vb, "Username:", usr, func() { v.copy(usr) }); err != nil {
		return nil, err
	}

	// Password
	pwd := v.d.Login().Password
	if err = v.packField(vb, "Password:", "********", func() { v.copy(pwd) }); err != nil {
		return nil, err
	}

	// tagNote
	if note := v.d.Login().Note; note != "" {
		if err = v.packField(vb, "tagNote:", note, nil); err != nil {
			return nil, err
		}
	}

	// URL
	if u := v.d.Login().URL; u != nil {
		lk, err := gtk.LinkButtonNewWithLabel(u.String(), u.Host)
		if err != nil {
			return nil, err
		}
		lk.SetHAlign(gtk.ALIGN_START)
		vb.Add(lk)
	}

	// Strength
	const (
		minStrength = 0
		maxStrength = 4
	)
	lb, err := NewLevelBar(v.d.Login().Strength(), minStrength, maxStrength)
	if err != nil {
		return nil, err
	}
	vb.Add(lb)

	// Last modification
	if err = v.packDate(vb, "Last updated:", v.d.LastUpdate()); err != nil {
		return nil, err
	}

	// Creation
	if err = v.packDate(vb, "Creation:", v.d.AddDate()); err != nil {
		return nil, err
	}
	return vb, nil
}

func (v *vaultInfo) packField(vb *gtk.Box, name, value string, copy Func) error {
	// Name
	l, err := NewLabel(name, font(defaultFontSize), defaultMargin, defaultMargin, 0, defaultMargin)
	if err != nil {
		return err
	}
	l.SetHAlign(gtk.ALIGN_START)
	vb.Add(l)

	// Value
	l, err = NewLabel(value, font(defaultFontSize+1), defaultMargin/2, defaultMargin, defaultMargin, defaultMargin)
	if err != nil {
		return err
	}
	l.SetHAlign(gtk.ALIGN_START)
	l.SetSelectable(true)

	// Groups both and adds copy action.
	if copy != nil {
		hb, err := NewHBox()
		if err != nil {
			return err
		}
		hb.PackStart(l, true, true, 0)

		// Easy copy to clipboard.
		b, err := NewButton("Copy", copy, 0, defaultMargin)
		if err != nil {
			return err
		}
		hb.PackEnd(b, false, true, 0)
		vb.Add(hb)
	} else {
		// Basic mode.
		vb.Add(l)
	}
	return nil
}

func (v *vaultInfo) packDate(vb *gtk.Box, lb string, t time.Time) error {
	l, err := NewLabel(lb+" "+elapsed.Time(t), font(defaultFontSize), defaultMargin, defaultMargin, 0, defaultMargin)
	if err != nil {
		return err
	}
	l.SetHAlign(gtk.ALIGN_START)
	vb.Add(l)

	return nil
}

func (v *vaultInfo) securityIcon() (*gtk.Image, error) {
	name := secLowIco
	if v.d.Login().Strength() > 2 {
		name = secMedIco
	}
	return gtk.ImageNewFromIconName(name, gtk.ICON_SIZE_DND)
}
