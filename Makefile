APP = app
VERSION ?= v0.0.0
IMAGE_BUILDER ?= podman
IMAGE ?= quay.io/rmarasch/app
IMAGE_TAG = $(IMAGE):latest
OUTPUT_DIR ?= _output
OUTPUT_BIN = $(OUTPUT_DIR)/bin
BIN = $(OUTPUT_BIN)/$(APP)

default: build

.PHONY: build
build:
	CGO_ENABLED=0 go build \
		-mod vendor \
		-ldflags="-X 'main.Version=$(VERSION)'" \
		-o $(BIN) \
		./cmd/$(APP)

.PHONY: image
image:
	$(IMAGE_BUILDER) build --tag=$(IMAGE_TAG) .

.PHONY: tests
tests:
	go test -mod vendor ./...

.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIR)
