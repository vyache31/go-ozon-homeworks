.PHONY: generate run-all test test-race test-api

generate:
	protoc -I ./ \
	--go_out ./pkg --go_opt paths=source_relative \
	--go-grpc_out ./pkg --go-grpc_opt paths=source_relative \
	api/product.proto

run-all:
	@go run ./cmd/product & \
	go run ./cmd/cart & \
	wait

test:
	go test ./...

test-race:
	go test -race ./...

test-api:
	go test -v -count=1 -tags=api ./test/
