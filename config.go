// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const configFilename = ".apidoc.json"

type config struct {
	Input  *input  `json:"input"`
	Output *output `json:"output"`
	Doc    *doc    `json:"doc"`

	lang *lang
}

type input struct {
	Type      string   `json:"type"`      // 输入的目标语言
	Dir       string   `json:"dir"`       // 源代码目录
	Exts      []string `json:"exts"`      // 需要扫描的文件扩展名
	Recursive bool     `json:"recursive"` // 是否查找Dir的子目录
}

type output struct {
	Dir string `json:"dir"`
	//Type string   `json:"type"` // 输出的语言格式
	//Groups     []string `json:"groups"`     // 需要打印的分组内容。
	//Timezone   string   `json:"timezone"`   // 时区
}

type doc struct {
	Version string `json:"version"` // 文档版本号
	Title   string `json:"title"`   // 文档的标题，默认为apidoc
	BaseURL string `json:"baseURL"` // api文档中url的前缀，不指定，则为空
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

	if err = initConfig(wd, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func initConfig(wd string, cfg *config) error {
	if len(cfg.Doc.Title) == 0 {
		cfg.Doc.Title = "apidoc"
	}

	if len(cfg.Input.Dir) == 0 {
		cfg.Input.Dir = wd
	}
	cfg.Input.Dir += string(os.PathSeparator)

	if len(cfg.Input.Exts) > 0 {
		exts := make([]string, 0, len(cfg.Input.Exts))
		for _, ext := range cfg.Input.Exts {
			if len(ext) == 0 {
				continue
			}

			if ext[0] != '.' {
				ext = "." + ext
			}
			exts = append(exts, ext)
		}
	}

	// 若没有指定Type，则根据exts和当前目录下的文件检测来确定其值
	if len(cfg.Input.Type) == 0 {
		var err error
		if len(cfg.Input.Exts) == 0 {
			cfg.Input.Type, err = detectDirLangType(cfg.Input.Dir)
		} else {
			cfg.Input.Type, err = detectLangType(cfg.Input.Exts)
		}

		if err != nil {
			return err
		}
	}
	cfg.Input.Type = strings.ToLower(cfg.Input.Type)

	l, found := langs[cfg.Input.Type]
	if !found {
		return fmt.Errorf("暂不支持该类型[%v]的语言", cfg.Input.Type)
	}
	cfg.lang = l

	if len(cfg.Input.Exts) == 0 {
		cfg.Input.Exts = l.exts
	}

	return nil
}

// 从扩展名检测其所属的语言名称。
// 以第一个匹配extsIndex的文件扩展名为准。
func detectLangType(exts []string) (string, error) {
	for _, ext := range exts {
		if lang, found := extsIndex[ext]; found {
			return lang, nil
		}
	}
	return "", fmt.Errorf("无法找到与这些扩展名[%v]相匹配的代码扫描函数", exts)
}

// 检测目录下的文件类型。
// 以第一个匹配extsIndex的文件扩展名为准。
func detectDirLangType(dir string) (string, error) {
	var lang string

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() || len(lang) > 0 {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		lang, _ = extsIndex[ext]
		return nil
	}

	if err := filepath.Walk(dir, walk); err != nil {
		return "", err
	}

	if len(lang) == 0 {
		return lang, fmt.Errorf("无法检测到[%v]目录下的文件类型", dir)
	}

	return lang, nil
}

// 根据recursive值确定是否递归查找paths每个目录下的子目录。
func recursivePath(dir string, recursive bool, exts ...string) ([]string, error) {
	paths := []string{}
	dir += string(os.PathSeparator)

	extIsEnabled := func(ext string) bool {
		for _, v := range exts {
			if ext == v {
				return true
			}
		}
		return false
	}

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() && !recursive && path != dir {
			return filepath.SkipDir
		} else if extIsEnabled(filepath.Ext(path)) {
			paths = append(paths, path)
		}
		return nil
	}

	if err := filepath.Walk(dir, walk); err != nil {
		return nil, err
	}

	return paths, nil
}
