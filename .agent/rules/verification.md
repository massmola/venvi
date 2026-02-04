---
trigger: always_on
description: Mandatory verification workflow after any code modification.
---

# Mandatory Verification Rule

To ensure the integrity of the project, the following rule is strictly enforced:

## Verification After Modification
After making code changes, including refactoring, new features, or bug fixes, you **MUST** run the project's verification workflow before considering the task complete and notifying the user.

### Required Steps:
1. **Run the Verification Workflow**: Execute the `/verify_changes` slash command (or the steps defined in `.agent/workflows/verify_changes.md`).
2. **Handle Failures**: If any verification step fails, you must fix the issues and re-run the entire verification workflow until it passes.
3. **Include Proof**: When calling `notify_user` or creating a `walkthrough.md`, explicitly state that the verification workflow was run and is passing.

> [!IMPORTANT]
> A task is NOT complete until all checks (Formatting, Linting, Type Checking, Tests, and Documentation Build) pass successfully in the project's development environment.