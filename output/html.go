// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"html/template"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/caixw/apidoc/app"
	"github.com/caixw/apidoc/doc"
	"github.com/caixw/apidoc/output/static"
)

// 输出页面的一些自定义项
const (
	groupPrefix = "group_" // 分组文件的前缀
	indexName   = "index"  // 索引文件名
	suffix      = ".html"  // 文件后缀名
)

// 用于页首和页脚的附加信息
type page struct {
	Content        string                // 索引文件的其它内容
	Groups         map[string][]*doc.API // 按组名形式组织的文档集合
	GroupName      string                // 当前分组名称
	Group          []*doc.API            // 当前组的文档集合
	IndexGroup     []*doc.API            // 索引页的文档列表
	Date           string                // 生成日期
	Version        string                // 文档版本
	AppVersion     string                // apidoc 的版本号
	AppName        string                // 程序名称
	AppRepoURL     string                // 仓库地址
	AppOfficialURL string                // 官网地址
	Title          string                // 标题
	Elapsed        time.Duration         // 生成文档所用的时间
}

// 将 docs 的内容以 html 格式输出。
func html(docs *doc.Doc, opt *Options) error {
	p := &page{
		Content:        docs.Content,
		Title:          docs.Title,
		Version:        docs.Version,
		AppVersion:     app.Version,
		AppName:        app.Name,
		AppRepoURL:     app.RepoURL,
		AppOfficialURL: app.OfficialURL,
		Elapsed:        opt.Elapsed,
		Date:           time.Now().Format(time.RFC3339), // TODO 可以自定义时间格式？
		Groups:         make(map[string][]*doc.API, 100),
		IndexGroup:     make([]*doc.API, 0, 100),
	}

	// 按分组名称进行分类
	for _, api := range docs.Apis {
		if len(api.Group) == 0 { // 未指定分组名称，则归类到索引页的文档
			p.IndexGroup = append(p.IndexGroup, api)
			continue
		}

		if p.Groups[api.Group] == nil {
			p.Groups[api.Group] = []*doc.API{}
		}
		p.Groups[api.Group] = append(p.Groups[api.Group], api)
	}

	// 编译模板
	t := template.New("html").
		Funcs(template.FuncMap{
			"groupURL": p.groupURL,
		})

	if len(opt.Template) > 0 { // 自定义模板
		if _, err := t.ParseGlob(path.Join(opt.Template, "*.html")); err != nil {
			return err
		}
	} else { // 系统模板
		for _, content := range static.Templates {
			template.Must(t.Parse(content))
		}
	}

	return p.render(t, opt.Dir)
}

// 根据分组名称获取相应的 url 地址。
func (p *page) groupURL(groupName string) string {
	return path.Join(".", groupPrefix+groupName+suffix)
}

// 根据分组名称获取相应的文件地址。
func (p *page) groupPath(parent, groupName string) string {
	return filepath.Join(parent, groupPrefix+groupName+suffix)
}

// 根据模板 t 将页面输出到 destDir 目录
func (p *page) render(t *template.Template, destDir string) error {
	if err := p.renderIndex(t, destDir); err != nil {
		return err
	}

	if err := p.renderGroup(t, destDir); err != nil {
		return err
	}

	// 输出static
	return static.Output(destDir)
}

// 输出索引页
func (p *page) renderIndex(t *template.Template, destDir string) error {
	index, err := os.Create(filepath.Join(destDir, indexName+suffix))
	if err != nil {
		return err
	}
	defer index.Close()

	return t.ExecuteTemplate(index, "index", p)
}

// 按分组输出内容页
func (p *page) renderGroup(t *template.Template, destDir string) error {
	for name, group := range p.Groups {
		file, err := os.Create(p.groupPath(destDir, name))
		if err != nil {
			return err
		}
		defer file.Close()

		p.GroupName = name
		p.Group = group
		if err = t.ExecuteTemplate(file, "group", p); err != nil {
			return err
		}
	}
	return nil
}
