# Action: Archive Change

When archiving an implemented change, follow these steps:

1. Locate the folder `.specrate/changes/{change-id}/` with the specified `change-id`.
   - Identify the `change-id` based on the user's description if not provided accurately, and confirm it with the user.
   - If the folder does not exist, inform the user and suggest proposing a new change instead, then abort the current action.
   - If the change is not in the `IMPLEMENTED` state, inform the user and abort the current action.
2. Ensure `.specrate/specs/` exists; create it if missing.
3. Integrate the spec delta in `spec-delta.md` into the relevant specs.
   - If a referenced spec does not exist, create it under `.specrate/specs/{spec-id}/spec.md` using the [SPECS-SPEC.md](../../assets/templates/SPECS-SPEC.md) template and then apply the delta.
   - Prefer editing the authoritative spec(s) rather than leaving requirements only in the change delta.
4. Update the `state` file to contain the text `ARCHIVED` only.
5. Summarize the archived change and provide any next steps if applicable.
