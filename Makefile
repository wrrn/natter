.DEFAULT_GOAL := help
PROTO_FILES := post.proto likes.proto rate-limit.proto
.PHONY: generate-grpc
generate-grpc: ## generate the gRPC code from the protobuf definitions
	@for proto in $(PROTO_FILES); do \
		protoc -I /usr/local/include -I api/proto --go_out=plugins=grpc:$$GOPATH/src $${proto}; \
	done



generate-protoset:
	@for proto in $(PROTO_FILES); do \
		protoc \
			-I /usr/local/include \
			-I api/proto \
			--descriptor_set_out=api/proto/protosets/$${proto}set \
			--include_imports \
			$${proto}; \
	done

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
