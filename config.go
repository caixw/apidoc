// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/caixw/apidoc/app"
	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/output"
)

// 项目的配置内容，分别引用到了 input.Options 和 output.Options
// 所有可能改变输出的表现形式的，应该添加到 config 中；
// 而如果只是改变输出内容的，则应该直接以标签的形式出现在代码中，
// 比如文档的版本号、标题等，都是直接使用 `@apidoc`来指定的，而不是出现在配置文件中。
type config struct {
	Version string          `json:"version"` // 产生该配置文件的程序版本号，主版本号不同，表示不兼容
	Input   *input.Options  `json:"input"`
	Output  *output.Options `json:"output"`
}

// 从配置文件中加载配置项。
func loadConfig() (*config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filepath.Join(wd, app.ConfigFilename))
	if err != nil {
		return nil, err
	}

	cfg := &config{}
	if err = json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if err := cfg.Input.Init(); err != nil {
		return nil, err
	}

	if err := cfg.Output.Init(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// 在当前目录下产生个默认的配置文件。
func genConfigFile() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	fi, err := os.Create(filepath.Join(wd, app.ConfigFilename))
	if err != nil {
		return err
	}
	defer fi.Close()

	lang, err := input.DetectDirLang(wd)
	if err != nil { // 不中断，仅作提示用。
		app.Warn(err)
	}

	cfg := &config{
		Version: app.Version,
		Input: &input.Options{
			Dir:       "./",
			Recursive: true,
			Lang:      lang,
		},
		Output: &output.Options{
			Type: "html",
			Dir:  "./apidoc",
		},
	}
	data, err := json.MarshalIndent(cfg, "", "    ")
	_, err = fi.Write(data)
	return err
}
