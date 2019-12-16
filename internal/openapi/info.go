// SPDX-License-Identifier: MIT

package openapi

import (
	"github.com/issue9/is"
	"github.com/issue9/version"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
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

func (info *Info) sanitize() *message.SyntaxError {
	if info.Title == "" {
		return message.NewLocaleError("", "title", 0, locale.ErrRequired)
	}

	if !version.SemVerValid(info.Version) {
		return message.NewLocaleError("", "version", 0, locale.ErrInvalidFormat)
	}

	if info.TermsOfService != "" && !is.URL(info.TermsOfService) {
		return message.NewLocaleError("", "termsOfService", 0, locale.ErrInvalidFormat)
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

func (l *License) sanitize() *message.SyntaxError {
	if l.URL != "" && !is.URL(l.URL) {
		return message.NewLocaleError("", "url", 0, locale.ErrInvalidFormat)
	}

	return nil
}

func newLicense(l *doc.Link) *License {
	return &License{
		Name: l.Text,
		URL:  l.URL,
	}
}

func newContact(c *doc.Contact) *Contact {
	return &Contact{
		Name:  c.Name,
		URL:   c.URL,
		Email: c.Email,
	}
}

func (c *Contact) sanitize() *message.SyntaxError {
	if c.URL != "" && !is.URL(c.URL) {
		return message.NewLocaleError("", "url", 0, locale.ErrInvalidFormat)
	}

	if c.Email != "" && !is.Email(c.Email) {
		return message.NewLocaleError("", "email", 0, locale.ErrInvalidFormat)
	}

	return nil
}
