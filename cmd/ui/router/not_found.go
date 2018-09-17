// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package router

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type notFoundPage struct{}

// NotFound ...
func NotFound() Handler {
	return &notFoundPage{}
}

// Handle implements the Handler interface.
func (h *notFoundPage) Handle(c *gin.Context) {
	c.JSON(http.StatusNotFound, toErr(errors.New(http.StatusText(http.StatusNotFound))))
}

func toErr(err error) gin.H {
	return gin.H{"error": err.Error()}
}
