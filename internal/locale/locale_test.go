// SPDX-License-Identifier: MIT

package locale

import (
	"testing"

	"github.com/issue9/assert"
	"golang.org/x/text/language"
)

var _ error = &Err{}

func TestVarsDefaultLocaleID(t *testing.T) {
	a := assert.New(t)

	// 必须保证第一个元素是默认值
	a.Equal(tags[0].String(), DefaultLocaleID)
}

func TestTranslate(t *testing.T) {
	a := assert.New(t)
	a.Equal(Translate("cmn-hans", ErrInvalidUTF8Character), cmnHans[ErrInvalidUTF8Character])
	a.Equal(Translate("zh-hant", ErrInvalidUTF8Character), cmnHant[ErrInvalidUTF8Character])
	a.NotEqual(Translate("zh-hant", ErrInvalidUTF8Character), cmnHans[ErrInvalidUTF8Character])
	Translate("not-well-format", cmnHans[ErrInvalidUTF8Character]) // 无效的 tag 格式
}

func TestSetTag(t *testing.T) {
	a := assert.New(t)

	tag := language.MustParse("zh-Hans")
	SetTag(tag)
	a.NotEqual(Sprintf(ErrInvalidUTF8Character), cmnHant[ErrInvalidUTF8Character]).
		Equal(Sprintf(ErrInvalidUTF8Character), cmnHans[ErrInvalidUTF8Character]).
		Equal(NewError(ErrInvalidUTF8Character).Error(), cmnHans[ErrInvalidUTF8Character])

	// zh-cn 应该会转换到 zh-hans
	tag = language.MustParse("zh-CN")
	SetTag(tag)
	a.NotEqual(Sprintf(ErrInvalidUTF8Character), cmnHant[ErrInvalidUTF8Character]).
		Equal(Sprintf(ErrInvalidUTF8Character), cmnHans[ErrInvalidUTF8Character]).
		Equal(NewError(ErrInvalidUTF8Character).Error(), cmnHans[ErrInvalidUTF8Character])

	tag = language.MustParse("zh-Hant")
	SetTag(tag)
	a.Equal(Sprintf(ErrInvalidUTF8Character), cmnHant[ErrInvalidUTF8Character]).
		Equal(NewError(ErrInvalidUTF8Character).Error(), cmnHant[ErrInvalidUTF8Character])
}
