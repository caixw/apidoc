// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/issue9/cmdopt"
	"golang.org/x/text/width"

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
	var maxName, maxID int
	for _, l := range langs {
		calcMaxWidth(l.DisplayName, &maxName)
		calcMaxWidth(l.ID, &maxID)
	}
	maxName += tail
	maxID += tail

	for _, l := range langs {
		id := l.ID + strings.Repeat(" ", maxID-textWidth(l.ID))
		name := l.DisplayName + strings.Repeat(" ", maxName-textWidth(l.DisplayName))
		if _, err := fmt.Fprintln(w, id, name, strings.Join(l.Exts, " ")); err != nil {
			return err
		}
	}

	return nil
}

func calcMaxWidth(content string, max *int) {
	if w := textWidth(content); w > *max {
		*max = w
	}
}

func textWidth(text string) (w int) {
	for _, r := range text {
		switch width.LookupRune(rune(r)).Kind() {
		case width.EastAsianFullwidth, width.EastAsianWide:
			w += 2
		default:
			w++
		}
	}
	return
}
