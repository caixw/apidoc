// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package locale

import (
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/issue9/assert"
)

func TestInitLocales(t *testing.T) {
	a := assert.New(t)

	a.NotError(initLocales())
	a.True(len(locales) > 0)

	p := message.NewPrinter(language.MustParse("zh-Hans"))
	a.Equal(p.Sprintf(FlagHUsage), locales["zh-Hans"][FlagHUsage])
}
