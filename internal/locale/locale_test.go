// SPDX-License-Identifier: MIT

package locale

import (
	"testing"

	"github.com/issue9/assert"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func TestDisplayNames(t *testing.T) {
	a := assert.New(t)

	a.Equal(len(DisplayNames()), len(locales))
	for tag := range DisplayNames() {
		_, found := locales[tag]
		a.True(found)
	}
}

func TestLocale(t *testing.T) {
	a := assert.New(t)

	l := NewLocale(FlagHUsage)
	a.NotNil(l)

	l1 := l.String(message.NewPrinter(language.MustParse("zh-hans")))
	l2 := l.String(message.NewPrinter(language.MustParse("zh-hant")))
	a.NotEqual(l1, l2)
}
