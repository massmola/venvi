---
description: Run the Autonomous Agent Loop (Perception, Reasoning, Action, Reflection) using the venvi-agent CLI.
---

# Autonomous Agent Loop Workflow

This workflow describes how to run the self-improving agent loop using the `venvi-agent` tool.

## Prerequisites
- `venvi-agent` binary must be built (`go build -o venvi-agent agent/main.go`).
- Familiarize yourself with `.agent/rules/coding_standards.md` for strict environment and safety rules.

## The Loop

### 1. Goal Setting (Orchestrator)
**Action**: Define the high-level goal.
**Command**:
```bash
./venvi-agent prompt orchestrator "YOUR_GOAL_HERE"
```
**Step**: Copy the prompt output and paste it into the IDE agent chat. The Orchestrator will generate a "Ralph Wiggum" task list (named after the "I'm helping" meme, referring to extremely granular, simple, and verifiable steps) with **mandatory verification steps**. This technique ensures the agent doesn't get ahead of itself and verifies every small change.

### 2. Session Start
**Action**: Start a log for this task.
**Command**:
```bash
./venvi-agent log start "session-id"
```

### 3. Perception (Memory Retrieval)
**Action**: Check if we have done this before.
**Command**:
```bash
./venvi-agent memory search "keywords"
```
**Step**: Read any relevant skills/lessons found. Specifically look for **Rules** or **Skills** that failed previously.

### 4. Action (The Loop with Verification)
**Action**: Execute the micro-tasks from Step 1.
**Rule**: Perform verification after **every** significant change.
**Logging**: After major steps, log the action.
```bash
./venvi-agent log append "session-id" "Agent" "Executed step 1: <details>"
```

### 5. Checkpoint (Auto-Commit)
**Action**: Save progress after a successful task and **passing tests**.
**Command**:
```bash
./venvi-agent commit "Implemented login validation"
```

### 6. Reflection (Critic)
**Action**: When the task is done (or failed), run the Critic.
**Command**:
```bash
./venvi-agent prompt critic "session-id"
```
**Step**: 
1. Copy the output prompt (which includes the logs) to the IDE agent.
2. The agent will analyze errors and suggest **Systematic Improvements** to Rules, Skills, or Workflows.
3. **Save** the lesson:
```bash
./venvi-agent memory add "Topic" "Content" "tag1" "tag2"
```

## Example Cycle
```bash
# 1. Goal Setting
./venvi-agent prompt orchestrator "Implement user registration"

# 2. Session Start
./venvi-agent log start "user-reg-task"

# 3. Perception
./venvi-agent memory search "registration"

# 4. Action (Work & Logging)
# ... work done ...
./venvi-agent log append "user-reg-task" "Agent" "Implemented validation logic"

# 5. Checkpoint (Auto-Commit)
./venvi-agent commit "Implemented user registration with validation"

# 6. Reflection & Learning
./venvi-agent prompt critic "user-reg-task"
# Agent says: "Validation regex was too strict."
./venvi-agent memory add "Registration" "Use loose regex for names" "regex" "ui"
```
