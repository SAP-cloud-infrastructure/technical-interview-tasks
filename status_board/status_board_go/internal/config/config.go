package config

import (
	"log"
	"os"
	"strings"
	"time"
)

const (
	defaultPort         = "8080"
	defaultTargetsFile  = "/etc/status-board/targets.json"
	defaultCheckTimeout = 2 * time.Second
)

// Target is a single upstream that the /checks endpoint has to probe.
type Target struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Config is the runtime configuration of the service. Every value comes from
// the environment, so on Kubernetes it can be supplied by a ConfigMap and a
// Secret.
type Config struct {
	Port         string
	BoardName    string
	Environment  string
	APIToken     string
	CheckTimeout time.Duration
	TargetsFile  string
	Targets      []Target
	PodName      string
	NodeName     string
}

// Load reads the configuration from the environment.
func Load() Config {
	cfg := Config{
		Port:         envOr("PORT", defaultPort),
		BoardName:    os.Getenv("BOARD_NAME"),
		Environment:  os.Getenv("ENVIRONMENT"),
		APIToken:     os.Getenv("API_TOKEN"),
		CheckTimeout: defaultCheckTimeout,
		TargetsFile:  envOr("TARGETS_FILE", defaultTargetsFile),
		Targets:      parseTargets(os.Getenv("CHECK_TARGETS")),
		PodName:      os.Getenv("POD_NAME"),
		NodeName:     os.Getenv("NODE_NAME"),
	}

	if raw := os.Getenv("CHECK_TIMEOUT"); raw != "" {
		timeout, err := time.ParseDuration(raw)
		if err != nil {
			log.Printf("WARNING: CHECK_TIMEOUT %q is not a duration, using %s: %v", raw, defaultCheckTimeout, err)
		} else {
			cfg.CheckTimeout = timeout
		}
	}

	return cfg
}

// Missing returns the names of the required settings that have no value.
func (c Config) Missing() []string {
	var missing []string
	for _, setting := range []struct{ name, value string }{
		{"BOARD_NAME", c.BoardName},
		{"ENVIRONMENT", c.Environment},
		{"API_TOKEN", c.APIToken},
	} {
		if strings.TrimSpace(setting.value) == "" {
			missing = append(missing, setting.name)
		}
	}
	return missing
}

// TargetsFileFound reports whether the configured targets file exists. The
// service does not read it yet (see the /checks TODO), but /settings reports it
// so that a targets file which never reached the container is easy to spot.
func (c Config) TargetsFileFound() bool {
	if c.TargetsFile == "" {
		return false
	}
	_, err := os.Stat(c.TargetsFile)
	return err == nil
}

// parseTargets parses CHECK_TARGETS, a comma-separated list of name=url pairs.
func parseTargets(raw string) []Target {
	var targets []Target
	for _, entry := range strings.Split(raw, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		name, url, found := strings.Cut(entry, "=")
		if !found {
			name, url = entry, entry
		}

		targets = append(targets, Target{
			Name: strings.TrimSpace(name),
			URL:  strings.TrimSpace(url),
		})
	}
	return targets
}

func envOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
