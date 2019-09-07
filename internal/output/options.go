// SPDX-License-Identifier: MIT

package output

import (
	"encoding/xml"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/errors"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/output/openapi"
	opt "github.com/caixw/apidoc/v5/options"
)

type marshaler func(v *doc.Doc) ([]byte, error)

type options struct {
	opt.Output
	marshal marshaler
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

func buildOptions(o *opt.Output) (*options, *errors.Error) {
	if o.Path == "" {
		return nil, errors.New("", "path", 0, locale.ErrRequired)
	}

	if o.Type == "" {
		o.Type = opt.ApidocXML
	}

	var marshal marshaler
	switch o.Type {
	case opt.ApidocXML:
		marshal = xmlMarshal
	case opt.OpenapiJSON:
		marshal = openapi.JSON
	case opt.OpenapiYAML:
		marshal = openapi.YAML
	case opt.RAMLYAML:
		// TODO
	default:
		return nil, errors.New("", "type", 0, locale.ErrInvalidValue)
	}

	return &options{
		Output:  *o,
		marshal: marshal,
	}, nil
}

func xmlMarshal(v *doc.Doc) ([]byte, error) {
	return xml.Marshal(v)
}
