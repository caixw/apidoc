// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package parser

import (
	"testing"

	"github.com/caixw/apidoc/locale"
	"github.com/issue9/assert"
)

var _ error = &syntaxError{}

func TestNewWarn(t *testing.T) {
	a := assert.New(t)

	err := newWarn("file.go", 1, locale.ErrDirIsEmpty)
	a.Error(err)
	a.Contains(err.Error(), "file.go")

	synErr, ok := err.(*syntaxError)
	a.True(ok)
	a.Equal(synErr.Type, typeWarn)
}

func TestNewError(t *testing.T) {
	a := assert.New(t)

	err := newError("file.go", 1, locale.ErrDirIsEmpty)
	a.Error(err)
	a.Contains(err.Error(), "file.go")

	synErr, ok := err.(*syntaxError)
	a.True(ok)
	a.Equal(synErr.Type, typeError)
}
