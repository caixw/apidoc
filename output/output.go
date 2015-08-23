// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"html/template"
	"os"
	"strconv"
	"time"

	"github.com/caixw/apidoc/core"
	"github.com/caixw/apidoc/output/static"
)

// 将docs的内容以html格式输出。
func Html(docs []*core.Doc, opt *Options) error {
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
		Groups:     make(map[string]string, len(docs)),
	}

	groups := map[string][]*core.Doc{}
	for _, v := range docs {
		i.Groups[v.Group] = "./group_" + v.Group + ".html"
		if groups[v.Group] == nil {
			groups[v.Group] = []*core.Doc{}
		}
		groups[v.Group] = append(groups[v.Group], v)
	}

	if err := outputIndex(t, i, opt.DocDir); err != nil {
		return err
	}

	if err := outputGroup(groups, t, i, opt.DocDir); err != nil {
		return err
	}

	// 输出static
	return static.Output(opt.DocDir)
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
func outputGroup(docs map[string][]*core.Doc, t *template.Template, i *info, destDir string) error {
	for k, v := range docs {
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
