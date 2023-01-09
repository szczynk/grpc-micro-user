# ==============================================================================
# Proto
# Generate protobuf in Go, gRPC for Go, gRPC Gateway for Go, and Swagger Documentation in JSON
proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto \
	--go_out=pb --go_opt=paths=source_relative \
	--validate_out=pb --validate_opt=lang=go,paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative,logtostderr=true \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=auth,logtostderr=true \
	proto/*.proto

# ==============================================================================
# Main
run:
	go run cmd/server/main.go

build:
	go build cmd/server/main.go

test:
	go test -v -cover ./...

lint:
	@echo Starting golangci-lint
	golangci-lint run ./...
	@echo Starting staticcheck
	staticcheck -checks U1000 ./...

migrate:
	go run migration/migration.go

.PHONY: proto run build test lint migrate