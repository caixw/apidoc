// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"io"

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

var (
	detectFlagSet   *flag.FlagSet
	detectRecursive bool
	detectDir       = uri("./")
)

func initDetect() {
	detectFlagSet = command.New("detect", locale.Sprintf(locale.CmdDetectUsage), detect)
	detectFlagSet.BoolVar(&detectRecursive, "r", true, locale.Sprintf(locale.FlagDetectRecursive))
	detectFlagSet.Var(&buildDir, "d", locale.Sprintf(locale.FlagDetectDirUsage))
}

func detect(io.Writer) error {
	h := core.NewMessageHandler(messageHandle)
	defer h.Stop()

	path := core.URI(detectDir)
	cfg, err := build.DetectConfig(path, detectRecursive)
	if err != nil {
		return err
	}

	if err = cfg.Save(path); err != nil {
		return err
	}

	h.Locale(core.Succ, locale.ConfigWriteSuccess, path)
	return nil
}
