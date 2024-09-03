LOCAL_BIN:=$(CURDIR)/bin

.PHONY: .bin-deps
.bin-deps:
	$(info Installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

PHONY: .protoc-generate
.protoc-generate: .bin-deps
	protoc \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
		--go_out=. --go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
      	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
      	metrics.proto
