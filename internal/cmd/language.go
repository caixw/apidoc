// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/caixw/apidoc/v5/internal/lang"
	"github.com/caixw/apidoc/v5/internal/locale"
)

var langFlagSet *flag.FlagSet

func init() {
	langFlagSet = command.New("lang", language, langUsage)
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
		width := len(l.DisplayName)
		if width > maxDisplay {
			maxDisplay = width
		}
		width = len(l.Name)
		if width > maxName {
			maxName = width
		}
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
	_, err := fmt.Fprintln(w, locale.Sprintf(locale.FlagLUsage))
	return err
}
