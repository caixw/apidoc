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

	a.Equal(len(DisplayNames()), 3)
}

func TestInit(t *testing.T) {
	a := assert.New(t)

	tag := language.MustParse("zh-Hans")
	a.NotError(Init(tag)).
		Equal(localeTag, tag).
		NotEqual(Sprintf(ErrRequired), zhHant[ErrRequired]).
		Equal(Sprintf(ErrRequired), zhHans[ErrRequired])

	// zh-cn 应该会转换到 zh-hans
	tag = language.MustParse("zh-CN")
	a.NotError(Init(tag)).
		Equal(localeTag, tag).
		NotEqual(Sprintf(ErrRequired), zhHant[ErrRequired]).
		Equal(Sprintf(ErrRequired), zhHans[ErrRequired])

	tag = language.MustParse("zh-Hant")
	a.NotError(Init(tag)).
		Equal(localeTag, tag).
		Equal(Sprintf(ErrRequired), zhHant[ErrRequired])

	// 设置为系统语言
	systag, err := syslocale.Get()
	a.NotError(err)

	// 设置为 Und，依然会采用系统语言
	a.NotError(Init(language.Und)).
		Equal(systag, localeTag)
}
