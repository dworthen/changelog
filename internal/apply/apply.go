package apply

import (
	"fmt"
	"strings"

	"github.com/dworthen/changelog/internal/config"
	"github.com/dworthen/changelog/internal/gitmanage"
	"github.com/mailgun/raymond/v2"
)

func Apply() error {
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}
	changelog, err := NewChangelog()
	if err != nil {
		return err
	}
	err = changelog.BumpFiles()
	if err != nil {
		return err
	}
	err = changelog.Save()
	if err != nil {
		return err
	}
	if conf.OnApply.CommitFiles {
		err = gitmanage.CommitFiles([]string{"."}, "chore: apply changelog")
		if err != nil {
			return err
		}
		tagFormat := strings.TrimSpace(conf.OnApply.TagFormat)
		if conf.OnApply.TagCommit && tagFormat != "" {
			tag, err := raymond.Render(tagFormat, changelog)
			if err != nil {
				return err
			}
			err = gitmanage.Tag(strings.TrimSpace(tag))
			if err != nil {
				return err
			}
		}
	}

	fmt.Printf("Changelog updates applied! Updated from version %s to %s\n", changelog.OldVersion, changelog.Version)

	return nil
}
