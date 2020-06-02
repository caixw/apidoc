// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/issue9/cmdopt"

	"github.com/caixw/apidoc/v7/internal/lang"
	"github.com/caixw/apidoc/v7/internal/locale"
)

func initLang(command *cmdopt.CmdOpt) {
	command.New("lang", locale.Sprintf(locale.CmdLangUsage), doLang)
}

func doLang(w io.Writer) error {
	ls := lang.Langs()
	langs := make([]*lang.Language, 1, len(ls)+1)
	langs[0] = &lang.Language{
		ID:          locale.Sprintf(locale.LangID),
		DisplayName: locale.Sprintf(locale.LangName),
		Exts:        []string{locale.Sprintf(locale.LangExts)},
	}
	langs = append(langs, ls...)

	// 计算各列的最大长度值
	var maxDisplay, maxName int
	for _, l := range langs {
		calcMaxWidth(l.DisplayName, &maxDisplay)
		calcMaxWidth(l.ID, &maxName)
	}
	maxDisplay += tail
	maxName += tail

	for _, l := range langs {
		n := l.ID + strings.Repeat(" ", maxName-len(l.ID))
		d := l.DisplayName + strings.Repeat(" ", maxDisplay-len(l.DisplayName))
		if _, err := fmt.Fprintln(w, n, d, strings.Join(l.Exts, " ")); err != nil {
			return err
		}
	}

	return nil
}

func calcMaxWidth(content string, max *int) {
	width := len(content)
	if width > *max {
		*max = width
	}
}
