# INTERVIEW SOLUTION - Complete Implementation

This directory contains the **complete solution** for the Go web application interview task.

## 📁 Files Overview

### Solution Files (NEW)
- **`handlers_solution.go`** - Complete implementation with detailed comments
- **`main_solution.go`** - Test runner that uses the solution
- **`test_solution.sh`** - Automated test script
- **`SOLUTION_README.md`** - This file

### Original Files (provided to applicant)
- **`handlers.go`** - Original skeleton with empty `GetRepoList` function
- **`main.go`** - Original main file
- **`auth.go`** - BasicAuth middleware (already complete)

---

## ✅ What Was Implemented

The solution implements the missing `/repo-list/{org_name}` endpoint that:

1. **Extracts path parameters** using Go 1.22's `r.PathValue("org_name")`
2. **Extracts query parameters** for optional filtering with `r.URL.Query().Get("repo_filter")`
3. **Calls GitHub API** at `https://api.github.com/orgs/{org}/repos`
4. **Handles all errors properly**:
   - Network failures → 503 Service Unavailable
   - Organization not found → 404 Not Found
   - Invalid requests → 400 Bad Request
   - GitHub API errors → 502 Bad Gateway
5. **Filters repositories** by name (case-insensitive substring match)
6. **Returns JSON response** with proper structure

---

## 🚀 Quick Test

### Option 1: Run the automated test script

```bash
./test_solution.sh
```

This will:
- Start the server
- Test all endpoints
- Verify error handling
- Check response formats
- Display test results

### Option 2: Manual testing

```bash
# Start the server with the solution
go run main_solution.go

# In another terminal, test the endpoints:

# Test hello-world (already working)
curl http://localhost:8080/hello-world

# Test repo-list without filter
curl http://localhost:8080/repo-list/golang

# Test repo-list with filter
curl "http://localhost:8080/repo-list/golang?repo_filter=go"

# Test non-existent organization (should return 404)
curl http://localhost:8080/repo-list/nonexistentorg12345

# Test protected endpoint (will need AUTH_USERNAME and AUTH_PASSWORD set)
curl http://localhost:8080/protected
```

---

## 📊 Expected Response Format

### Success Response for `/repo-list/golang?repo_filter=go`

```json
{
  "organization": "golang",
  "count": 2,
  "repositories": [
    {
      "name": "go",
      "full_name": "golang/go",
      "description": "The Go programming language",
      "html_url": "https://github.com/golang/go",
      "stargazers_count": 132767,
      "language": "Go"
    },
    {
      "name": "gofrontend",
      "full_name": "golang/gofrontend",
      "description": "Go compiler frontend (gccgo)",
      "html_url": "https://github.com/golang/gofrontend",
      "stargazers_count": 888,
      "language": "Go"
    }
  ]
}
```

### Error Response for Non-existent Organization

```json
{
  "error": "Organization 'nonexistentorg12345' not found"
}
```

---

## 🧪 Test Results

All tests passed:
- ✅ `/hello-world` endpoint works
- ✅ `/repo-list/{org_name}` fetches repositories
- ✅ `?repo_filter=` parameter filters results
- ✅ Returns 404 for non-existent organizations
- ✅ Returns proper error messages
- ✅ JSON response has correct structure
- ✅ HTTP status codes are appropriate
- ✅ Response time is acceptable (<10s)

---

## 🔍 Key Implementation Details

### 1. HTTP Client with Timeout
```go
client := &http.Client{
    Timeout: 10 * time.Second,
}
```
**Why?** Without timeout, requests can hang forever if GitHub is slow.

### 2. Required Headers
```go
req.Header.Set("User-Agent", "simple-web-app")
req.Header.Set("Accept", "application/vnd.github.v3+json")
```
**Why?** GitHub API requires User-Agent header or it returns 403.

### 3. Always Close Response Body
```go
defer resp.Body.Close()
```
**Why?** Prevents memory leaks and connection exhaustion.

### 4. Proper Error Handling
- Network errors → 503
- Not found → 404
- Bad request → 400
- GitHub errors → 502

### 5. Case-Insensitive Filtering
```go
if strings.Contains(strings.ToLower(repo.Name), filterLower) {
    filteredRepos = append(filteredRepos, repo)
}
```

---

## 🐳 Docker Testing

