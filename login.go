// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package safe

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/rvflash/safe/crypto"
)

type hashLogin struct {
	hash crypto.Hash
	*Login
}

// MarshalJSON implements tje json.Marshaler interface.
func (h *hashLogin) MarshalJSON() ([]byte, error) {
	// Changes the last modification date.
	h.Login.LastUpdate = time.Now()
	// To JSON string.
	plain, err := json.Marshal(h.Login)
	if err != nil {
		return nil, err
	}
	cipher, err := h.hash.Sign(plain)
	if err != nil {
		return nil, err
	}
	return json.Marshal(cipher)
}

// UnmarshalJSON implements tje json.Unmarshaler interface.
func (h *hashLogin) UnmarshalJSON(s []byte) error {
	// Given a JSON string
	b, err := b64Decode(bytes.Trim(s, `"`))
	if err != nil {
		return err
	}
	plain, err := h.hash.Decrypt(b)
	if err != nil {
		return err
	}
	return json.Unmarshal(plain, &h.Login)
}

func b64Decode(v []byte) ([]byte, error) {
	b := make([]byte, base64.StdEncoding.DecodedLen(len(v)))
	n, err := base64.StdEncoding.Decode(b, v)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}

// Login represents the couple of username / password and any other information to sign in.
type Login struct {
	LastUpdate time.Time `json:"since"`
	Name       string    `json:"name"`
	Note       string    `json:"note,omitempty"`
	Password   string    `json:"pass"`
	URL        *url.URL  `json:"url,omitempty"`
}

// NewLogin returns a new instance of Login and set the last update date.
func NewLogin(name, pass string) *Login {
	return &Login{
		Name:       strings.TrimSpace(name),
		Password:   strings.TrimSpace(pass),
		LastUpdate: time.Now(),
	}
}

// Valid returns in success if the the Login has all mandatory data to be store.
func (l *Login) Valid() bool {
	return strings.TrimSpace(l.Name) != "" && strings.TrimSpace(l.Password) != ""
}

// Safe indicates if the Login seems safe or not.
func (l *Login) Safe() (ok bool, err error) {
	if ok = l.Valid(); !ok {
		err = ErrInvalid
		return
	}
	if ok = len(l.Password) >= MinSize; !ok {
		err = ErrTooShort
		return
	}
	if ok = l.LastUpdate.After(time.Now().Add(-MaxDuration)); !ok {
		err = ErrOutdated
		return
	}
	return
}
