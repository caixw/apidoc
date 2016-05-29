// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/output"
)

type config struct {
	Version string          `json:"version"` // 兼容的 apidoc 版本
	Input   *input.Options  `json:"input"`
	Output  *output.Options `json:"output"`
}

// 从配置文件中加载配置项。
func loadConfig() (*config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(wd + "/" + configFilename)
	if err != nil {
		return nil, err
	}

	cfg := &config{}
	if err = json.Unmarshal(data, cfg); err != nil {
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

	path := wd + string(os.PathSeparator) + configFilename
	fi, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fi.Close()

	cfg := &config{
		Input:  &input.Options{Dir: "./", Recursive: true},
		Output: &output.Options{},
	}
	data, err := json.MarshalIndent(cfg, "", "    ")
	_, err = fi.Write(data)
	return err
}
