// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
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
	path := getPath(detectFlagSet)
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
	return cmdUsage(w, locale.CmdDetectUsage)
}
