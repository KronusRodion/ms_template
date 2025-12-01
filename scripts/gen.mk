# Generator for Protocol Buffers
PROTO_DIR = proto
GEN_GO_DIR = ./gen/go

# Find all .proto files recursively
PROTO_FILES = $(shell find $(PROTO_DIR) -name '*.proto')

# Generate Go targets from proto files
GO_PB_TARGETS = $(PROTO_FILES:$(PROTO_DIR)/%.proto=$(GEN_GO_DIR)/%.pb.go)
GO_GRPC_TARGETS = $(PROTO_FILES:$(PROTO_DIR)/%.proto=$(GEN_GO_DIR)/%_grpc.pb.go)

# Generate Go code
.PHONY: gen-go
gen-go: $(GO_PB_TARGETS) $(GO_GRPC_TARGETS)

# Pattern rule for Go files
$(GEN_GO_DIR)/%.pb.go $(GEN_GO_DIR)/%_grpc.pb.go: $(PROTO_DIR)/%.proto
	@mkdir -p $(dir $@)
	protoc -I $(PROTO_DIR) $< \
		--go_out=$(GEN_GO_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_GO_DIR) \
		--go-grpc_opt=paths=source_relative
	@echo "Generated: $@"

# Clean generated files
.PHONY: clean
clean:
	rm -rf $(GEN_GO_DIR)
	@echo "Cleaned generated files"

# Show found proto files
.PHONY: list-proto
list-proto:
	@echo "Found proto files:"
	@for file in $(PROTO_FILES); do \
		echo "  $$file"; \
	done

# Show targets that will be generated
.PHONY: list-targets
list-targets:
	@echo "Go PB targets:"
	@for target in $(GO_PB_TARGETS); do \
		echo "  $$target"; \
	done
	@echo "Go gRPC targets:"
	@for target in $(GO_GRPC_TARGETS); do \
		echo "  $$target"; \
	done

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Generate all code (default)"
	@echo "  gen-go       - Generate Go code only"
	@echo "  clean        - Remove all generated files"
	@echo "  list-proto   - List all found proto files"
	@echo "  list-targets - List all targets that will be generated"
	@echo "  help         - Show this help"