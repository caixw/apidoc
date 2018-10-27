// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/doc/lexer"
)

func TestNewLink(t *testing.T) {
	a := assert.New(t)

	// 格式不够长
	l, err := newLink(&lexer.Tag{Data: []byte("text")})
	a.Error(err).Nil(l)

	// 格式不正确
	l, err = newLink(&lexer.Tag{Data: []byte("text https://")})
	a.Error(err).Nil(l)

	l, err = newLink(&lexer.Tag{Data: []byte("text  https://example.com")})
	a.NotError(err).
		NotNil(l).
		Equal(l.Text, "text").
		Equal(l.URL, "https://example.com")
}

func TestNewContact(t *testing.T) {
	a := assert.New(t)

	// 格式不够长
	c, err := newContact(&lexer.Tag{Data: []byte("name")})
	a.Error(err).Nil(c)

	// 格式不正确
	c, err = newContact(&lexer.Tag{Data: []byte("name name@")})
	a.Error(err).Nil(c)

	// 格式不正确
	c, err = newContact(&lexer.Tag{Data: []byte("name name@example.com https://")})
	a.Error(err).Nil(c)

	c, err = newContact(&lexer.Tag{Data: []byte("name name@example.com")})
	a.NotError(err).
		NotNil(c).
		Equal(c.Name, "name").
		Equal(c.Email, "name@example.com").
		Empty(c.URL)

	c, err = newContact(&lexer.Tag{Data: []byte("name name@example.com https://example.com")})
	a.NotError(err).
		NotNil(c).
		Equal(c.Name, "name").
		Equal(c.Email, "name@example.com").
		Equal(c.URL, "https://example.com")

	c, err = newContact(&lexer.Tag{Data: []byte("name https://example.com name@example.com")})
	a.NotError(err).
		NotNil(c).
		Equal(c.Name, "name").
		Equal(c.Email, "name@example.com").
		Equal(c.URL, "https://example.com")
}

func TestCheckContactType(t *testing.T) {
	a := assert.New(t)

	a.Equal(1, checkContactType("https://example.com"))
	a.Equal(2, checkContactType("user@example.com"))
	a.Equal(0, checkContactType("xxxx"))
}
