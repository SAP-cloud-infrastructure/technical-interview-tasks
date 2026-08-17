#!/bin/bash
set -o errexit

# Creates a kind cluster that ENFORCES NetworkPolicy (the default kind CNI,
# kindnet, does not), wired to a local registry.
#
# Based on the kind local-registry example, with the default CNI disabled and Calico installed.

reg_name='kind-registry'
reg_port='5001'
cluster_name='simple-web-app'
calico_version='v3.28.2'

# 1. Create registry container unless it already exists
if [ "$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)" != 'true' ]; then
  docker run \
    -d --restart=always -p "127.0.0.1:${reg_port}:5000" --network bridge --name "${reg_name}" \
    registry:2
fi

# 2. Create kind cluster with the default CNI DISABLED so a policy-enforcing CNI
#    (Calico) can be installed. A pod subnet is set for Calico's IPAM.
cat <<EOF | kind create cluster --name "${cluster_name}" --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  disableDefaultCNI: true
  podSubnet: 192.168.0.0/16
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry]
    config_path = "/etc/containerd/certs.d"
EOF

# 3. Add the registry config to the nodes
REGISTRY_DIR="/etc/containerd/certs.d/localhost:${reg_port}"
for node in $(kind get nodes --name "${cluster_name}"); do
  docker exec "${node}" mkdir -p "${REGISTRY_DIR}"
  cat <<EOF | docker exec -i "${node}" cp /dev/stdin "${REGISTRY_DIR}/hosts.toml"
[host."http://${reg_name}:5000"]
EOF
done

# 4. Connect the registry to the cluster network if not already connected
if [ "$(docker inspect -f='{{json .NetworkSettings.Networks.kind}}' "${reg_name}")" = 'null' ]; then
  docker network connect "kind" "${reg_name}"
fi

# 5. Install Calico as the CNI (enforces NetworkPolicy)
kubectl apply -f "https://raw.githubusercontent.com/projectcalico/calico/${calico_version}/manifests/calico.yaml"

echo "Waiting for nodes to become Ready (Calico rollout)..."
kubectl wait --for=condition=Ready nodes --all --timeout=180s
kubectl -n kube-system rollout status daemonset/calico-node --timeout=180s

# 6. Document the local registry
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${reg_port}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF

echo
echo "Cluster '${cluster_name}' is ready with Calico enforcing NetworkPolicy."
