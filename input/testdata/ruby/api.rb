# Copyright 2016 by caixw, All rights reserved.
# Use of this source code is governed by a MIT
# license that can be found in the LICENSE file.

module Testdata


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
def login(w, r)
	puts "=begin login";
	# TODO
end

# 123
# 123
=begin
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
=end
def logout(w , r )
    puts 'logout  =end';
    # TODO
end

end
