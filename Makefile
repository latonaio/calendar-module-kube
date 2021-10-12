GO_SRCS := $(shell find . -type f -name '*.go')

build-protobuf:
	protoc -I./cmd/calendar-server/proto/ \
		--go_out=plugins=grpc:./cmd/calendar-server \
		--go_opt=module=bitbucket.org/latonaio/calendar-module-kube/cmd/calendar-server \
		./cmd/calendar-server/proto/calendarpb/calendar.proto

docker-build: $(GO_SRCS)
	bash ./builders/docker-build-calendar-server.sh

go-build: $(GO_SRCS)
	go build ./cmd/calendar-server

go-test-build: $(GO_SRCS)
	go build ./test/testclient

docker-test-build:
	bash ./test/builders/docker-build-calendar-testclient.sh

go-test: $(GO_SRCS)
	go test ./...

