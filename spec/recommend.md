# Recommend: task-aware model selection, plug-and-play (v1.1)

Status: **DECIDED for implementation**, 2026-10-04. Ivan's instruction, in his words: «прими наиболее подходящие
решения и запустим имплементацию… нужна уже библиотека + её соединения с task-board спавнами… желательно
плаг-н-плей». It resolves F1–F8 for v1 (§7). It builds on the landed contracts (`pkg/routing`, `pkg/evidence`,
`spec/headroom.md`, `spec/slice-b.md`) and does not change them.

## 0. The pains it removes

1. **Hand-written routing rules.** Today rules say which model and effort to use for each role, and they are
   edited by hand when a quota or budget runs out.
2. **Wrong reasoning effort.** Hard-coded rules spend strong reasoning on trivial tasks, or use too little.
   For example, Muse is usable only at `max`, and the gap from `xhigh` is large.
3. **No quality–cost frontier.** Decisions do not use an index-based Pareto curve of quality against cost,
   such as Artificial Analysis provides.
4. **No escalation.** Very hard or delicate tasks need the strongest model. Some research needs a pipeline
   that runs several of the strongest models.

## 1. Unit of choice and the catalog

The unit is a **configuration**: `runtime + model + effort`. Every effort is its own row, so choosing an
effort is choosing a row.

The **catalog** (`catalog/catalog.json`, versioned, with provenance) has one row per configuration. A row
carries:
- `runtime`, `model`, `effort`, `family` (the vendor), `billing` (`subscription | local | metered`);
- `quality`: the indices `{overall, coding, review}`, each 0–100 or absent. Every value has
  `{value, kind: measured | vendor | estimate, source, as_of, stderr?, n?}`. Optional `review`
  stays in catalog-v1 and requires `kind=measured`, `stderr` and positive integer run count `n`;
- `cost`: the relative cost per typical task, `{usd_per_task?, tokens_per_task?, kind, source, as_of}`;
- `latency`: `fast | medium | slow`, or absent;
- `constraints`: for example `{only_effort: "max"}` for Muse, or `not_for: ["windows"]`. Constraints are
  hard filters that the catalog declares.

**Tiers, computed** from `quality.coding`, falling back to `overall` when coding is absent:

| Tier | Coding index |
| --- | --- |
| `S` | ≥ 62 |
| `A` | ≥ 56 |
| `B` | ≥ 45 |
| `C` | < 45 |
| `U` | unknown |

- The default band edges are calibrated to the researched catalog’s index scale and live in the policy,
  where operators can override them.
- A row whose `stderr` places its interval across an edge takes the **lower** tier.

For `review.*`, a measured `quality.review` supplies both ranking and tiering, with separate
policy `review_tier_edges` (S ≥40, A ≥30, B ≥20). These rounded bands distinguish the pinned
Bug Hunt distribution's strongest max-effort results (~40–49), high-effort results (~30–35),
and moderate results (~20–29); they are operator-overridable routing bands, not validated
review success thresholds. Tiering uses the lower endpoint `value−stderr` as above.
Without a review measurement, retain the previous task index (coding then overall for
`review.code`; overall then coding for `review.spec`) and previous coding/overall index tiers.
The explanation carries `review_quality_from_index` and names `quality_index` explicitly.
Fallback indices and Bug Hunt percentages are different scales. In review quality
comparisons, measurements always rank above index fallbacks, including Pareto
dominance, best-quality objectives, burn and fan-out. Compare fallback indices only
among themselves; tiering retains the separate scales described above.

For review quality objectives (hard, critical, delicate and fan-out), means differing by
**less than** the larger SE tie for best quality; exact equality also ties. Ties use catalog
cost then canonical configuration key, including critical/delicate objectives. Only two
review measurements use this noise comparison; index fallbacks tie only on exact equality.
Measurements and fallbacks never share a tie band, even with equal numeric values.
Noise ties are non-transitive: anchor each band at the highest remaining mean, admit rows
that tie with that anchor, sort that band by cost/key, and repeat without chaining bands.
Balanced headroom and soft standing preferences retain this quality/cost/key order,
with matching scoped high-stakes reviewer preferences breaking ties inside each
anchored quality band as described in §3.1. Burn ties
prefer the subscription with the most known headroom before cost/key. Hard rules still apply.
Cheapest/standard objectives use provenance-first review quality in Pareto dominance
and their existing cost/headroom rules. Selector version `recommend-v3` identifies the legacy unscoped balanced
behavior; older decision records retain their version and are not silently replayed as v3.

