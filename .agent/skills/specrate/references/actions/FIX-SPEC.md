# Action: Fix Spec

When fixing a spec according to the codebase, follow these steps:

1. Check if the codebase contains the implementation necessary to fix the spec.
   - If not, inform the user that fixing is not possible due to insufficient information, and abort the current action.
2. Ensure `.specrate/specs/` exists; create it if missing.
3. Locate the folder `.specrate/specs/{spec-id}/` with the specified `spec-id` that describes the spec.
   - Propose a `spec-id` based on the user's description if not provided, and confirm it with the user.
   - If the folder does not exist, confirm with the user to create a new spec folder.
4. Check if the following required file exists inside the spec folder: `spec.md`.
   - If the file exists, update the file according to the current implementation in the codebase.
     Make sure to reflect any changes that are not aligned with the existing spec.
   - Otherwise, create the file using the [SPECS-SPEC.md](../../assets/templates/SPECS-SPEC.md) template.
     Gather information from existing implementations, documentation, and other relevant sources.
     Fill out the `spec.md` file with accurate and comprehensive details about the spec.
5. Summarize the fixed spec and provide any next steps if applicable.
