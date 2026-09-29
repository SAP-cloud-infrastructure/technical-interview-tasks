# Status Board - Go Version

## API Endpoints

The application provides the following endpoints:

1. **GET /hello-world** - Returns a simple "Hello World" message. Needs no configuration.
2. **GET /settings** - Returns the configuration the service is running with.
   Answers `503` with the list of missing settings until `BOARD_NAME`,
   `ENVIRONMENT` and `API_TOKEN` are set. The value of `API_TOKEN` is never
   returned, only whether it is set.
3. **GET /checks** - Skeleton to be implemented: probe every configured target
   and report the result.
4. **GET /healthz** - Returns `ok` as long as the process serves requests. Cheap
   enough for liveness and readiness probes.

## Steps:

1. Clone the repository

2. Run the project locally:
   ```bash
   go run .
   ```
   The server will start on `http://localhost:8080`

3. Test the endpoints:
   ```bash
   # Test hello-world endpoint
   curl http://localhost:8080/hello-world

   # Test settings endpoint (503 until the service is configured)
   curl -i http://localhost:8080/settings

   # Test checks endpoint
   curl http://localhost:8080/checks
   ```

4. Build the docker image with tag `local/status-board`:
   ```bash
   make docker-build
   ```

5. Run the docker image:
   ```bash
   docker run --rm -p 8080:8080 local/status-board
   ```

## Environment Variables

| Variable         | Required | Default                            | Used for                                             |
|------------------|----------|------------------------------------|------------------------------------------------------|
| `PORT`           | no       | `8080`                             | Port the server listens on                           |
| `BOARD_NAME`     | yes      | —                                  | `/settings`                                          |
| `ENVIRONMENT`    | yes      | —                                  | `/settings`                                          |
| `API_TOKEN`      | yes      | —                                  | `/settings`; supply it from a Secret, not a ConfigMap |
| `CHECK_TARGETS`  | no       | empty                              | `/checks`: comma-separated `name=url` pairs          |
| `CHECK_TIMEOUT`  | no       | `2s`                               | `/checks`: timeout per target                        |
| `TARGETS_FILE`   | no       | `/etc/status-board/targets.json`   | `/settings` reports whether this file is readable    |
| `POD_NAME`       | no       | empty                              | `/settings`; from the downward API on Kubernetes     |
| `NODE_NAME`      | no       | empty                              | `/settings`; from the downward API on Kubernetes     |

Example:
```bash
export BOARD_NAME="My Board"
export ENVIRONMENT=dev
export API_TOKEN=not-a-real-token
export CHECK_TARGETS="example=https://example.com,local=http://localhost:8080/hello-world"
go run .
```

Or with Docker:
```bash
docker run --rm -p 8080:8080 \
  -e BOARD_NAME="My Board" \
  -e ENVIRONMENT=dev \
  -e API_TOKEN=not-a-real-token \
  -e CHECK_TARGETS="local=http://localhost:8080/hello-world" \
  local/status-board
```

## Task: implement `/checks`

`Handler.Checks` in `internal/handlers/handlers.go` is a skeleton that answers
`501 Not Implemented`. The `TODO` above it states the requirements: probe every
configured target concurrently, apply the configured timeout per target, and
aggregate the results into `ChecksResponse`. The response types are already
defined; adjust them if you have a better idea.

`make run` and `curl "http://localhost:8080/checks"` are enough to try it out.
