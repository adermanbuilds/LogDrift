package main

import (
	"encoding/json"
	"os"
	"regexp"
	"slices"
)

// Config holds all user-configurable options for logdrift.
// Keep defaults sensible so the tool works out-of-the-box.
type Config struct {
	// Detection patterns
	ErrorPatterns      []string          `json:"error_patterns"`      // Regexes matching error-like log lines
	SlowPatterns       []string          `json:"slow_patterns"`       // Regexes matching slow/timeout indicators
	SuspiciousPatterns []string          `json:"suspicious_patterns"` // Regexes for security-related or suspicious events
	CustomPatterns     map[string]string `json:"custom_patterns"`     // User-defined named patterns

	// Thresholds
	AnomalyThreshold int `json:"anomaly_threshold"` // Percent or score threshold to flag anomalies
	SlowThresholdMs  int `json:"slow_threshold_ms"` // Response-time threshold in milliseconds

	// Output preferences
	ShowDebug     bool `json:"show_debug"`     // Include debug-level output when true
	ShowInfo      bool `json:"show_info"`      // Include info-level output when true
	ColorOutput   bool `json:"color_output"`   // Enable colored terminal output
	CompactOutput bool `json:"compact_output"` // Compact output format for denser logs

	// Filtering
	IncludeComponents []string `json:"include_components"` // If set, only these components are shown
	ExcludeComponents []string `json:"exclude_components"` // Components to always hide
	MinLevel          string   `json:"min_level"`          // Minimum level to display (e.g., "WARN")
}

// DefaultConfig returns sensible defaults so users get useful results without a config file.
func DefaultConfig() *Config {
	return &Config{
		ErrorPatterns: []string{
			`(?i)(error|exception|fatal|panic|traceback)`, // common error keywords (case-insensitive)
			`5\d{2}\s`,                             // HTTP 5xx status codes
			`(?i)(failed|failure|timeout)`,         // failure/timeout words
			`(?i)(connection refused|broken pipe)`, // network-related failures
			`(?i)(out of memory|oom)`,              // OOM indicators
			`(?i)(segmentation fault|core dump)`,   // native crash signs
		},
		SlowPatterns: []string{
			`(?i)(\d+(?:\.\d+)?)(ms|s)\s.*?(slow|timeout)`, // explicit "slow" with a duration
			`response_time[=:]\s*([5-9]\d{2,}|\d{4,})`,     // numeric response_time > ~500ms
			`latency[=:]\s*([5-9]\d{2,}|\d{4,})`,           // numeric latency > ~500ms
		},
		SuspiciousPatterns: []string{
			`(?i)(DROP|REJECT|blocked?|denied|unauthorized)`, // access/DB denial words
			`(?i)(brute.?force|injection|xss|csrf)`,          // common attack vectors
			`(?i)(malware|virus|exploit)`,                    // malware/attack terms
		},
		CustomPatterns: make(map[string]string),

		AnomalyThreshold: 40,  // score threshold for anomaly alerts
		SlowThresholdMs:  500, // default slow threshold of 500ms

		ShowDebug:     false,
		ShowInfo:      false,
		ColorOutput:   true,
		CompactOutput: false,

		MinLevel: "WARN", // conservative default to reduce noise
	}
}

// LoadConfig reads a JSON config file and merges it with defaults.
// If the file doesn't exist, defaults are returned (no error).
func LoadConfig(path string) (*Config, error) {
	// Start with defaults
	cfg := DefaultConfig()

	// If config file doesn't exist, return defaults
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	// Read and parse config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Merge with defaults (Unmarshal fills fields present in JSON)
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// SaveConfig writes the current config to a JSON file with indentation.
func (c *Config) SaveConfig(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// CompilePatterns converts string regex patterns into compiled *regexp.Regexp slices.
// Returns compiled regex slices for errors, slow events, and suspicious events.
func (c *Config) CompilePatterns() ([]*regexp.Regexp, []*regexp.Regexp, []*regexp.Regexp, error) {
	var errorRegexes []*regexp.Regexp
	var slowRegexes []*regexp.Regexp
	var suspiciousRegexes []*regexp.Regexp

	// Compile error patterns
	for _, pattern := range c.ErrorPatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, nil, nil, err // bubble up first compile error
		}
		errorRegexes = append(errorRegexes, re)
	}

	// Compile slow patterns
	for _, pattern := range c.SlowPatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, nil, nil, err
		}
		slowRegexes = append(slowRegexes, re)
	}

	// Compile suspicious patterns
	for _, pattern := range c.SuspiciousPatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, nil, nil, err
		}
		suspiciousRegexes = append(suspiciousRegexes, re)
	}

	return errorRegexes, slowRegexes, suspiciousRegexes, nil
}

// ShouldShow decides whether a log entry should be displayed based on level and component filters.
func (c *Config) ShouldShow(entry LogEntry) bool {
	// Check minimum level
	minLevel := c.parseMinLevel()
	if entry.Level < minLevel {
		return false
	}

	// Check component filters
	component := extractComponent(entry.Raw)

	// If IncludeComponents is non-empty, only show listed components
	if len(c.IncludeComponents) > 0 {
		found := slices.Contains(c.IncludeComponents, component)
		if !found {
			return false
		}
	}

	// ExcludeComponents always hides matching components
	return !slices.Contains(c.ExcludeComponents, component)
}

// parseMinLevel maps the MinLevel string to the LogLevel enum.
// Unknown values default to DEBUG for maximum verbosity.
func (c *Config) parseMinLevel() LogLevel {
	switch c.MinLevel {
	case "FATAL":
		return FATAL
	case "ERROR":
		return ERROR
	case "WARN":
		return WARN
	case "INFO":
		return INFO
	default:
		return DEBUG
	}
}

// extractComponent pulls a [component] token from structured log lines.
// Returns an empty string if no component bracket is present.
func extractComponent(line string) string {
	re := regexp.MustCompile(`\[([^\]]+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
