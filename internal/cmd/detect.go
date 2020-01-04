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

func initDetect() {
	detectFlagSet = command.New("detect", detect, buildUsage(locale.CmdDetectUsage))
}

func detect(w io.Writer) error {
	path := getPath(detectFlagSet)
	h := message.NewHandler(newHandlerFunc())
	defer h.Stop()

	if err := apidoc.Detect(path, true); err != nil {
		return err
	}

	h.Message(message.Succ, locale.ConfigWriteSuccess, path)
	return nil
}
