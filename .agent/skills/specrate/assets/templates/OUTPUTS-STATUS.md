# Status

## Summary

- Date: [YYYY-MM-DD]
- Enabled: [YES if `.specrate/` folder exists; NO otherwise]

## Specs

[State "No artifacts" and skip the table if no specs found or spec folder missing]

[A table of specs with ID, summary, and sanity columns filled out for each spec]

| ID | Summary | Sanity |
| -- | ------- | ------ |

[Summary is a concise summary extracted from the spec.md file]
[Sanity is a combination of "MISSING spec.md" or any other issues found; or "OK" if no issues]

## Changes

[State "No artifacts" and skip the table if no changes found or changes folder missing]

[A table of changes with ID, state, summary, affected specs, and sanity columns filled out for each change]

| ID | State | Summary | Affected Specs | Sanity |
| -- | ----- | ------- | -------------- | ------ |

[State is one of "PROPOSED", "PLANNED", "IMPLEMENTED", "ARCHIVED", "MISSING", "INVALID" from STATE file]
[Summary is a concise summary extracted from the proposal.md file]
[Affected Specs is a comma-separated list of spec IDs by checking proposal.md and spec-delta.md]
[Sanity is a combination of "OK", "MISSING {files}", "PENDING tasks in IMPLEMENTED state", or any other issues found; or "OK" if no issues]

## Suggestions

- [Suggestions if applicable]
- [Otherwise, remove this section]

## Notes

- [Additional observations or comments if applicable]
- [Otherwise, remove this section]
