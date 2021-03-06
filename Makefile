# Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.

test:
	GO_TESTING=1 go test --cover -covermode=count -coverprofile=coverage.out ./...

start:
	GO111MODULE=on go run cmd/user/main.go start

build-backend:
	bash build.sh

build-frontend:
	cd ./frontend && npm run build

build-docker:
	make build-docker-backend
	make build-docker-frontend

build-docker-backend:
	docker build --tag axetroy/wsm-backend:latest ./

build-docker-frontend:
	docker build --tag axetroy/wsm-frontend:latest ./frontend