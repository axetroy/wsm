test:
	GO_TESTING=1 go test --cover -covermode=count -coverprofile=coverage.out ./...

build:
	bash build.sh
