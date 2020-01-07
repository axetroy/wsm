![License](https://img.shields.io/github/license/axetroy/wsm.svg)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/wsm.svg)

Backend: [![Build Status](https://github.com/axetroy/wsm/workflows/backen/badge.svg)](https://github.com/axetroy/wsm/actions)
[![Docker Build Status](https://img.shields.io/docker/cloud/build/axetroy/wsm-backend)](https://hub.docker.com/r/axetroy/wsm-backend/builds)
[![Coverage Status](https://coveralls.io/repos/github/axetroy/wsm/badge.svg?branch=master)](https://coveralls.io/github/axetroy/wsm?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/wsm)](https://goreportcard.com/report/github.com/axetroy/wsm)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/axetroy/wsm)

Frontend:
[![Build Status](https://github.com/axetroy/wsm/workflows/frontend/badge.svg)](https://github.com/axetroy/wsm/actions)
[![Docker Build Status](https://img.shields.io/docker/cloud/build/axetroy/wsm-frontend)](https://hub.docker.com/r/axetroy/wsm-frontend/builds)
[![DeepScan grade](https://deepscan.io/api/teams/6484/projects/8581/branches/105883/badge/grade.svg)](https://deepscan.io/dashboard#view=project&tid=6484&pid=8581&bid=105883)

## Web Server Manager

通过 Web 来管理远端服务器

管理员邀请成员加入团队，管理服务器。

成员无需关心服务器的`地址/帐号/密码/密钥`等敏感信息，即可连接服务器进行操作

在成员离职后，在团队中移除成员即可

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

部分部署分为两部分

- 数据库
- 程序

### 数据库

使用 docker-compose 部署数据库

```yaml
version: "3"
services:
  # 数据库
  pg:
    image: postgres:9.6.16-alpine
    restart: always
    volumes:
      - "./volumes/pg:/var/lib/postgresql/data"
    ports:
      - 54321:5432 # 本机端口:容器端口
    environment:
      - POSTGRES_USER=terminal # 用户名
      - POSTGRES_PASSWORD=terminal # 数据库密码
      - POSTGRES_DB=terminal # 数据库名

  # 缓存
  redis:
    image: redis:5.0.7-alpine
    restart: always
    ports:
      - 6379:6379
    volumes:
      - "./volumes/redis:/data"
    environment:
      - REDIS_PASSWORD=password
    command: ["redis-server", "--requirepass", "password"]
```

### 程序

部署应用程序，使用 `Nginx` + `前端镜像` + `后端镜像` 进行部署

需要使用 [nginx.conf](nginx.conf) 文件和 [docker-compose.yml](docker-compose.yml) 文件

```yaml
version: "3"
services:
  # 网关
  nginx:
    image: nginx:1.17.6-alpine
    restart: always
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf # 映射 nginx 配置文件
      - ./logs/nginx:/var/log/nginx # 日志文件
    ports:
      - 9000:80 # 本机端口:容器端口
    links:
      - frontend
      - backend

  # 前端
  frontend:
    image: axetroy/wsm-frontend:latest
    restart: always
    environment:
      - PORT=80
      - HOST=0.0.0.0

  # 后端接口
  backend:
    image: axetroy/wsm-backend:latest
    restart: always
    environment:
      - USER_HTTP_PORT=80
      - DB_HOST=192.168.3.15 # 数据库的IP地址
      - DB_PORT=54321 # 数据库的端口
```

## 技术栈

- Golang
- Node.js + Nuxt

## TODO

- [ ] 一次性分享终端
- [ ] 终端操作记录
- [ ] 操作记录回放

## 许可协议

[Apache License 2.0](LICENSE)
