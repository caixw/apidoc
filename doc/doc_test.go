// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/input"
)

func TestDoc_parseBlock(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	a.NotError(d.parseBlock(input.Block{Data: []byte("@api GET /path summary")}))
	a.NotError(d.parseBlock(input.Block{Data: []byte("@apidoc title")}))
	// 任意其它内容
	a.NotError(d.parseBlock(input.Block{Data: []byte("xxxxx")}))
}

func TestDoc_parseapidoc(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	// 不能为空
	a.Error(d.parseapidoc(nil, newTag("")))

	// 正常
	a.NotError(d.parseapidoc(nil, newTag("title of doc")))
	a.Equal(d.Title, "title of doc")

	// 不能多次调用
	a.Error(d.parseapidoc(nil, newTag("xxx")))
}

func TestDoc_parseContent(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	// 不能为空
	a.Error(d.parseContent(nil, newTag("")))

	// 正常
	a.NotError(d.parseContent(nil, newTag("xxx\nyyy\nzzz")))
	a.Equal(d.Content, Markdown("xxx\nyyy\nzzz"))

	// 不能多次调用
	a.Error(d.parseContent(nil, newTag("xxx")))
}

func TestDoc_parseVersion(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	// 不能为空
	a.Error(d.parseVersion(nil, newTag("")))

	// 正常
	a.NotError(d.parseVersion(nil, newTag("3.2.1")))

	// 不能多次调用
	a.Error(d.parseVersion(nil, newTag("3.2.1")))
}

func TestDoc_parseContact(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	a.NotError(d.parseContact(nil, newTag("name name@example.com https://example.com")))

	// 不能重复调用
	a.Error(d.parseContact(nil, newTag("name name@example.com https://example.com")))
}

func TestDoc_parseTag(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	a.NotError(d.parseTag(nil, newTag("tag1 标签 1 的描述内容")))
	a.Equal(len(d.Tags), 1)
	tag := d.Tags[0]
	a.Equal(tag.Name, "tag1").
		Equal(tag.Description, "标签 1 的描述内容")

	// 格式错误
	a.Error(d.parseTag(nil, newTag("tag1")))

	// 重复的标签名
	a.Error(d.parseTag(nil, newTag("tag1 desc")))
}

func TestDoc_parseServer(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	a.NotError(d.parseServer(nil, newTag("admin https://admin.api.example.com admin api")))
	a.Equal(len(d.Servers), 1)
	srv := d.Servers[0]
	a.Equal(srv.Name, "admin").
		Equal(srv.URL, "https://admin.api.example.com").
		Equal(srv.Description, "admin api")

	a.NotError(d.parseServer(nil, newTag("client https://client.api.example.com client api")))
	a.Equal(len(d.Servers), 2)
	srv = d.Servers[1]
	a.Equal(srv.Name, "client").
		Equal(srv.URL, "https://client.api.example.com").
		Equal(srv.Description, "client api")

	// 少内容
	a.Error(d.parseServer(nil, newTag("client")))

	// 格式不正确
	a.Error(d.parseServer(nil, newTag("client https://url")))

	// 重复的内容
	a.Error(d.parseServer(nil, newTag("client https://example.com desc")))
}

func TestDoc_parseLicense(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	// 长度不够
	a.Error(d.parseLicense(nil, newTag("MIT")))

	// 非 URL
	a.Error(d.parseLicense(nil, newTag("MIT https://")))

	a.NotError(d.parseLicense(nil, newTag("MIT https://opensources.org/licenses/MIT")))
	a.NotNil(d.License).
		Equal(d.License.Text, "MIT").
		Equal(d.License.URL, "https://opensources.org/licenses/MIT")

	// d.License 已经存在，再次添加会出错。
	a.Error(d.parseLicense(nil, newTag("MIT https://opensources.org/licenses/MIT")))
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
