// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package syslocale

import (
	"testing"

	"golang.org/x/text/language"

	"github.com/issue9/assert"
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
