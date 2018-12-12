IMAGE ?= chickenzord/kube-annotate

build:
	docker build -t $(IMAGE) .

push:
	docker push $(IMAGE)
