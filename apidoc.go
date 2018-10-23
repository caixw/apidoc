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
	"github.com/caixw/apidoc/input"
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

// Parse 分析代码并输出
//
// erro 用于输出语法日志错误内容；
// output 输出设置项；
// input 输入设置项。
func Parse(erro *log.Logger, output *o.Options, input ...*input.Options) error {
	if len(input) == 0 {
		return errors.New("参数 input 不能为空")
	}

	if output == nil {
		return errors.New("参数 output 不能为空")
	}

	for _, opt := range input {
		if err := opt.Sanitize(); err != nil {
			return err
		}
	}

	if err := output.Sanitize(); err != nil {
		return err
	}

	start := time.Now()
	docs, err := docs.Parse(erro, input...)
	if err != nil {
		return err
	}

	output.Elapsed = time.Now().Sub(start)
	return o.Render(docs, output)
}
