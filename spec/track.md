# Track: model routing

Status: **DRAFT**, 2026-09-23. Specification: `spec/model-routing.md` (draft, R1–R13, independent review folded in). Roadmap: slice A0 inside milestone M4 (it is the agent-selection resolver), the rest in milestone M8.

## Why

Today a model is chosen for a task from a static preference list per workload class. The router chooses among already admitted candidates from evidence and records why; it never admits anything, never launches anything and never switches accounts or egress.

## How it fits agent selection

The agent-selection resolver (`skill-project-management/.specs/drafts/agent-selection.md` §7) is the router's CandidateResolver: its hard filters answer "allowed", "compatible" and "available"; its ranking step is where a routing strategy runs. Role needs from playbooks are the router's requirement features; role preferences are priors; hard requirements are never violated. Evidence is kept per role or requirement class and execution fingerprint.

## Milestones

| Slice | Delivers | When |
| --- | --- | --- |
| A0 | one resolver, `config-order` parity, decision record, contract appendix | M4 (with agent selection S1) |
| A1 | machine half: catalog bindings, `bound_on_this_machine`, dependent settings | M4 |
| B | evidence snapshots (Bug Hunt adapter, own evals), rule assessor, deterministic estimator | M8, after M5 availability |
| C | semantic assessment in shadow mode (Jev or a structured LLM, chosen explicitly) | M8 |
| D | opt-in selection at spawn, two-phase validation, effective-configuration proof | M8 |
| E | calibration and measured gain on own tasks | M8 |

## Repository

This repository, `curator-model-router` (created 2026-09-24; private until the operator decides its visibility, proposed public; consumed by task-board and the public `curator-run`). The draft specification is `spec/model-routing.md`.
