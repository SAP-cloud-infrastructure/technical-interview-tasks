#!/bin/bash
set -o errexit

# Creates a kind cluster that ENFORCES NetworkPolicy (the default CNI of kind,
# kindnet, does not): the default CNI is disabled and Calico is installed
# instead.

cluster_name='status-board'
calico_version='v3.28.2'

# 1. Create the kind cluster with the default CNI DISABLED so a policy-enforcing
#    CNI (Calico) can be installed. A pod subnet is set for Calico's IPAM.
cat <<EOM | kind create cluster --name "${cluster_name}" --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  disableDefaultCNI: true
  podSubnet: 192.168.0.0/16
EOM

# 2. Install Calico as the CNI (enforces NetworkPolicy)
kubectl apply -f "https://raw.githubusercontent.com/projectcalico/calico/${calico_version}/manifests/calico.yaml"

echo "Waiting for nodes to become Ready (Calico rollout)..."
kubectl wait --for=condition=Ready nodes --all --timeout=180s
kubectl -n kube-system rollout status daemonset/calico-node --timeout=180s

echo
echo "Cluster '${cluster_name}' is ready with Calico enforcing NetworkPolicy."
echo "Load the image with: make kind-load (in ../status_board_go)"
