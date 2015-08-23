apidoc [![Build Status](https://travis-ci.org/caixw/apidoc.svg?branch=master)](https://travis-ci.org/caixw/apidoc)
======

apidoc是一个简单的RESTful api文档生成工具，它从代码注释中提取特定格式的内容，生成文档。
目前支持c系列注释风格的语言和ruby。
具体文档可参考:https://caixw.github.io/apidoc/

```c
/*
 @api get /users 获取所有的用户信息
 @apiGroup users
 @apiQuery page int 显示第几页的内容
 @apiQuery size int 每页显示的数量
 @apiSuccess 200 ok
 @apiParam count int 符合条件的所有用户数量
 @apiParam users array 用户列表。
 @apiExample json
 {
     "count": 500,
     "users": [
        {"id":1, "username": "admin1", "name": "管理员2"},
        {"id":2, "username": "admin2", "name": "管理员2"}
     ],
 }
 @apiExmaple xml
 <users count="500">
     <user id="1" username="admin1" name="管理员1" />
     <user id="2" username="admin2" name="管理员2" />
 </users>
 */
```

### 安装

```shell
go get github.com/caixw/apidoc
```


### 版权

本项目采用[MIT](http://opensource.org/licenses/MIT)开源授权许可证，完整的授权说明可在[LICENSE](LICENSE)文件中找到。
