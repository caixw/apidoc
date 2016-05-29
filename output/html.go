// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"html/template"
	"os"
	"strconv"
	"time"

	"github.com/caixw/apidoc/doc"
	"github.com/caixw/apidoc/output/static"
)

// 用于页首和页脚的附加信息
type info struct {
	Groups     map[string]string // 分组名称与文件的对照表
	CurrGroup  string            // 当前所在的分组页，若为空，表示在列表页
	Date       string
	Version    string
	AppVersion string
	Title      string
	Elapsed    string
}

// 将docs的内容以html格式输出。
func html(docs *doc.Doc, opt *Options) error {
	t := template.New("core")
	for _, content := range static.Templates {
		template.Must(t.Parse(content))
	}

	i := &info{
		Title:      opt.Title,
		Version:    opt.Version,
		AppVersion: opt.AppVersion,
		Elapsed:    strconv.FormatFloat(float64(opt.Elapsed)/1000000, 'f', 2, 32),
		Date:       time.Now().Format(time.RFC3339),
		Groups:     make(map[string]string, len(docs.Apis)),
	}

	groups := map[string][]*doc.API{}
	for _, v := range docs.Apis {
		i.Groups[v.Group] = "./group_" + v.Group + ".html"
		if groups[v.Group] == nil {
			groups[v.Group] = []*doc.API{}
		}
		groups[v.Group] = append(groups[v.Group], v)
	}

	if err := outputIndex(t, i, opt.Dir); err != nil {
		return err
	}

	if err := outputGroup(groups, t, i, opt.Dir); err != nil {
		return err
	}

	// 输出static
	return static.Output(opt.Dir)
}

// 输出索引页
func outputIndex(t *template.Template, i *info, destDir string) error {
	index, err := os.Create(destDir + "index.html")
	if err != nil {
		return err
	}
	defer index.Close()

	err = t.ExecuteTemplate(index, "header", i)
	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(index, "index", i)
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(index, "footer", i)
}

// 按分组输出内容页
func outputGroup(apis map[string][]*doc.API, t *template.Template, i *info, destDir string) error {
	for k, v := range apis {
		group, err := os.Create(destDir + "group_" + k + ".html")
		if err != nil {
			return err
		}
		defer group.Close()

		i.CurrGroup = k
		err = t.ExecuteTemplate(group, "header", i)
		if err != nil {
			return err
		}
		for _, d := range v {
			err = t.ExecuteTemplate(group, "group", d)
			if err != nil {
				return err
			}
		}
		err = t.ExecuteTemplate(group, "footer", i)
		if err != nil {
			return err
		}
	}
	return nil
}
