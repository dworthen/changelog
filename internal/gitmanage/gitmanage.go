package gitmanage

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/dworthen/changelog/internal/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type GitNotInitializedError struct {
}

func (e *GitNotInitializedError) Error() string {
	return "Git is not initialized."
}

func NewGitNotInitializedError() error {
	return &GitNotInitializedError{}
}

func IsGitNotInitializedError(err error) bool {
	return errors.Is(err, &GitNotInitializedError{})
}

var gitRepoInstance *git.Repository
var directory string
var getGitRepoOnce sync.Once
var getGitRepoErr error

func GetGitRepo() (*git.Repository, string, error) {
	getGitRepoOnce.Do(func() {
		cwd := utils.GetCWD()

		for {
			repo, err := git.PlainOpen(cwd)
			if err == nil {
				gitRepoInstance = repo
				directory = cwd
				return
			}
			if err.Error() == "repository does not exist" {
				if dir := filepath.Dir(cwd); dir != cwd {
					cwd = dir
					continue
				}
				getGitRepoErr = NewGitNotInitializedError()
				return
			}
			getGitRepoErr = fmt.Errorf("Error opening git repository: %w", err)
			return
		}

	})
	if getGitRepoErr != nil {
		return nil, "", getGitRepoErr
	}
	return gitRepoInstance, directory, nil
}

func CommitFiles(files []string, description string) error {
	cwd := utils.GetCWD()
	repo, repoDirectory, err := GetGitRepo()
	if err != nil {
		return err
	}

	relDir, err := filepath.Rel(repoDirectory, cwd)
	if err != nil {
		return fmt.Errorf("Error committing files. Failed to get relative file location for %s from %s. Error: %w", repoDirectory, cwd, err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("Error committing files. Failed to get worktree: %w", err)
	}

	for _, file := range files {
		fileLocation := utils.JoinPaths(relDir, file)
		_, err = worktree.Add(fileLocation)
		if err != nil {
			return fmt.Errorf("Error committing files. Failed to add file to worktree. %s. Error: %w", fileLocation, err)
		}
	}

	_, err = worktree.Commit(description, &git.CommitOptions{})
	if err != nil {
		return fmt.Errorf("Error committing files. Failed to commit files: %w", err)
	}
	return nil
}

func GetCommitHashForFile(filePath string) (string, error) {
	cwd := utils.GetCWD()
	repo, repoDirectory, err := GetGitRepo()
	if err != nil {
		return "", err
	}

	relDir, err := filepath.Rel(repoDirectory, cwd)
	if err != nil {
		return "", fmt.Errorf("Error getting commit hash for file. Failed to get relative file location for %s from %s. Error: %w", repoDirectory, cwd, err)
	}

	filePath = utils.JoinPaths(relDir, filePath)
	filePath = filepath.ToSlash(filePath)

	commits, err := repo.Log(&git.LogOptions{
		FileName: &filePath,
	})
	if err != nil {
		return "", fmt.Errorf("Error getting commit hash for file. Failed to get commits for file %s: %w", filePath, err)
	}

	commit, err := commits.Next()
	if err != nil {
		return "", fmt.Errorf("Error getting commit hash for file. Failed to get next commit: %w", err)
	}

	return commit.Hash.String(), nil
}

func tagExists(tag string, r *git.Repository) (bool, error) {
	tags, err := r.TagObjects()
	if err != nil {
		return false, fmt.Errorf("Error getting tags: %w", err)
	}
	res := false
	err = tags.ForEach(func(t *object.Tag) error {
		if t.Name == tag {
			res = true
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return res, nil
}

func Tag(tag string) error {
	repo, _, err := GetGitRepo()
	if err != nil {
		return err
	}

	tagAlreadyExists, err := tagExists(tag, repo)
	if err != nil {
		return err
	}

	if tagAlreadyExists {
		return fmt.Errorf("Tag %s already exists", tag)
	}

	head, err := repo.Head()

	if err != nil {
		return fmt.Errorf("Error tagging commit. Failed to get HEAD: %w", err)
	}
	_, err = repo.CreateTag(tag, head.Hash(), &git.CreateTagOptions{
		Message: tag,
	})
	if err != nil {
		return fmt.Errorf("Error tagging commit. Failed to create tag %s. Error: %w", tag, err)
	}
	return nil
}

func getMainBranchCommit() (*object.Commit, error) {
	repo, _, err := GetGitRepo()
	if err != nil {
		return nil, err
	}

	branches, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("Error getting main branch commit. Failed to get branches: %w", err)
	}

	var mainReference *plumbing.Reference
	branches.ForEach(func(branch *plumbing.Reference) error {
		switch branch.Name() {
		case plumbing.Main, plumbing.Master:
			mainReference = branch
		}
		return nil
	})

	if mainReference == nil {
		head, err := repo.Head()
		if err != nil {
			return nil, fmt.Errorf("Error getting main branch commit. Failed to get HEAD: %w", err)
		}
		mainReference = head
	}

	commitObject, err := repo.CommitObject(mainReference.Hash())
	if err != nil {
		return nil, fmt.Errorf("Error getting main branch commit. Failed to get commit object for hash: %s. Error: %w", mainReference.Hash(), err)
	}
	return commitObject, nil
}

func LastCommitContainsChangelogEntry() (bool, error) {
	repo, repoDir, err := GetGitRepo()
	if err != nil {
		return false, err
	}

	relChangelogDir, err := filepath.Rel(repoDir, utils.GetChangelogDirPath())
	if err != nil {
		return false, fmt.Errorf("Error getting last commit changelog entry. Failed to get relative path for %s from %s. Error: %w", utils.GetChangelogDirPath(), repoDir, err)
	}

	pattern := filepath.ToSlash(filepath.Join(relChangelogDir, "*.md"))

	head, err := repo.Head()
	if err != nil {
		return false, fmt.Errorf("Error getting last commit changelog entry. Failed to get HEAD: %w", err)
	}

	headCommit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return false, fmt.Errorf("Error getting last commit changelog entry. Failed to get commit object for hash: %s. Error: %w", head.Hash(), err)
	}

	parentCommit, err := getMainBranchCommit()
	if err != nil {
		return false, fmt.Errorf("Error getting last commit changelog entry. Failed to get main branch commit: %w", err)
	}

	if headCommit.Hash.String() == parentCommit.Hash.String() {
		parentHash := parentCommit.ParentHashes[0]
		newParentCommit, err := repo.CommitObject(parentHash)
		if err != nil {
			return false, fmt.Errorf("Error getting last commit changelog entry. Failed to get parent commit object for hash: %s. Error: %w", parentHash, err)
		}
		parentCommit = newParentCommit
	}

	patch, err := parentCommit.Patch(headCommit)
	if err != nil {
		return false, err
	}

	matches := false
	for _, filePatch := range patch.FilePatches() {
		_, newFile := filePatch.Files()
		if newFile == nil {
			continue
		}
		match, err := filepath.Match(pattern, newFile.Path())
		if err != nil {
			return false, fmt.Errorf("Error matching file %s with pattern %s: %w", newFile.Path(), pattern, err)
		}
		if match {
			matches = true
			break
		}
	}
	return matches, nil

}
