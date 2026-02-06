# Action: Implement Change

When implementing a planned change, follow these steps:

1. Locate the folder `.specrate/changes/{change-id}/` with the specified `change-id`.
   - Identify the `change-id` based on the user's description if not provided accurately, and confirm it with the user.
   - If the folder does not exist, inform the user and suggest proposing a new change instead, then abort the current action.
   - If the change is not in the `PLANNED` state, inform the user and abort the current action.
2. Read the `tasks.md` file to understand the tasks to be completed.
3. For each task listed in `tasks.md`, perform the following:
   - Do the task as per the user's instructions.
   - If more information or clarification is needed, ask the user, and then amend relevant files if necessary.
   - Mark the task as completed in `tasks.md` using task list grammar (e.g., `- [x]` for completed tasks).
   - Keep the workspace consistent after each task (avoid leaving it half-broken unless the user explicitly accepts that).
   - Repeat this step until all tasks are completed.
4. Update the `state` file to contain the text `IMPLEMENTED` only.
5. Summarize the implemented change and provide any next steps if applicable.
