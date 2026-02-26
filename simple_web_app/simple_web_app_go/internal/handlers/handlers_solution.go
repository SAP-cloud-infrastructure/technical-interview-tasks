package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ============================================================================
// SOLUTION FILE - Complete implementation of the /repo-list endpoint
// This file contains the correct answer for the interview task
// ============================================================================

// GitHubRepo represents a repository from the GitHub API response
// We only include the fields we're interested in, but GitHub returns many more
type GitHubRepo struct {
	Name        string `json:"name"`         // Repository name (e.g., "go")
	FullName    string `json:"full_name"`    // Full name with org (e.g., "golang/go")
	Description string `json:"description"`  // Repository description
	HTMLURL     string `json:"html_url"`     // Web URL to the repository
	StarCount   int    `json:"stargazers_count"` // Number of stars
	Language    string `json:"language"`     // Primary programming language
}

// RepoListResponse is the structure we return to the client
type RepoListResponse struct {
	Organization string       `json:"organization"` // The organization name queried
	Count        int          `json:"count"`        // Number of repositories returned
	Repositories []GitHubRepo `json:"repositories"` // List of repositories
}

// ErrorResponse represents an error message returned to the client
type ErrorResponse struct {
	Error string `json:"error"`
}

// GetRepoListSolution is the complete implementation of the /repo-list/{org_name}[?repo_filter=filter] endpoint
// This endpoint fetches repositories from a GitHub organization and optionally filters them
//
// To use this solution, replace the GetRepoList function in handlers.go with this implementation
func GetRepoListSolution(w http.ResponseWriter, r *http.Request) {
	// ========================================================================
	// STEP 1: Extract the organization name from the URL path parameter
	// ========================================================================
	// Using Go 1.22's new routing feature: PathValue extracts the {org_name} from the route
	// Example: /repo-list/golang -> orgName = "golang"
	orgName := r.PathValue("org_name")

	// Validate that the organization name is not empty
	// This is important for security - we don't want to make requests with empty values
	// An empty org name would result in a malformed GitHub API URL
	if orgName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Organization name is required"})
		return
	}

	// ========================================================================
	// STEP 2: Extract the optional query parameter for filtering
	// ========================================================================
	// Example: ?repo_filter=go will filter repos whose names contain "go"
	// If not provided, all repositories will be returned
	repoFilter := r.URL.Query().Get("repo_filter")

	// ========================================================================
	// STEP 3: Create an HTTP client with a timeout
	// ========================================================================
	// IMPORTANT: Never use http.Get() or default client without timeout!
	// Without a timeout, if GitHub is slow or unresponsive, the goroutine will hang forever
	// This can lead to resource exhaustion and eventually crash your application
	//
	// Why 10 seconds?
	// - GitHub API typically responds in < 1 second
	// - 10 seconds allows for network latency and slow responses
	// - Not too long that users wait forever for errors
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// ========================================================================
	// STEP 4: Construct the GitHub API URL
	// ========================================================================
	// GitHub's REST API endpoint for listing organization repositories
	// Documentation: https://docs.github.com/en/rest/repos/repos#list-organization-repositories
	//
	// Note: This endpoint returns public repositories only (without authentication)
	// It returns up to 30 repositories per page by default
	url := fmt.Sprintf("https://api.github.com/orgs/%s/repos", orgName)

	// ========================================================================
	// STEP 5: Create the HTTP request
	// ========================================================================
	// We use NewRequest instead of client.Get() so we can add custom headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// This error is rare - usually means URL construction failed
		// Could happen if orgName contains invalid characters
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to create request"})
		return
	}

	// ========================================================================
	// STEP 6: Add required headers
	// ========================================================================
	// GitHub API REQUIRES a User-Agent header or it will reject the request with 403
	// This is GitHub's way of tracking which applications are using their API
	req.Header.Set("User-Agent", "simple-web-app")

	// The Accept header tells GitHub we want the v3 API format (JSON)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// OPTIONAL: Add authentication for higher rate limits
	// Without auth: 60 requests/hour per IP
	// With auth: 5000 requests/hour per user
	//
	// To enable authentication, set GITHUB_TOKEN environment variable:
	// githubToken := os.Getenv("GITHUB_TOKEN")
	// if githubToken != "" {
	//     req.Header.Set("Authorization", "token " + githubToken)
	// }

	// ========================================================================
	// STEP 7: Execute the HTTP request
	// ========================================================================
	resp, err := client.Do(req)
	if err != nil {
		// Network error: GitHub is unreachable, DNS failed, timeout exceeded, etc.
		// This is different from HTTP errors (which have status codes)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to connect to GitHub API"})
		return
	}

	// CRITICAL: Always close the response body to prevent memory leaks
	// defer ensures this runs even if we return early due to errors
	// If you don't close the body, the HTTP connection stays open and resources leak
	defer resp.Body.Close()

	// ========================================================================
	// STEP 8: Check the HTTP status code from GitHub
	// ========================================================================
	// Handle different error cases with appropriate HTTP status codes

	// 404 Not Found: Organization doesn't exist
	if resp.StatusCode == http.StatusNotFound {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("Organization '%s' not found", orgName)})
		return
	}

	// Any other non-200 status code
	// Could be:
	// - 403 Forbidden: Rate limit exceeded or missing User-Agent
	// - 500 Internal Server Error: GitHub is having issues
	// - 502 Bad Gateway: GitHub proxy issues
	if resp.StatusCode != http.StatusOK {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("GitHub API returned status: %d", resp.StatusCode)})
		return
	}

	// ========================================================================
	// STEP 9: Parse the JSON response from GitHub
	// ========================================================================
	// GitHub returns an array of repository objects
	// json.Decoder is efficient for reading from io.Reader (http response body)
	var repos []GitHubRepo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		// JSON parsing failed - malformed response from GitHub
		// This shouldn't happen unless GitHub API changes or network corruption
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to parse GitHub response"})
		return
	}

	// ========================================================================
	// STEP 10: Apply filtering if the repo_filter parameter was provided
	// ========================================================================
	// We do case-insensitive substring matching on repository names
	// Example: filter="go" will match "go", "golang", "go-tools", etc.
	//
	// Alternative approaches:
	// - Exact match: repo.Name == repoFilter
	// - Prefix match: strings.HasPrefix(repo.Name, repoFilter)
	// - Regex match: regexp.MatchString(pattern, repo.Name)
	filteredRepos := repos
	if repoFilter != "" {
		// Pre-allocate slice with zero length but capacity of original
		// This is more efficient than appending to nil slice
		filteredRepos = make([]GitHubRepo, 0, len(repos))

		// Convert filter to lowercase once (not in the loop)
		filterLower := strings.ToLower(repoFilter)

		for _, repo := range repos {
			// Check if the filter string appears anywhere in the repository name
			// Case-insensitive comparison
			if strings.Contains(strings.ToLower(repo.Name), filterLower) {
				filteredRepos = append(filteredRepos, repo)
			}
		}
	}

	// ========================================================================
	// STEP 11: Prepare the response structure
	// ========================================================================
	response := RepoListResponse{
		Organization: orgName,
		Count:        len(filteredRepos),
		Repositories: filteredRepos,
	}

	// ========================================================================
	// STEP 12: Send the JSON response to the client
	// ========================================================================
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ============================================================================
// PRODUCTION IMPROVEMENTS TO DISCUSS
// ============================================================================
//
// 1. RATE LIMITING AWARENESS
//    - Check X-RateLimit-Remaining header
//    - Return 429 Too Many Requests when GitHub rate limit is hit
//    - Implement exponential backoff for retries
//
// 2. CACHING
//    - Cache GitHub responses for a few minutes
//    - Use Redis or in-memory cache (like go-cache)
//    - Reduces GitHub API calls and improves response time
//
// 3. PAGINATION
//    - GitHub returns max 30 repos per page by default
//    - For orgs with >30 repos, you need to fetch multiple pages
//    - Look for "Link" header in response
//
// 4. LOGGING
//    - Add structured logging (zerolog, zap)
//    - Log all GitHub API calls with timing
//    - Log errors with context (org name, filter, etc.)
//
// 5. METRICS
//    - Add Prometheus metrics
//    - Track: request count, latency, error rate, cache hit rate
//
// 6. CIRCUIT BREAKER
//    - If GitHub is down, don't keep hammering it
//    - Use circuit breaker pattern (like gobreaker)
//    - Fail fast when GitHub is known to be down
//
// 7. INPUT VALIDATION
//    - Validate org name format (alphanumeric, hyphens only)
//    - Prevent injection attacks
//    - Limit filter length to prevent abuse
//
// 8. GRACEFUL DEGRADATION
//    - Return cached data if GitHub is down
//    - Partial results if some pages fail
//
// ============================================================================
