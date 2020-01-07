test:
	GO_TESTING=1 go test --cover -covermode=count -coverprofile=coverage.out ./...

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