// SPDX-License-Identifier: MIT

package build

import (
	"testing"

	"github.com/issue9/assert"
)

func TestDetectInput(t *testing.T) {
	a := assert.New(t)

	o, err := detectInput("./testdata", true)
	a.NotError(err).NotEmpty(o)
	a.Equal(len(o), 2). // c and php
				Equal(o[0].Lang, "c++").
				Equal(o[1].Lang, "php")
}

func TestDetectLanguage(t *testing.T) {
	a := assert.New(t)
	exts := map[string]int{
		".h":     2,
		".c":     3,
		".swift": 1,
		".php":   2,
	}

	langs := detectLanguage(exts)
	a.Equal(len(langs), 3) // c++,php,swift
	a.Equal(langs[0].ID, "c++").
		Equal(langs[0].count, 5)
	a.Equal(langs[1].ID, "php").
		Equal(langs[1].count, 2)
	a.Equal(langs[2].ID, "swift").
		Equal(langs[2].count, 1)
}

func TestDetectExts(t *testing.T) {
	a := assert.New(t)

	files, err := detectExts("./testdata", false)
	a.NotError(err)
	a.Equal(len(files), 5)
	a.Equal(files[".php"], 1).Equal(files[".c"], 1)

	files, err = detectExts("./testdata", true)
	a.NotError(err)
	a.Equal(len(files), 6)
	a.Equal(files[".php"], 1).Equal(files[".1"], 3)
}
