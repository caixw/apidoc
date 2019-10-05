// SPDX-License-Identifier: MIT

// Package apidoc RESTful API 文档生成工具
//
// 可以从代码文件的注释中提取文档内容，生成 API 文档，
// 支持大部分的主流的编程语言。
package apidoc

import (
	"io/ioutil"
	"net/http"
	"path"

	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v5/input"
	"github.com/caixw/apidoc/v5/internal/html"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
	"github.com/caixw/apidoc/v5/output"
)

// Init 初始化包
//
// 如果传递了 language.Und，则采用系统当前的本地化信息。
// 如果获取系统的本地化信息依然失败，则会失放 zh-Hans 作为默认值。
func Init(tag language.Tag) error {
	return locale.Init(tag)
}

// Version 当前程序的版本号
//
// 为一个正常的 semver(https://semver.org/lang/zh-CN/) 格式字符串。
func Version() string {
	return vars.Version()
}

// Do 解析文档并输出文档内容
//
// 如果需要控制详细的操作步骤，可以自行调用 input 和 output 的相关函数实现。
//
// 如果是文档语法错误，则相关的错误信息会反馈给 h，由 h 处理错误信息；
// 如果是配置项（o 和 i）有问题，则以 *message.SyntaxError 类型返回错误信息。
func Do(h *message.Handler, o *output.Options, i ...*input.Options) error {
	doc, err := input.Parse(h, i...)
	if err != nil {
		return err
	}

	return output.Render(doc, o)
}

// Handle 处理 apidoc 相关的依赖文件
//
// p 指定了 apidoc.xml 实际的文件路径；
func Handle(p string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		switch path.Base(r.URL.Path) {
		case "apidoc.xsl":
			w.Header().Set("Content-Type", "text/xsl")
			_, err = w.Write(html.XSL)
		case "apidoc.js":
			w.Header().Set("Content-Type", "application/javascript")
			_, err = w.Write(html.JS)
		case "apidoc.css":
			w.Header().Set("Content-Type", "text/css")
			_, err = w.Write(html.CSS)
		case "apidoc.xml":
			// TODO 替换掉 xml-stylesheet 为当前的 xsl
			data, err := ioutil.ReadFile(p)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "text/xml")
			_, err = w.Write(data)
		default:
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
		if err != nil {
			// TODO
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	})
}
