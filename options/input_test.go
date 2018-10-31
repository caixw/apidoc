// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package options

import (
	"testing"

	"github.com/issue9/assert"
)

func TestDetectExts(t *testing.T) {
	a := assert.New(t)

	files, err := detectExts("./testdata", false)
	a.NotError(err)
	a.Equal(len(files), 4)
	a.Equal(files[".php"], 1).Equal(files[".c"], 1)

	files, err = detectExts("./testdata", true)
	a.NotError(err)
	a.Equal(len(files), 5)
	a.Equal(files[".php"], 1).Equal(files[".1"], 3)
}

func TestDetect(t *testing.T) {
	a := assert.New(t)

	o, err := Detect("./testdata", true)
	a.NotError(err).NotEmpty(o)
	a.NotContains(o.Exts, ".1") // .1 不存在于已定义的语言中
}
