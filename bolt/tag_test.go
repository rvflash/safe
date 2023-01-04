// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package bolt_test

import (
	"testing"

	"github.com/rvflash/safe"
)

func newPassPhrase() *safe.Passphrase {
	return safe.NewPassPhrase("testOnlyDoNotUse")
}

func newDBWithOwner() (safe.Service, error) {
	db, err := newDB()
	if err != nil {
		return nil, err
	}
	if err = db.CreateOwner(newPassPhrase()); err != nil {
		return nil, err
	}
	return db, nil

}

func TestSafe_Tags(t *testing.T) {
	// Opens temporary DB
	db, err := newDBWithOwner()
	if err != nil {
		t.Fatalf("workspace: %s", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			t.Fatalf("db closing: %s", err)
		}
	}()
	// Sort order matters in this tests suite.
	var dt = []struct {
		tag                     *safe.Tag
		key                     string
		list                    []*safe.Tag
		errNew, errDel, errList error
	}{
		{tag: safe.NewTag(""), errNew: safe.ErrInvalid, errDel: safe.ErrNotFound},
		{tag: safe.NewTag("Business"), errDel: safe.ErrNotFound, list: []*safe.Tag{safe.NewTag("Business")}},
		{tag: safe.NewTag("Social"), key: "Unknown", list: []*safe.Tag{
			safe.NewTag("Business"),
			safe.NewTag("Social"),
		}},
		{tag: safe.NewTag("Jobs"), key: "Social#", list: []*safe.Tag{
			safe.NewTag("Business"),
			safe.NewTag("Jobs"),
		}},
	}
	var list []*safe.Tag
	for i, tt := range dt {
		if err = db.CreateTag(tt.tag); err != tt.errNew {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.errNew)
		}
		if err = db.DeleteTag(tt.key); err != tt.errDel {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.errDel)
		}
		list, err = db.Tags()
		if err != tt.errList {
			t.Errorf("%d. mismatch error: got=%q exp=%q", i, err, tt.errList)
		}
		if !equalTags(list, tt.list) {
			t.Errorf("%d. mismatch tags: got=%v exp=%v", i, list, tt.list)
		}
	}
}

func equalTags(l1, l2 []*safe.Tag) bool {
	if len(l1) != len(l2) {
		return false
	}
	if l1 == nil {
		return true
	}
	for i, t := range l1 {
		if t.Name() != l2[i].Name() {
			return false
		}
	}
	return true
}
