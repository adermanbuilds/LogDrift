package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// LogLevel represents severity of log entry
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// LogEntry represents a parsed log line and detection metadata
type LogEntry struct {
	Timestamp time.Time // time when the line was processed
	Level     LogLevel  // inferred severity
	Message   string    // trimmed log message
	Raw       string    // original line text
	IsAnomaly bool      // flagged by detector
	Score     int       // anomaly score (higher = more suspicious)
}

// AnomalyDetector holds compiled patterns, running counters, and config
type AnomalyDetector struct {
	errorPatterns []*regexp.Regexp // compiled error regexes
	slowPatterns  []*regexp.Regexp // compiled slow/latency regexes
	suspiciousIPs []*regexp.Regexp // compiled suspicious activity regexes
	errorKeywords []string         // simple keyword fallback
	lineCount     int              // total processed lines
	errorCount    int              // counted error-level lines
	warnCount     int              // counted warn-level lines
	anomalyCount  int              // flagged anomalies
	startTime     time.Time        // detector start time for stats
	config        *Config          // active configuration
}

// NewDetector initializes the anomaly detector with default config
func NewDetector() *AnomalyDetector {
	return NewDetectorWithConfig(DefaultConfig())
}

// NewDetectorWithConfig initializes with custom config, compiling patterns.
// If compilation fails, warns and uses empty pattern sets to keep running.
func NewDetectorWithConfig(cfg *Config) *AnomalyDetector {
	errorRegexes, slowRegexes, suspiciousRegexes, err := cfg.CompilePatterns()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[1;32mWarning: Error compiling patterns:\033[0m %v\n", err)
		// Use empty patterns if compilation fails so the program remains usable
		errorRegexes = []*regexp.Regexp{}
		slowRegexes = []*regexp.Regexp{}
		suspiciousRegexes = []*regexp.Regexp{}
	}

	return &AnomalyDetector{
		errorPatterns: errorRegexes,
		slowPatterns:  slowRegexes,
		suspiciousIPs: suspiciousRegexes,
		errorKeywords: []string{
			"error",
			"exception",
			"fatal",
			"panic",
			"failed",
			"timeout",
			"refused",
			"denied",
			"invalid",
			"corrupt",
		},
		startTime: time.Now(),
		config:    cfg,
	}
}

// ParseLine analyzes a single log line, infers level, applies patterns and scores.
// This is intentionally fast and heuristic-driven rather than perfectly accurate.
func (d *AnomalyDetector) ParseLine(line string) LogEntry {
	d.lineCount++

	entry := LogEntry{
		Raw:       line,
		Timestamp: time.Now(),
		Score:     0,
	}

	lineLower := strings.ToLower(line)

	// Prefer structured level tokens like "ERROR [component]" or similar.
	levelPattern := regexp.MustCompile(`(?i)\b(FATAL|CRITICAL|ERROR|WARN|INFO|DEBUG)\s+\[`)
	if matches := levelPattern.FindStringSubmatch(line); len(matches) > 1 {
		levelStr := strings.ToUpper(matches[1])
		switch levelStr {
		case "FATAL", "CRITICAL":
			entry.Level = FATAL
			entry.Score += 100
			d.errorCount++
		case "ERROR":
			entry.Level = ERROR
			entry.Score += 50
			d.errorCount++
		case "WARN":
			entry.Level = WARN
			entry.Score += 20
			d.warnCount++
		case "INFO":
			entry.Level = INFO
		default:
			entry.Level = DEBUG
		}
	} else {
		// Fallback: simple keyword checks when no structured level found.
		switch {
		case strings.Contains(lineLower, "fatal") || strings.Contains(lineLower, "critical"):
			entry.Level = FATAL
			entry.Score += 100
			d.errorCount++
		case strings.Contains(lineLower, "error") && !strings.Contains(lineLower, "errors=0"):
			entry.Level = ERROR
			entry.Score += 50
			d.errorCount++
		case strings.Contains(lineLower, "warn"):
			entry.Level = WARN
			entry.Score += 20
			d.warnCount++
		case strings.Contains(lineLower, "info"):
			entry.Level = INFO
		default:
			entry.Level = DEBUG
		}
	}

	// Increase score and mark anomaly if any error patterns match.
	for _, pattern := range d.errorPatterns {
		if pattern.MatchString(line) {
			entry.Score += 30
			entry.IsAnomaly = true
		}
	}

	// Slow/latency indicators also increase score.
	for _, pattern := range d.slowPatterns {
		if pattern.MatchString(line) {
			entry.Score += 25
			entry.IsAnomaly = true
		}
	}

	// Suspicious activity (security/attack indicators) increases score more.
	for _, pattern := range d.suspiciousIPs {
		if pattern.MatchString(line) {
			entry.Score += 40
			entry.IsAnomaly = true
		}
	}

	// Final anomaly decision based on configured threshold.
	if entry.Score >= d.config.AnomalyThreshold {
		entry.IsAnomaly = true
		d.anomalyCount++
	}

	entry.Message = strings.TrimSpace(line)
	return entry
}

// ColorCode returns an ANSI color code for a log level.
func (l LogLevel) ColorCode() string {
	switch l {
	case FATAL:
		return "\033[1;35m" // Magenta
	case ERROR:
		return "\033[1;31m" // Red
	case WARN:
		return "\033[1;33m" // Yellow
	case INFO:
		return "\033[1;32m" // Green
	default:
		return "\033[0;37m" // White (DEBUG/default)
	}
}

// String returns the textual level name.
func (l LogLevel) String() string {
	switch l {
	case FATAL:
		return "FATAL"
	case ERROR:
		return "ERROR"
	case WARN:
		return "WARN"
	case INFO:
		return "INFO"
	default:
		return "DEBUG"
	}
}

// Print outputs a formatted log entry respecting color and compact settings.
// Anomalies are highlighted unless compact output is requested.
func (e LogEntry) Print(cfg *Config) {
	reset := "\033[0m"

	// Disable colors if requested
	levelColor := e.Level.ColorCode()
	if !cfg.ColorOutput {
		levelColor = ""
		reset = ""
	}

	if e.IsAnomaly && !cfg.CompactOutput {
		// Highlight anomalies prominently (red background when colors enabled).
		if cfg.ColorOutput {
			fmt.Printf("\033[41mANOMALY (score: %d)\033[0m ", e.Score)
		} else {
			fmt.Printf("ANOMALY (score: %d) ", e.Score)
		}
	}

	// Print simple level + message line to keep output parseable.
	fmt.Printf("%s[%s]%s %s\n",
		levelColor,
		e.Level.String(),
		reset,
		e.Message,
	)
}

// PrintStats prints a brief summary of runtime statistics at program end.
func (d *AnomalyDetector) PrintStats() {
	runtime := time.Since(d.startTime)

	fmt.Println("\n" + strings.Repeat("─", 60))
	fmt.Printf("\033[1;36mLogDrift Statistics\033[0m\n")
	fmt.Println(strings.Repeat("─", 60))
	fmt.Printf("Runtime:         %s\n", runtime.Round(time.Second))
	fmt.Printf("Lines Processed: %d\n", d.lineCount)
	fmt.Printf("Anomalies:       %d (%.1f%%)\n",
		d.anomalyCount,
		float64(d.anomalyCount)/float64(d.lineCount)*100,
	)
	fmt.Printf("Errors:          %d (%.1f%%)\n",
		d.errorCount,
		float64(d.errorCount)/float64(d.lineCount)*100,
	)
	fmt.Printf("Warnings:        %d (%.1f%%)\n",
		d.warnCount,
		float64(d.warnCount)/float64(d.lineCount)*100,
	)
	fmt.Printf("Rate:            %.0f lines/sec\n",
		float64(d.lineCount)/runtime.Seconds(),
	)
	fmt.Println(strings.Repeat("─", 60))
}

func main() {
	// Parse CLI options
	opts := ParseCLI()

	// Handle special commands (version, generate-config) and exit early if used
	if opts.HandleSpecialCommands() {
		return
	}

	// Load configuration (defaults used when no file provided)
	cfg := DefaultConfig()
	if opts.ConfigFile != "" {
		loadedCfg, err := LoadConfig(opts.ConfigFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		cfg = loadedCfg
	}

	// Apply CLI overrides on top of config file
	opts.ApplyToConfig(cfg)

	// Create detector with the active configuration
	detector := NewDetectorWithConfig(cfg)
	scanner := bufio.NewScanner(os.Stdin)

	// Increase buffer size to handle long log lines safely (1MB limit here).
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 1024*1024)

	if !cfg.CompactOutput {
		fmt.Fprintf(os.Stderr, "\033[1;36mLogDrift v0.1.0 - Monitoring stdin...\033[0m\n\n")
	}

	// Process each line from stdin
	for scanner.Scan() {
		line := scanner.Text()
		entry := detector.ParseLine(line)

		// Apply config-driven filtering (level/component)
		if !cfg.ShouldShow(entry) {
			continue
		}

		// Respect --anomalies-only flag; otherwise show anomalies and WARN+.
		if opts.OnlyAnomalies {
			if entry.IsAnomaly {
				entry.Print(cfg)
			}
		} else {
			// Default: display anomalies and WARN+ messages for visibility.
			if entry.IsAnomaly || entry.Level >= WARN {
				entry.Print(cfg)
			}
		}
	}

	// Handle any scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	// Print final statistics unless compact output was requested
	if !cfg.CompactOutput {
		detector.PrintStats()
	}
}
