---
description: Enforce immediate draft PR creation after branching.
glob: "**/*"
---

# Pull Request Standard

To ensure early visibility and collaboration, the following rule is strictly enforced:

## Immediate Draft PR Creation
When you create a new branch that is intended to be merged later, you **MUST** immediately create a pull request marked as a draft using the GitHub CLI.

### command
```bash
gh pr create --draft --fill
```

### Why?
- **Visibility**: Allows others to see what is being worked on.
- **Early Feedback**: Enables early discussions on the approach.
- **CI/CD**: Triggers CI/CD pipelines early (if configured) to catch issues sooner.
