// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package locale

import (
	"testing"

	"github.com/issue9/assert"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v5/internal/locale/syslocale"
)

func TestInit(t *testing.T) {
	a := assert.New(t)

	tag := language.MustParse("zh-Hans")
	a.NotError(Init(tag))
	a.Equal(localeTag, tag)

	tag = language.MustParse("zh-Hant")
	a.NotError(Init(tag))
	a.Equal(localeTag, tag)

	// 设置为系统语言
	systag, err := syslocale.Get()
	a.NotError(err)
	a.NotError(Init(language.Und))
	a.Equal(systag, localeTag)
}

func TestInitLocales(t *testing.T) {
	a := assert.New(t)

	a.NotError(initLocales())
	a.True(len(locales) > 0)
	a.Equal(len(locales), len(displayNames))

	p := message.NewPrinter(language.MustParse("zh-Hans"))
	a.Equal(p.Sprintf(FlagHUsage), locales[language.MustParse("zh-Hans")][FlagHUsage])
}
