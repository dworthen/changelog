package main

import (
	"io"
	"os"
	"path/filepath"
)

func main() {

	mapping := map[string][]string{
		"dist/changelog_windows_arm64_v8.0/changelog.exe": {"js/win32-arm64/bin/changelog.exe", "python/win32-arm64/src/changesets_win32_arm64/bin/changelog.exe"},
		"dist/changelog_windows_amd64_v1/changelog.exe":   {"js/win32-x64/bin/changelog.exe", "python/win32-x64/src/changesets_win32_x64/bin/changelog.exe"},
		"dist/changelog_darwin_arm64_v8.0/changelog":      {"js/darwin-arm64/bin/changelog", "python/darwin-arm64/src/changesets_darwin_arm64/bin/changelog"},
		"dist/changelog_darwin_amd64_v1/changelog":        {"js/darwin-x64/bin/changelog", "python/darwin-x64/src/changesets_darwin_x64/bin/changelog"},
		"dist/changelog_linux_arm64_v8.0/changelog":       {"js/linux-arm64/bin/changelog", "python/linux-arm64/src/changesets_linux_arm64/bin/changelog"},
		"dist/changelog_linux_amd64_v1/changelog":         {"js/linux-x64/bin/changelog", "python/linux-x64/src/changesets_linux_x64/bin/changelog"},
		"README.md": {"js/changelog/README.md", "js/darwin-arm64/README.md", "js/darwin-x64/README.md", "js/linux-arm64/README.md", "js/linux-x64/README.md", "js/win32-arm64/README.md", "js/win32-x64/README.md", "python/changelog/README.md", "python/darwin-arm64/README.md", "python/darwin-x64/README.md", "python/linux-arm64/README.md", "python/linux-x64/README.md", "python/win32-arm64/README.md", "python/win32-x64/README.md"},
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for src, destinations := range mapping {
		for _, dest := range destinations {
			srcPath := filepath.Join(cwd, src)
			destPath := filepath.Join(cwd, dest)
			destDir := filepath.Dir(destPath)

			err = os.MkdirAll(destDir, 0755)
			if err != nil {
				panic(err)
			}

			sourceFile, err := os.Open(srcPath)
			if err != nil {
				panic(err)
			}

			destFile, err := os.Create(destPath)
			if err != nil {
				panic(err)
			}
			defer destFile.Close()
			_, err = io.Copy(destFile, sourceFile)
			if err != nil {
				panic(err)
			}
		}
	}

}
