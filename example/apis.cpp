// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// @api GET /users/{id}/logs 获取用户信息
// @apigroup g1
// @apitags t1,t2
// @apideprecated desc
// @apiquery page int default desc
// @apiquery size int default desc
// @apiquery state array.string [normal,lock] 状态码
// @apiparam id int desc
// @apiparam id int desc
//
// @apirequest object * 通用的请求主体
// @apiheader name desc optional
// @apiheader name desc optional
// @apiparam count int optional desc
// @apiparam list array must desc
// @apiparam list.id int optional desc
// @apiparam list.name int must desc
// @apiparam list.groups array.string optional.xxxx desc markdown enum:
//  * xx: xxxxx
//  * xx: xxxxx
// @apiexample application/json summary
// {
//  count: 5,
//  list: [
//    {id:1, name: 'name1', 'groups': [1,2]},
//    {id:2, name: 'name2', 'groups': [1,2]}
//  ]
// }
//
// @apirequest object application/xml 特定的请求主体
//
// @apiresponse 200 array.object * 通用的返回内容定义
// @apiheader string xxx
// @apiparam id int desc
// @apiparam name string desc
// @apiparam group object desc
// @apiparam group.id int desc
//
// @apiresponse 404 object application/json
// @apiheader string xxx
// @apiparam code int desc
// @apiparam message string desc
// @apiparam detail array.object desc
// @apiparam detail.id string desc
// @apiparam detail.message string desc
void api() {
    // TODO
}
