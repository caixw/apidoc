// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"strings"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/is"
)

// 对一些堂量的基本检测。
func TestConsts(t *testing.T) {
	a := assert.New(t)

	a.True(is.URL(RepoURL))
	a.True(is.URL(OfficialURL))
	a.True(strings.IndexRune(Symbols, '@') < 0)
}
