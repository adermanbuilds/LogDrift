package main

import (
	"flag"
	"fmt"
	"os"
)

// CLIOptions holds command-line flags that adjust runtime behavior.
type CLIOptions struct {
	ConfigFile     string
	ShowAll        bool
	ShowDebug      bool
	ShowInfo       bool
	NoColor        bool
	Compact        bool
	Component      string
	MinLevel       string
	OnlyAnomalies  bool
	GenerateConfig bool
	Version        bool
}

// printInfo displays a compact, nicely formatted info block to stderr.
func printInfo() {
	fmt.Fprintln(os.Stderr, "\033[1;36mLogDrift v0.1.0 — Fast log anomaly detection\033[0m")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "    logdrift [options] < logfile.log")
	fmt.Fprintln(os.Stderr, "    tail -f /var/log/app.log | logdrift [options]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Options:")
	fmt.Fprintln(os.Stderr, "    --config string         Path to config file (JSON)")
	fmt.Fprintln(os.Stderr, "    --all                   Show all logs (including DEBUG/INFO)")
	fmt.Fprintln(os.Stderr, "    --debug                 Show DEBUG logs")
	fmt.Fprintln(os.Stderr, "    --info                  Show INFO logs")
	fmt.Fprintln(os.Stderr, "    --no-color              Disable colored output")
	fmt.Fprintln(os.Stderr, "    --compact               Compact output (no anomaly highlights)")
	fmt.Fprintln(os.Stderr, "    --component string      Filter by component (e.g., 'parser')")
	fmt.Fprintln(os.Stderr, "    --min-level string      Minimum log level (default: WARN)")
	fmt.Fprintln(os.Stderr, "    --anomalies-only        Show only anomalies")
	fmt.Fprintln(os.Stderr, "    --generate-config       Generate default config file")
	fmt.Fprintln(os.Stderr, "    --version               Show version")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "For more: https://github.com/adermanbuilds/logdrift")
	fmt.Fprintln(os.Stderr, "")
}

// ParseCLI parses command-line arguments and returns a populated CLIOptions.
func ParseCLI() *CLIOptions {
	opts := &CLIOptions{}

	flag.StringVar(&opts.ConfigFile, "config", "", "Path to config file (JSON)")
	flag.BoolVar(&opts.ShowAll, "all", false, "Show all logs (including DEBUG/INFO)")
	flag.BoolVar(&opts.ShowDebug, "debug", false, "Show DEBUG logs")
	flag.BoolVar(&opts.ShowInfo, "info", false, "Show INFO logs")
	flag.BoolVar(&opts.NoColor, "no-color", false, "Disable colored output")
	flag.BoolVar(&opts.Compact, "compact", false, "Compact output (no anomaly highlights)")
	flag.StringVar(&opts.Component, "component", "", "Filter by component (e.g., 'parser')")
	flag.StringVar(&opts.MinLevel, "min-level", "WARN", "Minimum log level (DEBUG/INFO/WARN/ERROR/FATAL)")
	flag.BoolVar(&opts.OnlyAnomalies, "anomalies-only", false, "Show only anomalies")
	flag.BoolVar(&opts.GenerateConfig, "generate-config", false, "Generate default config file")
	flag.BoolVar(&opts.Version, "version", false, "Show version")

	// Replace default usage with the neat, aligned block.
	flag.Usage = func() {
		printInfo()
		// flag.PrintDefaults() is intentionally omitted to keep the aligned format above.
	}

	flag.Parse()
	return opts
}

// ApplyToConfig mutates the provided Config to reflect CLI options.
func (o *CLIOptions) ApplyToConfig(cfg *Config) {
	if o.ShowAll {
		cfg.ShowDebug = true
		cfg.ShowInfo = true
		cfg.MinLevel = "DEBUG"
	}

	if o.ShowDebug {
		cfg.ShowDebug = true
		if cfg.MinLevel == "INFO" || cfg.MinLevel == "WARN" {
			cfg.MinLevel = "DEBUG"
		}
	}

	if o.ShowInfo {
		cfg.ShowInfo = true
		if cfg.MinLevel == "WARN" {
			cfg.MinLevel = "INFO"
		}
	}

	if o.NoColor {
		cfg.ColorOutput = false
	}

	if o.Compact {
		cfg.CompactOutput = true
	}

	if o.Component != "" {
		cfg.IncludeComponents = []string{o.Component}
	}

	if o.MinLevel != "" {
		cfg.MinLevel = o.MinLevel
	}
}

// HandleSpecialCommands processes flags that should cause immediate output and exit.
func (o *CLIOptions) HandleSpecialCommands() bool {
	if o.Version {
		printInfo()
		fmt.Println("Built with Go, optimized for speed")
		return true
	}

	if o.GenerateConfig {
		cfg := DefaultConfig()
		if err := cfg.SaveConfig("/dev/stdout"); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating config: %v\n", err)
			os.Exit(1)
		}
		return true
	}

	return false
}
