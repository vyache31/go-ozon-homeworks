.PHONY:generate
generate:
	protoc -I ./ \
	--go_out ./pkg --go_opt paths=source_relative \
	--go-grpc_out ./pkg --go-grpc_opt paths=source_relative \
	api/product.proto

.PHONY:run-all
run-all:
	@go run ./cmd/product & \
	go run ./cmd/cart & \
	wait

.PHONY:test-api
test-api:
	go test -v -count=1 -tags=api ./test/

.PHONY:test-race
test-race:
	go test -race ./...