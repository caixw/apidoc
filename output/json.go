// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/caixw/apidoc/doc"
)

type jsonPage struct {
	Title     string        `json:"title"`
	Version   string        `json:"version,omitempty"`
	Content   string        `json:"content,omitempty"`
	Date      time.Time     `json:"date"`
	Elapsed   time.Duration `json:"elapsed"`
	GroupName string        `json:"groupName"` // 当前分组的名称
	Apis      []*doc.API    `json:"apis"`      // 当前分组的 api 文档
}

func renderJSON(docs *doc.Doc, opt *Options) error {
	groups := make(map[string][]*doc.API, 100)
	for _, api := range docs.Apis {
		if groups[api.Group] == nil {
			groups[api.Group] = []*doc.API{}
		}
		groups[api.Group] = append(groups[api.Group], api)
	}

	page := &jsonPage{
		Title:   docs.Title,
		Version: docs.Version,
		Date:    time.Now(),
		Elapsed: opt.Elapsed,
		Content: docs.Content,
	}

	return renderJSONGroups(page, groups, opt.Dir)
}

func renderJSONGroups(p *jsonPage, groups map[string][]*doc.API, destDir string) error {
	for p.GroupName, p.Apis = range groups {
		path := filepath.Join(destDir, p.GroupName+".json")
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		data, err := json.MarshalIndent(p, "", "    ")
		if err != nil {
			return err
		}

		if _, err = file.Write(data); err != nil {
			return err
		}
	}
	return nil
}
