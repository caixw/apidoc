// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package apidoc RESTful API 文档生成工具。
package apidoc

import (
	"errors"
	"log"
	"time"

	"golang.org/x/text/language"

	"github.com/caixw/apidoc/docs"
	i "github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/internal/locale"
	"github.com/caixw/apidoc/internal/vars"
	o "github.com/caixw/apidoc/output"
)

// InitLocale 初始化语言环境
//
// NOTE: 必须保证在第一时间调用。
//
// 如果 tag 的值为 language.Und，则表示采用系统语言
func InitLocale(tag language.Tag) error {
	return locale.Init(tag)
}

// Version 获取版本号
func Version() string {
	return vars.Version()
}

// Do 分析输入信息，并最终输出到指定的文件。
//
// erro 用于输出语法错误内容；
// output 输出设置项；
// input 输入设置项。
func Do(erro *log.Logger, output *o.Options, input ...*i.Options) error {
	if output == nil {
		return errors.New("参数 output 不能为空")
	}
	if err := output.Sanitize(); err != nil {
		return err
	}

	docs, err := Parse(erro, input...)
	if err != nil {
		return err
	}

	return o.Render(docs, output)
}

// Parse 分析输入信息，并获取 docs.Docs 实例。
//
// erro 用于输出语法错误内容；
// input 输入设置项。
func Parse(erro *log.Logger, input ...*i.Options) (*docs.Docs, error) {
	if len(input) == 0 {
		return nil, errors.New("参数 input 不能为空")
	}

	for _, opt := range input {
		if err := opt.Sanitize(); err != nil {
			return nil, err
		}
	}

	start := time.Now()
	block := i.Parse(erro, input...)
	docs := docs.Parse(erro, block)
	docs.Elapsed = time.Now().Sub(start)

	return docs, nil
}
