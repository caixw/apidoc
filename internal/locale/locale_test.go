// SPDX-License-Identifier: MIT

package locale

import (
	"testing"

	"github.com/issue9/assert"
	"golang.org/x/text/language"
)

var _ error = &Err{}

func TestDisplayNames(t *testing.T) {
	a := assert.New(t)

	a.Equal(len(DisplayNames()), 2)

	tag := language.MustParse("zh-Hans")
	ds := DisplayNames()
	ds[tag] = "123"
	a.NotEqual(ds[tag], displayNames[tag])
}

func TestTranslate(t *testing.T) {
	a := assert.New(t)
	a.Equal(Translate("cmn-hans", ErrRequired), zhHans[ErrRequired])
	a.Equal(Translate("zh-hant", ErrRequired), zhHant[ErrRequired])
	a.NotEqual(Translate("zh-hant", ErrRequired), zhHans[ErrRequired])
	Translate("not-well-format", zhHans[ErrRequired]) // 无效的 tag 格式
}

func TestSetLocale(t *testing.T) {
	a := assert.New(t)

	tag := language.MustParse("zh-Hans")
	SetLanguageTag(tag)
	a.NotEqual(Sprintf(ErrRequired), zhHant[ErrRequired]).
		Equal(Sprintf(ErrRequired), zhHans[ErrRequired]).
		Equal(NewError(ErrRequired).Error(), zhHans[ErrRequired])

	// zh-cn 应该会转换到 zh-hans
	tag = language.MustParse("zh-CN")
	SetLanguageTag(tag)
	a.NotEqual(Sprintf(ErrRequired), zhHant[ErrRequired]).
		Equal(Sprintf(ErrRequired), zhHans[ErrRequired]).
		Equal(NewError(ErrRequired).Error(), zhHans[ErrRequired])

	tag = language.MustParse("zh-Hant")
	SetLanguageTag(tag)
	a.Equal(Sprintf(ErrRequired), zhHant[ErrRequired]).
		Equal(NewError(ErrRequired).Error(), zhHant[ErrRequired])
}
