// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// Commons errors.
var (
	// ErrKey is returned is the given key is invalid.
	ErrKey = errors.New("invalid key")
	// ErrData is returned is the data is mal-formatted.
	ErrData = errors.New("invalid data")
)

// Signer must be implemented by any hash to encrypt a data.
type Signer interface {
	Sign(plain []byte) ([]byte, error)
}

// Decrypter must be implemented by any hash to decrypt a data.
type Decrypter interface {
	Decrypt(hash []byte) ([]byte, error)
}

// Hash must be implemented by any hashed to sign and decrypt data.
type Hash interface {
	Signer
	Decrypter
}

// New returns a private key.
func New(key []byte) *PrivateKey {
	return &PrivateKey{value: key}
}

// PrivateKey is a private key.
type PrivateKey struct {
	value []byte
}

// Decrypt decrypts the given hash, returns the plain text behind it or a error if it fails to do it.
func (k *PrivateKey) Decrypt(hash []byte) ([]byte, error) {
	if len(k.value) == 0 {
		return nil, ErrKey
	}
	if len(hash) < aes.BlockSize {
		return nil, ErrData
	}
	block, err := aes.NewCipher(k.value)
	if err != nil {
		return nil, err
	}
	plain := hash[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, hash[:aes.BlockSize])
	stream.XORKeyStream(plain, plain)

	return plain, nil
}

// Sign signs the given plain text, returns a hash or a error if it fails to do it.
func (k *PrivateKey) Sign(plain []byte) ([]byte, error) {
	if len(k.value) == 0 {
		return nil, ErrKey
	}
	if len(plain) == 0 {
		return nil, ErrData
	}
	block, err := aes.NewCipher(k.value)
	if err != nil {
		return nil, err
	}
	hash := make([]byte, aes.BlockSize+len(plain))
	iv := hash[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(hash[aes.BlockSize:], plain)

	return hash, nil
}
