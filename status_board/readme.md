# Technical Interview Task - Status Board

This skeleton of a small HTTP service (`status-board`) is provided to help you get
started with the technical interview task. A Dockerfile and Kubernetes manifests
are provided to containerize and deploy it.

The service exposes three endpoints:

* `GET /hello-world` — works out of the box, needs no configuration.
* `GET /settings` — reports the configuration the service runs with. It fails
  until the service is configured.
* `GET /checks` — a skeleton to be implemented: probe the configured upstream
  targets and report their status.

(`GET /healthz` exists as well, as a cheap target for health probes.)

## Repository Structure

This folder contains the following:

* `status_board_go`
  * The service implemented in Go using the plain `net/http` package, with a
    Dockerfile for containerization.
* `kubernetes_manifests`
  * Kubernetes manifests to deploy the service, plus a troubleshooting variant
    and scripts to create a cluster whose CNI enforces `NetworkPolicy`.

## The Tasks

1. **Implement** the `/checks` endpoint (see `status_board_go/README.md`).
2. **Deploy** the service to Kubernetes and get the workload healthy (see
   `kubernetes_manifests/README.md`).

The two tasks are independent — the deployment task does not need `/checks` to
be implemented.
