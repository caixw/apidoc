// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"
	"path"

	"github.com/caixw/apidoc/v6/build"
	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/internal/vars"
)

var detectFlagSet *flag.FlagSet

var detectRecursive bool

func initDetect() {
	detectFlagSet = command.New("detect", detect, detectUsage)
	detectFlagSet.BoolVar(&detectRecursive, "r", true, locale.Sprintf(locale.FlagDetectRecursive))
}

func detect(io.Writer) error {
	p := getPath(detectFlagSet)
	h := core.NewMessageHandler(newHandlerFunc())
	defer h.Stop()

	cfg, err := build.DetectConfig(p, detectRecursive)
	if err != nil {
		return err
	}

	if err = cfg.SaveToFile(path.Join(p, vars.AllowConfigFilenames[0])); err != nil {
		return err
	}

	h.Message(core.Succ, locale.ConfigWriteSuccess, p)
	return nil
}

func detectUsage(w io.Writer) error {
	_, err := fmt.Fprintln(w, locale.Sprintf(locale.CmdDetectUsage, getFlagSetUsage(detectFlagSet)))
	return err
}
