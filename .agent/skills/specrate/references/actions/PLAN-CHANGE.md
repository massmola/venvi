# Action: Plan Change

When planning a proposed change, follow these steps:

1. Locate the folder `.specrate/changes/{change-id}/` with the specified `change-id`.
   - Identify the `change-id` based on the user's description if not provided accurately, and confirm it with the user.
   - If the folder does not exist, inform the user and suggest proposing a new change instead, then abort the current action.
   - If the change is not in the `PROPOSED` state, inform the user and abort the current action.
2. Create the following required files inside the change folder:
   - `tasks.md`: using the [CHANGES-TASKS.md](../../assets/templates/CHANGES-TASKS.md) template.
   - Optionally, `design.md`: using the [CHANGES-DESIGN.md](../../assets/templates/CHANGES-DESIGN.md) template.
3. Fill out the `tasks.md` file (and `design.md` if present) with the relevant information about the change.
   - If `design.md` is present, fill it out first since it may inform the tasks.
   - Create actionable tasks in `tasks.md` using task list grammar (e.g., `- [ ]` for pending tasks).
4. Update the `state` file to contain the text `PLANNED` only.
5. Summarize the planned change and provide any next steps if applicable.
