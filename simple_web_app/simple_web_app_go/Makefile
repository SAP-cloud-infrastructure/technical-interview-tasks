# Makefile for the Go simple-web-app.
#
# The image name/tag matches what the Kubernetes manifests expect. After
# building, load it into your cluster with `make kind-load` (kind) or
# `make minikube-load` (minikube).

IMAGE_NAME ?= local/simple-web-server
IMAGE_TAG  ?= latest
IMAGE      ?= $(IMAGE_NAME):$(IMAGE_TAG)
CLUSTER    ?= simple-web-app

.PHONY: docker-build
docker-build: ## Build the Docker image
	docker build -t $(IMAGE) .

.PHONY: kind-load
kind-load: docker-build ## Build and load the image into the kind cluster
	kind load docker-image $(IMAGE) --name $(CLUSTER)

.PHONY: minikube-load
minikube-load: docker-build ## Build and load the image into the minikube cluster
	minikube image load $(IMAGE) --profile $(CLUSTER)

.PHONY: help
help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  %-16s %s\n", $$1, $$2}'
