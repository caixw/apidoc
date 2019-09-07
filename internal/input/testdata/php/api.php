<?php
// SPDX-License-Identifier: MIT

namespace \testdata;

// @api POST /users/login 登录
// group users
// tags: [t1,t2]
//
// request:
//   description: request body
//   content:
//     application/json:
//       schema:
//         type: object
//         properties:
//           username:
//             type: string
//             description: 登录账号
//           password:
//             type: string
//             description: 密码
function login(ResponseWriter $w, Request $r) {
    echo '/******login';
	// TODO
}

// 123
// 123
/* @api DELETE /users/login 注销登录
group users
tags: [t1,t2]

request:
  description: request body
  content:
    application/json:
      schema:
        type: object
        properties:
          username:
            type: string
            description: 登录账号
          password:
            type: string
            description: 密码
*/
function logout(ResponseWriter $w, Request $r) {
    echo 'logout*********/';
	// TODO
}
