// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"

	"github.com/caixw/apidoc/v5"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

var detectFlagSet *flag.FlagSet

func init() {
	detectFlagSet = command.New("detect", detect, detectUsage)
}

func detect(w io.Writer) error {
	var path = "./"
	if 0 != detectFlagSet.NArg() {
		path = detectFlagSet.Arg(0)
	}

	h := message.NewHandler(newHandlerFunc())
	if err := apidoc.Detect(path, true); err != nil {
		printLine(erroOut, erroColor, err)
		return nil
	}
	printLocale(succOut, succColor, locale.ConfigWriteSuccess, path)

	h.Stop()
	return nil
}

func detectUsage(w io.Writer) error {
	_, err := fmt.Fprintln(w, locale.Sprintf(locale.CmdDetectUsage))
	return err
}
