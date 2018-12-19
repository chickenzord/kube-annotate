TAG ?= latest
IMAGE ?= chickenzord/kube-annotate:$(TAG)

test:
	go test ./...

build:
	docker build -t $(IMAGE) .

push:
	docker push $(IMAGE)
