# Copyright 2016 by caixw, All rights reserved.
# Use of this source code is governed by a MIT
# license that can be found in the LICENSE file.


# @api POST /users/login 登录
# group users
# tags: [t1,t2]
#
# request:
#   description: request body
#   content:
#     application/json:
#       schema:
#         type: object
#         properties:
#           username:
#             type: string
#             description: 登录账号
#           password:
#             type: string
#             description: 密码
sub login()
{
	print("/**********", "login")
	# TODO
}

# 123
# 123
=pod
@api DELETE /users/login 注销登录
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
=cut
sub logout()
{
	print("logout", "**********/")
	# TODO
}
