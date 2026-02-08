# Autonomous Loop Rule

**Mandatory Rule**: You MUST use the **Autonomous Loop** pattern for any task that involves multiple steps, ambiguity, or significant code changes.

## The Loop
1.  **Plan**: Break the task into micro-tasks using the Orchestrator prompt.
2.  **Execute**: Perform the task, logging each step.
3.  **Verify**: **MANDATORY**. Run tests and `verify_changes` after every significant step.
4.  **Reflect**: Use the Critic prompt to analyze the session, specifically looking for ways to improve **Rules**, **Skills**, and **Workflows**.

## When to use
- New feature implementation.
- Refactoring.
- Debugging complex issues.

## Error Handling
- If an error occurs, do NOT just fix it. **Reflect** on why it happened and generate a "Lesson Learned" to update your system prompts or rules.
