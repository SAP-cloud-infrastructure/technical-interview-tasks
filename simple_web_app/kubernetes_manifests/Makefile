CLUSTER_NAME           ?= interview
REGISTRY_NAME          ?= registry
REGISTRY_PORT          ?= 5000
REGISTRY_INTERNAL_PORT := 5000
REGISTRY_CONFIG        := $(abspath registry-config.json)
KIND_CONFIG            := $(abspath kind-config.yaml)

REGISTRY_IMAGE := ghcr.io/project-zot/zot:latest

.PHONY: all cluster registry connect-registry teardown

## Bring up cluster + registry in one shot
all: cluster registry connect-registry

## Create the kind cluster
cluster:
	@if kind get clusters 2>/dev/null | grep -qw "$(CLUSTER_NAME)"; then \
	  echo "Cluster $(CLUSTER_NAME) already exists, skipping."; \
	else \
	  echo "Creating kind cluster $(CLUSTER_NAME)..."; \
	  kind create cluster --name "$(CLUSTER_NAME)" --config "$(KIND_CONFIG)"; \
	fi

## Start the zot registry container
registry:
	@if docker inspect "$(REGISTRY_NAME)" &>/dev/null; then \
	  echo "Registry container already exists, skipping."; \
	else \
	  echo "Starting zot registry ($(REGISTRY_IMAGE))..."; \
	  docker run -d \
	    --name "$(REGISTRY_NAME)" \
	    --restart=unless-stopped \
	    -p "$(REGISTRY_PORT):5000" \
	    -v "$(REGISTRY_CONFIG):/etc/zot/config.json:ro" \
	    "$(REGISTRY_IMAGE)"; \
	fi

## Connect registry to kind network and configure containerd on each node
connect-registry:
	@echo "Connecting registry to kind network..."
	@if docker network inspect kind --format '{{range .Containers}}{{.Name}} {{end}}' 2>/dev/null \
	    | grep -qw "$(REGISTRY_NAME)"; then \
	  echo "Registry already on kind network."; \
	else \
	  docker network connect kind "$(REGISTRY_NAME)"; \
	fi
	@echo "Configuring containerd on cluster nodes..."
	@for node in $$(kind get nodes --name "$(CLUSTER_NAME)"); do \
	  echo "  Configuring $$node..."; \
	  if ! docker exec "$$node" grep -q 'config_path.*certs\.d' /etc/containerd/config.toml 2>/dev/null; then \
	    docker exec "$$node" bash -c \
	      'printf "\n[plugins.\"io.containerd.grpc.v1.cri\".registry]\n  config_path = \"/etc/containerd/certs.d\"\n" \
	       >> /etc/containerd/config.toml'; \
	    docker exec "$$node" systemctl restart containerd; \
	  fi; \
	  docker exec "$$node" mkdir -p "/etc/containerd/certs.d/$(REGISTRY_NAME):$(REGISTRY_INTERNAL_PORT)"; \
	  printf 'server = "http://$(REGISTRY_NAME):$(REGISTRY_INTERNAL_PORT)"\n\n[host."http://$(REGISTRY_NAME):$(REGISTRY_INTERNAL_PORT)"]\n  capabilities = ["pull", "resolve", "push"]\n' \
	    | docker exec -i "$$node" tee "/etc/containerd/certs.d/$(REGISTRY_NAME):$(REGISTRY_INTERNAL_PORT)/hosts.toml" >/dev/null; \
	done
	@printf 'apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: local-registry-hosting\n  namespace: kube-public\ndata:\n  localRegistryHosting.v1: |\n    host: "$(REGISTRY_NAME):$(REGISTRY_INTERNAL_PORT)"\n    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"\n' \
	  | kubectl --context "kind-$(CLUSTER_NAME)" apply -f -

## Tear down cluster and registry
teardown:
	@echo "Deleting kind cluster $(CLUSTER_NAME)..."
	@kind delete cluster --name "$(CLUSTER_NAME)" || true
	@echo "Removing registry container..."
	@docker rm -f "$(REGISTRY_NAME)" || true