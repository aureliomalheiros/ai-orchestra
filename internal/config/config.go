package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type App struct {
	Name       string `yaml:"name"`
	WorkingDir string `yaml:"working_dir"`
	LogLevel   string `yaml:"log_level"`
}

type ClaudeConfig struct {
	Binary       string        `yaml:"binary"`
	Model        string        `yaml:"model"`
	Timeout      time.Duration `yaml:"timeout"`
	OutputFormat string        `yaml:"output_format"`
	Interactive  bool          `yaml:"interactive"`
	MaxTurns     int           `yaml:"max_turns"`
}

type CodexConfig struct {
	Binary       string        `yaml:"binary"`
	Model        string        `yaml:"model"`
	Timeout      time.Duration `yaml:"timeout"`
	ApprovalMode string        `yaml:"approval_mode"`
	SandboxMode  string        `yaml:"sandbox_mode"`
}

type AgentsConfig struct {
	Claude ClaudeConfig `yaml:"claude"`
	Codex  CodexConfig  `yaml:"codex"`
}

type Stage struct {
	Name  string `yaml:"name"`
	Agent string `yaml:"agent"`
	Role  string `yaml:"role"`
}

type ApprovalConfig struct {
	BeforePlanning  bool `yaml:"before_planning"`
	BeforeExecution bool `yaml:"before_execution"`
	BeforeReview    bool `yaml:"before_review"`
	ShowDiff        bool `yaml:"show_diff"`
}

type RetryConfig struct {
	MaxAttempts int    `yaml:"max_attempts"`
	Backoff     string `yaml:"backoff"`
}

type WorkflowConfig struct {
	Stages   []Stage        `yaml:"stages"`
	Approval ApprovalConfig `yaml:"approval"`
	Retry    RetryConfig    `yaml:"retry"`
}

type SkillsConfig struct {
	Directories []string `yaml:"directories"`
}

type StateConfig struct {
	Directory   string `yaml:"directory"`
	SaveHistory bool   `yaml:"save_history"`
	SavePlans   bool   `yaml:"save_plans"`
}

type LoggingConfig struct {
	Directory string `yaml:"directory"`
	Level     string `yaml:"level"`
	Format    string `yaml:"format"`
}

type Config struct {
	App      App            `yaml:"app"`
	Agents   AgentsConfig   `yaml:"agents"`
	Workflow WorkflowConfig `yaml:"workflow"`
	Skills   SkillsConfig   `yaml:"skills"`
	State    StateConfig    `yaml:"state"`
	Logging  LoggingConfig  `yaml:"logging"`
}

type Loader struct {
	configPath string
}

func NewLoader() *Loader {
	return &Loader{}
}

func NewLoaderWithPath(path string) *Loader {
	return &Loader{configPath: path}
}

func (l *Loader) Load() (*Config, error) {
	path, err := l.findConfigFile()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	expanded := os.ExpandEnv(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	cfg.App.WorkingDir = expandPath(cfg.App.WorkingDir)
	cfg.State.Directory = expandPath(cfg.State.Directory)
	cfg.Logging.Directory = expandPath(cfg.Logging.Directory)
	for i, dir := range cfg.Skills.Directories {
		cfg.Skills.Directories[i] = expandPath(dir)
	}

	return &cfg, nil
}

func (l *Loader) findConfigFile() (string, error) {
	if l.configPath != "" {
		return l.configPath, nil
	}

	candidates := []string{
		"./configs/config.yaml",
		"./config.yaml",
	}

	home, err := os.UserHomeDir()
	if err == nil {
		candidates = append(candidates, filepath.Join(home, ".ai-orchestra", "config.yaml"))
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("config file not found in any of: %v", candidates)
}

func (c *Config) Validate() error {
	if c.Agents.Claude.Binary == "" {
		return fmt.Errorf("agents.claude.binary is required")
	}
	if c.Agents.Codex.Binary == "" {
		return fmt.Errorf("agents.codex.binary is required")
	}
	if c.Agents.Claude.Timeout <= 0 {
		return fmt.Errorf("agents.claude.timeout must be positive")
	}
	if c.Agents.Codex.Timeout <= 0 {
		return fmt.Errorf("agents.codex.timeout must be positive")
	}
	return nil
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}
