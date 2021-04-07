#!/bin/sh
# SPDX-License-Identifier: MIT

# 指定工作目录
wd=$(dirname $0)/../cmd/apidoc

# 指定编译日期
date=`date -u '+%Y%m%d'`

# 获取最后一条前交的 hash
hash=`git rev-parse HEAD`

# 需要修改变量的地址，若为 main，则指接使用 main，而不是全地址
path=github.com/caixw/apidoc/v7/core

ldflags="-X ${path}.metadata=${date}.${hash}"
