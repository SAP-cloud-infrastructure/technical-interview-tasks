# Deploy the status board to Kubernetes

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

   Both create a cluster named `status-board`.

3. Build and load the Docker image for the Go app from `status_board_go/`:

   **kind:**

   ```bash
   cd ../status_board_go
   make kind-load
   ```

   **minikube:**

   ```bash
   cd ../status_board_go
   make minikube-load
   ```

   These targets build the image (`docker build`) and load it into the cluster.
   To only build the image without loading it, run `make docker-build`.

4. Deploy manifests to Kubernetes

## Manifests

- `kubernetes-manifests-working.yaml` — a clean deployment that comes up and
  stays healthy as-is. It carries no `NetworkPolicy`, so it also works on a
  cluster whose CNI does not enforce them. Use this to get the workload up
  quickly.
- `kubernetes-manifests-troubleshoot.yaml` — troubleshooting exercise. It needs a
  cluster with a policy-enforcing CNI (step 2).

Both deploy three workloads:

| Workload              | Purpose                                                                 |
|-----------------------|-------------------------------------------------------------------------|
| `status-board`        | The app, configured from a ConfigMap and a Secret, behind a `Service`.  |
| `echo`                | The same image without configuration, a stand-in upstream for `/checks`.|
| `status-board-tester` | Calls the app through its `Service` in a loop and exits when a call fails. |

## Troubleshooting exercise

Deploy `kubernetes-manifests-troubleshoot.yaml` and get the workload running,
healthy and reachable:

```bash
kubectl apply -f kubernetes-manifests-troubleshoot.yaml
```

The manifest deploys the same application as the working variant, but it does not
come up. Several independent problems are in your way; fixing one uncovers the
next. Change whatever you need to — manifests, objects in the cluster, or both —
and say what you changed and why.

You are done when:

- every pod is `Running` and ready,
- `status-board-tester` stops restarting (it keeps calling `/hello-world` and
  `/settings` through the `Service` and exits as soon as a call fails),
- `GET /settings` answers `200` through the `Service`.

Note that `/checks` answers `501` until it is implemented — that is the other
task, not a problem to fix here.

## Clean up

**kind:**

```bash
kind delete cluster --name status-board
```

**minikube:**

```bash
minikube delete --profile status-board
```
