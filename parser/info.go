// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package parser

import (
	"bytes"

	"github.com/issue9/is"
	"github.com/issue9/version"

	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/openapi"
)

// @apidoc 的格式下如：
//
// @apidoc title of doc
// @group g1
// @tag t1 desc
// @tag t2 desc
// @license name url
// @contact name url
// @contact name email
// @version v1
// @terms url
// @server name url
// @server name url
// @external name url
// @description description markdown

type info struct {
	title       string
	group       string
	tags        []*openapi.Tag
	license     *openapi.License
	contracts   []*openapi.Contact
	version     string
	terms       string
	servers     []*openapi.Server
	description openapi.Description
	externaldoc *openapi.ExternalDocumentation
}

func (p *parser) parseAPIDoc(l *lexer) error {
	i := &info{}
	for tag, eof := l.tag(); !eof; tag, eof = l.tag() {
		switch string(bytes.ToLower(tag.name)) {
		case "@apidoc":
			if len(tag.data) == 0 {
				return tag.syntaxError(locale.Sprintf(locale.ErrTagArgNotEnough, "@apidoc"))
			}
			if i.title != "" {
				return tag.syntaxError(locale.Sprintf(locale.ErrDuplicateTag, "@apidoc"))
			}
			i.title = string(tag.data)
		case "@apigroup":
			if i.group != "" {
				return tag.syntaxError(locale.Sprintf(locale.ErrDuplicateTag, "@apiGroup"))
			}
			i.group = string(tag.data)
		case "@apitag":
			data := split(tag.data, 2)
			if len(data) != 2 {
				return tag.syntaxError(locale.Sprintf(locale.ErrInvalidFormat, "@apiTag"))
			}
			if i.tags == nil {
				i.tags = make([]*openapi.Tag, 0, 10)
			}
			i.tags = append(i.tags, &openapi.Tag{
				Name:        string(data[0]),
				Description: openapi.Description(data[1]),
			})
		case "@apilicense":
			if i.license != nil {
				return tag.syntaxError(locale.Sprintf(locale.ErrDuplicateTag, "@apiLicense"))
			}

			data := split(tag.data, 2)
			if len(data) != 2 {
				return tag.syntaxError(locale.Sprintf(locale.ErrInvalidFormat, "@apiLicense"))
			}
			if !is.URL(data[1]) {
				return tag.syntaxError(locale.Sprintf(locale.ErrInvalidFormat, "@apiLicense"))
			}
			i.license = &openapi.License{
				Name: string(data[0]),
				URL:  string(data[1]),
			}
		case "@apicontract":
			if err := i.parseContract(tag); err != nil {
				return err
			}
		case "@apiversion":
			if i.version != "" {
				return tag.syntaxError(locale.Sprintf(locale.ErrDuplicateTag, "@apiVersion"))
			}
			i.version = string(tag.data)

			if !version.SemVerValid(i.version) {
				return tag.syntaxError(locale.Sprintf(locale.ErrInvalidFormat, "@apiVersion"))
			}
		case "@apiterms":
			if i.terms != "" {
				return tag.syntaxError(locale.Sprintf(locale.ErrDuplicateTag, "@apiTerms"))
			}
			i.terms = string(tag.data)
		case "@apiservers":
			data := split(tag.data, 2)
			if len(data) != 2 {
				return tag.syntaxError(locale.Sprintf(locale.ErrInvalidFormat, "@apiServer"))
			}
			if i.servers == nil {
				i.servers = make([]*openapi.Server, 0, 10)
			}
			i.servers = append(i.servers, &openapi.Server{
				URL:         string(data[0]),
				Description: openapi.Description(data[1]),
			})
		case "@apidescription":
			if i.description == "" {
				return tag.syntaxError(locale.Sprintf(locale.ErrDuplicateTag, "@apiDescription"))
			}
			i.description = openapi.Description(tag.data)
		case "@apiexternaldoc":
			if i.externaldoc != nil {
				return tag.syntaxError(locale.Sprintf(locale.ErrDuplicateTag, "@apiExternalDoc"))
			}

			data := split(tag.data, 2)
			if len(data) != 2 {
				return tag.syntaxError(locale.Sprintf(locale.ErrInvalidFormat, "@apiExternalDoc"))
			}

			if !is.URL(data[0]) {
				return tag.syntaxError(locale.Sprintf(locale.ErrInvalidFormat, "@apiExternalDoc"))
			}

			i.externaldoc = &openapi.ExternalDocumentation{
				URL:         string(data[0]),
				Description: openapi.Description(data[1]),
			}
		default:
			return tag.syntaxError(locale.Sprintf(locale.ErrInvalidTag, string(tag.name)))
		}
	}

	// TODO p.getDoc(i.group).
	return nil
}

func (i *info) parseContract(tag *tag) error {
	if i.contracts == nil {
		i.contracts = make([]*openapi.Contact, 0, 10)
	}

	data := split(tag.data, 3)

	if len(data) < 2 || len(data) > 3 {
		return tag.syntaxError(locale.Sprintf(locale.ErrInvalidFormat, "@apiContract"))
	}

	c := &openapi.Contact{Name: string(data[0])}
	v := string(data[1])
	switch checkContractType(v) {
	case 1:
		c.URL = v
	case 2:
		c.Email = v
	case 3:
		return tag.syntaxError(locale.Sprintf(locale.ErrInvalidFormat, "@apiContract"))
	}

	if len(data) == 3 {
		v := string(data[2])
		switch checkContractType(v) {
		case 1:
			c.URL = v
		case 2:
			c.Email = v
		case 3:
			return tag.syntaxError(locale.Sprintf(locale.ErrInvalidFormat, "@apiContract"))
		}
	}

	i.contracts = append(i.contracts, c)
	return nil
}

func checkContractType(v string) int8 {
	switch {
	case is.URL(v):
		return 1
	case is.Email(v):
		return 2
	default:
		return 0
	}
}
