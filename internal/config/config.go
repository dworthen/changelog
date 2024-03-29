package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type BumpInfo struct {
	FileName string `json:"file"`
	JsonPath string `json:"path"`
}

type OnAdd struct {
	ComitFiles bool `json:"commitFiles"`
}

type OnApply struct {
	CommitFiles bool   `json:"commitFiles"`
	TagCommit   bool   `json:"tagCommit"`
	TagFormat   string `json:"tagFormat"`
}

type Config struct {
	Version   string     `json:"version"`
	BumpFiles []BumpInfo `json:"bumpFiles"`
	OnAdd     OnAdd      `json:"onAdd"`
	OnApply   OnApply    `json:"onApply"`
}

func newConfig() *Config {
	return &Config{
		Version:   "",
		BumpFiles: []BumpInfo{},
		OnAdd: OnAdd{
			ComitFiles: true,
		},
		OnApply: OnApply{
			CommitFiles: true,
			TagCommit:   true,
			TagFormat:   "v{{version}}",
		},
	}
}

func (c *Config) Save() error {
	dir, err := filepath.Abs(".changelog")
	if err != nil {
		return err
	}

	stats, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(".changelog directory does not exist. Run changelog init.")
		} else {
			return err
		}
	}
	if !stats.IsDir() {
		return fmt.Errorf(".changelog is not a directory")
	}

	filename := filepath.Join(dir, "config.json")

	configContents, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, configContents, 0644)
}

var conf *Config
var once sync.Once

func GetConfig() (*Config, error) {
	var errMsg string = ""
	once.Do(func() {
		conf = newConfig()
		fileLocation, err := filepath.Abs(".changelog/config.json")
		if err != nil {
			errMsg = err.Error()
			return
		}

		_, err = os.Stat(fileLocation)
		if err != nil {
			return
		}

		fileContents, err := os.ReadFile(fileLocation)
		if err != nil {
			errMsg = err.Error()
			return
		}

		err = json.Unmarshal(fileContents, &conf)
		if err != nil {
			errMsg = err.Error()
			return
		}
	})

	if errMsg != "" {
		return nil, fmt.Errorf("Failed to load config. %s", errMsg)
	}

	return conf, nil
}

func (bumpInfo *BumpInfo) Bump(newVersion string) error {
	filePath, err := filepath.Abs(bumpInfo.FileName)
	if err != nil {
		return err
	}

	_, err = os.Stat(filePath)
	if err != nil {
		return err
	}

	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(fileContents, &data)
	if err != nil {
		return err
	}
	var currentInterface interface{} = data
	pathWalk := strings.Split(bumpInfo.JsonPath, ".")

	for i := 0; i < len(pathWalk)-1; i++ {
		key := pathWalk[i]
		switch typedInterface := currentInterface.(type) {
		case map[string]interface{}:
			currentInterface = typedInterface[key]
		case []interface{}:
			intKey, err := strconv.Atoi(key)
			if err != nil {
				return err
			}
			if intKey < 0 || intKey >= len(typedInterface) {
				return fmt.Errorf("Unable to parse json path %s for %s. Array index falls outside of array", bumpInfo.JsonPath, bumpInfo.FileName)
			}
			currentInterface = typedInterface[intKey]
		default:
			return fmt.Errorf("Unable to parse json path %s for %s", bumpInfo.JsonPath, bumpInfo.FileName)
		}
	}

	key := pathWalk[len(pathWalk)-1]
	switch currentInterface := currentInterface.(type) {
	case map[string]interface{}:
		currentInterface[key] = newVersion
	case []interface{}:
		intKey, err := strconv.Atoi(key)
		if err != nil {
			return err
		}
		if intKey < 0 || intKey >= len(currentInterface) {
			return fmt.Errorf("Unable to parse json path %s for %s. Array index falls outside of array", bumpInfo.JsonPath, bumpInfo.FileName)
		}
		currentInterface[intKey] = newVersion
	default:
		return fmt.Errorf("Unable to parse json path %s for %s", bumpInfo.JsonPath, bumpInfo.FileName)
	}

	newFileContents, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, newFileContents, 0644)
}
