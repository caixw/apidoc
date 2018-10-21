#!/bin/sh
# Copyright 2017 by caixw, All rights reserved.
# Use of this source code is governed by a MIT
# license that can be found in the LICENSE file.

# 指定工作目录
wd=$(dirname $0)/../cmd

# 指定编译日期
date=`date -u '+%Y%m%d'`

# 获取最后一条前交的 hash
hash=`git rev-parse HEAD`

# 需要修改变量的名名，若为 main，则指接使用 main，而不是全地址
path=github.com/caixw/apidoc/vars

cd ${wd}

echo '开始编译'
go build -o ./apidoc -ldflags "-X ${path}.buildDate=${date} -X ${path}.commitHash=${hash}" -v
