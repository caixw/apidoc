// SPDX-License-Identifier: MIT

package openapi

import (
	"github.com/issue9/is"
	"github.com/issue9/version"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// Info 接口文档的基本信息
type Info struct {
	Title          string   `json:"title" yaml:"title"`
	Description    string   `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty" json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *License `json:"license,omitempty" yaml:"license,omitempty"`
	Version        string   `json:"version" yaml:"version"`
}

// Contact 描述联系方式
type Contact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

// License 授权信息
type License struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url,omitempty" yaml:"url,omitempty"`
}

func (info *Info) sanitize() *core.SyntaxError {
	if info.Title == "" {
		return core.NewLocaleError(core.Location{}, "title", locale.ErrRequired)
	}

	if !version.SemVerValid(info.Version) {
		return core.NewLocaleError(core.Location{}, "version", locale.ErrInvalidFormat)
	}

	if info.TermsOfService != "" && !is.URL(info.TermsOfService) {
		return core.NewLocaleError(core.Location{}, "termsOfService", locale.ErrInvalidFormat)
	}

	if info.Contact != nil {
		if err := info.Contact.sanitize(); err != nil {
			err.Field = "contact." + err.Field
			return err
		}
	}

	if info.License != nil {
		if err := info.License.sanitize(); err != nil {
			err.Field = "license." + err.Field
			return err
		}
	}

	return nil
}

func (l *License) sanitize() *core.SyntaxError {
	if l.URL != "" && !is.URL(l.URL) {
		return core.NewLocaleError(core.Location{}, "url", locale.ErrInvalidFormat)
	}

	return nil
}

func newLicense(l *ast.Link) *License {
	if l == nil {
		return nil
	}

	return &License{
		Name: l.Text.V(),
		URL:  l.URL.V(),
	}
}

func newContact(c *ast.Contact) *Contact {
	if c == nil {
		return nil
	}

	return &Contact{
		Name:  c.Name.V(),
		URL:   c.URL.V(),
		Email: c.Email.V(),
	}
}

func (c *Contact) sanitize() *core.SyntaxError {
	if c.URL != "" && !is.URL(c.URL) {
		return core.NewLocaleError(core.Location{}, "url", locale.ErrInvalidFormat)
	}

	if c.Email != "" && !is.Email(c.Email) {
		return core.NewLocaleError(core.Location{}, "email", locale.ErrInvalidFormat)
	}

	return nil
}
