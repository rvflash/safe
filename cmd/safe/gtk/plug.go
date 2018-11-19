// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package gtk

import "github.com/rvflash/safe/app"

// Plugin ...
type Plugin interface {
	App() *app.Safe
	Log(format string, args ...interface{})
}

// Logger must be implemented by any logger.
type Logger interface {
	// Printf logs a message at level Info on the standard logger.
	Printf(format string, args ...interface{})
}

// Plug ...
type Plug struct {
	*Object
	s *app.Safe
	l Logger
}

// NewPlug ...
func NewPlug(id, xml string, app *app.Safe, log Logger) (*Plug, error) {
	o, err := NewObject(id, xml)
	if err != nil {
		return nil, err
	}
	return &Plug{Object: o, s: app, l: log}, nil
}

// App ...
func (p *Plug) App() *app.Safe {
	return p.s
}

// Log ...
func (p *Plug) Log(format string, args ...interface{}) {
	if p.l != nil {
		p.l.Printf(format, args)
	}
}
