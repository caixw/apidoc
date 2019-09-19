// SPDX-License-Identifier: MIT

package locale

import (
	"testing"

	"github.com/issue9/assert"
	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v5/internal/locale/syslocale"
)

func TestDisplayNames(t *testing.T) {
	a := assert.New(t)

	a.Equal(len(DisplayNames()), 2)
}

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
