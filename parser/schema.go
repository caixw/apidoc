// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package parser

import (
	"bytes"

	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/openapi"
)

// @param list.groups array.string optional desc markdown
//  * xx: xxxxx
//  * xx: xxxxx
func buildSchema(schema *openapi.Schema, name, typ, optional, desc []byte) *syntaxError {
	type0, type1, err := parseType(typ)
	if err != nil {
		return err
	}

	var p *openapi.Schema
	var last []byte // 最后的名称
	if len(name) > 0 {
		names := bytes.Split(name, seqaratorDot)
		for _, name := range names {
			if schema.Properties == nil {
				schema.Properties = make(map[string]*openapi.Schema, 2)
			}

			ss := schema.Properties[string(name)]
			if ss == nil {
				ss = &openapi.Schema{}
				schema.Properties[string(name)] = ss
			}
			p = schema
			last = name
			schema = ss
		}
	}

	schema.Type = type0
	schema.Description = openapi.Description(desc)
	if type0 == "array" {
		schema.Items = &openapi.Schema{Type: type1}
	}

	if p != nil && isRequired(string(optional)) {
		if p.Required == nil {
			p.Required = make([]string, 0, 10)
		}
		p.Required = append(p.Required, string(last))
	}

	return nil
}

func parseType(typ []byte) (t1, t2 string, err *syntaxError) {
	types := bytes.SplitN(typ, seqaratorDot, 2)
	if len(types) == 0 {
		return "", "", &syntaxError{MessageKey: locale.ErrInvalidFormat}
	}

	type0 := string(types[0])
	if type0 != "array" {
		return type0, "", nil
	}

	if len(types) == 1 {
		return "", "", &syntaxError{MessageKey: locale.ErrInvalidFormat}
	}

	return type0, string(types[1]), nil
}

func isRequired(optional string) bool {
	switch optional {
	case "required":
		return true
	case "optional":
		return false
	default:
		return false
	}
}

func parseEnum(desc []byte) []string {
	// TODO
	return nil
}
