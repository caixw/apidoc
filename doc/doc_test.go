// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"testing"

	"github.com/issue9/assert"
)

func TestDoc_parseTag(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	a.NotError(d.parseTag(newTag("tag1 标签 1 的描述内容")))
	a.Equal(len(d.Tags), 1)
	tag := d.Tags[0]
	a.Equal(tag.Name, "tag1").
		Equal(tag.Description, "标签 1 的描述内容")

	// 格式错误
	a.Error(d.parseTag(newTag("tag1")))
}

func TestDoc_parseServer(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	a.NotError(d.parseServer(newTag("admin https://admin.api.example.com admin api")))
	a.Equal(len(d.Servers), 1)
	srv := d.Servers[0]
	a.Equal(srv.Name, "admin").
		Equal(srv.URL, "https://admin.api.example.com").
		Equal(srv.Description, "admin api")

	a.NotError(d.parseServer(newTag("client https://client.api.example.com client api")))
	a.Equal(len(d.Servers), 2)
	srv = d.Servers[1]
	a.Equal(srv.Name, "client").
		Equal(srv.URL, "https://client.api.example.com").
		Equal(srv.Description, "client api")

	a.Error(d.parseServer(newTag("client")))
	a.Error(d.parseServer(newTag("client https://url")))
}

func TestDoc_parseLicense(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	a.NotError(d.parseLicense(newTag("MIT https://opensources.org/licenses/MIT")))
	a.NotNil(d.License).
		Equal(d.License.Text, "MIT").
		Equal(d.License.URL, "https://opensources.org/licenses/MIT")

	// d.License 已经存在，再次添加会出错。
	a.Error(d.parseLicense(newTag("MIT https://opensources.org/licenses/MIT")))
}

func TestNewLink(t *testing.T) {
	a := assert.New(t)

	// 格式不够长
	l, err := newLink(newTag("text"))
	a.Error(err).Nil(l)

	// 格式不正确
	l, err = newLink(newTag("text https://"))
	a.Error(err).Nil(l)

	l, err = newLink(newTag("text  https://example.com"))
	a.NotError(err).
		NotNil(l).
		Equal(l.Text, "text").
		Equal(l.URL, "https://example.com")
}

func TestNewContact(t *testing.T) {
	a := assert.New(t)

	// 格式不够长
	c, err := newContact(newTag("name"))
	a.Error(err).Nil(c)

	// 格式不正确
	c, err = newContact(newTag("name name@"))
	a.Error(err).Nil(c)

	// 格式不正确
	c, err = newContact(newTag("name name@example.com https://"))
	a.Error(err).Nil(c)

	c, err = newContact(newTag("name name@example.com"))
	a.NotError(err).
		NotNil(c).
		Equal(c.Name, "name").
		Equal(c.Email, "name@example.com").
		Empty(c.URL)

	c, err = newContact(newTag("name name@example.com https://example.com"))
	a.NotError(err).
		NotNil(c).
		Equal(c.Name, "name").
		Equal(c.Email, "name@example.com").
		Equal(c.URL, "https://example.com")

	c, err = newContact(newTag("name https://example.com name@example.com"))
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
