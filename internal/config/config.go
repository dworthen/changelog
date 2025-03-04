package config

import (
	"log/slog"
	"os"
	"sync"

	"github.com/dworthen/changelog/internal/utils"
	"gopkg.in/yaml.v3"
)

type VersionFile struct {
	Path    string `yaml:"path"`
	Pattern string `yaml:"pattern"`
}

type OnAdd struct {
	CommitFiles bool `yaml:"commitFiles"`
}

type OnApply struct {
	CommitFiles bool     `yaml:"commitFiles"`
	TagCommit   bool     `yaml:"tagCommit"`
	TagFormat   string   `yaml:"tagFormat"`
	Commands    []string `yaml:"commands"`
}

type Config struct {
	m             sync.RWMutex  `yaml:"-"`
	Version       string        `yaml:"version"`
	ChangelogFile string        `yaml:"changelogFile"`
	Files         []VersionFile `yaml:"files"`
	OnAdd         OnAdd         `yaml:"onAdd"`
	OnApply       OnApply       `yaml:"onApply"`
}

func NewConfig() *Config {
	return &Config{
		Version:       "0.0.0",
		ChangelogFile: "CHANGELOG.md",
		Files:         []VersionFile{},
		OnAdd: OnAdd{
			CommitFiles: true,
		},
		OnApply: OnApply{
			CommitFiles: true,
			TagCommit:   true,
			TagFormat:   "v{{version}}",
		},
	}
}

func (c *Config) Save() error {
	c.m.Lock()
	defer c.m.Unlock()

	configPath := utils.GetConfigFilePath()

	fileContents, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, fileContents, 0644)
}

func (c *Config) SetVersion(version string) {
	c.m.Lock()
	defer c.m.Unlock()

	c.Version = version
}

func (c *Config) GetVersion() string {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.Version
}

func (c *Config) SetChangelogFile(changelogFile string) {
	c.m.Lock()
	defer c.m.Unlock()
	c.ChangelogFile = changelogFile
}

func (c *Config) GetChangelogFile() string {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.ChangelogFile
}

func (c *Config) AddVersionFile(path string, pattern string) {
	c.m.Lock()
	defer c.m.Unlock()
	c.Files = append(c.Files, VersionFile{
		Path:    path,
		Pattern: pattern,
	})
}

func (c *Config) GetVersionFiles() []VersionFile {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.Files
}

func (c *Config) SetOnAddCommitFiles(commitFiles bool) {
	c.m.Lock()
	defer c.m.Unlock()
	c.OnAdd.CommitFiles = commitFiles
}

func (c *Config) GetOnAddCommitFiles() bool {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.OnAdd.CommitFiles
}

func (c *Config) SetOnApplyCommands(commands []string) {
	c.m.Lock()
	defer c.m.Unlock()
	c.OnApply.Commands = commands
}

func (c *Config) GetOnApplyCommands() []string {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.OnApply.Commands
}

func (c *Config) SetOnApplyCommitFiles(commitFiles bool) {
	c.m.Lock()
	defer c.m.Unlock()
	c.OnApply.CommitFiles = commitFiles
}

func (c *Config) GetOnApplyCommitFiles() bool {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.OnApply.CommitFiles
}

func (c *Config) SetOnApplyTagCommit(tagCommit bool) {
	c.m.Lock()
	defer c.m.Unlock()
	c.OnApply.TagCommit = tagCommit
}

func (c *Config) GetOnApplyTagCommit() bool {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.OnApply.TagCommit
}

func (c *Config) SetOnApplyTagFormat(tagFormat string) {
	c.m.Lock()
	defer c.m.Unlock()
	c.OnApply.TagFormat = tagFormat
}

func (c *Config) GetOnApplyTagFormat() string {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.OnApply.TagFormat
}

var conf *Config
var configOnce sync.Once
var configLoadingErr error

func LoadConfig() (*Config, error) {
	configOnce.Do(func() {
		conf = NewConfig()
		configPath := utils.GetConfigFilePath()
		_, err := os.Stat(configPath)
		if err != nil {
			return
		}

		fileContents, err := os.ReadFile(configPath)
		if err != nil {
			configLoadingErr = err
			return
		}

		err = yaml.Unmarshal(fileContents, conf)
		if err != nil {
			configLoadingErr = err
			return
		}
		slog.Info("Loaded config", "config", conf)
	})
	if configLoadingErr != nil {
		return nil, configLoadingErr
	}
	return conf, nil
}
