// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package types

import "github.com/caixw/apidoc/types/openapi"

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
	Title          string              `yaml:"title"`
	Description    openapi.Description `yaml:"description,omitempty"`
	TermsOfService string              `json:"termsOfService,omitempty"`
	Contact        *openapi.Contact    `yaml:"contact,omitempty"`
	License        *openapi.License    `yaml:"license,omitempty"`
	Version        string              `yaml:"version"`
}

func (doc *Doc) parseInfo(info *Info) error {
	// TODO
	return nil
}
