// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"

	"github.com/issue9/logs/writers"
	"github.com/issue9/term/colors"

	"github.com/caixw/apidoc/internal/locale"
)

// 控制台的输出颜色
const (
	infoColor = colors.Green
	warnColor = colors.Cyan
	erroColor = colors.Red
)

// 日志信息输出
var (
	info = newLog(os.Stdout, infoColor, "[INFO] ")
	warn = newLog(os.Stderr, warnColor, "[WARN] ")
	erro = newLog(os.Stderr, erroColor, "[ERRO] ")
)

func initLogsLocale() {
	info.SetPrefix(locale.Sprintf(locale.InfoPrefix))
	warn.SetPrefix(locale.Sprintf(locale.WarnPrefix))
	erro.SetPrefix(locale.Sprintf(locale.ErrorPrefix))
}

func newLog(out *os.File, color colors.Color, prefix string) *log.Logger {
	return log.New(writers.NewConsole(out, color, colors.Default), prefix, 0)
}
