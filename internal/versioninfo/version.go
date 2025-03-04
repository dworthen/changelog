package versioninfo

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/dworthen/updater"
)

var update *updater.Updater

func getUpdater() (*updater.Updater, error) {
	if update == nil {
		currentVersion, err := GetVersion()
		if err != nil {
			return nil, err
		}
		update = updater.New(&updater.UpdaterConfig{
			CurrentVersion: currentVersion,
			BaseUrl:        "https://github.com/dworthen/changelog/releases/latest/download",
			UpdaterConfig:  "updater.config.json",
		})
	}

	return update, nil
}

//go:embed updater.config.json
var versionFileContents string

type VersionInfo struct {
	Version string `json:"version"`
}

func GetVersion() (string, error) {
	var versionInfo VersionInfo

	err := json.Unmarshal([]byte(versionFileContents), &versionInfo)
	if err != nil {
		return "", err
	}

	return versionInfo.Version, nil
}

func CheckForUpdate() (bool, string, error) {
	update, err := getUpdater()
	if err != nil {
		return false, "", err
	}
	return update.CheckForAvailableUpdate()
}

func PrintAvailableUpdate() error {
	isUpdate, newVersion, err := CheckForUpdate()
	if err != nil {
		return err
	}

	if isUpdate {
		fmt.Printf("A new version of changelog is available, %s. Run the `changelog update` to upgrade to the latest version.\n", newVersion)
	}

	return nil
}

func Update() error {
	update, err := getUpdater()
	if err != nil {
		return err
	}
	return update.Update()
}
