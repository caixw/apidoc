// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"

	"github.com/caixw/apidoc/v5/internal/locale"
)

var localeFlagSet *flag.FlagSet

func init() {
	localeFlagSet = command.New("locale", doLocale, localeUsage)
}

func doLocale(w io.Writer) error {
	for tag, name := range locale.DisplayNames() {
		if _, err := fmt.Fprintln(w, tag, name); err != nil {
			return err
		}
	}
	return nil
}

func localeUsage(w io.Writer) error {
	_, err := fmt.Fprintln(w, locale.Sprintf(locale.CmdLocaleUsage))
	return err
}
