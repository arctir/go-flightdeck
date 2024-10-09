OAPI_CODEGEN_VERSION ?= v1.16.3
API_SPEC ?= "./api/generated/v1/api.gen.yaml"

.PHONY: all
all: package 

.PHONY: gen
gen: deps
	mkdir -p pkg/api/v1/
	oapi-codegen -config oapi-config.yaml ${API_SPEC} > pkg/api/v1/api.gen.go

.PHONY: deps
deps:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)

.PHONY: package
package: gen
	go build ./...
