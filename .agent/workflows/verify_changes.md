---
description: Verify that code changes meet the project's quality standards.
---
# Verify Changes Workflow

// turbo-all

1. **Enter Nix Development Environment**:
   ```bash
   nix develop
   ```

2. **Run Full Validation**:
   ```bash
   ./scripts/validate.sh
   ```

3. **Report**:
   If the validation script fails, fix the issues and restart the workflow. If it passes, the code is ready for review.
