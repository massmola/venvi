---
description: Verify that code changes meet the project's quality standards.
---
# Verify Changes Workflow

// turbo-all

1. **Run Full Validation in Nix Environment**:
   ```bash
   nix develop --command ./scripts/validate.sh
   ```

3. **Report**:
   If the validation script fails, fix the issues and restart the workflow. If it passes, the code is ready for review.
