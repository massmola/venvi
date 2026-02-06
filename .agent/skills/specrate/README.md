# Specrate

Specrate is an agent skill that manages specs and changes in a structured way. It provides a framework for proposing, planning, implementing, and archiving changes to the system's specs.

## Installation

Ensure your agent supports [agent skills](https://agentskills.io) and install the skill by git-cloning into your repo's `.github/skills` folder. Replace `.github/skills` with the appropriate path if your agent uses a different location for skills.

Note: the skill itself is installed under `.github/skills/specrate`, while the documents it manages are created/updated under `.specrate/` at your repository root.

### Bash

```sh
git clone --depth 1 https://github.com/rickygao/specrate .github/skills/specrate
rm -rf .github/skills/specrate/.git
```

### PowerShell

```pwsh
& git clone --depth 1 https://github.com/rickygao/specrate .github/skills/specrate
Remove-Item -Path .github/skills/specrate/.git -Recurse -Force
```

## Usage

Specrate manages spec and change artifacts under `.specrate/` at the repository root (see [references/PRINCIPLE.md](./references/PRINCIPLE.md) for conventions).

1. Show current status: "Show the current status of specs and changes"
2. Propose a new change: "Propose a new change to {description of the change}"
3. Amend an existing change: "Amend the change {change-id}: {description of the amendment}"
4. Plan a proposed change: "Plan the proposed {change-id}"
5. Implement a planned change: "Implement the planned {change-id}"
6. Archive an implemented change: "Archive the implemented {change-id}"
7. Fix a spec from codebase: "Fix the spec for {spec-id, or description of the spec scope} from the codebase"

Refer to the [SKILL.md](./SKILL.md) file for how the skill works in detail.

## Acknowledgements

This skill was heavily inspired by the [OpenSpec](https://openspec.dev).