The reproducible local CSV importer is documented in [tools/review-import](../tools/review-import/README.md).
Bug Hunt is a judge-graded TypeScript defect-finding/fixing proxy, not a measured code-review
finding rate. Its list cost and wall time are reported separately; catalog cost remains the
existing resource proxy. Unmeasured configurations receive no review value.

**The frontier.** For each task class, the catalog yields the Pareto frontier of
(quality ↑, cost ↓). A row is dominated when another row is at least as good on both axes and better on one.

## 2. Task profile (input)

```text
TaskProfile
  role         developer | reviewer | researcher | orchestrator | <any playbook role>
  task_class   code.implement | code.fix | code.refactor | code.test | review.code | review.spec |
               docs.write | research | planning | orchestration | tool-use | ops | routine   (taxonomy v1)
  difficulty   trivial | routine | standard | hard | critical          (default standard)
  sensitivity  normal | delicate                                     (default normal)
  language     optional facet (go, swift, …)
  pipeline     single | fanout                                       (default single)
```

Defaults and inference:
- An explicit `--difficulty` always wins.
- Otherwise `difficulty` comes from `task_class`:
  - `routine` and `docs.write` → `routine`;
  - `research`, `planning` and `orchestration` → `hard`;
  - everything else → `standard`.
- `sensitivity: delicate` is set by the caller, for example for security, data-loss paths or a landing on a
  shared main.

## 3. Policy (operator file, sensible defaults built in)

Budget policy `budget_mode = "economy" | "balanced" | "burn"` is independent of the
spawn execution `mode`. `--budget` overrides it for both recommend and spawn. Balanced
is the default and retains the following table and legacy unscoped-policy decisions exactly.
The policy and its override are frozen into the decision and content ID. An absent
`budget_mode` is the legacy wire representation of balanced; recommendation JSON omits
that default to preserve existing golden bytes. Human explanations always print the mode.

| Budget | Objective at every difficulty | Fan-out |
| --- | --- | --- |
| `economy` | Cheapest qualified; standard skips the ≤1.5× A upgrade | Existing distinct-family cap |
| `balanced` | Existing difficulty table below | `fanout_k`, default 3 |
| `burn` | Best quality; equal-quality ties prefer the largest known subscription headroom, then cost/key | Every qualified distinct A+ family up to `burn_fanout_cap`, default 4 |

Budget modes keep admission, locks, constraints, billing order, standing hard rules and
tier floors unchanged, including the existing delicate S-to-A fallback. Economy does
not let headroom move a more expensive row ahead of the cheapest row. Burn keeps
quality ahead of headroom, with anchored review noise bands as above. Non-review
standing preferences retain their existing bias inside each qualified tier. Unknown or stale
usage is neutral and cannot count as spare quota. The cap is configurable up to 100;
zero/omitted `burn_fanout_cap` uses four. Automatic per-runtime mode switching is deferred.
Non-balanced decisions and difficulty/sensitivity-scoped rules use `recommend-v4`;
legacy unscoped balanced decisions retain `recommend-v3`. Replay recomputes frozen
inputs and refuses any prior record whose result differs under the measurement-first
review quality rule.

The policy lives at `~/.curator/routing.toml`. A missing file means the built-in defaults. Every value is
versioned and digested into the decision.

| Difficulty | Minimum tier | Objective |
| --- | --- | --- |
| `trivial` | C | cheapest on the frontier |
| `routine` | B | cheapest on the frontier |
| `standard` | B | cheapest that meets the tier, preferring A when A costs ≤ 1.5× the B pick |
| `hard` | A | best quality among A or above; cost breaks ties |
| `critical` | S | best quality, cost ignored |

Rules:
- **Delicate tasks** require tier S and maximise quality. They never fall back below A. If no A-or-above
  candidate is admitted, the router abstains with `no_qualified_candidate`.
- **Fan-out** (`pipeline: fanout`) returns the top K = 3 configurations at tier A or above. The rows must
  come from **distinct families**, ranked by quality, so that one research question goes to several strong
  models.
- **Billing order: `subscription`, then `local`, then `metered`.**
  - Metered is used only when no subscription or local candidate qualifies.
  - When `allow_metered = false` (the default), metered is never used.
- **Headroom.** Among candidates that meet the tier, with the same tier and cost within ±25% (the
  interchangeable set), order by subscription headroom per `spec/headroom.md` §5:
  - prefer capacity that expires soon;
  - keep the reserve (`reserve_bp` 2000) for orchestrators and reviewers;
  - spread runs that are in flight;
  - treat a stale or unknown value as neutral;
  - put an over-quota candidate last.

  Headroom never moves a candidate across tiers.
