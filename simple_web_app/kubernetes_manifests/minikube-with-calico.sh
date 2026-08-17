#!/bin/bash
set -o errexit

# Creates a minikube cluster that ENFORCES NetworkPolicy. minikube's default CNI
# does not enforce NetworkPolicy, so it is started with Calico (--cni=calico) —
# the equivalent of kind-with-calico.sh for minikube users.

profile='simple-web-app'

minikube start --profile "${profile}" --cni=calico

echo "Waiting for Calico to be ready..."
kubectl --context "${profile}" -n kube-system rollout status daemonset/calico-node --timeout=180s

echo
echo "Cluster '${profile}' is ready with Calico enforcing NetworkPolicy."
echo "Load the image with: make minikube-load"
