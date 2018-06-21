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

func TestNewSyntaxError(t *testing.T) {
	a := assert.New(t)

	err := newSyntaxError("file.go", 1, locale.ErrDirIsEmpty)
	a.Error(err)
	a.Contains(err.Error(), "file.go")
}
