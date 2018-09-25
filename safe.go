// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package safe

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/rvflash/safe/crypto"
)

// Default bounds.
const (
	// MinSize is the minimum number of bytes accepted for a pass.
	MinSize = 16
	// MaxDuration is the duration before to warn to update a data: 90 days.
	MaxDuration = time.Hour * 24 * 90
)

// List of common errors.
var (
	// ErrTooShort is returned is the pass phrase is too weak.
	ErrTooShort = fmt.Errorf("too short, minimum required: %d characters", MinSize)
	// ErrMissing is returned is the mandatory data is missing.
	ErrMissing = errors.New("missing data")
	// ErrInvalid is returned if the data doesn't respect the minimum requirement.
	ErrInvalid = errors.New("invalid data")
	// ErrOutdated is returned if the data is deprecated.
	ErrOutdated = errors.New("outdated data")
	// ErrNotFound is the data doesn't exist.
	ErrNotFound = errors.New("not found")
	// ErrStrength is returned if the password is not safe.
	ErrStrength = errors.New("low password strength")
)

// Keyer returns the key of the data.
type Keyer interface {
	Key() []byte
}

// Validator returns in success if the data can be store.
type Validator interface {
	Valid() bool
}

// Data must be implement by any data to store.
type Data interface {
	Keyer
	Validator
}

// Service must be implements by any data source.
type Service interface {
	VaultService
	OwnerService
	TagService
	io.Closer
}

// VaultService must be implements by any service to manipulate the Vaults.
type VaultService interface {
	// CreateVault stores a Vault in database.
	CreateVault(v *Vault) error
	// DeleteVault deletes a Vault in database.
	DeleteVault(key string) error
	// Vaults lists the vaults in the given tag.
	Vaults(hash crypto.Hash, tag *Tag, prefix string) ([]*Vault, error)
	// Vault returns the requested Vault.
	Vault(hash crypto.Hash, key string) (*Vault, error)
	// UpdateVault updates the given Vault.
	UpdateVault(v *Vault) error
}

type timeUnix struct {
	time time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *timeUnix) UnmarshalJSON(b []byte) error {
	sec, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	t.time = time.Unix(sec, 0)

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (t *timeUnix) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.time.Unix())
}

type jsonVault struct {
	AddDate    *timeUnix  `json:"add_ts"`
	LastUpdate *timeUnix  `json:"upd_ts"`
	Login      *hashLogin `json:"login"`
	Name       string     `json:"name"`
	Tag        *Tag       `json:"tag"`
}

// Vault stores the data (login etc.) to be protected by encryption.
type Vault struct {
	v jsonVault
}

// EmptyVault returns a empty Vault based on the given hash to sign data.
func EmptyVault(hash crypto.Hash) *Vault {
	return &Vault{v: jsonVault{Login: &hashLogin{hash: hash}}}
}

// NewVault returns a new instance of Vault for the given data.
func NewVault(hash crypto.Hash, name string, tag *Tag, login *Login) *Vault {
	return &Vault{v: jsonVault{
		Name:       strings.TrimSpace(name),
		Tag:        tag,
		Login:      &hashLogin{hash: hash, Login: login},
		AddDate:    &timeUnix{time: time.Now()},
		LastUpdate: &timeUnix{time: time.Now()},
	}}
}

// CreationDate returns the creation date of the Vault.
func (v *Vault) AddDate() time.Time {
	return v.v.AddDate.time
}

// Key implements the Keyer interface.
func (v *Vault) Key() []byte {
	if !v.hasBasics() {
		// Avoids panic
		return nil
	}
	return bytes.Join([][]byte{v.Tag().Key(), []byte(v.Name())}, []byte(""))
}

// LastUpdate returns the last update of the Vault.
func (v *Vault) LastUpdate() time.Time {
	return v.v.LastUpdate.time
}

// Login returns the Login stored inside the Vault.
func (v *Vault) Login() *Login {
	if v.v.Login == nil {
		// Avoids panic.
		return nil
	}
	return v.v.Login.Login
}

// SignLogin adds a login to the vault after to have sign it.
func (v *Vault) SignLogin(hash crypto.Hash, l *Login) error {
	if l == nil || !l.Valid() {
		return ErrInvalid
	}
	v.v.Login = &hashLogin{hash: hash, Login: l}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (v *Vault) MarshalJSON() ([]byte, error) {
	// Changes the last modification date.
	v.v.LastUpdate = &timeUnix{time: time.Now()}
	return json.Marshal(v.v)
}

// Name returns the Vault's name.
func (v *Vault) Name() string {
	return v.v.Name
}

// Tag returns the Tag where the Vault is stored.
func (v *Vault) Tag() *Tag {
	if v.v.Tag == nil {
		// Avoids panic.
		return nil
	}
	return v.v.Tag
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (v *Vault) UnmarshalJSON(b []byte) (err error) {
	return json.Unmarshal(b, &v.v)
}

// Valid implements the Validator interface.
func (v *Vault) Valid() bool {
	if !v.hasBasics() {
		return false
	}
	if v.Login() == nil || !v.Login().Valid() {
		return false
	}
	if v.AddDate().IsZero() || v.AddDate().After(time.Now()) {
		return false
	}
	return true
}

func (v *Vault) hasBasics() bool {
	if v.Tag() == nil || !v.Tag().Valid() {
		return false
	}
	return v.Name() != ""
}

// OwnerService must be implemented by any service to manipulate the database owner.
type OwnerService interface {
	// CreateOwner creates and stores the owner of this database.
	CreateOwner(p *Passphrase) error
	// HasOwner returns in success if the database has a owner.
	HasOwner() bool
	// IsOwner returns in success if the given Passphrase matches to that of the base.
	IsOwner(p *Passphrase) bool
}

// TagService must be implemented by any service to manipulate the tags.
type TagService interface {
	// CreateTag creates a tag.
	CreateTag(t *Tag) error
	// DeleteTag deletes a tag.
	DeleteTag(key string) error
	// Tags lists all the tags.
	Tags() ([]*Tag, error)
}

// Used to join tag and entry names
const tagSep = "#"

// Tag is a tag.
type Tag struct {
	name string
}

// NewTag returns a new instance of Tag.
func NewTag(s string) *Tag {
	return &Tag{name: strings.TrimSpace(s)}
}

// Key implements the Keyer interface.
func (t *Tag) Key() []byte {
	if !t.Valid() {
		return nil
	}
	return []byte(t.name + tagSep)
}

// MarshalJSON implements tje json.Marshaler interface.
func (t *Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.name)
}

// Name returns the tag's name.
func (t *Tag) Name() string {
	return t.name
}

// UnmarshalJSON implements tje json.Unmarshaler interface.
func (t *Tag) UnmarshalJSON(b []byte) (err error) {
	return json.Unmarshal(b, &t.name)
}

// Valid implements the Validator interface.
func (t *Tag) Valid() bool {
	return t.name != "" && !strings.Contains(t.name, tagSep)
}
