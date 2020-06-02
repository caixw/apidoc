// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/language/display"

	"github.com/caixw/apidoc/v7"
	"github.com/caixw/apidoc/v7/internal/locale"
)

func initLocale() {
	command.New("locale", locale.Sprintf(locale.CmdLocaleUsage), doLocale)
}

func doLocale(w io.Writer) error {
	locales := make(map[string]string, len(apidoc.Locales()))

	// 计算各列的最大长度值
	var maxID int
	for _, tag := range apidoc.Locales() {
		id := tag.String()
		calcMaxWidth(id, &maxID)
		locales[id] = display.Self.Name(tag)
	}
	maxID += tail

	for k, v := range locales {
		id := k + strings.Repeat(" ", maxID-len(k))
		if _, err := fmt.Fprintln(w, id, v); err != nil {
			return err
		}
	}

	return nil
}
