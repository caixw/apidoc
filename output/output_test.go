// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"os"
	"testing"

	"github.com/caixw/apidoc/core"
	"github.com/issue9/assert"
)

var (
	block1 = `
@api get /api/users  获取所有用户信息
获取所有用户的详细信息，包括用户所属的权限组，昵称等。
若没有权限，则返回空对象。
@apiGroup users
@apiQuery status string optional 只获取指定状态的数据，可用值为normal, locked
@apiSuccess 200 json 成功获取用户信息
@apiParam users object 表示所有的用户列表
@apiExample json
{"users":[
	{"id":1, "name": "n1", "group": 1},
	{"id":2, "name": "n2", "group": 1},
	{"id":3, "name": "n3", "group": 1},
]}
`
	block2 = `
@api get /api/users/{id} 获取指定用户的详细信息
@apiParam id int 用户的ID值
@apiGroup users
@apiSuccess 200 json 成功获取信息
@apiParam id int 用户的ID
@apiParam name string 用户名称
@apiParam group int 用户所在的权限组ID
@apiExample json
{"id":1, "name": "n1", "group": 1}
`
	block3 = `
@api post /api/auth/login 请求登录用户
@apiGroup auth
@apiRequest json
@apiParam username string 登录用户名
@apiParam password string 登录密码
@apiExample json
{"username": "admin", "password": "admin"}
@apiSuccess 200 成功登录
@apiHeader token xxx
@apiError 200 账号或是密码错误
@apiParam message string 描述错误信息
@apiExample json
{"message":"账号或是密码错误"}
`
)

func TestLexer_Html(t *testing.T) {
	testdir := "./testdir"
	a := assert.New(t)

	// 创建测试目录
	a.NotError(os.MkdirAll(testdir, os.ModePerm), "无法创建测试目录")

	docs := []*core.Doc{}

	doc, err := core.Scan([]rune(block1), 1, "test1.go")
	a.NotError(err).NotNil(doc)
	docs = append(docs, doc)

	doc, err = core.Scan([]rune(block2), 100, "test1.go")
	a.NotError(err).NotNil(doc)
	docs = append(docs, doc)

	doc, err = core.Scan([]rune(block3), 1, "test2.go")
	a.NotError(err).NotNil(doc)
	docs = append(docs, doc)

	opt := &Options{
		DocDir:     testdir,
		Version:    "doc v0.1",
		AppVersion: "appver 0.1",
		Title:      "TestDoc",
	}
	a.NotError(Html(docs, opt))
}
