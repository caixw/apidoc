// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"html/template"
	"os"
	"time"

	"github.com/caixw/apidoc/app"
	"github.com/caixw/apidoc/doc"
	"github.com/caixw/apidoc/output/static"
)

// 用于页首和页脚的附加信息
type page struct {
	Groups         map[string]string // 分组名称与文件的对照表
	CurrGroup      string            // 当前所在的分组页，若为空，表示在列表页
	Date           string            // 生成日期
	Version        string            // 文档版本
	AppVersion     string            // apidoc 的版本号
	AppName        string            // 程序名称
	AppRepoURL     string            // 仓库地址
	AppOfficialURL string            // 官网地址
	Title          string            // 标题
	Elapsed        time.Duration     // 生成文档所用的时间
}

// 将 docs 的内容以 html 格式输出。
func html(docs *doc.Doc, opt *Options) error {
	t := template.New("html")
	for _, content := range static.Templates {
		template.Must(t.Parse(content))
	}

	p := &page{
		Title:          opt.Title,
		Version:        opt.Version,
		AppVersion:     app.Version,
		AppName:        app.Name,
		AppRepoURL:     app.RepoURL,
		AppOfficialURL: app.OfficialURL,
		Elapsed:        opt.Elapsed,
		Date:           time.Now().Format(time.RFC3339),
		Groups:         make(map[string]string, len(docs.Apis)),
	}

	groups := map[string][]*doc.API{}
	for _, v := range docs.Apis {
		p.Groups[v.Group] = "./group_" + v.Group + ".html"
		if groups[v.Group] == nil {
			groups[v.Group] = []*doc.API{}
		}
		groups[v.Group] = append(groups[v.Group], v)
	}

	if err := outputIndex(t, p, opt.Dir); err != nil {
		return err
	}

	if err := outputGroup(groups, t, p, opt.Dir); err != nil {
		return err
	}

	// 输出static
	return static.Output(opt.Dir)
}

// 输出索引页
func outputIndex(t *template.Template, p *page, destDir string) error {
	index, err := os.Create(destDir + "index.html")
	if err != nil {
		return err
	}
	defer index.Close()

	err = t.ExecuteTemplate(index, "header", p)
	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(index, "index", p)
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(index, "footer", p)
}

// 按分组输出内容页
func outputGroup(apis map[string][]*doc.API, t *template.Template, p *page, destDir string) error {
	for k, v := range apis {
		group, err := os.Create(destDir + "group_" + k + ".html")
		if err != nil {
			return err
		}
		defer group.Close()

		p.CurrGroup = k
		err = t.ExecuteTemplate(group, "header", p)
		if err != nil {
			return err
		}
		for _, d := range v {
			err = t.ExecuteTemplate(group, "group", d)
			if err != nil {
				return err
			}
		}
		err = t.ExecuteTemplate(group, "footer", p)
		if err != nil {
			return err
		}
	}
	return nil
}
