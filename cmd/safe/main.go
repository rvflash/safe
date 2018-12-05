// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package main

//go:generate esc -o static/static.go -pkg static static/ui

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/rvflash/safe/app"
	"github.com/rvflash/safe/bolt"
	"github.com/rvflash/safe/cmd/safe/gtk"
)

func main() {
	// Parses flags
	root := flag.String("dir", "", "directory path of the database")
	salt := flag.String("salt", "Sh0u!ldN0t8eUs3d", "public key")
	user := flag.String("user", "default", "name of the database")
	test := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	// Application engine.
	logger := log.New(os.Stdout, "safe: ", log.Lshortfile|log.LstdFlags)
	if *root == "" {
		exec, err := os.Executable()
		if err != nil {
			log.Fatalf("os: %s", err)
		}
		*root = filepath.Dir(exec)
	}
	if *salt == "" {
		logger.Fatal("flag: missing salt")
	}
	if *user == "" {
		logger.Fatal("flag: missing username")
	}

	// Bolt database.
	db, err := bolt.Open(*root, *user)
	if err != nil {
		log.Fatalf("db: %s", err)
	}
	engine := app.New(db, *salt, *root, nil)

	// Launches the application.
	safe, err := gtk.Init(&os.Args, engine, logger, *test)
	if err != nil {
		log.Fatal(err)
	}
	safe.Run()
}
