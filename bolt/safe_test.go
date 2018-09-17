// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package bolt_test

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/rvflash/safe/bolt"
)

func TestOpen(t *testing.T) {
	dir, err := unknownDir()
	if err != nil {
		t.Fatalf("workspace: %s", err)
	}
	var dt = []struct {
		path, name string
		err        error
	}{
		{path: os.TempDir()},
		{path: os.TempDir(), name: "file"},
		{path: dir, name: "file", err: errors.New("no such file or directory")},
	}
	for i, tt := range dt {
		if _, err = bolt.Open(tt.path, tt.name); !checkErr(err, tt.err, strings.Contains) {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.err)
		}
	}
}

func unknownDir() (string, error) {
	// Creates a new one in default  temporary directory.
	name, err := ioutil.TempDir("", "safe")
	if err != nil {
		return "", err
	}
	// Removes it (safely)
	if err = os.Remove(name); err != nil {
		return "", err
	}
	return name, nil
}

func checkErr(e1, e2 error, f func(a, b string) bool) bool {
	if e1 == nil {
		return e2 == nil
	}
	if e2 == nil {
		// fact: e1 is not nil
		return false
	}
	return f(e1.Error(), e2.Error())
}
