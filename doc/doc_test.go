// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/errors"
	"github.com/caixw/apidoc/internal/input"
)

func TestDoc_parseBlock(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}
	h := &errors.Handler{}

	d.parseBlock(input.Block{Data: []byte("@api GET /path summary")}, h)
	a.Equal(len(d.Apis), 1)

	d.parseBlock(input.Block{Data: []byte("@apidoc title")}, h)
	a.Equal(d.Title, "title")

	// 任意其它内容
	d.parseBlock(input.Block{Data: []byte("xxxxx")}, h)
}

func TestDoc_parseapidoc(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	// 不能为空
	tag := newTag("@apidoc")
	d.parseapidoc(nil, tag)

	// 正常
	tag = newTag("@apidoc title of doc")
	d.parseapidoc(nil, tag)
	a.Equal(d.Title, "title of doc")

	// 不能多次调用
	tag = newTag("xxx")
	d.parseapidoc(nil, tag)
}

func TestDoc_parseContent(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	// 不能为空
	tag := newTag("@apiContent")
	d.parseContent(nil, tag)

	// 正常
	tag = newTag("@apiContent xxx\nyyy\nzzz")
	d.parseContent(nil, tag)
	a.Equal(d.Content, Markdown("xxx\nyyy\nzzz"))

	// 不能多次调用
	tag = newTag("@apiContent xxx")
	d.parseContent(nil, tag)
}

func TestDoc_parseVersion(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	// 不能为空
	tag := newTag("@apiVersion ")
	d.parseVersion(nil, tag)
	a.Empty(d.Tags)

	// 正常
	tag = newTag("@apiVersion 3.2.1")
	d.parseVersion(nil, tag)
	a.Equal(d.Version, "3.2.1")

	// 不能多次调用
	tag = newTag("@apiVersion 4.2.1")
	d.parseVersion(nil, tag)
	a.Equal(d.Version, "3.2.1")
}

func TestDoc_parseContact(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	tag := newTag("@apiContact name name@example.com https://example.com")
	d.parseContact(nil, tag)
	a.Equal(d.Contact.Name, "name")

	// 不能重复调用
	tag = newTag("@apiContact name1 name@example.com https://example.com")
	d.parseContact(nil, tag)
	a.Equal(d.Contact.Name, "name")
}

func TestDoc_parseTag(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	tag := newTag("@apiTag tag1 标签 1 的描述内容")
	d.parseTag(nil, tag)
	a.Equal(len(d.Tags), 1)
	tag0 := d.Tags[0]
	a.Equal(tag0.Name, "tag1").
		Equal(tag0.Description, "标签 1 的描述内容")

	// 格式错误
	tag = newTag("@apiTag tag1")
	d.parseTag(nil, tag)

	// 重复的标签名
	tag = newTag("@apiTag tag1 desc")
	d.parseTag(nil, tag)
}

func TestDoc_parseServer(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	tag := newTag("@apiServer admin https://admin.api.example.com admin api")
	d.parseServer(nil, tag)
	a.Equal(len(d.Servers), 1)
	srv := d.Servers[0]
	a.Equal(srv.Name, "admin").
		Equal(srv.URL, "https://admin.api.example.com").
		Equal(srv.Description, "admin api")

	tag = newTag("@apiServer client https://client.api.example.com client api")
	d.parseServer(nil, tag)
	a.Equal(len(d.Servers), 2)
	srv = d.Servers[1]
	a.Equal(srv.Name, "client").
		Equal(srv.URL, "https://client.api.example.com").
		Equal(srv.Description, "client api")

	// 少内容
	tag = newTag("@apiServer client")
	d.parseServer(nil, tag)

	// 格式不正确
	tag = newTag("@apiServer client https://url")
	d.parseServer(nil, tag)

	// 重复的内容
	tag = newTag("@apiServer client https://example.com desc")
	d.parseServer(nil, tag)
}

func TestDoc_parseLicense(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	// 长度不够
	tag := newTag("MIT")
	d.parseLicense(nil, tag)

	// 非 URL
	tag = newTag("@apiLicense MIT https://")
	d.parseLicense(nil, tag)

	tag = newTag("@apiLicense MIT https://opensources.org/licenses/MIT")
	d.parseLicense(nil, tag)
	a.NotNil(d.License).
		Equal(d.License.Text, "MIT").
		Equal(d.License.URL, "https://opensources.org/licenses/MIT")

	// d.License 已经存在，再次添加会出错。
	tag = newTag("@apiLicense MIT https://opensources.org/licenses/MIT")
	d.parseLicense(nil, tag)
}

func TestNewContact(t *testing.T) {
	a := assert.New(t)

	// 格式不够长
	tag := newTag("@apiContact name")
	c, ok := newContact(tag)
	a.False(ok).Nil(c)

	// 格式不正确
	tag = newTag("@apiContact name name@")
	c, ok = newContact(tag)
	a.False(ok).Nil(c)

	// 格式不正确
	tag = newTag("@apiContact name name@example.com https://")
	c, ok = newContact(tag)
	a.False(ok).Nil(c)

	tag = newTag("@apiContact name name@example.com")
	c, ok = newContact(tag)
	a.True(ok).
		NotNil(c).
		Equal(c.Name, "name").
		Equal(c.Email, "name@example.com").
		Empty(c.URL)

	tag = newTag("@apiContact name name@example.com https://example.com")
	c, ok = newContact(tag)
	a.True(ok).
		NotNil(c).
		Equal(c.Name, "name").
		Equal(c.Email, "name@example.com").
		Equal(c.URL, "https://example.com")

	tag = newTag("@apiContact name https://example.com name@example.com")
	c, ok = newContact(tag)
	a.True(ok).
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
