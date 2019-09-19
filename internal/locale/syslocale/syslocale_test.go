// SPDX-License-Identifier: MIT

package syslocale

import (
	"testing"

	"github.com/issue9/assert"
	"golang.org/x/text/language"
)

func TestGet(t *testing.T) {
	a := assert.New(t)

	lang, err := Get()
	if err != nil {
		a.Equal(lang, language.Und)
	} else {
		a.NotEqual(lang, language.Und)
	}
}

func TestGetLocaleName(t *testing.T) {
	a := assert.New(t)

	name, err := getLocaleName()
	a.NotError(err).True(len(name) > 0)
}
