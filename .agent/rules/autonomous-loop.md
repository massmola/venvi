---
trigger: always_on
---

# Autonomous Loop Rule

**Mandatory Rule**: You MUST use the **Autonomous Loop** pattern for any task that involves multiple steps, ambiguity, or significant code changes. you have a dedicated workflow

## When to use
- New feature implementation.
- Refactoring.
- Debugging complex issues.

## Error Handling
- If an error occurs, do NOT just fix it. **Reflect** on why it happened and generate a "Lesson Learned" to update your system prompts or rules.
