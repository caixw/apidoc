// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/caixw/apidoc/app"
	"github.com/caixw/apidoc/doc"
)

type page struct {
	Title       string            `json:"title"`
	Version     string            `json:"version,omitempty"`
	BaseURL     string            `json:"baseURL"`
	LicenseName string            `json:"licenseName"`
	LicenseURL  string            `json:"licenseURL"`
	Content     string            `json:"content,omitempty"`
	Date        time.Time         `json:"date"`
	Elapsed     time.Duration     `json:"elapsed"`
	Groups      map[string]string `json:"groups"` // 组名与文件名的对应关系
}

type group struct {
	path string // 相对路径名

	Name string     `json:"groupName"` // 当前分组的名称
	Apis []*doc.API `json:"apis"`      // 当前分组的 api 文档
}

func render(docs *doc.Doc, opt *Options) error {
	groups := make(map[string]*group, 100)

	for _, api := range docs.Apis {
		name := strings.ToLower(api.Group)

		if groups[name] == nil {
			groups[name] = &group{
				path: filepath.Join(opt.Dir, name),
				Name: api.Group, // 名称区分大小写，不采用 name 变量
				Apis: make([]*doc.API, 0, 100),
			}
		}
		groups[name].Apis = append(groups[name].Apis, api)
	}

	names := make(map[string]string, len(groups))
	for _, group := range groups {
		names[group.Name] = group.path
	}

	page := &page{
		Title:       docs.Title,
		Version:     docs.Version,
		BaseURL:     docs.BaseURL,
		LicenseName: docs.LicenseName,
		LicenseURL:  docs.LicenseURL,
		Content:     docs.Content,
		Date:        time.Now(),
		Elapsed:     opt.Elapsed,
		Groups:      names,
	}

	if err := renderPage(page, opt.Dir); err != nil {
		return err
	}

	return renderGroups(groups, opt.Dir)
}

func renderPage(p *page, destDir string) error {
	destDir = filepath.Join(destDir, app.JSONDataDirName)

	path := filepath.Join(destDir, "page.json")
	if err := renderJSON(p, path); err != nil {
		return err
	}

	return nil
}

func renderGroups(groups map[string]*group, destDir string) error {
	destDir = filepath.Join(destDir, app.JSONDataDirName)

	for _, g := range groups {
		if err := renderJSON(g, g.path); err != nil {
			return err
		}
	}

	return nil
}

func renderJSON(obj interface{}, path string) error {
	data, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, os.ModePerm)
}
