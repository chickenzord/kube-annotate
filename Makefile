DOCKER_TAG ?= local
DOCKER_IMAGE ?= chickenzord/kube-annotate:$(DOCKER_TAG)
PACKAGE ?= github.com/chickenzord/kube-annotate
BUILD_OUTPUT ?= bin/kube-annotate
BUILD_GOOS ?= darwin linux
BUILD_GOARCH ?= 386 amd64

clean:
	mkdir -p bin
	rm -f bin/*

deps:
	go get -v github.com/stretchr/testify
	go get -v github.com/ahmetb/govvv
	dep ensure -v -vendor-only

test:
	go test -v -cover -coverprofile=coverage.txt -covermode=atomic ./...

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) \
	go build -o $(BUILD_OUTPUT) -ldflags="$$(govvv -flags -pkg $(PACKAGE)/pkg/config)" ./cmd/kube-annotate

build-all-platforms:
	for GOOS in $(BUILD_GOOS); do \
		for GOARCH in $(BUILD_GOARCH); do \
			GOOS=$$GOOS \
			GOARCH=$$GOARCH \
			BUILD_OUTPUT="bin/kube-annotate-$$GOOS-$$GOARCH" \
			CGO_ENABLED=0 \
			$(MAKE) build; \
		done; \
	done;

run:
	go run ./cmd/kube-annotate

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-push:
	docker push $(DOCKER_IMAGE)
