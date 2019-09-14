// SPDX-License-Identifier: MIT

package locale

import (
	"fmt"
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

var _ fmt.Stringer = &Locale{}

func TestLocale(t *testing.T) {
	a := assert.New(t)

	l := NewLocale(FlagHUsage)
	a.NotNil(l)

	a.NotError(Init(language.MustParse("zh-hans")))
	l1 := l.String()
	a.NotError(Init(language.MustParse("zh-hant")))
	l2 := l.String()
	a.NotEqual(l1, l2)
}
