// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package lexer

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/internal/locale"
)

func TestNewError(t *testing.T) {
	a := assert.New(t)

	err := newError("file.go", "@api", 1, locale.ErrDirIsEmpty)
	a.Error(err)
	a.Contains(err.Error(), "file.go")
}
