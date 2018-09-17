// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package safe

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"hash"
)

// Passphrase is the phrase used to protect the database.
type Passphrase struct {
	value []byte
}

// NewPassPhrase returns a new instance of Passphrase.
func NewPassPhrase(s string) *Passphrase {
	return &Passphrase{value: []byte(s)}
}

// Compare returns in error if the given hash doesn't match with the encrypted Passphrase.
func (p *Passphrase) Compare(hashed []byte) error {
	if len(hashed) == 0 {
		// Missing encoded pass phrase to compare.
		return ErrMissing
	}
	hash, err := p.sign()
	if err != nil {
		return err
	}
	if bytes.Compare(b64Encode(hash), hashed) != 0 {
		return ErrInvalid
	}
	return nil
}

func b64Encode(v []byte) []byte {
	b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
	base64.StdEncoding.Encode(b, v)
	return b
}

// Key implements the Keyer interface.
func (p *Passphrase) Key() []byte {
	return []byte("pass")
}

// MarshalJSON implements the json.Marshaler interface.
func (p *Passphrase) MarshalJSON() ([]byte, error) {
	b, err := p.sign()
	if err != nil {
		return nil, err
	}
	return json.Marshal(b)
}

// NewCipher returns a hash of 32 bytes to use as AES key to encrypt data.
// This key is not stored.
func (p *Passphrase) NewCipher(salt string) ([]byte, error) {
	if len(salt) == 0 {
		// Missing secret key to sign the pass phrase.
		return nil, ErrMissing
	}
	h := hmac.New(sha256.New, []byte(salt))
	return p.hash(h)
}

func (p *Passphrase) sign() (b []byte, err error) {
	if !p.Valid() {
		return nil, ErrTooShort
	}
	h := sha256.New()
	return p.hash(h)
}

func (p *Passphrase) hash(h hash.Hash) ([]byte, error) {
	_, err := h.Write(p.value)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// Valid returns in success if the key is long enough.
// It implements the Validator interface.
func (p *Passphrase) Valid() bool {
	return len(p.value) >= MinSize
}
