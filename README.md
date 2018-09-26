# Safe

[![GoDoc](https://godoc.org/github.com/rvflash/safe?status.svg)](https://godoc.org/github.com/rvflash/safe)
[![Build Status](https://img.shields.io/travis/rvflash/safe.svg)](https://travis-ci.org/rvflash/safe)
[![Code Coverage](https://img.shields.io/codecov/c/github/rvflash/safe.svg)](http://codecov.io/github/rvflash/safe?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/safe)](https://goreportcard.com/report/github.com/rvflash/safe)

Safe is a tool to store all yours passwords in a local encrypted storage.

The owner uses a passphrase to sign in (sha256 hash).
This passphrase combined with a salt given on the application's launching is used to generate a HMAC hash.
This hash will be used to sign all the data with AES encryption.

![Web view](https://raw.githubusercontent.com/rvflash/safe/master/doc/webview.jpeg)


## Installation

Safe uses the Go Modules coming from the 1.11 version of the language.

```bash
$ git clone https://github.com/rvflash/safe.git
$ cd safe/cmd/ui
$ GO111MODULE=on go build -o safe
$ ./safe -salt="whatever-you-want-as-salt"
```
The web server starts on `http://localhost:7233/` by default. Opens it in private mode in your favorite browser. 


### Features

- [x] Local storage using boltDB.
- [x] Web view based on local version of Bootstrap v4.1.3 (only CSS) and Vue.js v2.5.17.
- [ ] Migrate the Vue.js application to a Qt GUI in order to not use a web browser (avoids HTTP, unsafe extension, etc.).
- [ ] Historic of password's modifications.
- [ ] Notification center with alerts on outdated or low strength password.