- **Unknown quality** (tier U) is never selected when a known candidate meets the tier. It is shown in
  `explain`.
- **Notes do not score.** Operator notes become catalog `constraints` or policy overrides.

### 3.1 Standing rules (operator rulings as hard constraints)

The orchestrators already follow model rulings from the coordinating keeper. Until Ivan retires them, the
recommender applies them as **hard constraints**, and every recommendation names each ruling it applied. The
rules live in the policy file as a `rules` list. Each entry has:
- `id` and `source` (the ruling, for example `tb-R164`);
- `when`, exact conjunctive matches on `host`, `role`, `task_class`, `runtime`, `model`, `story`,
  plus `difficulty` (a list of trivial/routine/standard/hard/critical) and `sensitivity`
  (a list of normal/delicate). Absent lists match any value; present lists must be
  nonempty, unique and contain only valid values. Task defaults are resolved first;
  scoped-rule explanations record `rule_not_applicable: difficulty` or `sensitivity`
  when the task falls outside the rule's scope;
- one or more actions:
  - `require` / `forbid` a runtime, a model, a family or an effort;
  - `prefer` an order of runtimes or families (a soft bias within qualified tiers, with the
    high-stakes review quality-tie behavior below);
  - `effort` (a pin, or a minimum and maximum);
  - `quota_stop_bp`: a subscription whose tightest applicable known window has headroom at or below this many
    basis points is excluded. Unknown usage never triggers the stop;
  - `cross_provider_review`: a reviewer from a family other than the producer's, when `producer_family` is given;
    same-provider is allowed only as the rule's stated fallback.

Rules are evaluated in file order:
- `forbid` beats `require`;
- a rule that leaves no qualified candidate yields `no_qualified_candidate`, naming the rule.

The host is `--host`, or the policy's `host`, or the machine's short hostname.

The repository ships `docs/standing-rules.toml` with the rulings in force on 2026-10-05 (copy it into
`~/.curator/routing.toml`):
- **tb-R164 (Mac mini):** codex first everywhere; only normal trivial/routine/standard
  reviews are pinned to `gpt-6-astra` at `medium`. Hard, critical and delicate reviews
  retain their tier floors and choose the best qualified configuration, preferring
  Codex among candidates in the same anchored best-quality band. The Codex-first
  preference explicitly lists all difficulties; only matching preferences with
  their own difficulty/sensitivity scope get this demanding-review tie behavior.
  Nonmatching rules have no effect on selection or other rules. Preferences never
  move a lower-quality band or an index fallback ahead of a better measurement.
  Economy retains its cheapest-qualified objective and existing within-tier preferences.
- **tb-R139:** codex weekly window stop when headroom is 25% or less.
- **tb-R121:** cross-provider review; same-provider only as the last resort, with R164 as the operator
  exception on the mini.
- **tb-R119:** on a refusal, rotate the model. The caller passes `--exclude <runtime/model>` and the rule
  records why.
- **tb-R147:** the ivmbp pairs.
- **Per-Story operator exceptions,** for example cocoa's 380dwz on sol high, as `story` matches.

The recommender never changes a ruling by inference: replacing one is a card for Ivan.

Host limits stay with task-board and the host (disk lines, ceilings, the shared FIFO). `cmr spawn` refuses
nothing on their behalf and passes through task-board's own refusal.

## 4. Candidates (hard filter: admission is never widened)

The admitted set comes from, in order:
1. `--candidates FILE`, a JSON list of `{runtime, model, effort}`;
2. task-board's spawn-preflight for the role. `cmr` calls
   `task-board q 'project_config(view=spawn-preflight, role=R, agent=A)'` for each allowed agent and reads
   the admitted pairs;
3. if neither is available, the catalog rows whose runtime binary is on `PATH`. In that case the decision
   says so.

Every recommendation is a subset of the admitted set. Locked values are respected (R7): if the caller passes
`--agent`, `--model` or `--reasoning-effort`, the router only fills the rest in select/recommend. Shadow advisory decisions ignore forwarded selection locks; cmr-side locks still apply and forwarded role still defines admission.

## 5. Usage facts (quota)

`cmr usage refresh` runs the providerquota plans from skill-agents-management v0.5.40:
- codex app-server `account/rateLimits/read`;
- claude `-p /usage`;
- muse MSP `usage/read`;
- agy `/usage`.

These are read-only, make no inference call, and run only through the harness binaries. The command runs each
plan at most once per `(runtime, home)` per TTL (10 min), and caches the records under
`$XDG_STATE_HOME/curator/model-router/usage/`. `cmr recommend` reads the cache, refreshes stale entries when
`--refresh` is passed, and projects the records into `routing.UsageFact`, freezing `as_of` once.