```bash
# Build the Docker image
docker build -t simple-web-server .

# Run the container
docker run --rm -p 8080:8080 simple-web-server

# Test it
curl http://localhost:8080/hello-world
curl "http://localhost:8080/repo-list/golang?repo_filter=go"
```

---

## ☸️ Kubernetes Testing

```bash
cd ../kubernetes_manifests

# Create Kind cluster with registry
./kind-with-registry.sh

# Tag and push image
docker tag simple-web-server localhost:5001/simple-web-app
docker push localhost:5001/simple-web-app

# IMPORTANT: Fix the NetworkPolicy first!
# Edit kubernetes-manifests.yaml and change ingress: [] to allow traffic
# OR remove the NetworkPolicy section entirely for testing

# Deploy
kubectl apply -f kubernetes-manifests.yaml

# Check pods
kubectl get pods

# Port-forward to test
kubectl port-forward svc/simple-web-app 8080:80

# Test
curl http://localhost:8080/hello-world
```

### ⚠️ Known Issue: NetworkPolicy Blocks Traffic

The provided `kubernetes-manifests.yaml` has a NetworkPolicy with `ingress: []` which **blocks all incoming traffic**.

This is intentional to test if the applicant can debug Kubernetes networking issues!

**Solution:** Update the NetworkPolicy to allow traffic:

```yaml
ingress:
- from:
  - podSelector: {}  # Allow from all pods in namespace
  ports:
  - protocol: TCP
    port: 8080
```

---

## 📝 Interview Evaluation Points

### What the Solution Demonstrates

#### Go Knowledge
- ✅ Uses Go 1.22 routing features (`PathValue`)
- ✅ Proper HTTP client usage with timeout
- ✅ Struct tags for JSON marshaling
- ✅ Error handling patterns
- ✅ defer for resource cleanup

#### HTTP/REST API Knowledge
- ✅ Proper HTTP status codes
- ✅ Content-Type headers
- ✅ Request/response structure
- ✅ External API integration

#### Security Awareness
- ✅ Input validation (empty org name)
- ✅ Timeout to prevent DoS
- ✅ Uses existing BasicAuth middleware properly
- ✅ Aware of timing attacks (`subtle.ConstantTimeCompare` in auth.go)

#### Code Quality
- ✅ Clear variable names
- ✅ Comprehensive comments
- ✅ Proper error messages
- ✅ Clean code structure

#### Docker/Kubernetes
- ✅ Understands multi-stage builds
- ✅ Can debug NetworkPolicy issues
- ✅ Understands Kubernetes networking

---

## 🚀 Production Improvements to Discuss

### 1. Rate Limiting
- Check `X-RateLimit-Remaining` header from GitHub
- Implement circuit breaker pattern
- Return 429 when rate limited

### 2. Caching
- Cache GitHub responses for 5-10 minutes
- Use Redis or in-memory cache
- Reduces API calls and improves latency

### 3. Pagination
- GitHub returns max 30 repos per page
- Implement pagination for large organizations
- Parse "Link" header for next page

### 4. Observability
- Structured logging (zerolog, zap)
- Prometheus metrics (request count, latency, errors)
- Distributed tracing (OpenTelemetry)

### 5. Authentication
- Add GitHub token via environment variable
- Increases rate limit from 60 to 5000 req/hour
- Required for private repositories

### 6. Input Validation
- Validate org name format (alphanumeric + hyphens)
- Limit filter length
- Prevent injection attacks

### 7. Graceful Shutdown
- Handle SIGTERM/SIGINT
- Drain existing connections
- Important for Kubernetes deployments

### 8. Health Checks
- Add `/health` and `/ready` endpoints
- Kubernetes liveness/readiness probes
- Return status of dependencies

---

## 📚 Additional Resources

- [GitHub REST API Docs](https://docs.github.com/en/rest)
- [Go 1.22 Routing](https://go.dev/blog/routing-enhancements)
- [Go HTTP Client Best Practices](https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779)
- [Kubernetes NetworkPolicy](https://kubernetes.io/docs/concepts/services-networking/network-policies/)

---

## ✨ Summary

This solution demonstrates a **senior-level** understanding of:
- Go web development
- External API integration
- Error handling and edge cases
- Security best practices
- Docker and Kubernetes
- Production-ready considerations

The code is well-commented, tested, and ready for discussion in the interview!
