// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/vars"
)

var detectFlagSet *flag.FlagSet

var detectRecursive bool

func initDetect() {
	detectFlagSet = command.New("detect", detect, detectUsage)
	detectFlagSet.BoolVar(&detectRecursive, "r", true, locale.Sprintf(locale.FlagDetectRecursive))
}

func detect(io.Writer) error {
	h := core.NewMessageHandler(newHandlerFunc())
	defer h.Stop()

	uri, err := getPath(detectFlagSet)
	if err != nil {
		return err
	}

	cfg, err := build.DetectConfig(uri, detectRecursive)
	if err != nil {
		return err
	}

	path := uri.Append(vars.AllowConfigFilenames[0])
	if err = cfg.SaveToFile(path); err != nil {
		return err
	}

	h.Message(core.Succ, locale.ConfigWriteSuccess, uri)
	return nil
}

func detectUsage(w io.Writer) error {
	_, err := fmt.Fprintln(w, locale.Sprintf(locale.CmdDetectUsage, getFlagSetUsage(detectFlagSet)))
	return err
}
