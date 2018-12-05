# Safe

[![GoDoc](https://godoc.org/github.com/rvflash/safe?status.svg)](https://godoc.org/github.com/rvflash/safe)
[![Build Status](https://img.shields.io/travis/rvflash/safe.svg)](https://travis-ci.org/rvflash/safe)
[![Code Coverage](https://img.shields.io/codecov/c/github/rvflash/safe.svg)](http://codecov.io/github/rvflash/safe?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/safe)](https://goreportcard.com/report/github.com/rvflash/safe)

Safe is a tool to store all yours passwords in a local encrypted storage.

The owner uses a passphrase to sign in (sha256 hash).
This passphrase combined with a salt given on the application's launching is used to generate a HMAC hash.
This hash will be used to sign all the data with AES encryption.
![Application GTK3](https://raw.githubusercontent.com/rvflash/safe/master/doc/app.png)

## Installation

Safe uses the Go Modules coming from the 1.11 version of the language and GTK+3 as GUI.

Since the version 0.2.0, Safe is not anymore a web application but a application powered by GTK+3.
Thanks to the [gotk3](https://github.com/gotk3/gotk3)'s project for the bindings.

See the [installation instructions](https://github.com/gotk3/gotk3/wiki#installation) regarding our OS before going to the next step.

Finally, build and launch it:

```bash
$ git clone https://github.com/rvflash/safe.git
$ cd safe/cmd/safe
$ GO111MODULE=on go build
$ ./safe -salt="whatever-you-want-as-salt"
```

### Features

- [x] Local storage using boltDB.
- [x] ~~Web view based on local version of Bootstrap v4.1.3 (only CSS) and Vue.js v2.5.17.~~
- [x] Migrate the Vue.js application to a GTK+3 GUI in order to not use a web browser (avoids HTTP, unsafe extension, etc.).
- [ ] Historic of password's modifications.
- [ ] Notification center with alerts on outdated or low strength password.