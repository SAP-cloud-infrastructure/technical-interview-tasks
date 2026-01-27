# Simple Web App - Go Version

## API Endpoints

The application provides the following endpoints:

1. **GET /hello-world** - Returns a simple "Hello World" message
2. **GET /repo-list/{org_name}[?repo_filter={filter}**] - Skeleton to be implemented
3. **GET /protected** - Protected endpoint requiring HTTP Basic Authentication

## Steps:

1. Clone the repository
2. Run the project locally:
   ```bash
   poetry install
   poetry run simple_web_app
   ```
   The server will start on `http://localhost:8080`

3. Test the endpoints:
   ```bash
   # Test hello-world endpoint
   curl http://localhost:8080/hello-world
   
   # Test repo-list endpoint
   curl "http://localhost:8080/repo-list/golang?repo_filter=go"
   ```

4. Build the docker image with tag `simple-web-server`:
   ```bash
   docker build -t simple-web-server .
   ```

5. Run the docker image:
   ```bash
   docker run --rm -p 8080:8080 simple-web-server
   ```

## Environment Variables

For the `/protected` endpoint, set the following environment variables:
- `AUTH_USERNAME` - Username for basic authentication
- `AUTH_PASSWORD` - Password for basic authentication

Example:
```bash
export AUTH_USERNAME=admin
export AUTH_PASSWORD=secret
poetry run simple_web_app
```

Or with Docker:
```bash
docker run --rm -p 8080:8080 \
  -e AUTH_USERNAME=admin \
  -e AUTH_PASSWORD=secret \
  simple-web-server
```