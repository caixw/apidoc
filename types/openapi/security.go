// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package openapi

// SecurityScheme.IN 的可选值
const (
	SecurityInQuery  = "query"
	SecurityInHeader = "header"
	SecurityInCookie = "cookie"
)

// Security.Type 的可选值
const (
	SecurityTypeAPIKey        = "apikey"
	SecurityTypeHTTP          = "http"
	SecurityTypeOAuth2        = "oauth2"
	SecurityTypeOpenIDConnect = "openIdConnect"
)

// SecurityRequirement Object
//
// 键名指向的是 Components.SecuritySchemes 中的名称。
// 若 SecurityScheme.Type 是 oauth2 或是 openIDConnect，
// 则 SecurityRequirement 的键值必须是个空值，否则键值为一个 scope 列表。
type SecurityRequirement map[string][]string

// SecurityScheme Object
type SecurityScheme struct {
	Type             string      `json:"type" yaml:"type"`
	Description      Description `json:"description,omitempty" yaml:"description,omitempty"`
	Name             string      `json:"name" yaml:"name"` // 报头或是 cookie 的名称
	IN               string      `json:"in" yaml:"in"`     // 位置, header, query 和 cookie
	Scheme           string      `json:"scheme" yaml:"scheme"`
	BearerFormat     string      `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	Flows            *OAuthFlows `json:"flows" yaml:"flows"`
	OpenIDConnectURL string      `json:"openIdConnectUrl" yaml:"openIdConnectUrl"`

	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
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
