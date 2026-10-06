# Native task-board routing seam

Status: proposal for task-board owners; this change adds a design note and a
library benchmark only. Source inspection used task-board `origin/main` at
`e4e971871be6c38595c4a490d8f487bd3d643a1c`. References below are relative to that
repository. The library contract is [recommend](recommend.md); the existing
external CLI connection is documented in [integration](../docs/integration.md).

## 1. Goal and non-goals

Let an opted-in project choose a runtime/model/effort with an in-process call to
`github.com/relux-works/curator-model-router/pkg/recommend`, pinned at a tag.
Add a decision ID to the RUN record. Keep the integration small: task-board
supplies an admitted set, receives a recommendation, and still validates and
claims the eventual launch.

Admission remains entirely task-board's. Explicit `--agent`, `--model`, or
`--reasoning-effort` always takes precedence over routing. The initial seam
moves no admission, limit state, confirmation, lease, scheduling, execution, or
RUN persistence code into the library. It introduces no harness probing,
network calls, subprocess routing, fan-out, or automatic policy migration.
Library difficulty floors and standing rules restrict the recommendation;
they do not become task-board admission rules. Because this seam is fail-open,
a router refusal is not a launch prohibition.

## 2. Current selection map

| Layer | File:function or type | Responsibility |
| --- | --- | --- |
| Command inputs and seed | `tools/board-cli/cmd/spawn.go:runSpawnPipeline`, `resolveSpawnAgentForCommandAt`, `applySpawnRoleSelectionDefaults` | Captures typed model/effort before defaults; resolves schedule/provider policy and affinity; seeds omitted dimensions after canonical role loading. |
| Per-role defaults | `pkg/remoteconfig/spawn_role_default.go:SpawnRoleSelectionDefaultFor`; `tools/board-cli/internal/spawn/selection_policy.go:SpawnRoleSelectionDefaultForAgentRole`, `ValidateSpawnRoleSelectionDefault` | Choice: default model/effort, validated against the role's own bounds. A default must not immediately be adjusted by its ceiling. |
| Admission and resolution | `tools/board-cli/internal/spawn/selection_policy.go:SpawnCeilingForAgentRole`, `ValidateSpawnPolicyForLaunch`, `ResolveSelectionForCeiling`, `validateSpawnAllowSet`; `tools/board-cli/internal/spawn/admission_registry.go:ExpandAdmittedPairs` | Authority: role ceilings, explicit allow-sets and supported efforts. Resolution also supplies a bounded pair. |
| Canonical admission | `pkg/remoteconfig/spawn_admission.go:AdmittedPairSet`, `Contains`, `Pairs`, `Without` | Authority: exact model/effort membership, including entry unions and pinned expansion. Bare model lists are insufficient. |
| Confirmation and rationale | `tools/board-cli/internal/spawn/selection_policy.go:ValidateSelectionConfirmation`, `validateAllowSetSelectionConfirmation`; `tools/board-cli/cmd/spawn.go:resolveSpawnSelectionPreflightWithLimits`, `appendSpawnSelectionRationaleTuple` | Authority/state: required one-line rationale, exact resolved acknowledgement and fresh rationale history; `none` forbids confirmation flags. Limit degradation occurs after this gate. |
| Basis ladder | `tools/board-cli/cmd/spawn_limits.go:resolveSpawnLimitSelection`; `tools/board-cli/internal/spawn/limit_selection.go:SeedClassForPath`, `ResolveSelectionBasis` | Choice: explicit user, configuration default, spread, affinity, within the policy's provider set. |
| Bounded fallback | `tools/board-cli/internal/spawn/limit_selection.go:resolveUnderLimits`, `chooseProvider`, `restrictToBasisGroup`, `choosePair`, `chooseEffort` | Choice inside surviving admission: honor basis provider/group, then provider fallback; closest ranked model at/below seed, otherwise least overshoot; analogous effort ordering. |
| Broker ranks | `pkg/remoteconfig/brokerrank/brokerrank.go:Rank.Compare`; `tools/board-cli/internal/spawn/limit_selection.go:registryPolicyRank` | Capability evidence: ordinals compare only inside the same broker. They are not cross-broker quality scores or independent permission. |
| Re-rank protection | `tools/board-cli/internal/spawn/catalog_rerank.go:RerankModelForTest`; `pkg/remoteconfig/spawn_admission.go:ExpandSpawnAdmission` | Test seam for vendor rank changes; admission expansion pins detect drift rather than silently changing permissions. |
| Workload classes | `tools/board-cli/internal/spawn/workload_class.go:DeriveWorkloadClass`, `WorkloadClassSnapshotDigest`; `pkg/remoteconfig/spawn_workload.go:Entries`, `SidecarSpawnProvider.LoadWorkloadClasses` | Choice metadata: role/task-derived or explicit class, advisory ordered pairs, unified fallback; snapshot provenance stays task-board's. |
| Available workload pairs | `tools/board-cli/cmd/workload_projection.go:projectWorkloadRecommendation`, `resolveWorkloadPairs`, `readWorkloadLimitSnapshot`, `validateSpawnWorkloadSelection` | Authority evidence plus advisory choice: canonical admitted universe intersected with registry/providers and limit availability. Emits `available_pairs` and separate `claimable_pairs`; configured recommendation entries do not grant permission. |
| Limit groups | `tools/board-cli/internal/spawn/limits_groups.go:SpawnLimitsGroupTable`, `ValidateSpawnLimitsGroups`; `tools/board-cli/cmd/spawn_limits.go:spawnLimitPartition`; `tools/board-cli/internal/spawn/limit_selection.go:GroupPartition`, `gateIndeterminateGroups` | Authority/state: effective group membership, free/gated/dropped classification, including unreadable-state gating. Overrides replace group rows and validate ownership/uniqueness. |
| Claims and observations | `tools/board-cli/internal/spawn/limit_selection.go:SelectUnderLimits`; `tools/board-cli/cmd/spawn_limits.go:claimForSelection`, `spawnLimitClaim`, `spawnProbeLeaseHolder.release`, `transfer`; `tools/board-cli/internal/spawn/provider_limits.go:ClassifyProviderLimit`, `RecordProviderQuota`, `RecordProviderSuccess` | Authority/state: one atomic probe claim for the selected gated group, lease cleanup/transfer, quota suppression and success recovery. |

