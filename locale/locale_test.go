// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package locale

import (
	"testing"

	"github.com/issue9/assert"
)

func TestInit(t *testing.T) {
	a := assert.New(t)

	_, err := Init()
	a.NotError(err)
	a.True(len(locales) > 0)
}

func TestGetLocaleName(t *testing.T) {
	a := assert.New(t)

	name, err := getLocaleName()
	a.NotError(err).True(len(name) > 0)
}
