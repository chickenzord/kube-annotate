TAG ?= latest
IMAGE ?= chickenzord/kube-annotate:$(TAG)

build:
	docker build -t $(IMAGE) .

push:
	docker push $(IMAGE)