The ladder's descriptive order is `explicit_user → config_default → spread →
affinity`. Its actual `BasisRung` values are `explicit_user`,
`legacy_preferred_agent`, `spread`, `affinity`; `config_default` is a
`SeedClass`, not a same-named rung. Legacy `preferred_agent` occupies rung two
only when `preferred_agentic_system` is absent. Exclusive systems constrain the
provider universe; mixed systems enable configured spread (preference or
round-robin), which outranks affinity. Role defaults seed the pair separately.
Provider allow-sets and schedule windows remain authority even when their
preferred order is replaced by routing.

`free` means ordinarily available, or probing under this RUN's held lease.
`gated` means eligible to attempt a claim, not permission to launch without it.
`dropped` means suppressed, probing under a live foreign lease, or refused a
claim in this selection. An unreadable limit store is not an empty healthy
store: task-board gates the admitted groups and the store decides claimability.

## 3. Exact seam and request mapping

The native call site is
`tools/board-cli/cmd/spawn_limits.go:resolveSpawnLimitSelection`, after steps
1–3 construct `providers`, `admitted`, `groupTable` and the read-only
`partition`, at step 5 immediately before `spawn.ResolveSelectionBasis`.
`runSpawnPipeline → resolveSpawnSelectionPreflightWithLimits →
resolveSpawnLimitSelection` is the launch path. Route once per frozen preflight
snapshot, not again on each confirmation retry or probe-claim retry. The
queued launch preparation must carry routing inputs/audit through its manifest;
reservation currently has a nil limit input, so a hook only in the preparation
path must not be presented as routing the earlier reservation.

Any typed selection dimension makes the launch explicit. In select, skip the
router for that spawn and preserve the existing handling of missing dimensions.
In shadow, an advisory call may use `Request.Locks` for those typed values, but
must never change the launched pair. Typed means the original CLI values, not
role defaults subsequently copied into the selection request.

| task-board input | `recommend.Request` mapping |
| --- | --- |
| `map[string]remoteconfig.AdmittedPairSet` within schedule/provider allow-set | Flatten each `Pairs()` into `Candidates`: provider → `Runtime`, pair.Model → `Model`, pair.Effort → `Effort`. Deduplicate exact triples. Empty board effort maps to library `none`; reverse on output. |
| Workload projection, when configured | Intersect canonical admission with `available_pairs ∪ claimable_pairs` from the same validated snapshot. Do not restrict to advisory `recommendations`. An indeterminate projection is a routing error, not permission to expand the set. Unconfigured ceilings are expanded from task-board's live registry, never from the router catalog or PATH. |
| Effective limit group partition | Dropped groups remove all corresponding triples and add exact `runtime/model` entries to `Exclude`. Gated triples remain candidates with `probe_required` metadata in the task-board routing envelope. Free triples remain available. Mirror indeterminate-group gating before constructing this envelope. |
| Existing cached quota measurements | `Usage.AsOf` is one frozen valid UTC Unix timestamp. Copy genuine window measurements into `Usage.Facts`, their exact scope mapping into `ScopeMap`, exact candidate keys into `UsageKeys`, and known in-flight counts into `Inflight`. Missing/stale measurements stay unknown/neutral. |
| Canonical loaded role | `Task.Role`, unchanged, including custom roles. |
| Board task class / resolved workload class | Translate to `Task.TaskClass` with the table below; retain the original class and snapshot provenance alongside the decision. |
| Optional `--difficulty`, `--delicate` | `Task.Difficulty`; `Task.Sensitivity = delicate` when true, otherwise `normal`. Omitted difficulty uses library inference. |
| This single spawn | `Task.Pipeline = single`, `Task.Platform` from task-board's platform; known language is optional. |
| Embedded catalog and decoded policy | `Catalog = LoadCatalog(DefaultCatalog())` initialized once; policy defaults plus the explicitly configured project policy file and budget override. |
| Native source | `AdmissionSource = task-board-native`; optional `Admission.RationaleRequired` describes task-board's confirmation policy. No board CLI flags are needed for this in-process caller. |

`Request` has no group-class field. A group class is a boolean availability/
claim fact, not a numeric quota measurement: never invent 0% or 100% usage from
free/gated/dropped. The task-board envelope carries gate annotations, group
snapshot digest and drop reasons beside the library explanation. The immutable
library record retains its original decision ID; adding annotations must not
rewrite its content under that ID. Host/story/rule facets, when configured,
come from caller metadata and are not inferred by consulting credential roots.

Proposed deterministic task-class translation:

| Workload class | Router task class |
| --- | --- |
| architecture | planning |
| implementation | code.implement |
| mechanical | routine |
| debugging | code.fix |
| testing | code.test |
| review | review.code for code, review.spec otherwise |
| research | research |
| documentation | docs.write |
| migration | code.refactor |
| operations | orchestration for orchestrator, ops otherwise |
| unified | code → code.implement; metadata → routine; docs → docs.write; research → research |

An explicit valid workload-class override wins over a derived class. Unknown or
unavailable derivation is a recorded routing error with ladder fallback, not an
invented profile. Native opt-in may reuse board class metadata in memory; it
must preserve existing workload snapshot/rationale validation, including stale
snapshot refusal, and must not silently synthesize caller acknowledgement.

### Selection, claims and fallback

Use `recommend.BuildDecision(request)` so the output contains a replayable
record as well as `Recommendation.Selected`, `DecisionID`, `Explanation`, and
ordered `Alternatives`. `Recommend(request)` is the equivalent pair-only API
and is what the benchmark measures. Reverse-map the selected triple into the
native pair. A normal `Recommendation.Refusal` (including no qualified
candidate) is handled like a routing error for this seam.

- **off (default):** execute the existing ladder and launch path unchanged. No
  router call, policy/cache read, extra audit field, or cursor effect.
- **shadow:** the ladder decides and only its selected group may be claimed.
  Record the router proposal/decision or error; never claim its hypothetical
  group, advance the cursor twice, or change the pair.
- **select:** for a non-explicit spawn, use router order in place of
  configuration-default, spread and affinity choice. The provider universe,
  admitted membership and limit stage still constrain every attempt. Do not
  label the router choice `explicit_user`.

A seed-only substitution is insufficient: `choosePair` may choose another
model at the same broker rank. The small task-board change should pass an
optional ordered list of native triples into `LimitSelectionInput` (proposed
field `RoutingChoices`). In `resolveUnderLimits`, after subtraction/nonempty
validation and before legacy provider/pair choice, choose the first exact
router triple that still belongs to the provider's surviving admitted set.
Empty choices execute the existing ladder. This changes a choice branch only;
`SelectUnderLimits` continues to own its existing claim loop and callbacks.

Recheck exact membership, supported effort, workload snapshot and the chosen
provider's confirmation policy before any claim or launch. In select without
typed flags, selection-specific confirmation must be validated against the
routed seed, rather than stopping on a now-replaced default beforehand. The
opt-in command adapter can defer that specific check until choice is available;
policy validity and admission checks remain early. Preserve caller rationale;
a generated decision-ID rationale is allowed only where task-board policy
accepts/requires it, and never supplies an acknowledgement or fresh human
rationale when an actual adjustment requires those. Shadow retains the old
confirmation order. Failures of task-board's gates remain task-board failures.

If a claim race drops the router pick's group, retry the next ranked candidate
from this same decision, skipping every pair in that group. Once ranked choices
are exhausted, use the original ladder over current survivors. Do not reroute
or reinstate dropped groups. Store both the original router proposal and final
selected pair; alternative/fallback execution does not change the decision ID.
Claims, errors, lease release/transfer and no-surviving-model refusal remain
in task-board. A claim I/O error is not a router error and follows existing
claim handling.

**Both shadow and select fail open on every router error:** catalog/policy
load or validation, task mapping, usage-cache read/decode, invalid request,
abstention, malformed/out-of-set output, or decision serialization/storage.
Record a sanitized error on the RUN and use the existing ladder. Fail-open
never widens admission, bypasses confirmation, or makes a refused claim free.
If the old ladder itself cannot form a seed, preserve its original refusal;
this first seam need not invent defaults for previously invalid commands.
Complete removal of early default dependence needs an owner-approved follow-up.

## 4. Configuration and CLI

Proposed addition in `pkg/remoteconfig/spawn.go:SpawnConfig`:

```go
type SpawnRoutingConfig struct {
    Mode   string `json:"mode,omitempty"`   // off | shadow | select; omitted = off
    Policy string `json:"policy,omitempty"` // project-relative or absolute file path
    Budget string `json:"budget,omitempty"` // economy | balanced | burn
}
// SpawnConfig.Routing *SpawnRoutingConfig `json:"routing,omitempty"`
```

```json
{"spawn":{"routing":{"mode":"shadow","policy":"routing/policy.toml","budget":"balanced"}}}
```

Remote spawn policy uses the same whole-section selection and config-key
rebasing as other spawn fields; an active `remote.spawn` must not unexpectedly
inherit root routing. Resolve relative policy paths against the selected
project config's directory, never the shell cwd. Omitted policy uses built-in
defaults; an explicitly missing file is a routing error. Decode JSON with
`recommend.LoadPolicy`; for TOML convert a decoded object to JSON and use the
same strict library validator. `internal/cmrio` is not an importable public
API. Do not copy the CLI's implicit user-home policy lookup into task-board.

| Flag | Proposed behavior |
| --- | --- |
| `--route[=off\|shadow\|select]` | String flag with no-argument value `select`; explicit value wins over project mode. |
| `--difficulty trivial\|routine\|standard\|hard\|critical` | Overrides task-derived library difficulty. |
| `--delicate[=true\|false]` | Explicit sensitivity; omitted means normal. |
| `--routing-budget economy\|balanced\|burn` | Overrides routing config budget, then policy `budget_mode`, then balanced. This is an objective, not a spend cap. |

Routing execution mode belongs to task-board; policy `mode` cannot enable
routing or override `--route`. Parse/read mode independently before advisory
policy work. Malformed routing policy/settings become audited ladder fallback
in an enabled mode. CLI syntax errors retain normal command validation.

The routing objective deliberately does not reuse `--budget`: `spawn.go:init`
already registers a numeric `--budget` spend cap and its deprecated
`--max-budget-usd` alias. Proposal: a separate `--routing-budget` flag, leaving
the existing numeric flags, their precedence and their goldens untouched.
Named routing objectives and numeric spend caps are separate fields. (The
alternative, one string `--budget` that accepts either form, needs a
compatibility transport and owner approval; it is not proposed for v1.)

## 5. RUN audit

Add an optional `routing` object to the existing selection audit, manifest and
RUN state/response. Carry it through queue preparation and terminalization.
Omit it entirely in off mode to preserve old serialized bytes.

| Field | Meaning |
| --- | --- |
| `mode` | Effective shadow/select mode. |
| `decision_id` | Library content ID when a decision was produced, including a normal refusal; absent on pre-decision error or explicit select bypass. |
| `selected_pair` | Actual final native runtime/model/effort after task-board validation/claims. |
| `router_pair` | Original router proposal, separate from the final launch. |
| `ladder_pair` | Old ladder's result on the frozen snapshot; in shadow, the actual ladder outcome. |
| `agreement` | Nullable boolean: router_pair equals ladder_pair exactly; null when either is unavailable. |
| `error` | Nullable sanitized `{code, message}` describing router failure/abstention. |
| `source`, `fallback`, `attempts` | Distinguish explicit bypass, router, ranked alternative and legacy fallback; retain rejected group reasons. |
| `decision`, `limit_snapshot_digest`, `gate_annotations` | Frozen record or durable reference, plus board-owned evidence needed to explain/replay admission and availability. |

For select's counterfactual `ladder_pair`, add a task-board-only read-only
preview helper over `resolveUnderLimits` with a cloned, indeterminate-gated
partition and the original seed/basis. It must not call `SelectUnderLimits`
with a nil claim callback: that would falsely drop all gated groups. A select
preview makes no promise that a hypothetical future claim would succeed.
Record it as a pre-claim preview. Shadow records the real post-claim ladder
pair. Audit/decision-write failure must not convert an advisory error into a
launch blocker; attach an error to the normal RUN write where possible. A
failure of task-board's existing required RUN persistence still follows its
existing semantics.

## 6. Per-spawn cost and benchmark

The routing path is an in-process function call with a versioned embedded
catalog decoded once. Each enabled spawn reads at most one configured policy
file and one existing usage cache snapshot, freezes time once, and computes
ranking, explanations, canonical input copies and decision ID. Missing usage
is neutral; an unreadable/malformed cache falls back with an error. No network,
subprocess, harness startup or quota refresh is part of native recommendation.
Existing task-board admission reads and probe claims are separate costs.

`pkg/recommend/benchmark_test.go:BenchmarkRecommendEmbeddedCatalog` supplies all
46 embedded catalog rows as already admitted, developer/code.implement,
balanced/default policy, one fixed clock, unknown usage and a single pipeline.
The timed loop calls `Recommend`, including normalization, input copying,
ranking, explanations and decision-ID hashing. Embedded JSON decoding/setup
and policy/cache/RUN I/O are excluded. This measures selection CPU/allocation
cost, not total spawn latency or a populated usage-cache workload.

One successful run, via the supplied minilock wrapper with `GOFLAGS=-mod=mod`
and `GOPROXY=off`, used:

```sh
go test -p 1 ./pkg/recommend/... -run . \
  -bench '^BenchmarkRecommendEmbeddedCatalog$' -benchmem -count=1
