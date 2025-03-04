package add

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/dworthen/changelog/internal/config"
	"github.com/dworthen/changelog/internal/gitmanage"
	"github.com/dworthen/changelog/internal/utils"
	"github.com/mailgun/raymond/v2"
)

type changelogEntry struct {
	m           sync.RWMutex
	Change      string
	Description string
	timestamp   time.Time
}

func (c *changelogEntry) SetChange(change string) {
	c.m.Lock()
	defer c.m.Unlock()
	c.Change = change
}

func (c *changelogEntry) GetChange() string {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.Change
}

func (c *changelogEntry) SetDescription(description string) {
	c.m.Lock()
	defer c.m.Unlock()
	c.Description = description
}

func (c *changelogEntry) GetDescription() string {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.Description
}

//go:embed changelogEntryTemplate.hbs
var changelogEntryTemplate string

func (c *changelogEntry) Save() error {
	c.m.Lock()
	defer c.m.Unlock()
	filename := fmt.Sprintf("%d.md", c.timestamp.UTC().Unix())
	fileLocation := utils.JoinPaths(utils.GetChangelogDirPath(), filename)

	dir := filepath.Dir(fileLocation)

	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s directory does not exist in working directory. Please run `changelog init` first.", dir)
		} else {
			return err
		}
	}

	fileContents, err := raymond.Render(changelogEntryTemplate, c)
	if err != nil {
		return fmt.Errorf("Error saving changelog entry. Failed to render template: %w", err)
	}

	err = os.WriteFile(fileLocation, []byte(fileContents), 0644)
	if err != nil {
		return fmt.Errorf("Error saving changelog entry. Failed to write file: %w", err)
	}

	configData, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("Error saving changelog entry. Failed to load config: %w", err)
	}

	commitFiles := configData.GetOnAddCommitFiles()
	if commitFiles {
		relFileLocation, err := filepath.Rel(utils.GetCWD(), fileLocation)
		if err != nil {
			return fmt.Errorf("Error saving changelog entry. Failed to get relative file location for changelog entry: %w", err)
		}
		err = gitmanage.CommitFiles([]string{relFileLocation}, c.Description)
		if err != nil {
			if gitmanage.IsGitNotInitializedError(err) {
				return fmt.Errorf("Error saving changelog entry. This project is configured to commit files after `changelog add` but git is not initialized. Please run `git init` first.")
			}
			return fmt.Errorf("Error saving changelog entry. Failed to commit files: %w", err)
		}
	}

	return nil
}

func NewChangelogEntry() *changelogEntry {
	return &changelogEntry{
		timestamp: time.Now(),
	}
}

var changelogEntryOnce sync.Once
var changelogEntryInstance *changelogEntry

func GetChangeLogEntry() *changelogEntry {
	changelogEntryOnce.Do(func() {
		changelogEntryInstance = NewChangelogEntry()
	})
	return changelogEntryInstance
}
