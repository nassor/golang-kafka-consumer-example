all: test

test:
	@echo "#unit tests"
	@go test -mod=vendor -race -cover -short ./...

test-int-func:
	@echo "#unit tests"
	@go test -mod=vendor -race -cover -short ./...

proto:
	mkdir -p internal/pkg/pb
	protoc -I/usr/local/include --go_out=plugins=grpc:./internal/pkg/pb --proto_path=./api/protobuf-spec/ ./api/protobuf-spec/*.proto