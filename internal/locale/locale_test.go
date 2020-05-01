// SPDX-License-Identifier: MIT

package locale

import (
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/utils"
	"golang.org/x/text/language"
)

var _ error = &Err{}

func TestDisplayNames(t *testing.T) {
	a := assert.New(t)

	a.Equal(len(DisplayNames()), 3)

	tag := language.MustParse("zh-Hans")
	ds := DisplayNames()
	ds[tag] = "123"
	a.NotEqual(ds[tag], displayNames[tag])
}

func TestTranslate(t *testing.T) {
	a := assert.New(t)
	a.Equal(Translate("zh-hans", ErrRequired), zhHans[ErrRequired])
	a.Equal(Translate("zh-hant", ErrRequired), zhHant[ErrRequired])
	a.NotEqual(Translate("zh-hant", ErrRequired), zhHans[ErrRequired])
	a.Panic(func() {
		Translate("not-well-format", ErrRequired)
	})
}

func TestSetLocale(t *testing.T) {
	a := assert.New(t)

	tag := language.MustParse("zh-Hans")
	a.True(SetLanguageTag(tag)).
		Equal(localeTag, tag).
		NotEqual(Sprintf(ErrRequired), zhHant[ErrRequired]).
		Equal(Sprintf(ErrRequired), zhHans[ErrRequired]).
		Equal(NewError(ErrRequired).Error(), zhHans[ErrRequired])

	// zh-cn 应该会转换到 zh-hans
	tag = language.MustParse("zh-CN")
	a.True(SetLanguageTag(tag)).
		Equal(LanguageTag(), tag).
		NotEqual(Sprintf(ErrRequired), zhHant[ErrRequired]).
		Equal(Sprintf(ErrRequired), zhHans[ErrRequired]).
		Equal(NewError(ErrRequired).Error(), zhHans[ErrRequired])

	tag = language.MustParse("zh-Hant")
	a.True(SetLanguageTag(tag)).
		Equal(LanguageTag(), tag).
		Equal(Sprintf(ErrRequired), zhHant[ErrRequired]).
		Equal(NewError(ErrRequired).Error(), zhHant[ErrRequired])

	// 设置为系统语言
	systemTag, err := utils.GetSystemLanguageTag()
	a.NotError(err)
	a.True(SetLanguageTag(systemTag)).
		Equal(systemTag, LanguageTag())
}
