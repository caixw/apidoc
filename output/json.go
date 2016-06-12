// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	j "encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/caixw/apidoc/doc"
)

type jsonData struct {
	Title     string        `json:"title"`
	Version   string        `json:"version"`
	Date      time.Time     `json:"date"`
	Elapsed   time.Duration `json:"elapsed"`
	Content   string        `json:"content"`
	GroupName string        `json:"groupName"`
	Apis      []*doc.API    `json:"apis"`
}

func json(docs *doc.Doc, opt *Options) error {
	indexGroup := make([]*doc.API, 0, 100)
	groups := make(map[string][]*doc.API, 100)

	for _, api := range docs.Apis {
		if len(api.Group) == 0 { // 未指定分组名称，则归类到索引页的文档
			indexGroup = append(indexGroup, api)
			continue
		}

		if groups[api.Group] == nil {
			groups[api.Group] = []*doc.API{}
		}
		groups[api.Group] = append(groups[api.Group], api)
	}

	page := &jsonData{
		Title:   docs.Title,
		Version: docs.Version,
		Date:    time.Now(),
		Elapsed: opt.Elapsed,
		Content: docs.Content,
	}
	if len(indexGroup) > 0 {
		file, err := os.Create(filepath.Join(opt.Dir, indexName+".json"))
		if err != nil {
			return err
		}
		page.Apis = indexGroup
		page.GroupName = ""
		data, err := j.MarshalIndent(page, "", "    ")
		if err != nil {
			return err
		}
		if _, err := file.Write(data); err != nil {
			return err
		}
	}

	for name, apis := range groups {
		file, err := os.Create(filepath.Join(opt.Dir, groupPrefix+name+".json"))
		if err != nil {
			return err
		}

		page.Apis = apis
		page.GroupName = name
		data, err := j.MarshalIndent(page, "", "    ")
		if err != nil {
			return err
		}
		if _, err := file.Write(data); err != nil {
			return err
		}
	}

	return nil
}
