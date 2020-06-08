// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"io"

	"github.com/issue9/cmdopt"
	"gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

var (
	detectRecursive bool
	detectWrite     bool
	detectDir       = uri("./")
)

func initDetect(command *cmdopt.CmdOpt) {
	fs := command.New("detect", locale.Sprintf(locale.CmdDetectUsage), detect)
	fs.BoolVar(&detectRecursive, "r", true, locale.Sprintf(locale.FlagDetectRecursive))
	fs.BoolVar(&detectWrite, "w", false, locale.Sprintf(locale.FlagDetectWrite))
	fs.Var(&buildDir, "d", locale.Sprintf(locale.FlagDetectDirUsage))
}

func detect(w io.Writer) error {
	h := core.NewMessageHandler(messageHandle)
	defer h.Stop()

	dir := core.URI(detectDir)
	cfg, err := build.DetectConfig(dir, detectRecursive)
	if err != nil {
		return err
	}

	if !detectWrite {
		data, err := yaml.Marshal(cfg)
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(w, string(data))
		return err
	}

	if err = cfg.Save(dir); err != nil {
		return err
	}
	h.Locale(core.Succ, locale.ConfigWriteSuccess, dir)
	return nil
}