The executor lives in `cmd/` (internal), never in `pkg/`; the library stays pure. This is a deliberate v1
placement of Lane 2's executor beside the router, so the router works stand-alone. Task-board's own executor
(W5) can replace it later.

## 6. Output and the task-board connection (plug-and-play)

`cmr recommend --role R --task-class C [--difficulty D] [--budget economy|balanced|burn] [--delicate] [--fanout] [--json]` prints:
- the selected configuration;
- the ranked alternatives, each with its tier, quality, cost, headroom and reason codes;
- the exact task-board spawn flags (`--agent X --model Y --reasoning-effort Z`);
- a content-addressed decision id.

The decision record is written to `$XDG_STATE_HOME/curator/model-router/decisions/`.

`cmr spawn [recommend flags] -- <task-board spawn args>` is the plug-and-play connection:
1. It queries admission with the forwarded board/root selectors, freezes the confirmation policy, and computes the recommendation, injecting only the flags the caller did not set.
2. It appends `--selection-rationale "cmr:<decision-id> <short reason>"` only when the selected provider/role preflight policy accepts or requires it (v2 ordered criteria or v3/v4 `adjustment_confirmation = "required"`). Caller rationales are preserved. For absent/equal criteria, no effective allow-set, `none`, or candidate-file/PATH admission, no rationale is injected. The decision id is always stored and printed by cmr.
3. It execs `task-board spawn` with the result.

Modes, set by `--mode` or the policy:
- `select` (the default for `cmr spawn`, since using it is the opt-in);
- `recommend`, which prints the flags and does not exec;
- `shadow`, which is fail-open: execute the original spawn args unchanged even when
  recommendation, admission, preflight, policy, usage or logging fails. Log the
  would-be decision or error as `cmr:shadow ...` on stderr and return task-board's
  exit status. Only wrapper parse errors before `--` and a missing task-board
  binary stop launch. Resolve mode from the explicit flag first, then a minimal
  independent read of the top-level policy `mode` when full loading fails. A
  readable shadow declaration survives syntax and conversion errors elsewhere;
  an unreadable policy supplies no mode and select/recommend refuse. Explicit
  `--mode shadow` always launches. Shadow preserves the completed child’s exit
  status even if copying stdout or stderr fails; failure to start still refuses.
  Append a `shadow-observation-v1` JSON line after completion, comparing the
  unconstrained advisory pick to explicit forwarded selection flags. Missing
  caller dimensions remain null and yield `unknown_caller_default`, never agreement.
  Observation-write failures warn without changing status. `cmr shadow report`
  (`--since RFC3339`, `--json`) aggregates divergence groups and lists advisory
  errors separately as FAIL-OPEN launches. Stable observation/report schemas and
  comparison denominators are documented in `docs/integration.md`.

On `no_qualified_candidate`, select/recommend refuse and never launch anything;
shadow still forwards the original launch.

A fan-out recommendation prints one spawn line per model, and `cmr spawn --fanout` launches each.

A native `task-board spawn --route` flag (W5) can replace the wrapper later without changing `cmr`.

## 7. F1–F8 for v1 (decided; reversible by editing the policy or catalog)

| | Decision |
| --- | --- |
| **F1** floors | Tier floors per difficulty (§3), from the catalog indices. No per-class empirical floors yet. |
| **F2** interchangeable | Same tier, and cost within ±25%. Headroom moves work only inside that set. |
| **F3** billing | `subscription`, then `local`, then `metered`; metered is off by default. |
| **F4** notes | Notes never score; they become constraints or policy overrides. |
| **F5–F7** | Unchanged: the internal review set is planted plus adjudicated, the metric is confirmed majors, the first class is Go review. Its results feed `quality` rows of kind `measured` later. |
| **F8** Muse | Pull-only `usage/read`; no key minting. |

## 8. Acceptance

- Deterministic: the same catalog, policy, candidates and usage snapshot give the same decision. Replay
  reproduces it.
- Golden scenarios:
  1. a trivial task picks a cheap, low-effort row;
  2. a hard task picks tier A or above;
  3. a delicate task picks the strongest row;
  4. Muse is only ever proposed at `max`;
  5. an exhausted subscription moves a standard task to an equivalent row with headroom;
  6. fan-out picks three distinct families;
  7. metered is never chosen when off;
  8. locked flags are respected;
  9. an empty qualified set refuses.
- `cmr spawn` with a fake `task-board` binary on `PATH` passes the exact args.
