# Migrating

If your project already has a changelog file, you can adopt `changelog` without losing your existing history.

## Steps

### 1. Initialize

Run `changelog init` and set `currentVersion` to your project's current version:

```bash
changelog init
```

```
? Current version: 2.3.1
? Changelog file path: CHANGELOG.md
? Main git branch: main
Enter post-apply commands (leave blank to finish):
? Post-apply command:
✓ Initialized .changelog directory
```

### 2. Preserve Existing History

Copy the contents of your existing changelog file into `.changelog/templates/footer.md`:

```bash
cp CHANGELOG.md .changelog/templates/footer.md
```

The footer template is appended after all generated release sections. This means your historical entries will appear below any new releases, preserving the full project history.

### 3. Regenerate the Changelog

Run `changelog view` and write the output to your changelog file:

```bash
changelog view > CHANGELOG.md
```

This regenerates the changelog using the header template, any existing releases (none yet), and the footer containing your historical content. The result should look very similar to your original file, with the header from `templates/header.md` at the top.

### 4. Commit

Commit the new `.changelog/` directory and the updated changelog:

```bash
git add .changelog/ CHANGELOG.md
git commit -m "chore: adopt changelog tool"
```

## Going Forward

From this point on, use `changelog` as designed:

1. Run [`changelog add`](/commands?id=changelog-add) when making notable changes.
2. Run [`changelog apply`](/commands?id=changelog-apply) to cut releases.

New releases will appear **above** the historical footer content, maintaining a natural reverse-chronological order.

?> **Tip:** You can edit `.changelog/templates/header.md` if you want to customize the header that appears at the top of your changelog.
