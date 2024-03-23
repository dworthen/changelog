package gitmanage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CommitFiles(files []string, description string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	for _, file := range files {
		if err != nil {
			return err
		}
		_, err = worktree.Add(file)
		if err != nil {
			return err
		}
	}

	_, err = worktree.Commit(description, &git.CommitOptions{})
	return err
}

func GetFileCommit(filePath string) (string, error) {
	filePath = filepath.ToSlash(filePath)
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return "", err
	}

	commits, err := repo.Log(&git.LogOptions{
		FileName: &filePath,
	})
	if err != nil {
		return "", err
	}

	commit, err := commits.Next()
	if err != nil {
		return "", err
	}

	return commit.Hash.String(), nil
}

func tagExists(tag string, r *git.Repository) bool {
	tagFoundErr := "tag was found"
	tags, err := r.TagObjects()
	if err != nil {
		return false
	}
	res := false
	err = tags.ForEach(func(t *object.Tag) error {
		if t.Name == tag {
			res = true
			return fmt.Errorf(tagFoundErr)
		}
		return nil
	})
	if err != nil && err.Error() != tagFoundErr {
		log.Printf("iterate tags error: %s", err)
		return false
	}
	return res
}

func Tag(tag string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	if tagExists(tag, repo) {
		return fmt.Errorf("Tag %s already exists", tag)
	}

	head, err := repo.Head()
	if err != nil {
		return err
	}
	_, err = repo.CreateTag(tag, head.Hash(), &git.CreateTagOptions{
		Message: tag,
	})
	return err
}
