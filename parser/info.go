// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package parser

import (
	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/openapi"
)

// Info 文档的信息
type Info struct {
	Group string `yaml:"group,omitempty"`

	// openapi 根元素
	Servers      []*openapi.Server              `yaml:"servers,omitempty"`
	Components   *openapi.Components            `yaml:"components,omitempty"`
	Security     []*openapi.SecurityRequirement `yaml:"security,omitempty"`
	Tags         []*openapi.Tag                 `yaml:"tags,omitempty"`
	ExternalDocs *openapi.ExternalDocumentation `yaml:"externalDocs,omitempty"`

	// openapi.Info 元素内容
	Title          string              // 从 @apidoc 中获取
	Description    openapi.Description `yaml:"description,omitempty"`
	TermsOfService string              `json:"termsOfService,omitempty"`
	Contact        *openapi.Contact    `yaml:"contact,omitempty"`
	License        *openapi.License    `yaml:"license,omitempty"`
	Version        string              `yaml:"version"`
}

func (doc *doc) parseInfo(info *Info) error {
	doc.locker.Lock()
	defer doc.locker.Unlock()

	if doc.OpenAPI.Info != nil {
		return locale.Errorf(locale.ErrApidocExists)
	}

	doc.OpenAPI.Servers = info.Servers
	doc.OpenAPI.Components = info.Components
	doc.OpenAPI.Security = info.Security
	doc.OpenAPI.Tags = info.Tags
	doc.OpenAPI.ExternalDocs = info.ExternalDocs

	doc.OpenAPI.Info = &openapi.Info{
		Title:          info.Title,
		Description:    info.Description,
		TermsOfService: info.TermsOfService,
		Contact:        info.Contact,
		License:        info.License,
		Version:        info.Version,
	}

	return nil
}
