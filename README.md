# curator-model-router

Task-aware model routing for coding agents. Given a description of the work, its role and its constraints, the router chooses a launch configuration (runtime binding + model + reasoning effort) from the candidates that are already admitted and available. It records why in a decision bound to an evidence snapshot and hands the result to the caller's existing spawn gate. It never admits anything, never launches anything, and never switches accounts or egress.

**Status: DRAFT design. No code yet; nothing here is normative.**

## Layout

| Path | Content |
| --- | --- |
| `spec/model-routing.md` | draft specification (decisions R1–R13) |
| `spec/track.md` | delivery track (slices A0 to E) |

## How it fits

- A Go library with caller adapters for task-board and `curator-run`; it creates no processes and imports none of its callers.
- Its CandidateResolver is task-board's agent-selection resolver (`skill-project-management/.specs/drafts/agent-selection.md`): hard filters answer "allowed", "compatible" and "available"; a routing strategy only ranks what passes them.
- The caller launches the chosen candidate through the agents-management launch plane (`skill-agents-management`), and the existing spawn gate re-validates it.
- Roadmap: slices A0/A1 with milestone M4, slices B–E with milestone M8, in `wiki/roadmap/ecosystem-roadmap.md`.
