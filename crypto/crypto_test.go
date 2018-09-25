// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package crypto_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/rvflash/safe/crypto"
)

func hash() *crypto.PrivateKey {
	return crypto.New([]byte("6368616e676520746869732070617373"))
}

func TestNew(t *testing.T) {
	var dt = []struct {
		h       *crypto.PrivateKey
		size    int
		in, out []byte
		err     error
	}{
		{h: crypto.New(nil), err: crypto.ErrKey},
		{h: hash(), in: []byte(""), err: crypto.ErrData},
		{h: hash(), in: []byte("secret data"), size: 27},
	}
	var (
		err  error
		out  []byte
		size int
	)
	for i, tt := range dt {
		out, err = tt.h.Sign(tt.in)
		if err != tt.err {
			t.Fatalf("%d. unexpected error: got=%q exp=%q", i, err, tt.err)
		}
		if size = len(out); size != tt.size {
			t.Errorf("%d. unexpected result: got=%q (%d)", i, hex.EncodeToString(out), size)
		}
		out, err = tt.h.Decrypt(out)
		if err != tt.err {
			t.Fatalf("%d. unexpected error: got=%q exp=%q", i, err, tt.err)
		}
		if !bytes.Equal(out, tt.in) {
			t.Errorf("%d. unexpected result: got=%s exp=%s", i, out, tt.out)
		}
	}
}
