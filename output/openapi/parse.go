// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package openapi

import (
	"encoding/json"

	yaml "gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/docs"
)

func parse(docs *docs.Docs) (*OpenAPI, error) {
	panic("该功能未实现")
	return nil, nil
}

// JSON 输出 JSON 格式数据
func JSON(docs *docs.Docs) ([]byte, error) {
	openapi, err := parse(docs)
	if err != nil {
		return nil, err
	}

	return json.Marshal(openapi)
}

// YAML 输出 YAML 格式数据
func YAML(docs *docs.Docs) ([]byte, error) {
	openapi, err := parse(docs)
	if err != nil {
		return nil, err
	}

	return yaml.Marshal(openapi)
}
