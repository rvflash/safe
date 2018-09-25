// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package safe_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/rvflash/safe"
)

func passPhrase() *safe.Passphrase {
	return safe.NewPassPhrase("what-a-beautiful-pass-phrase")
}

func TestPassphrase_Compare(t *testing.T) {
	var ph = []byte("xP/8nsOlAjlRXQ54vsY1amzecvnuJjmHMo+OIRXs/Ys=")
	var dt = []struct {
		pass   *safe.Passphrase
		hashed []byte
		err    error
	}{
		{pass: passPhrase(), err: safe.ErrMissing},
		{pass: passPhrase(), hashed: []byte("euh"), err: safe.ErrInvalid},
		{pass: safe.NewPassPhrase("what-a-"), hashed: ph, err: safe.ErrTooShort},
		{pass: passPhrase(), hashed: ph},
	}
	var err error
	for i, tt := range dt {
		err = tt.pass.Compare(tt.hashed)
		if !reflect.DeepEqual(err, tt.err) {
			t.Fatalf("%d. mismatch result: got=%q exp=%q", i, err, tt.err)
		}
		if err == nil {
			b, err := json.Marshal(tt.pass) //tt.pass.MarshalJSON()
			if err != nil {
				t.Fatalf("%d. unexpected error: %q", i, err)
			}
			if b = bytes.Trim(b, `"`); !bytes.Equal(b, tt.hashed) {
				t.Fatalf("%d. mismatch hash: got=%q exp=%q", i, b, tt.hashed)
			}
		}
	}
}

func TestPassphrase_Key(t *testing.T) {
	k := []byte("pass")
	if b := passPhrase().Key(); !bytes.Equal(b, k) {
		t.Fatalf("mismatch content: got=%q exp=%q", b, k)
	}
}

func TestPassphrase_NewCipher(t *testing.T) {
	_, err := passPhrase().NewCipher("")
	if err != safe.ErrMissing {
		t.Fatalf("unexpected error: %q", err)
	}
	b, err := passPhrase().NewCipher("say-hello")
	if err != nil {
		t.Fatalf("unexpected error: %q", err)
	}
	exp := "teP8g3YVWDHMzYgyuODZvWYoYgF8GRtGM1hKlAMpnqk="
	if s := base64.StdEncoding.EncodeToString(b); s != exp {
		t.Fatalf("mismatch content: got=%q exp=%q", s, exp)
	}
}
