// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"encoding/json"

	"golang.org/x/text/message"
	yaml "gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/doc"
	"github.com/caixw/apidoc/errors"
	"github.com/caixw/apidoc/internal/locale"
	"github.com/caixw/apidoc/internal/output/openapi"
	opt "github.com/caixw/apidoc/options"
)

type marshaler func(v *doc.Doc) ([]byte, error)

type options struct {
	opt.Output
	marshal marshaler
}

func newError(field string, key message.Reference, args ...interface{}) *errors.Error {
	return &errors.Error{
		Field: field,
		LocaleError: errors.LocaleError{
			MessageKey:  key,
			MessageArgs: args,
		},
	}
}

func (o *options) contains(tags ...string) bool {
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

func buildOptions(o *opt.Output) (*options, error) {
	// TODO 改用默认值
	if o.Path == "" {
		return nil, newError("path", locale.ErrRequired)
	}

	if o.Type == "" {
		o.Type = opt.ApidocJSON
	}

	var marshal marshaler
	switch o.Type {
	case opt.ApidocJSON:
		marshal = apidocJSONMarshal
	case opt.ApidocYAML:
		marshal = apidocYAMLMarshal
	case opt.OpenapiJSON:
		marshal = openapi.JSON
	case opt.OpenapiYAML:
		marshal = openapi.YAML
	case opt.RamlJSON:
		// TODO
	default:
		return nil, newError("type", locale.ErrInvalidValue)
	}

	return &options{
		Output:  *o,
		marshal: marshal,
	}, nil
}

func apidocJSONMarshal(v *doc.Doc) ([]byte, error) {
	return json.Marshal(v)
}

func apidocYAMLMarshal(v *doc.Doc) ([]byte, error) {
	return yaml.Marshal(v)
}
