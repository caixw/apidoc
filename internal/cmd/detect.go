// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"

	"github.com/caixw/apidoc/v6"
	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/message"
)

var detectFlagSet *flag.FlagSet

var detectRecursive bool

func initDetect() {
	detectFlagSet = command.New("detect", detect, detectUsage)
	detectFlagSet.BoolVar(&detectRecursive, "r", true, locale.Sprintf(locale.FlagDetectRecursive))
}

func detect(w io.Writer) error {
	path := getPath(detectFlagSet)
	h := message.NewHandler(newHandlerFunc())
	defer h.Stop()

	if err := apidoc.Detect(path, detectRecursive); err != nil {
		return err
	}

	h.Message(message.Succ, locale.ConfigWriteSuccess, path)
	return nil
}

func detectUsage(w io.Writer) error {
	_, err := fmt.Fprintln(w, locale.Sprintf(locale.CmdDetectUsage, getFlagSetUsage(detectFlagSet)))
	return err
}
