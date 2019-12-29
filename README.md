[![Build Status](https://travis-ci.com/axetroy/wsm.svg?token=QMG6TLRNwECnaTsy6ssj&branch=master)](https://travis-ci.com/axetroy/wsm)
[![Coverage Status](https://coveralls.io/repos/github/axetroy/wsm/badge.svg?branch=master)](https://coveralls.io/github/axetroy/wsm?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/wsm)](https://goreportcard.com/report/github.com/axetroy/wsm)
![License](https://img.shields.io/github/license/axetroy/wsm.svg)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/wsm.svg)

## Web Server Manager

通过 Web 来管理远端服务器

特性:

- [x] 团队支持，创建/管理团队
- [x] Web 终端

## 使用

```shell
$ git clone https://github.com/axetroy/wsm.git $GOPATH/src/github.com/axetroy/wsm
$ cd $GOPATH/src/github.com/axetroy/wsm
```

### 启动后端 API

```shell
$ go run cmd/user/main.go start
```

### 启动前端页面

```shell
$ cd ./frontend
$ yarn
$ npm run dev
```

## 部署

TODO: 在往后会打包成 Docker 镜像进行部署

## 许可协议

[Apache License 2.0](LICENSE)
