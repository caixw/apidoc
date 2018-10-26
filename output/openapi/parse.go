// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package openapi

import (
	"encoding/json"

	yaml "gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/doc"
)

func parse(doc *doc.Doc) (*OpenAPI, error) {
	openapi := &OpenAPI{
		OpenAPI: doc.APIDoc,
		Info: &Info{
			Title:       doc.Title,
			Description: doc.Content,
			Contact:     newContact(doc.Contact),
			License:     newLicense(doc.License),
			Version:     doc.Version,
		},
		Servers: make([]*Server, 0, len(doc.Servers)),
		Tags:    make([]*Tag, 0, len(doc.Tags)),

		// TODO Paths
	}

	for _, srv := range doc.Servers {
		openapi.Servers = append(openapi.Servers, newServer(srv))
	}

	for _, tag := range doc.Tags {
		openapi.Tags = append(openapi.Tags, newTag(tag))
	}

	return openapi, nil
}

// JSON 输出 JSON 格式数据
func JSON(doc *doc.Doc) ([]byte, error) {
	openapi, err := parse(doc)
	if err != nil {
		return nil, err
	}

	return json.Marshal(openapi)
}

// YAML 输出 YAML 格式数据
func YAML(doc *doc.Doc) ([]byte, error) {
	openapi, err := parse(doc)
	if err != nil {
		return nil, err
	}

	return yaml.Marshal(openapi)
}
