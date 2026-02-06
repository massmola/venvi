# Principle

## Key Conventions

- Favor straightforward, minimal implementations first and add complexity only when it is requested or clearly required.
- Keep changes tightly scoped to the requested outcome. Avoid scope creep.
- Specs describe "what" the system should do, while changes describe how the system should evolve.
- During the lifecycle of a change (`PROPOSED` → `PLANNED` → `IMPLEMENTED`), do **NOT** update specs. All proposed requirements and work-in-progress intent must live under the change folder (especially `spec-delta.md`).
- Specs are updated only when archiving a change (`ARCHIVED`), by integrating the change's `spec-delta.md` into the relevant spec(s).

## Folder Structure

All specrate-managed artifacts are organized in the `.specrate/` folder at the project root. The key folders under it are `specs/` and `changes/`.
Specs are the current state of the system, while changes are proposed modifications to it.

### Specs

The folder `.specrate/specs/` contains the specs that define various aspects of the system.
Specs are the authoritative source of truth for the system's behavior.
They are about "what" the system should do, not "how" it should be implemented.

Each spec is stored in its own subfolder named after its unique spec id (`.specrate/specs/{spec-id}/`).
Spec ids are kebab-case, noun phrases, unique identifiers, e.g., `access-control`, `payment-processing`, or `project-risk-management`.
The corresponding spec names are of the same wording but in Title Case, e.g., `Access Control`, `Payment Processing`, or `Project Risk Management`.

- `.specrate/specs/{spec-id}/spec.md`: (required, follows [SPECS-SPEC.md](../assets/templates/SPECS-SPEC.md))
  The main spec document for the spec `{spec-id}`.
  This document **MUST** outline the functionality-central, user-oriented, implementation-agnostic requirements, behaviors, and constraints.
  - functionality-central: focus on what the system should do from a functional perspective.
  - user-oriented: describe the system's behavior in terms of user needs and experiences.
  - implementation-agnostic: avoid specifying how the system should be implemented technically.
- `.specrate/specs/{spec-id}/*`: (optional)
  Any additional files related to the spec `{spec-id}`, such as diagrams or supporting documents.

### Changes

The folder `.specrate/changes/` contains the changes that describe proposed modifications to the system.
Changes are about "how" the system should evolve to meet new requirements or fix issues.

Each change is stored in its own subfolder named after its unique change id (`.specrate/changes/{change-id}/`).
Change ids are kebab-case, verb-led, unique identifiers, e.g., `add-multi-factor-auth`, `improve-payment-latency`, or `add-project-dashboard`. Prefer verb prefixes like `add-`, `remove-`, `update-`, `improve-`, *etc.*, to indicate the action being proposed.
The corresponding change names are of the same wording but in Title Case, e.g., `Add Multi-Factor Auth`, `Improve Payment Latency`, or `Add Project Dashboard`.

- `.specrate/changes/{change-id}/state`: (required)
  The state document for the change `{change-id}`, containing its current state only.
  Possible states are literally `PROPOSED`, `PLANNED`, `IMPLEMENTED`, and `ARCHIVED`.
  See policies in the **Change lifecycle** section of this document.
- `.specrate/changes/{change-id}/proposal.md`: (required, follows [CHANGES-PROPOSAL.md](../assets/templates/CHANGES-PROPOSAL.md))
  The main proposal document for the change `{change-id}`, outlining the motivation, goals, and high-level approach for the change `{change-id}`.
- `.specrate/changes/{change-id}/spec-delta.md`: (required, follows [CHANGES-SPEC-DELTA.md](../assets/templates/CHANGES-SPEC-DELTA.md))
  The spec delta for the change `{change-id}`.
  This document describes the specific modifications to be made to existing specs or new specs to be created as part of implementing the change `{change-id}`.
- `.specrate/changes/{change-id}/design.md`: (optional, follows [CHANGES-DESIGN.md](../assets/templates/CHANGES-DESIGN.md))
  The design document for the change `{change-id}`.
  This document provides a detailed design and rationale for the change `{change-id}`.
- `.specrate/changes/{change-id}/tasks.md`: (required after `PLANNED`, follows [CHANGES-TASKS.md](../assets/templates/CHANGES-TASKS.md))
  The task list for the change `{change-id}`.
  Tasks are actionable items to be completed as part of implementing the change `{change-id}`.
  Use task list grammar (e.g., `- [ ]` for pending tasks and `- [x]` for completed tasks) to track progress.
- `.specrate/changes/{change-id}/*`: (optional)
  Any additional files related to the change `{change-id}`, such as diagrams or supporting documents.

## Change Lifecycle

Each change goes through a defined lifecycle with specific states. The state must be included in the `state` file within each change folder. The possible states are: `PROPOSED`, `PLANNED`, `IMPLEMENTED`, and `ARCHIVED`. They happen in this order, and a change cannot skip any state.

Policy: treat specs as stable and do not edit them for an in-flight change. The only point where a change is merged into specs is the `ARCHIVED` step.

To manage the state of a change, use the following PowerShell commands or equivalent operations:

```pwsh
# To check the current state of a change, read the content of its state file:
Get-Content .specrate/changes/{change-id}/state

# for example:
Get-Content .specrate/changes/add-multi-factor-auth/state
```

```pwsh
# To update the state of a change, overwrite the content of its state file with the new state:
Set-Content .specrate/changes/{change-id}/state {NEW_STATE}

# for example:
Set-Content .specrate/changes/improve-payment-latency/state IMPLEMENTED
```

### Proposed

In this state, the state file contains the text `PROPOSED` only.

The state indicates that the change has been proposed and is ready for planning.
The change folder must contain at least the `state`, `proposal.md` and `spec-delta.md` files in this state.

By reading the `proposal.md` and `spec-delta.md` files, stakeholders can understand the motivation and the specific modifications proposed for the system. Then, designated engineers can plan the implementation of the change and create the `tasks.md` and `design.md` files accordingly to move the change to the next state of `PLANNED`.

### Planned

In this state, the state file contains the text `PLANNED` only.

The state indicates that the change has been planned with a defined set of tasks and is ready for implementation.
The change folder must contain the `tasks.md` file in this state and optionally the `design.md` file for a more comprehensive approach, in addition to the files required in the `PROPOSED` state.

By reading the `tasks.md` and `design.md` files, engineers can understand the specific tasks to be completed and the design approach for implementing the change. Then, designated engineers can proceed with the implementation of the change to move it to the next state of `IMPLEMENTED`.

### Implemented

In this state, the state file contains the text `IMPLEMENTED` only.

The state indicates that the change has been implemented and is ready for archival.
The change folder must contain the `tasks.md` file where all tasks are marked completed in this state. Task list grammar is used for tasks, e.g., `- [x]` for completed tasks and `- [ ]` for pending tasks.

Specs should still remain unchanged in this state. If the implementation reveals mismatches, update the change documents (especially `spec-delta.md`) and proceed to archival to merge into specs.

By reviewing the implementation and the behavior of the system, designated engineers can verify that the change has been successfully integrated. Once verified, the change can be moved to the final state of `ARCHIVED`.

### Archived

In this state, the state file contains the text `ARCHIVED` only.

The state indicates that the change has been completed and archived for historical reference.
The spec delta in `spec-delta.md` must have been fully integrated into the relevant specs in this state. This is the point where specs are updated.
