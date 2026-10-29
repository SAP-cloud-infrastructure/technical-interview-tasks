package handlers

import (
	"encoding/json"
	"net/http"

	"status-board/internal/config"
)

// Handler serves the endpoints that depend on the service configuration.
type Handler struct {
	cfg config.Config
}

// New returns a Handler for the given configuration.
func New(cfg config.Config) *Handler {
	return &Handler{cfg: cfg}
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string   `json:"error"`
	Missing []string `json:"missing,omitempty"`
	Hint    string   `json:"hint,omitempty"`
}

// SettingsResponse is what GET /settings returns once the service is configured.
type SettingsResponse struct {
	BoardName        string          `json:"board_name"`
	Environment      string          `json:"environment"`
	APIToken         string          `json:"api_token"`
	Port             string          `json:"port"`
	CheckTimeout     string          `json:"check_timeout"`
	Targets          []config.Target `json:"targets"`
	TargetsFile      string          `json:"targets_file"`
	TargetsFileFound bool            `json:"targets_file_found"`
	Pod              string          `json:"pod,omitempty"`
	Node             string          `json:"node,omitempty"`
}

// CheckResult is the outcome of probing a single target.
type CheckResult struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Status     string `json:"status"`
	HTTPStatus int    `json:"http_status,omitempty"`
	LatencyMS  int64  `json:"latency_ms"`
	Error      string `json:"error,omitempty"`
}

// ChecksResponse is what GET /checks returns.
type ChecksResponse struct {
	Board   string        `json:"board"`
	Results []CheckResult `json:"results"`
	Up      int           `json:"up"`
	Down    int           `json:"down"`
}

// Healthz handles the /healthz endpoint. It only reports that the process is
// serving requests, which makes it the right target for liveness and readiness
// probes.
func Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, MessageResponse{Message: "ok"})
}

// HelloWorld handles the /hello-world endpoint. It needs no configuration.
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, MessageResponse{Message: "Hello World"})
}

// Settings handles the /settings endpoint. It reports the configuration the
// service is running with and fails with 503 for as long as one of the required
// settings (BOARD_NAME, ENVIRONMENT, API_TOKEN) is missing. The value of
// API_TOKEN is never returned, only whether it is set.
func (h *Handler) Settings(w http.ResponseWriter, r *http.Request) {
	if missing := h.cfg.Missing(); len(missing) > 0 {
		writeJSON(w, http.StatusServiceUnavailable, ErrorResponse{
			Error:   "incomplete configuration",
			Missing: missing,
			Hint:    "supply the missing settings via the environment",
		})
		return
	}

	writeJSON(w, http.StatusOK, SettingsResponse{
		BoardName:        h.cfg.BoardName,
		Environment:      h.cfg.Environment,
		APIToken:         "configured",
		Port:             h.cfg.Port,
		CheckTimeout:     h.cfg.CheckTimeout.String(),
		Targets:          h.cfg.Targets,
		TargetsFile:      h.cfg.TargetsFile,
		TargetsFileFound: h.cfg.TargetsFileFound(),
		Pod:              h.cfg.PodName,
		Node:             h.cfg.NodeName,
	})
}

// Checks handles the /checks endpoint.
//
// TODO(candidate): implement this endpoint. It has to probe every configured
// target and report the result:
//
//   - probe each target in h.cfg.Targets with an HTTP GET
//   - probe the targets concurrently, not one after another
//   - give every probe at most h.cfg.CheckTimeout (see context.WithTimeout)
//   - a target counts as up when it answers with a 2xx status code
//   - fill in and return ChecksResponse as JSON, with the HTTP status code
//     200 when every target is up, 207 when only some are, and 503 when none
//     is up (200 as well when nothing is configured to probe)
//   - a slow or unreachable target must not fail the whole response
func (h *Handler) Checks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, ErrorResponse{
		Error: "not implemented",
		Hint:  "implement Handler.Checks in internal/handlers/handlers.go",
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
