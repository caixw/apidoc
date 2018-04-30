// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package openapi

// SecurityScheme.IN 的值定义
const (
	SecurityInQuery  = "query"
	SecurityInHeader = "header"
	SecurityInCookie = "cookie"
)

// SecurityRequirement Object
type SecurityRequirement map[string][]string

// SecurityScheme Object
type SecurityScheme struct {
	Type             string      `json:"type" yaml:"type"`
	Description      string      `json:"description,omitempty" yaml:"description,omitempty"`
	Name             string      `json:"name" yaml:"name"` // 报头或是 cookie 的名称
	IN               string      `json:"in" yaml:"in"`     // 位置, header, query 和 cookie
	Scheme           string      `json:"scheme" yaml:"scheme"`
	BearerFormat     string      `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	Flows            *OAuthFlows `json:"flows" yaml:"flows"`
	OpenIDConnectURL string      `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`
}

// OAuthFlows Object
type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty" yaml:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
}

// OAuthFlow Object
type OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`
	RefreshURL       string            `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty" yaml:"scopes,omitempty"`
}
