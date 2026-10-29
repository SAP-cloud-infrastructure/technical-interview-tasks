# Deploy workloads to Kubernetes

## Steps:

1. Clone the repository

2. Create the cluster with a CNI that enforces `NetworkPolicy` (the default CNI
   of kind and minikube does not). Use whichever tool you have:

   **kind:**

   ```bash
   ./kind-with-calico.sh
   ```

   **minikube:**

   ```bash
   ./minikube-with-calico.sh
   ```

3. Build and load the Docker image for the Go app from `simple_web_app_go/`:

   **kind:**

   ```bash
   cd ../simple_web_app_go
   make kind-load
   ```

   **minikube:**

   ```bash
   cd ../simple_web_app_go
   make minikube-load
   ```

   These targets build the image (`docker build`) and load it into the cluster.
   To only build the image without loading it, run `make docker-build`.

4. Deploy manifests to Kubernetes

## Manifests

- `kubernetes-manifests.yaml` — troubleshooting exercise for a cluster **without** `NetworkPolicy` enforcement (Docker Desktop, or kind/minikube with their default CNI).
  `NetworkPolicy`.
- `kubernetes-manifests-troubleshoot-kind.yaml` — troubleshooting exercise (kind).
- `kubernetes-manifests-troubleshoot-minikube.yaml` — troubleshooting exercise (minikube).

## Troubleshooting exercise

The `kubernetes-manifests*.yaml` files are deployments used for the troubleshooting
exercise. They run the same application; the goal is to get the workload running and
healthy.

`kubernetes-manifests.yaml` targets a cluster without `NetworkPolicy` enforcement, so
it does not need the Calico setup from the steps above — a plain Docker Desktop
cluster (or a default kind/minikube cluster) is enough. Load the image, then deploy:

```bash
kubectl apply -f kubernetes-manifests.yaml
```

The `kubernetes-manifests-troubleshoot-*.yaml` files target a cluster that enforces
`NetworkPolicy`. Follow the steps above to create the cluster and load the image, then
deploy the manifest for your tool:

```bash
# kind
kubectl apply -f kubernetes-manifests-troubleshoot-kind.yaml

# minikube
kubectl apply -f kubernetes-manifests-troubleshoot-minikube.yaml
```
