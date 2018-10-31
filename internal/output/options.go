// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"encoding/json"

	"golang.org/x/text/message"
	yaml "gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/doc"
	"github.com/caixw/apidoc/internal/errors"
	"github.com/caixw/apidoc/internal/locale"
	"github.com/caixw/apidoc/internal/output/openapi"
	"github.com/caixw/apidoc/options"
)

type marshaler func(v *doc.Doc) ([]byte, error)

// Options 输出配置项
type Options struct {
	options.Output
	marshal marshaler
}

func newError(field string, key message.Reference, args ...interface{}) *errors.Error {
	return &errors.Error{
		Field:       field,
		MessageKey:  key,
		MessageArgs: args,
	}
}

func (o *Options) contains(tags ...string) bool {
	if len(o.Tags) == 0 {
		return true
	}

	for _, t := range o.Tags {
		for _, tag := range tags {
			if tag == t {
				return true
			}
		}
	}
	return false
}

// Sanitize 对 Options 作一些初始化操作。
func (o *Options) Sanitize() error {
	// TODO 改用默认值
	if o.Path == "" {
		return newError("path", locale.ErrRequired)
	}

	if o.Type == "" {
		o.Type = options.ApidocJSON
	}

	switch o.Type {
	case options.ApidocJSON:
		o.marshal = apidocJSONMarshal
	case options.ApidocYAML:
		o.marshal = apidocYAMLMarshal
	case options.OpenapiJSON:
		o.marshal = openapi.JSON
	case options.OpenapiYAML:
		o.marshal = openapi.YAML
	case options.RamlJSON:
		// TODO
	default:
		return newError("type", locale.ErrInvalidValue)
	}

	return nil
}

func apidocJSONMarshal(v *doc.Doc) ([]byte, error) {
	return json.Marshal(v)
}

func apidocYAMLMarshal(v *doc.Doc) ([]byte, error) {
	return yaml.Marshal(v)
}
