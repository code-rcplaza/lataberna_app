# Skill Registry — FORGE_RPG

Generated: 2026-03-22
Project: FORGE_RPG

---

## User-Level Skills (`~/.claude/skills/`)

| Skill | Description | Trigger | Path |
|-------|-------------|---------|------|
| `go-testing` | Go testing patterns including Bubbletea TUI testing | When writing Go tests, using teatest, or adding test coverage | `~/.claude/skills/go-testing/SKILL.md` |
| `skill-creator` | Creates new AI agent skills following the Agent Skills spec | When user asks to create a new skill, add agent instructions, or document patterns for AI | `~/.claude/skills/skill-creator/SKILL.md` |
| `sdd-apply` | Implement tasks from the change, writing actual code following the specs and design | When the orchestrator launches you to implement one or more tasks from a change | `~/.claude/skills/sdd-apply/SKILL.md` |
| `sdd-archive` | Sync delta specs to main specs and archive a completed change | When the orchestrator launches you to archive a change after implementation and verification | `~/.claude/skills/sdd-archive/SKILL.md` |
| `sdd-design` | Create technical design document with architecture decisions and approach | When the orchestrator launches you to write or update the technical design for a change | `~/.claude/skills/sdd-design/SKILL.md` |
| `sdd-explore` | Explore and investigate ideas before committing to a change | When the orchestrator launches you to think through a feature, investigate the codebase, or clarify requirements | `~/.claude/skills/sdd-explore/SKILL.md` |
| `sdd-init` | Initialize Spec-Driven Development context in any project | When user wants to initialize SDD in a project, or says "sdd init", "iniciar sdd", "openspec init" | `~/.claude/skills/sdd-init/SKILL.md` |
| `sdd-propose` | Create a change proposal with intent, scope, and approach | When the orchestrator launches you to create or update a proposal for a change | `~/.claude/skills/sdd-propose/SKILL.md` |
| `sdd-spec` | Write specifications with requirements and scenarios (delta specs for changes) | When the orchestrator launches you to write or update specs for a change | `~/.claude/skills/sdd-spec/SKILL.md` |
| `sdd-tasks` | Break down a change into an implementation task checklist | When the orchestrator launches you to create or update the task breakdown for a change | `~/.claude/skills/sdd-tasks/SKILL.md` |
| `sdd-verify` | Validate that implementation matches specs, design, and tasks | When the orchestrator launches you to verify a completed (or partially completed) change | `~/.claude/skills/sdd-verify/SKILL.md` |

## Project-Level Skills

None found at `/Users/gear-rex/Desktop/Projects/FORGE_RPG/.claude/skills/`.

## Project Convention Files

| File | Path |
|------|------|
| `CLAUDE.md` (project) | `/Users/gear-rex/Desktop/Projects/FORGE_RPG/CLAUDE.md` |

---

## Shared Conventions (`~/.claude/skills/_shared/`)

| File | Description |
|------|-------------|
| `engram-convention.md` | Engram naming conventions for SDD artifacts |
| `openspec-convention.md` | OpenSpec file-based conventions |
| `persistence-contract.md` | Persistence backend contract |
| `sdd-phase-common.md` | Common patterns across SDD phases |
