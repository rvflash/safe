// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/rvflash/safe/app"
	"github.com/rvflash/safe/bolt"
	"github.com/rvflash/safe/cmd/ui/router"
)

func main() {
	// Gets configuration
	port := flag.Int("port", 4433, "service port")
	salt := flag.String("salt", "", "public key")
	user := flag.String("user", "default", "name of the database")
	test := flag.Bool("test", false, "test mode")
	flag.Parse()

	// Current working space.
	exec, err := os.Executable()
	if err != nil {
		log.Fatalf("os: %s", err)
	}
	root := filepath.Dir(exec)

	// Application engine.
	if *salt == "" {
		fatal(errors.New("flag: missing salt"))
	}
	if *user == "" {
		fatal(errors.New("flag: missing username"))
	}
	// > Bolt database
	db, err := bolt.Open(root, *user)
	if err != nil {
		log.Fatalf("db: %s", err)
	}
	// > Session
	session := app.NewSession(time.Minute*20, time.Second)
	engine := app.New(db, *salt, root, session)

	// Launches the server.
	route := router.NewRouter(*port, engine, *test)
	srv := http.Server{
		Addr:    route.Addr(),
		Handler: route.Handler(),
	}
	go func() {
		log.Printf("server: listening on localhost%s\n", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server (listen): %s\n", err)
		}
	}()

	// Graceful shutdown.
	var failed int
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	// > Waiting for SIGINT
	<-exit
	// > Closes the server.
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	if err = srv.Shutdown(ctx); err != nil {
		log.Printf("server (shutdown): %s\n", err)
		failed = 1
	}
	// > Closes the database.
	if err = engine.Close(); err != nil {
		log.Printf("app: %s\n", err)
		failed = 1
	}
	log.Println("server: closed")
	os.Exit(failed)
}

func fatal(err error) {
	flag.Usage()
	log.Fatalf("flag: %s", err.Error())
}
