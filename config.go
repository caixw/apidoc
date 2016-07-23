// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/caixw/apidoc/app"
	"github.com/caixw/apidoc/app/locale"
	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/output"
	"github.com/issue9/version"
)

// 项目的配置内容，分别引用到了 input.Options 和 output.Options。
//
// 所有可能改变输出的表现形式的，应该添加此配置中；
// 而如果只是改变输出内容的，应该直接以标签的形式出现在代码中，
// 比如文档的版本号、标题等，都是直接使用 `@apidoc` 来指定的。
type config struct {
	Version string           `json:"version"` // 产生此配置文件的程序版本号
	Inputs  []*input.Options `json:"inputs"`
	Output  *output.Options  `json:"output"`
}

// 带色彩输出的控制台。
type syntaxWriter struct {
}

func (w *syntaxWriter) Write(bs []byte) (int, error) {
	app.Error().Print(string(bs))
	return len(bs), nil
}

// 加载 path 所指的文件内容到 *config 实例。
func loadConfig(path string) (*config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &config{}
	if err = json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if !version.SemVerValid(cfg.Version) {
		return nil, &app.OptionsError{Field: "version", Message: app.Sprintf(locale.ErrInvalidFormat)}
	}

	if len(cfg.Inputs) == 0 {
		return nil, &app.OptionsError{Field: "inputs", Message: app.Sprintf(locale.ErrRequired)}
	}

	if cfg.Output == nil {
		return nil, &app.OptionsError{Field: "output", Message: app.Sprintf(locale.ErrRequired)}
	}

	l := log.New(&syntaxWriter{}, "", 0)
	for i, opt := range cfg.Inputs {
		index := strconv.Itoa(i)
		if err := opt.Init(); err != nil {
			err.Field = "inputs[" + index + "]." + err.Field
			return nil, err
		}
		opt.SyntaxLog = l
	}

	if err := cfg.Output.Init(); err != nil {
		err.Field = "outputs." + err.Field
		return nil, err
	}

	return cfg, nil
}

// 生成一个默认的配置文件，并写入到 path 中。
func genConfigFile(path string) error {
	dir := filepath.Dir(path)
	lang, err := input.DetectDirLang(dir)
	if err != nil { // 不中断，仅作提示用。
		app.Warn().Print(err)
	}

	cfg := &config{
		Version: app.Version,
		Inputs: []*input.Options{
			&input.Options{
				Dir:       dir,
				Recursive: true,
				Lang:      lang,
			},
		},
		Output: &output.Options{
			Type: "html",
			Dir:  filepath.Join(dir, "doc"),
		},
	}
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	fi, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fi.Close()

	_, err = fi.Write(data)
	return err
}
