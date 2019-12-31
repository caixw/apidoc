// SPDX-License-Identifier: MIT

package cmd

import (
	"io"
	"strings"

	"github.com/caixw/apidoc/v5/internal/lang"
	"github.com/caixw/apidoc/v5/internal/locale"
)

func initLang() {
	command.New("lang", language, langUsage)
}

func language(w io.Writer) error {
	tail := 3

	ls := lang.Langs()
	langs := make([]*lang.Language, 1, len(ls)+1)
	langs[0] = &lang.Language{
		Name:        locale.Sprintf(locale.LangID),
		DisplayName: locale.Sprintf(locale.LangName),
		Exts:        []string{locale.Sprintf(locale.LangExts)},
	}
	langs = append(langs, ls...)

	// 计算各列的最大长度值
	var maxDisplay, maxName int
	for _, l := range langs {
		calcMaxWidth(l.DisplayName, &maxDisplay)
		calcMaxWidth(l.Name, &maxName)
	}
	maxDisplay += tail
	maxName += tail

	for _, l := range langs {
		n := l.Name + strings.Repeat(" ", maxName-len(l.Name))
		d := l.DisplayName + strings.Repeat(" ", maxDisplay-len(l.DisplayName))
		printLine(w, infoColor, n, d, strings.Join(l.Exts, " "))
	}

	return nil
}

func langUsage(w io.Writer) error {
	return cmdUsage(w, locale.CmdLangUsage)
}

func calcMaxWidth(content string, max *int) {
	width := len(content)
	if width > *max {
		*max = width
	}
}
