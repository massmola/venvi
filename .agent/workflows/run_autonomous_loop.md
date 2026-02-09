---
description: Run the Autonomous Agent Loop (Perception, Reasoning, Action, Reflection) using the venvi-agent CLI.
---

# Autonomous Agent Loop Workflow

This workflow describes how to run the self-improving agent loop using the `venvi-agent` tool.

## Prerequisites
- `venvi-agent` binary must be built (`go build -o venvi-agent agent/main.go`).
- Familiarize yourself with `coding_standards.md` for strict environment and safety rules.

## The Loop

### 1. Goal Setting (Orchestrator)
**Action**: Define the high-level goal.
**Command**:
```bash
./venvi-agent prompt orchestrator "YOUR_GOAL_HERE"
```
**Step**: Copy the prompt output and paste it into the IDE agent chat. The Orchestrator will generate a "Ralph Wiggum" task list (named after the "I'm helping" meme, implying simple, granular, and obvious steps) with **mandatory verification steps**.

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
# 1. Start
./venvi-agent log start "fix-login-bug"

# 2. Search
./venvi-agent memory search "login"

# 3. Work... (Agent does stuff)
./venvi-agent log append "fix-login-bug" "Agent" "Updated user variable"

# 4. Checkpoint (Auto-Commit)
./venvi-agent commit "Fixed login bug"

# 5. Reflect
./venvi-agent prompt critic "fix-login-bug"
# Agent says: "We forgot to hash the password."

# 6. Learn
./venvi-agent memory add "Authentication" "Always hash passwords before saving" "security" "auth"
```
