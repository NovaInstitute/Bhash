# Automation issue import

This directory contains an import file that turns the automation backlog from
[`docs/tooling/automation-requirements.md`](../automation-requirements.md) into
a set of GitHub issues. Each entry matches the AUT-00X identifier defined in
the backlog.

## Importing with GitHub CLI

Use the GitHub CLI `issue import` command to create the issues in the project
repository. The CLI accepts newline-delimited JSON (`.jsonl`) where each line is
an issue definition.

```bash
# Authenticate gh with a token that can create issues
# gh auth login --scopes "repo"

# Create the automation issues inside the repository
# Replace <owner>/<repo> with the GitHub repository name.
gh issue import -R <owner>/<repo> -F docs/tooling/automation-issues/aut-issues.jsonl
```

The issues are created with the `automation` label and include the implementation
steps and definitions of done lifted directly from the automation requirements
document. Edit the generated issues as needed to add owners or milestones.

## Manual usage

If you prefer to create issues manually, open the `.jsonl` file and copy the
`title` and `body` content for the relevant AUT item into a new GitHub issue and
apply the `automation` label.