```

Result: **24,833,977 ns/op**, **6,887,556 B/op**, **210,257 allocs/op**
(55 iterations; approximately 24.8 ms per decision), Go 1.26.0, darwin/amd64,
Intel Core i5-8500B 3.00 GHz. Recommendation tests in this run passed. This
single measurement is an indicative baseline, not a latency guarantee; its
allocation volume is material and merits profiling before high-rate adoption.

## 7. Dependency step and compatibility

Before seam logic, task-board needs a separate dependency bump commit that
pins this module at an agreed tag in `tools/board-cli/go.mod`, updates
`go.sum`, regenerates `vendor/` and `vendor/modules.txt`, and passes its
`vendor-check` workflow plus existing tests. Because task-board vendors its
dependencies, its build needs no network access. The module repository is
currently private, so `vendor-check` and any `go mod download` in CI need read
access to it (GOPRIVATE plus a read credential on the runners; the credential
is an operator step, the runner wiring belongs to the CI owners). If the module
is published later, that requirement disappears. This note neither publishes a
tag nor changes task-board dependencies.

`pkg/recommend` imports only this module's `catalog`, `pkg/canonical` and
`pkg/routing`; it does not import skill-agents-management. Task-board's
`tools/board-cli/go.mod` already requires `skill-agents-management
v0.5.41-0.maint.1` while this module requires v0.5.40, so minimal version
selection keeps task-board's version and the routing code is unaffected.
Checked in a scratch copy with the newer version (`go get` there, then the full
suite): every package passes except one internal version-pin check in
`pkg/evidence` (`TestBugHuntPinnedVectors`: the evidence importer pins its source
to v0.5.40; its vectors did not change). `pkg/evidence` is not imported by
`pkg/recommend`. This module moves that pin to the newer version itself before
the integration tag, so the tagged module is tested on exactly task-board's
dependency graph.

## 8. Task-board test plan

1. Off: existing argv, response, manifest, RUN, warning and selection goldens
   remain byte-identical; no routing reads/call/extra writes/cursor advancement.
2. Shadow: across free/gated/dropped groups, errors and explicit flags, the
   launched pair equals the old ladder outcome. No hypothetical probe claim
   or duplicate cursor advancement; audit records proposal/agreement/error.
3. Select: each individual explicit flag, partial selections, combined flags
   and `--flag=value` bypass replacement. Existing ceiling adjustment/refusal
   semantics remain intact; role defaults are not mistaken for explicit input.
4. Property test: every router choice/alternative and every actual launched
   triple is in the admitted provider/pair set. Vary union entries, efforts,
   provider windows, empty admission, unknown catalog rows and registry drift.
5. Claim race: first router group loses a claim, all its pairs are skipped,
   next ranked candidate is tried, then ladder; no-survivor refusal persists.
   Validate leases release/transfer once and indeterminate state is gated.
6. Fail-open in both modes: missing/malformed policy, malformed cache, invalid
   clock/task/catalog/request, refusal, nil/out-of-set output, decision/audit
   write failure. Ladder succeeds when it normally would; admission,
   confirmation, workload and claim errors are never swallowed as router errors.
7. Determinism: shuffled admitted input with identical frozen catalog, policy,
   usage, task and clock yields identical decision ID/order. Replay the record;
   freeze cursor/limit evidence for board replay; distinguish post-claim outcome.
8. Queue/config/CLI: reservation and preparation preserve requested mode/profile
   and audit; preparation revalidates live admission and records a new decision
   if its frozen snapshot changed. Cover remote policy isolation, path resolution,
   budget precedence/numeric compatibility and required/forbidden rationales.

Use fake harnesses and injected claim/storage functions throughout.

## 9. Proposed staged end state

These are proposals for the task-board owners, not approved migrations:

1. Land this small dependency/call/audit seam, per-project off by default, and
   gather shadow evidence. Keep the existing ladder as the fallback.
2. Where enabled, use select to replace config-default, spread and affinity
   choice. Preserve explicit input and collect fallback/claim-race evidence.
3. Establish one headroom ranking in the library, fed by task-board's limit
   snapshot and truthful cached usage. Retire duplicate ranking decisions only
   after parity tests; claims and leases remain in task-board.
4. Migrate duplicated *choice* configuration (preferred system/order, spread,
   workload-class preferred pairs) into routing policy with versioned conversion,
   dry-run diff, explicit opt-in and rollback. Preserve provider restrictions,
   exclusive windows and admission ceilings as authority; never mechanically
   turn them into soft preferences.

Admission, confirmation, probe claims/leases, RUN audit and execution stay in
task-board for good.

## 10. Open owner questions

- Which tag and pinned catalog/policy version should task-board consume?
- Approve the separate `--routing-budget` spelling (proposed), or prefer a
  compatibility transport on the existing numeric `--budget`?
- Approve the task-class translation, especially review.spec, migration and
  custom roles, and define class provenance when no board class is available?
- Should v1 require an existing complete ladder seed, as proposed, or also
  reorder early validation to support routing-only commands without defaults?
- How should native opt-in transport configured workload digest/rationale
  evidence without replacing task-board's existing confirmation contract?
- Which cache supplies truthful UsageSnapshot facts, what freshness policy
  applies, and where should durable decision records referenced by RUN live?
- Approve pre-claim counterfactual ladder comparison in select, queue-time
  redecision rules, and retention of gate annotations outside the library ID?
- What shadow agreement/error/fallback evidence is sufficient to enable select,
  and who owns the later configuration migration and rollback?
