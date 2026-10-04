# Capability catalog notes — 2026-10-04

46 canonical configurations: 20 Codex, 16 Claude, 8 agy, 1 Muse and 1 local Qwen. The JSON follows the strict `catalog-v1` shape in `pkg/recommend/types.go` and `pkg/recommend/load.go`: extra metadata belongs here. It is a research catalog, not an admission grant or a claim that a binary, account entitlement, local engine or fallback setting was tested. No inference, credential access, private harness-home reads, or harness execution was performed. Only the two requested deliverables were written.

**Decision consequence:** the shipped default edges (S ≥62, A ≥56, B ≥45) produce S=5, A=7, B=12, C=21, U=1. Critical/S-only selection has five catalog configurations to consider, subject to admission, locks, constraints, and policy rules. Only OpenAI and Anthropic reach A or S, so the default up-to-three-family A-or-better fan-out returns at most two families. Scores and uncertainty allowances are unchanged.

**Vocabulary and scope.** [RegistryFacts](https://github.com/relux-works/skill-agents-management/blob/v0.5.40/pkg/vendorplugin/benchdata/bench.go) is the authority, together with [OpenAI](https://github.com/relux-works/skill-agents-management/blob/v0.5.40/pkg/vendorplugin/benchdata/openai.go), [Anthropic](https://github.com/relux-works/skill-agents-management/blob/v0.5.40/pkg/vendorplugin/benchdata/anthropic.go), [Google](https://github.com/relux-works/skill-agents-management/blob/v0.5.40/pkg/vendorplugin/benchdata/google.go), and [Muse](https://github.com/relux-works/skill-agents-management/blob/v0.5.40/pkg/vendorplugin/benchdata/muse.go). Read from the Go module cache, not reconstructed from public marketing names. RegistryFacts sorts each axis and changes EffortNone to explicit `none`; the underlying declarations identify each supported harness.

| Runtime | Canonical models | Efforts |
|---|---|---|
| codex | `gpt-6-astra` | `low`, `medium`, `high`, `xhigh`, `max` |
| codex | `gpt-6.1-sol` | `low`, `medium`, `high`, `xhigh`, `max` |
| codex | `gpt-6-sol` | `low`, `medium`, `high`, `xhigh`, `max` |
| codex | `gpt-6-luna` | `low`, `medium`, `high`, `xhigh`, `max` |
| claude | `claude-fable-5-1` | `low`, `medium`, `high`, `xhigh`, `max` |
| claude | `claude-opus-5-5` | `low`, `medium`, `high`, `xhigh`, `max` |
| claude | `claude-sonnet-5-5` | `low`, `medium`, `high`, `xhigh`, `max` |
| claude | `claude-haiku-4-5` | `none` |
| agy | `gemini-3.8-flash-high` | `none` |
| agy | `gemini-3.7-flash-high` | `none` |
| agy | `gemini-3.6-flash-high` | `none` |
| agy | `gemini-3.6-flash-medium` | `none` |
| agy | `gemini-3.6-flash-low` | `none` |
| agy | `gemini-3.5-flash-high` | `none` |
| agy | `gemini-3.5-flash-medium` | `none` |
| agy | `gemini-3.5-flash-low` | `none` |
| muse | `muse-spark-1.3-contributor` | `max` |
| local-qwen | `qwen-3.8-27b-mlx-8bit` | `none` (configured-model fixture) |

The registry’s CLI 0.159.0 capture includes GPT-6.1 Sol, GPT-6 Astra, GPT-6 Sol and GPT-6 Luna. GPT-6 Sol stays in scope even though newer Sol supersedes it. GPT-5.x, reserve and auto-review rows are outside the requested GPT-6 scope. `ultra` is intentionally retired from registry effort vocabularies; API-only `none` support does not add Codex rows. No fresh account-specific model listing was performed. Gemini 4 Argon appears on today’s public leaderboard but is absent from v0.5.40 and is not added.

Canonical aliases are not duplicate configurations: `astra` → `gpt-6-astra`, `sol` → `gpt-6.1-sol`, `luna` → `gpt-6-luna`, `opus` → `claude-opus-5-5`, `sonnet` → `claude-sonnet-5-5`, `gemini-flash` → `gemini-3.8-flash-high`, `muse-spark` → `muse-spark-1.3-contributor`. Resolve before matching. Agy suffixes encode effective low/medium/high; its separate `effort=none` means omit the CLI effort flag. Haiku also takes no effort flag.

Muse is restricted by `constraints.only_effort=max` under the brief’s operator policy, despite the registry allowing high/xhigh/max. It is not a claim that the API lacks lower efforts. The registry records older Muse 1.0.2 rejecting max and 1.3.0 accepting it. Public Bug Hunt means are 32.2/105 at max versus 20.333/105 at xhigh and 18.667/105 at high; AA coding is about 54 versus 48 at xhigh ([comparison](https://artificialanalysis.ai/agents/coding-agents/comparisons/antigravity-sdk-vs-muse-code)). The registry explicitly leaves the contributor runtime’s broker unresolved. Here family `meta` is the public model-publisher hypothesis, not verified billing-provider identity; both quality fields are estimates for that transfer.

Local is runtime `local-qwen`, system Pi, MLX, Qwen3.8-27B 8-bit, effort none, as recorded in the pinned [local-model fixture](https://github.com/relux-works/skill-agents-management/blob/v0.5.40/pkg/vendorplugin/vendors/local-models/testdata/local-qwen.toml) and [shipped-state documentation](https://github.com/relux-works/skill-agents-management/blob/v0.5.40/docs/shipped-state.md). Local configured models are explicitly outside compiled RegistryFacts; using this fixture is the necessary documented exception, not a hosted-Qwen registry substitution. Its Gate B deployment, weights digest, quantizer, effective thinking mode, availability and latency are unverified. Quality stays absent/U. The public hosted 8-bit Qwen score is not a local measurement. USD/task=0 means incremental API charge only, excluding electricity, hardware, contention and human time. Latency classification is omitted because no local measurement supports fast/medium/slow.

**Quality scales and rules.** `overall` uses AA Intelligence Index v4.3.2 points; `coding` uses AA Coding Agent Index v1.5 ×100. Neither is a probability of success on our projects. Bug Hunt percentages and SWE-bench percentages are retained separately below; they are never averaged with AA. Exact measured rows describe the publisher’s tested system, including its older harness versions, served aliases, fallbacks and reference scaffold—not a fresh test of our current installation. The supplied research reports remain methodological context rather than substituted current data.

AA does not supply an applicable standard error in the inspected index records; measured fields omit `stderr`, which means unknown, never zero uncertainty. No binomial SE is fabricated from correlated tasks/attempts. An estimate’s `stderr` is an authored uncertainty allowance, not an empirically calibrated sampling SE: base 5 points plus 2 for each effort step. It is wider than the base allowance for unreported anchor uncertainty. E1 copies the nearest measured same-model coding effort and adds the difference in that model’s measured overall index. Distance uses low < medium < high < xhigh < max; ties choose lower effort; clamp to [0,100]. This is a transparent prior, not validated cross-index calibration.

G1 estimates missing Gemini overall efforts from the same model’s measured high: subtract 3 points per step to medium/low; stderr 7/9. AA’s own starred 3.5-medium estimate (34) is not promoted to a measurement and is not the anchor. G2 and M1 preserve point values while widening uncertainty for harness/identity transfers. H1 preserves the non-reasoning Haiku model observation as an estimate because an absent effort control does not establish thinking mode. Missing entire coding anchors (Haiku and Gemini 3.5–3.7) remain absent; inventing a cross-model conversion from Bug Hunt or SWE-bench would mix scales.

Tier is computed from coding, otherwise overall, using the lower end of `[value−stderr,value+stderr]` when stderr exists; without it use the point value and flag unknown uncertainty. This matches the inspected v1 loader’s `TierFor` (one SE, not a 95% interval). No stored tier can become stale in the JSON. Repair/implementation observations do not establish Go-review, Swift, docs, planning or orchestration accuracy.

**Prices and task cost.** All remote rows are subscription-funded by instruction; no metered row or API-key path is assumed. `usd_per_task` is an estimated list-equivalent resource proxy, not a bill, quota debit or marginal subscription price. The subscription fee and remaining headroom are unknown. Actual subscription incremental cash cost alone would make remote costs tie and collapse the frontier to highest quality; the tables below deliberately use the declared list proxy.

| Model | Input $/1M | Output $/1M | Cache-read $/1M | Source |
|---|---:|---:|---:|---|
| `claude-fable-5-1` | 10 | 50 | 0.25 | [rate](https://platform.claude.com/docs/en/about-claude/pricing) |
| `claude-haiku-4-5` | 1 | 5 | 0.1 | [rate](https://platform.claude.com/docs/en/about-claude/pricing) |
| `claude-opus-5-5` | 4 | 20 | 0.2 | [rate](https://platform.claude.com/docs/en/about-claude/pricing) |
| `claude-sonnet-5-5` | 2 | 10 | 0.2 | [rate](https://platform.claude.com/docs/en/about-claude/pricing) |
| `gemini-3.5-flash-high` | 1.5 | 9 | 0.15 | [rate](https://ai.google.dev/gemini-api/docs/pricing) |
| `gemini-3.5-flash-low` | 1.5 | 9 | 0.15 | [rate](https://ai.google.dev/gemini-api/docs/pricing) |
| `gemini-3.5-flash-medium` | 1.5 | 9 | 0.15 | [rate](https://ai.google.dev/gemini-api/docs/pricing) |
| `gemini-3.6-flash-high` | 0.75 | 3.75 | 0.075 | [rate](https://ai.google.dev/gemini-api/docs/pricing) |
| `gemini-3.6-flash-low` | 0.75 | 3.75 | 0.075 | [rate](https://ai.google.dev/gemini-api/docs/pricing) |
| `gemini-3.6-flash-medium` | 0.75 | 3.75 | 0.075 | [rate](https://ai.google.dev/gemini-api/docs/pricing) |
| `gemini-3.7-flash-high` | 0.75 | 3.75 | 0.075 | [rate](https://ai.google.dev/gemini-api/docs/pricing) |
| `gemini-3.8-flash-high` | 0.75 | 3.75 | 0.075 | [rate](https://ai.google.dev/gemini-api/docs/pricing) |
| `gpt-6-astra` | 10 | 50 | 1 | [rate](https://developers.openai.com/api/docs/models/gpt-6-astra) |
| `gpt-6-luna` | 0.1 | 0.5 | 0.01 | [rate](https://developers.openai.com/api/docs/models/gpt-6-luna) |
| `gpt-6-sol` | 2 | 10 | 0.2 | [rate](https://developers.openai.com/api/docs/models/gpt-6-sol) |
| `gpt-6.1-sol` | 2 | 10 | 0.1 | [rate](https://developers.openai.com/api/docs/models/gpt-6.1-sol) |
| `muse-spark-1.3-contributor` | 1.25 | 4.25 | 0.15 | [rate](https://artificialanalysis.ai/models/comparisons/muse-spark-1-3-xhigh-vs-muse-spark-1-2) |

Prices are standard text API-equivalent rates observed 2026-10-04. Gemini 3.6–3.8 rates are introductory through 2026-12-31. Muse rates describe public Meta API Spark 1.3, not a proven contributor billing contract. No cache-write premium, cache-storage charges, tool charges, long-context premium, fast/batch/regional pricing or fallback-model repricing is included. These are material limitations, not measured savings.

For each row, `USD = ((I−C)×Pinput + C×Pcache + O×Poutput)/1e6`, where I includes cache reads C, and O includes reasoning output. `tokens_per_task=I+O`. All budgets are rounded to whole tokens before pricing. Exact AA coding rows provide workload mean token counts; unmeasured efforts scale the nearest same-model budget by effort factors low=.25, medium=.40, high=.60, xhigh=.80, max=1. These multipliers are assumptions, not quota conversion. Agy uses the 3.8/high SDK workload with low=.5, medium=.75, high=1, including cross-generation transfer. Haiku uses an authored 100k-in/8k-out small task. Cache tokens are assumed reads because the public aggregate does not split writes; therefore recomputed costs need not equal AA’s published cost. Per-request context length is unavailable even when task totals exceed a million tokens.

| Configuration | Input I | Cached C | Output O | Estimated $/task | Budget source/rule |
|---|---:|---:|---:|---:|---|
| `agy/gemini-3.5-flash-high/none` | 13395247 | 11871692 | 115826 | 5.108520 | Gemini 3.8/high SDK mean × 1; estimated CLI and generation transfer |
| `agy/gemini-3.5-flash-low/none` | 6697623 | 5935846 | 57913 | 2.554259 | Gemini 3.8/high SDK mean × 0.5; estimated CLI and generation transfer |
| `agy/gemini-3.5-flash-medium/none` | 10046435 | 8903769 | 86870 | 3.831394 | Gemini 3.8/high SDK mean × 0.75; estimated CLI and generation transfer |
| `agy/gemini-3.6-flash-high/none` | 13395247 | 11871692 | 115826 | 2.467391 | Gemini 3.8/high SDK mean × 1; estimated CLI and generation transfer |
| `agy/gemini-3.6-flash-low/none` | 6697623 | 5935846 | 57913 | 1.233695 | Gemini 3.8/high SDK mean × 0.5; estimated CLI and generation transfer |
| `agy/gemini-3.6-flash-medium/none` | 10046435 | 8903769 | 86870 | 1.850545 | Gemini 3.8/high SDK mean × 0.75; estimated CLI and generation transfer |
| `agy/gemini-3.7-flash-high/none` | 13395247 | 11871692 | 115826 | 2.467391 | Gemini 3.8/high SDK mean × 1; estimated CLI and generation transfer |
| `agy/gemini-3.8-flash-high/none` | 13395247 | 11871692 | 115826 | 2.467391 | Gemini 3.8/high SDK mean × 1; estimated CLI transfer |
| `claude/claude-fable-5-1/high` | 3348973 | 3123333 | 80179 | 7.046183 | AA coding mean at max × 0.6; estimated effort scaling |
| `claude/claude-fable-5-1/low` | 1395406 | 1301389 | 33408 | 2.935917 | AA coding mean at max × 0.25; estimated effort scaling |
| `claude/claude-fable-5-1/max` | 5581622 | 5205556 | 133632 | 11.743649 | AA coding mean at max × 1; exact effort budget |
| `claude/claude-fable-5-1/medium` | 2232649 | 2082222 | 53453 | 4.697476 | AA coding mean at max × 0.4; estimated effort scaling |
| `claude/claude-fable-5-1/xhigh` | 4465298 | 4164445 | 106905 | 9.394891 | AA coding mean at max × 0.8; estimated effort scaling |
| `claude/claude-haiku-4-5/none` | 100000 | 0 | 8000 | 0.140000 | Authored small-task assumption: 100000 input, 8000 output, no cache; no workload telemetry |
| `claude/claude-opus-5-5/high` | 9130924 | 8741354 | 199906 | 7.304671 | AA coding mean at max × 0.6; estimated effort scaling |
| `claude/claude-opus-5-5/low` | 3804552 | 3642231 | 83294 | 3.043610 | AA coding mean at max × 0.25; estimated effort scaling |
| `claude/claude-opus-5-5/max` | 15218206 | 14568924 | 333177 | 12.174453 | AA coding mean at max × 1; exact effort budget |
| `claude/claude-opus-5-5/medium` | 6087283 | 5827569 | 133271 | 4.869790 | AA coding mean at max × 0.4; estimated effort scaling |
| `claude/claude-opus-5-5/xhigh` | 12174565 | 11655139 | 266541 | 9.739552 | AA coding mean at max × 0.8; estimated effort scaling |
| `claude/claude-sonnet-5-5/high` | 2610884 | 2507283 | 41496 | 1.123619 | AA coding mean at high × 1; exact effort budget |
| `claude/claude-sonnet-5-5/low` | 961481 | 917183 | 15804 | 0.430073 | AA coding mean at low × 1; exact effort budget |
| `claude/claude-sonnet-5-5/max` | 27141369 | 26006028 | 600968 | 13.481568 | AA coding mean at max × 1; exact effort budget |
| `claude/claude-sonnet-5-5/medium` | 1232831 | 1173805 | 20521 | 0.558023 | AA coding mean at medium × 1; exact effort budget |
| `claude/claude-sonnet-5-5/xhigh` | 6961982 | 6693764 | 125696 | 3.132149 | AA coding mean at xhigh × 1; exact effort budget |
| `codex/gpt-6-astra/high` | 1979730 | 1884953 | 28237 | 4.244573 | AA coding mean at max × 0.6; estimated effort scaling |
| `codex/gpt-6-astra/low` | 824887 | 785397 | 11765 | 1.768547 | AA coding mean at max × 0.25; estimated effort scaling |
| `codex/gpt-6-astra/max` | 3299550 | 3141589 | 47061 | 7.074249 | AA coding mean at max × 1; exact effort budget |
| `codex/gpt-6-astra/medium` | 1319820 | 1256636 | 18824 | 2.829676 | AA coding mean at max × 0.4; estimated effort scaling |
| `codex/gpt-6-astra/xhigh` | 2639640 | 2513271 | 37649 | 5.659411 | AA coding mean at max × 0.8; estimated effort scaling |
| `codex/gpt-6-luna/high` | 6071837 | 5938388 | 59005 | 0.102231 | AA coding mean at max × 0.6; estimated effort scaling |
| `codex/gpt-6-luna/low` | 2529932 | 2474328 | 24585 | 0.042596 | AA coding mean at max × 0.25; estimated effort scaling |
| `codex/gpt-6-luna/max` | 10119728 | 9897313 | 98341 | 0.170385 | AA coding mean at max × 1; exact effort budget |
| `codex/gpt-6-luna/medium` | 4047891 | 3958925 | 39337 | 0.068154 | AA coding mean at max × 0.4; estimated effort scaling |
| `codex/gpt-6-luna/xhigh` | 8095782 | 7917850 | 78673 | 0.136308 | AA coding mean at max × 0.8; estimated effort scaling |
| `codex/gpt-6-sol/high` | 5867161 | 5759357 | 37305 | 1.740529 | AA coding mean at max × 0.6; estimated effort scaling |
| `codex/gpt-6-sol/low` | 2444651 | 2399732 | 15544 | 0.725224 | AA coding mean at max × 0.25; estimated effort scaling |
| `codex/gpt-6-sol/max` | 9778602 | 9598929 | 62176 | 2.900892 | AA coding mean at max × 1; exact effort budget |
| `codex/gpt-6-sol/medium` | 3911441 | 3839572 | 24870 | 1.160352 | AA coding mean at max × 0.4; estimated effort scaling |
| `codex/gpt-6-sol/xhigh` | 7822882 | 7679143 | 49740 | 2.320707 | AA coding mean at max × 0.8; estimated effort scaling |
| `codex/gpt-6.1-sol/high` | 2791163 | 2665415 | 28749 | 0.805527 | AA coding mean at high × 1; exact effort budget |
| `codex/gpt-6.1-sol/low` | 1699124 | 1627296 | 14695 | 0.453336 | AA coding mean at low × 1; exact effort budget |
| `codex/gpt-6.1-sol/max` | 3937843 | 3740552 | 61659 | 1.385227 | AA coding mean at max × 1; exact effort budget |
| `codex/gpt-6.1-sol/medium` | 2275867 | 2176437 | 22412 | 0.640624 | AA coding mean at medium × 1; exact effort budget |
| `codex/gpt-6.1-sol/xhigh` | 3177125 | 3036066 | 35403 | 0.939755 | AA coding mean at xhigh × 1; exact effort budget |
| `local-qwen/qwen-3.8-27b-mlx-8bit/none` | unknown | unknown | unknown | 0 | Local API charge only; latency unknown |
| `muse/muse-spark-1.3-contributor/max` | 16445971 | 16090656 | 263474 | 3.977507 | AA coding mean at max × 1; exact effort budget |

**Computed tiers and quality.** M=measured public result; E=estimate; ± is one stated uncertainty allowance. Unknown SE is not a confidence guarantee.

| Configuration | Overall | Coding | Tier |
|---|---:|---:|:---:|
| `agy/gemini-3.5-flash-high/none` | 33.000 M (SE unknown) | — | C |
| `agy/gemini-3.5-flash-low/none` | 27.000 E ±9 | — | C |
| `agy/gemini-3.5-flash-medium/none` | 30.000 E ±7 | — | C |
| `agy/gemini-3.6-flash-high/none` | 34.000 M (SE unknown) | — | C |
| `agy/gemini-3.6-flash-low/none` | 28.000 E ±9 | — | C |
| `agy/gemini-3.6-flash-medium/none` | 31.000 E ±7 | — | C |
| `agy/gemini-3.7-flash-high/none` | 39.000 M (SE unknown) | — | C |
| `agy/gemini-3.8-flash-high/none` | 41.000 M (SE unknown) | 41.863 E ±6 | C |
| `claude/claude-fable-5-1/high` | 51.000 M (SE unknown) | 60.222 E ±9 | B |
| `claude/claude-fable-5-1/low` | 47.000 M (SE unknown) | 56.222 E ±13 | C |
| `claude/claude-fable-5-1/max` | 53.000 M (SE unknown) | 62.222 M (SE unknown) | S |
| `claude/claude-fable-5-1/medium` | 49.000 M (SE unknown) | 58.222 E ±11 | B |
| `claude/claude-fable-5-1/xhigh` | 53.000 M (SE unknown) | 62.222 E ±7 | B |
| `claude/claude-haiku-4-5/none` | 15.000 E ±5 | — | C |
| `claude/claude-opus-5-5/high` | 54.000 M (SE unknown) | 61.989 E ±9 | B |
| `claude/claude-opus-5-5/low` | 42.000 M (SE unknown) | 49.989 E ±13 | C |
| `claude/claude-opus-5-5/max` | 58.000 M (SE unknown) | 65.989 M (SE unknown) | S |
| `claude/claude-opus-5-5/medium` | 51.000 M (SE unknown) | 58.989 E ±11 | B |
| `claude/claude-opus-5-5/xhigh` | 56.000 M (SE unknown) | 63.989 E ±7 | A |
| `claude/claude-sonnet-5-5/high` | 47.000 M (SE unknown) | 55.012 M (SE unknown) | B |
| `claude/claude-sonnet-5-5/low` | 36.000 M (SE unknown) | 42.059 M (SE unknown) | C |
| `claude/claude-sonnet-5-5/max` | 56.000 M (SE unknown) | 68.358 M (SE unknown) | S |
| `claude/claude-sonnet-5-5/medium` | 41.000 M (SE unknown) | 45.884 M (SE unknown) | B |
| `claude/claude-sonnet-5-5/xhigh` | 52.000 M (SE unknown) | 62.871 M (SE unknown) | S |
| `codex/gpt-6-astra/high` | 51.000 M (SE unknown) | 59.645 E ±9 | B |
| `codex/gpt-6-astra/low` | 46.000 M (SE unknown) | 54.645 E ±13 | C |
| `codex/gpt-6-astra/max` | 53.000 M (SE unknown) | 61.645 M (SE unknown) | A |
| `codex/gpt-6-astra/medium` | 50.000 M (SE unknown) | 58.645 E ±11 | B |
| `codex/gpt-6-astra/xhigh` | 52.000 M (SE unknown) | 60.645 E ±7 | B |
| `codex/gpt-6-luna/high` | 33.000 M (SE unknown) | 36.074 E ±9 | C |
| `codex/gpt-6-luna/low` | 22.000 M (SE unknown) | 25.074 E ±13 | C |
| `codex/gpt-6-luna/max` | 38.000 M (SE unknown) | 41.074 M (SE unknown) | C |
| `codex/gpt-6-luna/medium` | 30.000 M (SE unknown) | 33.074 E ±11 | C |
| `codex/gpt-6-luna/xhigh` | 35.000 M (SE unknown) | 38.074 E ±7 | C |
| `codex/gpt-6-sol/high` | 42.000 M (SE unknown) | 50.663 E ±9 | C |
| `codex/gpt-6-sol/low` | 34.000 M (SE unknown) | 42.663 E ±13 | C |
| `codex/gpt-6-sol/max` | 48.000 M (SE unknown) | 56.663 M (SE unknown) | A |
| `codex/gpt-6-sol/medium` | 40.000 M (SE unknown) | 48.663 E ±11 | C |
| `codex/gpt-6-sol/xhigh` | 44.000 M (SE unknown) | 52.663 E ±7 | B |
| `codex/gpt-6.1-sol/high` | 50.000 M (SE unknown) | 60.149 M (SE unknown) | A |
| `codex/gpt-6.1-sol/low` | 42.000 M (SE unknown) | 57.216 M (SE unknown) | A |
| `codex/gpt-6.1-sol/max` | 52.000 M (SE unknown) | 60.148 M (SE unknown) | A |
| `codex/gpt-6.1-sol/medium` | 48.000 M (SE unknown) | 61.415 M (SE unknown) | A |
| `codex/gpt-6.1-sol/xhigh` | 51.000 M (SE unknown) | 62.908 M (SE unknown) | S |
| `local-qwen/qwen-3.8-27b-mlx-8bit/none` | — | — | U |
| `muse/muse-spark-1.3-contributor/max` | 48.000 E ±5 | 54.303 E ±6 | B |

**Pareto frontier by task class.** Higher quality and lower list-equivalent USD/task; dominance requires at least as good on both, strictly better on one. No rounding before comparisons beyond stored six-decimal values; equal points both survive. Use coding for code.*, review.code, tool-use and ops, falling back to overall; use overall for the remaining classes, falling back to coding. This task-to-index map is an explicit v1 proxy assumption, not measured per-class performance. Admission, tier floors, platform filters and headroom must be applied by the consumer; these are unrestricted catalog frontiers. Every eligible row here is subscription-funded. Local/U is excluded because missing quality is incomparable, not zero. The conservative column repeats dominance using value−stderr; absent SE still cannot establish robustness.

| Task class | Index | Point frontier | Conservative frontier |
|---|---|---|---|
| code.implement | coding | `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `claude/claude-sonnet-5-5/low`; `claude/claude-sonnet-5-5/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` | `claude/claude-opus-5-5/max`; `claude/claude-sonnet-5-5/low`; `claude/claude-sonnet-5-5/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` |
| code.fix | coding | `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `claude/claude-sonnet-5-5/low`; `claude/claude-sonnet-5-5/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` | `claude/claude-opus-5-5/max`; `claude/claude-sonnet-5-5/low`; `claude/claude-sonnet-5-5/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` |
| code.refactor | coding | `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `claude/claude-sonnet-5-5/low`; `claude/claude-sonnet-5-5/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` | `claude/claude-opus-5-5/max`; `claude/claude-sonnet-5-5/low`; `claude/claude-sonnet-5-5/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` |
| code.test | coding | `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `claude/claude-sonnet-5-5/low`; `claude/claude-sonnet-5-5/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` | `claude/claude-opus-5-5/max`; `claude/claude-sonnet-5-5/low`; `claude/claude-sonnet-5-5/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` |
| review.code | coding | `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `claude/claude-sonnet-5-5/low`; `claude/claude-sonnet-5-5/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` | `claude/claude-opus-5-5/max`; `claude/claude-sonnet-5-5/low`; `claude/claude-sonnet-5-5/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` |
| review.spec | overall | `claude/claude-opus-5-5/high`; `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `codex/gpt-6-astra/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/high`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/max`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` | `claude/claude-opus-5-5/high`; `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `codex/gpt-6-astra/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/high`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/max`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` |
| docs.write | overall | `claude/claude-opus-5-5/high`; `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `codex/gpt-6-astra/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/high`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/max`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` | `claude/claude-opus-5-5/high`; `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `codex/gpt-6-astra/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/high`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/max`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` |
| research | overall | `claude/claude-opus-5-5/high`; `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `codex/gpt-6-astra/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/high`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/max`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` | `claude/claude-opus-5-5/high`; `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `codex/gpt-6-astra/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/high`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/max`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` |
| planning | overall | `claude/claude-opus-5-5/high`; `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `codex/gpt-6-astra/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/high`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/max`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` | `claude/claude-opus-5-5/high`; `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `codex/gpt-6-astra/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/high`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/max`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` |
| orchestration | overall | `claude/claude-opus-5-5/high`; `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `codex/gpt-6-astra/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/high`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/max`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` | `claude/claude-opus-5-5/high`; `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `codex/gpt-6-astra/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/high`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/max`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` |
| tool-use | coding | `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `claude/claude-sonnet-5-5/low`; `claude/claude-sonnet-5-5/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` | `claude/claude-opus-5-5/max`; `claude/claude-sonnet-5-5/low`; `claude/claude-sonnet-5-5/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` |
| ops | coding | `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `claude/claude-sonnet-5-5/low`; `claude/claude-sonnet-5-5/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` | `claude/claude-opus-5-5/max`; `claude/claude-sonnet-5-5/low`; `claude/claude-sonnet-5-5/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` |
| routine | overall | `claude/claude-opus-5-5/high`; `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `codex/gpt-6-astra/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/high`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/max`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` | `claude/claude-opus-5-5/high`; `claude/claude-opus-5-5/max`; `claude/claude-opus-5-5/xhigh`; `codex/gpt-6-astra/max`; `codex/gpt-6-luna/high`; `codex/gpt-6-luna/low`; `codex/gpt-6-luna/max`; `codex/gpt-6-luna/medium`; `codex/gpt-6-luna/xhigh`; `codex/gpt-6.1-sol/high`; `codex/gpt-6.1-sol/low`; `codex/gpt-6.1-sol/max`; `codex/gpt-6.1-sol/medium`; `codex/gpt-6.1-sol/xhigh` |

Frontier coordinates (same row may appear on both):

| Index | Configuration | Quality | $/task |
|---|---|---:|---:|
| coding | `claude/claude-opus-5-5/max` | 65.988580 | 12.174453 |
| coding | `claude/claude-opus-5-5/xhigh` | 63.988580 | 9.739552 |
| coding | `claude/claude-sonnet-5-5/low` | 42.059308 | 0.430073 |
| coding | `claude/claude-sonnet-5-5/max` | 68.357834 | 13.481568 |
| coding | `codex/gpt-6-luna/high` | 36.074389 | 0.102231 |
| coding | `codex/gpt-6-luna/low` | 25.074389 | 0.042596 |
| coding | `codex/gpt-6-luna/max` | 41.074389 | 0.170385 |
| coding | `codex/gpt-6-luna/medium` | 33.074389 | 0.068154 |
| coding | `codex/gpt-6-luna/xhigh` | 38.074389 | 0.136308 |
| coding | `codex/gpt-6.1-sol/low` | 57.216349 | 0.453336 |
| coding | `codex/gpt-6.1-sol/medium` | 61.414747 | 0.640624 |
| coding | `codex/gpt-6.1-sol/xhigh` | 62.907767 | 0.939755 |
| overall | `claude/claude-opus-5-5/high` | 54.000000 | 7.304671 |
| overall | `claude/claude-opus-5-5/max` | 58.000000 | 12.174453 |
| overall | `claude/claude-opus-5-5/xhigh` | 56.000000 | 9.739552 |
| overall | `codex/gpt-6-astra/max` | 53.000000 | 7.074249 |
| overall | `codex/gpt-6-luna/high` | 33.000000 | 0.102231 |
| overall | `codex/gpt-6-luna/low` | 22.000000 | 0.042596 |
| overall | `codex/gpt-6-luna/max` | 38.000000 | 0.170385 |
| overall | `codex/gpt-6-luna/medium` | 30.000000 | 0.068154 |
| overall | `codex/gpt-6-luna/xhigh` | 35.000000 | 0.136308 |
| overall | `codex/gpt-6.1-sol/high` | 50.000000 | 0.805527 |
| overall | `codex/gpt-6.1-sol/low` | 42.000000 | 0.453336 |
| overall | `codex/gpt-6.1-sol/max` | 52.000000 | 1.385227 |
| overall | `codex/gpt-6.1-sol/medium` | 48.000000 | 0.640624 |
| overall | `codex/gpt-6.1-sol/xhigh` | 51.000000 | 0.939755 |

**Row-specific provenance and estimates.** Quality source URLs are in each JSON metric; the following states the transfer or derivation.

- `agy/gemini-3.5-flash-high/none`: G1: model id fixes effective effort high; separate effort is none. Overall: exact high AA API-model observation. Coding absent: no same-model AA coding anchor found; Bug Hunt receipts stay on their own scale. Tier falls back to overall.
- `agy/gemini-3.5-flash-low/none`: G1: model id fixes effective effort low; separate effort is none. Overall estimate: same-model measured high 33 minus 3×2; stderr=9. Coding absent: no same-model AA coding anchor found; Bug Hunt receipts stay on their own scale. Tier falls back to overall.
- `agy/gemini-3.5-flash-medium/none`: G1: model id fixes effective effort medium; separate effort is none. Overall estimate: same-model measured high 33 minus 3×1; stderr=7. Coding absent: no same-model AA coding anchor found; Bug Hunt receipts stay on their own scale. Tier falls back to overall.
- `agy/gemini-3.6-flash-high/none`: G1: model id fixes effective effort high; separate effort is none. Overall: exact high AA API-model observation. Coding absent: no same-model AA coding anchor found; Bug Hunt receipts stay on their own scale. Tier falls back to overall.
- `agy/gemini-3.6-flash-low/none`: G1: model id fixes effective effort low; separate effort is none. Overall estimate: same-model measured high 34 minus 3×2; stderr=9. Coding absent: no same-model AA coding anchor found; Bug Hunt receipts stay on their own scale. Tier falls back to overall.
- `agy/gemini-3.6-flash-medium/none`: G1: model id fixes effective effort medium; separate effort is none. Overall estimate: same-model measured high 34 minus 3×1; stderr=7. Coding absent: no same-model AA coding anchor found; Bug Hunt receipts stay on their own scale. Tier falls back to overall.
- `agy/gemini-3.7-flash-high/none`: G1: model id fixes effective effort high; separate effort is none. Overall: exact high AA API-model observation. Coding absent: no same-model AA coding anchor found; Bug Hunt receipts stay on their own scale. Tier falls back to overall.
- `agy/gemini-3.8-flash-high/none`: G1: model id fixes effective effort high; separate effort is none. Overall: exact high AA API-model observation. G2: copy exact-model/high Antigravity SDK AA coding score with stderr 6 for SDK-to-CLI transfer; not a measured agy result.
- `claude/claude-fable-5-1/high`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Default fallback included. Coding E1: nearest measured same-model effort max, AA coding 62.222496 + overall delta -2; stderr=5+2×2=9.
- `claude/claude-fable-5-1/low`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Default fallback included. Coding E1: nearest measured same-model effort max, AA coding 62.222496 + overall delta -6; stderr=5+2×4=13.
- `claude/claude-fable-5-1/max`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Default fallback included. Coding: AA exact model/effort in the named harness; raw row 29a1cdad0d140781d05b9b4746cfdd0d.
- `claude/claude-fable-5-1/medium`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Default fallback included. Coding E1: nearest measured same-model effort max, AA coding 62.222496 + overall delta -4; stderr=5+2×3=11.
- `claude/claude-fable-5-1/xhigh`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Default fallback included. Coding E1: nearest measured same-model effort max, AA coding 62.222496 + overall delta +0; stderr=5+2×1=7.
- `claude/claude-haiku-4-5/none`: H1: copy AA non-reasoning Haiku overall 15 with stderr 5; registry none means no effort flag, not proof of disabled thinking. Coding absent: vendor SWE-bench used 128K thinking and a different scaffold.
- `claude/claude-opus-5-5/high`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Default fallback included. Coding E1: nearest measured same-model effort max, AA coding 65.988580 + overall delta -4; stderr=5+2×2=9.
- `claude/claude-opus-5-5/low`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Default fallback included. Coding E1: nearest measured same-model effort max, AA coding 65.988580 + overall delta -16; stderr=5+2×4=13.
- `claude/claude-opus-5-5/max`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Default fallback included. Coding: AA exact model/effort in the named harness; raw row ab5074e404be48bbd4ac3c488ad9a3e2.
- `claude/claude-opus-5-5/medium`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Default fallback included. Coding E1: nearest measured same-model effort max, AA coding 65.988580 + overall delta -7; stderr=5+2×3=11.
- `claude/claude-opus-5-5/xhigh`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Default fallback included. Coding E1: nearest measured same-model effort max, AA coding 65.988580 + overall delta -2; stderr=5+2×1=7.
- `claude/claude-sonnet-5-5/high`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Default fallback included. Coding: AA exact model/effort in the named harness; raw row 33e80f383dc644527161435b076798ad.
- `claude/claude-sonnet-5-5/low`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Default fallback included. Coding: AA exact model/effort in the named harness; raw row e9f9d3ae5b03aa060114dd4b3275d67d.
- `claude/claude-sonnet-5-5/max`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Default fallback included. Coding: AA exact model/effort in the named harness; raw row a2c87c062f3cef73d8525e7578f14742.
- `claude/claude-sonnet-5-5/medium`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Default fallback included. Coding: AA exact model/effort in the named harness; raw row 4d6aa9e9c8df931d54dba5956d94973c.
- `claude/claude-sonnet-5-5/xhigh`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Default fallback included. Coding: AA exact model/effort in the named harness; raw row 666605c66ad569c503ce53c01682e285.
- `codex/gpt-6-astra/high`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding E1: nearest measured same-model effort max, AA coding 61.645045 + overall delta -2; stderr=5+2×2=9.
- `codex/gpt-6-astra/low`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding E1: nearest measured same-model effort max, AA coding 61.645045 + overall delta -7; stderr=5+2×4=13.
- `codex/gpt-6-astra/max`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding: AA exact model/effort in the named harness; raw row 9e810ff3c6cce5ebaaff7c744c687a47.
- `codex/gpt-6-astra/medium`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding E1: nearest measured same-model effort max, AA coding 61.645045 + overall delta -3; stderr=5+2×3=11.
- `codex/gpt-6-astra/xhigh`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding E1: nearest measured same-model effort max, AA coding 61.645045 + overall delta -1; stderr=5+2×1=7.
- `codex/gpt-6-luna/high`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding E1: nearest measured same-model effort max, AA coding 41.074389 + overall delta -5; stderr=5+2×2=9.
- `codex/gpt-6-luna/low`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding E1: nearest measured same-model effort max, AA coding 41.074389 + overall delta -16; stderr=5+2×4=13.
- `codex/gpt-6-luna/max`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding: AA exact model/effort in the named harness; raw row ac8276cd5d15ac2daacb9956ad44c132.
- `codex/gpt-6-luna/medium`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding E1: nearest measured same-model effort max, AA coding 41.074389 + overall delta -8; stderr=5+2×3=11.
- `codex/gpt-6-luna/xhigh`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding E1: nearest measured same-model effort max, AA coding 41.074389 + overall delta -3; stderr=5+2×1=7.
- `codex/gpt-6-sol/high`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding E1: nearest measured same-model effort max, AA coding 56.662591 + overall delta -6; stderr=5+2×2=9.
- `codex/gpt-6-sol/low`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding E1: nearest measured same-model effort max, AA coding 56.662591 + overall delta -14; stderr=5+2×4=13.
- `codex/gpt-6-sol/max`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding: AA exact model/effort in the named harness; raw row 7a0dfd2a7a9d93056f860d4675586d66.
- `codex/gpt-6-sol/medium`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding E1: nearest measured same-model effort max, AA coding 56.662591 + overall delta -8; stderr=5+2×3=11.
- `codex/gpt-6-sol/xhigh`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding E1: nearest measured same-model effort max, AA coding 56.662591 + overall delta -4; stderr=5+2×1=7.
- `codex/gpt-6.1-sol/high`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding: AA exact model/effort in the named harness; raw row 4aee8aa2d4a0bad6068644f40edae235.
- `codex/gpt-6.1-sol/low`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding: AA exact model/effort in the named harness; raw row fa589dad83ffcdd6cce72a467e552063.
- `codex/gpt-6.1-sol/max`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding: AA exact model/effort in the named harness; raw row ea296ce785422d3b97cfff4427f2e0d2.
- `codex/gpt-6.1-sol/medium`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding: AA exact model/effort in the named harness; raw row 93acd59f31f8e5478686d650aa650f79.
- `codex/gpt-6.1-sol/xhigh`: Overall: AA model index at exact effort; API/reference-scaffold evidence, not a production-task success rate. Coding: AA exact model/effort in the named harness; raw row ed8c1d0b0a2fe4560bbfa0a9f47e8941.
- `muse/muse-spark-1.3-contributor/max`: M1: copy public Muse Spark 1.3 max observations; contributor-to-public execution binding remains unverified, so both are estimates (stderr 5 overall, 6 coding).

**Bug Hunt receipts.** [Pinned runs.csv](https://raw.githubusercontent.com/phuryn/bug-hunt-bench/1217192a6d04e89da3f6106ca3a304d2734882eb/results/runs.csv) at commit `1217192a6d04e89da3f6106ca3a304d2734882eb`, fetched 2026-10-04; live rows only, grouped by exact model, harness and effort. `mean/105×100` is a separate 0–100 bug-fix percentage. For n≥2, sample SE = sample SD/√n; for n=1 it is unknown. A separate pooled SE = 3.285553/√n uses all repeated live groups (81 residual degrees of freedom); it assumes homogeneous variance and is not a measured single-row SE. Both are in raw bugs; multiply by 100/105 for percentage SE. Repeated runs reuse the same two repositories/105 bugs, so these errors describe run variability, not new-repository generalization or 105 independent cases.

| Public configuration | Runs: scores (out of 105) | Mean | Percent | Sample SE | Pooled SE | Run ids | Last run date |
|---|---|---:|---:|---:|---:|---|---|
| Fable 5.1; Claude Code; high | 1: 33 | 33.000000 | 31.428571 | unknown | 3.285553 | fable51 | 2026-09-01 |
| Fable 5.1; Claude Code; low | 1: 29 | 29.000000 | 27.619048 | unknown | 3.285553 | fable51low | 2026-09-02 |
| Fable 5.1; Claude Code; max | 1: 43 | 43.000000 | 40.952381 | unknown | 3.285553 | fable51max | 2026-09-01 |
| Fable 5.1; Claude Code; medium | 1: 21 | 21.000000 | 20.000000 | unknown | 3.285553 | fable51medium | 2026-09-10 |
| Fable 5.1; Claude Code; xhigh | 1: 29 | 29.000000 | 27.619048 | unknown | 3.285553 | fable51xhigh | 2026-09-10 |
| GPT-6 Astra; Codex CLI; high | 1: 35 | 35.000000 | 33.333333 | unknown | 3.285553 | gpt6astrahigh | 2026-09-04 |
| GPT-6 Astra; Codex CLI; low | 1: 27 | 27.000000 | 25.714286 | unknown | 3.285553 | gpt6astralow | 2026-09-05 |
| GPT-6 Astra; Codex CLI; max | 3: 48, 40, 47 | 45.000000 | 42.857143 | 2.516611 | 1.896915 | gpt6astramax, gpt6astramaxr2, gpt6astramaxr3 | 2026-09-14 |
| GPT-6 Astra; Codex CLI; medium | 1: 34 | 34.000000 | 32.380952 | unknown | 3.285553 | gpt6astramedium | 2026-09-05 |
| GPT-6 Astra; Codex CLI; xhigh | 1: 43 | 43.000000 | 40.952381 | unknown | 3.285553 | gpt6astraxhigh | 2026-09-04 |
| GPT-6 Luna; Codex CLI; high | 1: 9 | 9.000000 | 8.571429 | unknown | 3.285553 | gpt6lunahigh | 2026-09-23 |
| GPT-6 Luna; Codex CLI; low | 1: 4 | 4.000000 | 3.809524 | unknown | 3.285553 | gpt6lunalow | 2026-09-23 |
| GPT-6 Luna; Codex CLI; max | 3: 17, 21, 17 | 18.333333 | 17.460317 | 1.333333 | 1.896915 | gpt6lunamax, gpt6lunamaxr2, gpt6lunamaxr3 | 2026-09-23 |
| GPT-6 Luna; Codex CLI; medium | 1: 4 | 4.000000 | 3.809524 | unknown | 3.285553 | gpt6lunamedium | 2026-09-23 |
| GPT-6 Luna; Codex CLI; xhigh | 1: 14 | 14.000000 | 13.333333 | unknown | 3.285553 | gpt6lunaxhigh | 2026-09-23 |
| GPT-6 Sol; Codex CLI; high | 1: 20 | 20.000000 | 19.047619 | unknown | 3.285553 | gpt6solhigh | 2026-09-23 |
| GPT-6 Sol; Codex CLI; low | 1: 6 | 6.000000 | 5.714286 | unknown | 3.285553 | gpt6sollow | 2026-09-23 |
| GPT-6 Sol; Codex CLI; max | 3: 32, 30, 26 | 29.333333 | 27.936508 | 1.763834 | 1.896915 | gpt6solmax, gpt6solmaxr2, gpt6solmaxr3 | 2026-09-23 |
| GPT-6 Sol; Codex CLI; medium | 1: 14 | 14.000000 | 13.333333 | unknown | 3.285553 | gpt6solmedium | 2026-09-23 |
| GPT-6 Sol; Codex CLI; xhigh | 1: 25 | 25.000000 | 23.809524 | unknown | 3.285553 | gpt6solxhigh | 2026-09-23 |
| GPT-6.1 Sol; Codex CLI; high | 2: 38, 35 | 36.500000 | 34.761905 | 1.500000 | 2.323237 | gpt61solhigh, gpt61solhighr2 | 2026-09-30 |
| GPT-6.1 Sol; Codex CLI; low | 2: 31, 14 | 22.500000 | 21.428571 | 8.500000 | 2.323237 | gpt61sollow, gpt61sollowr2 | 2026-10-01 |
| GPT-6.1 Sol; Codex CLI; max | 3: 44, 42, 47 | 44.333333 | 42.222222 | 1.452966 | 1.896915 | gpt61solmax, gpt61solmaxr2, gpt61solmaxr3 | 2026-09-30 |
| GPT-6.1 Sol; Codex CLI; medium | 2: 28, 30 | 29.000000 | 27.619048 | 1.000000 | 2.323237 | gpt61solmedium, gpt61solmediumr2 | 2026-09-30 |
| GPT-6.1 Sol; Codex CLI; xhigh | 3: 41, 43, 44 | 42.666667 | 40.634921 | 0.881917 | 1.896915 | gpt61solxhigh, gpt61solxhighr2, gpt61solxhighr3 | 2026-09-30 |
| Gemini 3.7 Flash; Antigravity CLI; high | 1: 16 | 16.000000 | 15.238095 | unknown | 3.285553 | gemini37flashseq | 2026-08-28 |
| Gemini 3.7 Flash; Gemini CLI (retired) + model-pinning gateway; high | 1: 22 | 22.000000 | 20.952381 | unknown | 3.285553 | gemini37flash | 2026-08-24 |
| Gemini 3.8 Flash; Antigravity CLI; high | 3: 20, 19, 15 | 18.000000 | 17.142857 | 1.527525 | 1.896915 | gemini38flash_agy, gemini38flash_agyr2, gemini38flash_agyr3 | 2026-09-15 |
| Muse Spark 1.3; Muse Code / Meta API; high | 3: 19, 17, 20 | 18.666667 | 17.777778 | 0.881917 | 1.896915 | musespark13high, musespark13highr2, musespark13highr3 | 2026-09-14 |
| Muse Spark 1.3; Muse Code / Meta API; low | 3: 12, 13, 4 | 9.666667 | 9.206349 | 2.848001 | 1.896915 | musespark13low, musespark13lowr2, musespark13lowr3 | 2026-09-14 |
| Muse Spark 1.3; Muse Code / Meta API; max | 5: 33, 33, 35, 29, 31 | 32.200000 | 30.666667 | 1.019804 | 1.469344 | musespark13max, musespark13maxr2, musespark13maxr3, musespark13maxr4, musespark13maxr5 | 2026-09-17 |
| Muse Spark 1.3; Muse Code / Meta API; medium | 3: 9, 15, 15 | 13.000000 | 12.380952 | 2.000000 | 1.896915 | musespark13medium, musespark13mediumr2, musespark13mediumr3 | 2026-09-14 |
| Muse Spark 1.3; Muse Code / Meta API; xhigh | 3: 20, 25, 16 | 20.333333 | 19.365079 | 2.603417 | 1.896915 | musespark13xhigh, musespark13xhighr2, musespark13xhighr3 | 2026-09-14 |
| Opus 5.5; Claude Code; high | 3: 32, 32, 31 | 31.666667 | 30.158730 | 0.333333 | 1.896915 | opus55high, opus55highr2, opus55highr3 | 2026-09-23 |
| Opus 5.5; Claude Code; low | 3: 23, 25, 19 | 22.333333 | 21.269841 | 1.763834 | 1.896915 | opus55low, opus55lowr2, opus55lowr3 | 2026-09-23 |
| Opus 5.5; Claude Code; max | 3: 43, 39, 43 | 41.666667 | 39.682540 | 1.333333 | 1.896915 | opus55max, opus55maxr2, opus55maxr3 | 2026-09-23 |
| Opus 5.5; Claude Code; medium | 3: 30, 30, 31 | 30.333333 | 28.888889 | 0.333333 | 1.896915 | opus55medium, opus55mediumr2, opus55mediumr3 | 2026-09-23 |
| Opus 5.5; Claude Code; xhigh | 3: 35, 37, 36 | 36.000000 | 34.285714 | 0.577350 | 1.896915 | opus55xhigh, opus55xhighr2, opus55xhighr3 | 2026-09-23 |
| Qwen3.8-27B 8-bit; Claude Code / OpenRouter; default | 3: 8, 12, 12 | 10.666667 | 10.158730 | 1.333333 | 1.896915 | qwen3827b8bitnc, qwen3827b8bitncr2, qwen3827b8bitncr3 | 2026-09-16 |
| Sonnet 5.5; Claude Code; high | 3: 32, 34, 37 | 34.333333 | 32.698413 | 1.452966 | 1.896915 | sonnet55high, sonnet55highr2, sonnet55highr3 | 2026-09-29 |
| Sonnet 5.5; Claude Code; low | 3: 20, 17, 20 | 19.000000 | 18.095238 | 1.000000 | 1.896915 | sonnet55low, sonnet55lowr2, sonnet55lowr3 | 2026-09-30 |
| Sonnet 5.5; Claude Code; max | 3: 57, 54, 43 | 51.333333 | 48.888889 | 4.255715 | 1.896915 | sonnet55max, sonnet55maxr2, sonnet55maxr3 | 2026-09-29 |
| Sonnet 5.5; Claude Code; medium | 3: 18, 28, 24 | 23.333333 | 22.222222 | 2.905933 | 1.896915 | sonnet55medium, sonnet55mediumr2, sonnet55mediumr3 | 2026-09-30 |
| Sonnet 5.5; Claude Code; xhigh | 3: 39, 34, 35 | 36.000000 | 34.285714 | 1.527525 | 1.896915 | sonnet55xhigh, sonnet55xhighr2, sonnet55xhighr3 | 2026-09-29 |

Do not pool the two Gemini 3.7 harnesses: the agy/Antigravity CLI score is 16, while 22 belongs to a retired Gemini CLI plus pinning gateway. Hosted Qwen 8-bit uses Claude Code/OpenRouter, not the Pi/MLX local engine. Bug Hunt is judged bug finding and fixing in two TypeScript repositories; counts are not deterministic unit-test pass rates. Superseded best runs are not added to live means.

**SWE-bench and vendor reports.** Inspected the [official SWE-bench board](https://www.swebench.com/) and [Haiku release methodology](https://www.anthropic.com/news/claude-haiku-4-5). Vendor observation: Haiku 4.5 SWE-bench Verified 73.3/100 (`kind=vendor`, source=the Haiku URL, as_of=2026-10-04, published=2025-10-15, stderr unknown). It uses 500 problems, 50 trials, a 128K thinking budget and a bash/edit scaffold; it does not establish this catalog’s no-effort-control Claude Code row and is excluded from primary coding and tiers. No exact supported-effort SWE-bench Verified result was established for the newer GPT-6/Claude-5.5/Muse/agy configurations.

The [Sonnet 5.5 vendor report](https://www.anthropic.com/claude-sonnet-5-5) reports Terminal-Bench 4.0 70.6% and Opus 5.5 66.4% (`kind=vendor`, as_of=2026-10-04, stderr unknown); those are neither SWE-bench Verified nor AA coding-index observations. They are not substituted into primary quality. Differences from AA are expected to depend on scaffold and protocol; no averaging or cherry-picking the higher number.

**Source index and freshness.** All primary source URLs in JSON were inspected through web search/open or direct public HTML. Overall release tables: [gpt-6.1-sol](https://artificialanalysis.ai/models/releases/gpt-6-1-sol), [gpt-6-astra](https://artificialanalysis.ai/models/releases/gpt-6-astra), [gpt-6-sol](https://artificialanalysis.ai/models/releases/gpt-6-sol), [gpt-6-luna](https://artificialanalysis.ai/models/releases/gpt-6-luna), [claude-fable-5-1](https://artificialanalysis.ai/models/releases/claude-fable-5-1), [claude-opus-5-5](https://artificialanalysis.ai/models/releases/claude-opus-5-5), [claude-sonnet-5-5](https://artificialanalysis.ai/models/releases/claude-sonnet-5-5). Other overall sources: [Gemini 3.8](https://artificialanalysis.ai/models/releases/gemini-3-8-flash), [3.7](https://artificialanalysis.ai/models/releases/gemini-3-7-flash), [3.6](https://artificialanalysis.ai/models/gemini-3-6-flash), [3.5](https://artificialanalysis.ai/models/releases/gemini-3-5-flash), [Muse 1.3](https://artificialanalysis.ai/models/releases/muse-spark-1-3), [Haiku non-reasoning](https://artificialanalysis.ai/models/claude-4-5-haiku). AA overall API indices are not direct measurements of Claude Code or agy task execution.

The AA coding page’s plain-text rendering omits plotted rows. The ledger below was extracted from its publicly served `self.__next_f.push` JSON payload, requiring no login. The retained records include the source labels and served slugs, not an inference from family names. Some slugs are pre-release aliases and versions predate the current model’s public release; no unsupported claim of exact current deployment parity is made. Anthropic coding results include fallback attempts; do not interpret them as isolated base-model scores. Recent raw data supersedes rounded article headlines; for example GPT-6.1 Sol xhigh scores higher than max and that reversal is preserved.

Source digests: Bug Hunt CSV SHA-256 `995ee322bab1019790fb58cc6c25af968e688f2952d8e28380e975f222407a74`; AA fetched HTML SHA-256 `25ace6d02df1cef8ff51e8c3836f13671a206e68df00be6a07b3ce5eb7333d85`; extracted ledger canonical SHA-256 `9952d8fee5498a09f1115f24b44cc5860dded05049c2ecf572a6424bcbfa506e`. Registry digest `5dfc2f74b9234a8b60c3e1d7e5c6ae9d9b0d1bb4b1b2b30db118ed90a6a25d9f` concatenates raw bench.go, openai.go, anthropic.go, google.go, muse.go in that order at v0.5.40. Mutable AA HTML itself is not retained; the minimal numerical/harness ledger is embedded here for replay, not represented as an immutable upstream release.

Recompute the ledger digest by parsing the JSON block below, then UTF-8 encoding `json.dumps(ledger, sort_keys=True, separators=(",", ":"))`. Primary coding values are `100*index_score`, rounded to six decimals. Source `mean` budgets support the cost table; source `mean.costUsd` is supplemental and not the catalog’s recomputed price.

```json
[
  {
    "id": "29a1cdad0d140781d05b9b4746cfdd0d",
    "label": "Claude Code - Fable 5.1 (max) (with fallback)",
    "index_score": 0.6222249615769457,
    "harness_versions": {
      "deep-swe-v1.1": {
        "min": {
          "version": "2.1.263",
          "dateReleased": "2026-09-06"
        },
        "max": {
          "version": "2.1.263",
          "dateReleased": "2026-09-06"
        }
      },
      "swe-atlas-qna": {
        "min": {
          "version": "2.1.259",
          "dateReleased": "2026-09-02"
        },
        "max": {
          "version": "2.1.259",
          "dateReleased": "2026-09-02"
        }
      },
      "terminal-bench-v4": {
        "min": {
          "version": "2.1.263",
          "dateReleased": "2026-09-06"
        },
        "max": {
          "version": "2.1.263",
          "dateReleased": "2026-09-06"
        }
      }
    },
    "host_model_slug": "anthropic_claude-fable-5-1",
    "mean": {
      "reward": 0.6222249615769457,
      "costUsd": 12.388092375137534,
      "agentWallTimeSec": 2089.900125412543,
      "steps": 36.65676567656766,
      "inputTokens": 5581622.259625964,
      "outputTokens": 133631.59955995603,
      "cacheTokens": 5205555.651265128,
      "cacheHitRate": 0.9159555838702423,
      "totalTokens": 5715253.85918592
    },
    "fallback_attempts": 74,
    "attempts": 909,
    "components": [
      {
        "metric": "deep-swe-v1.1",
        "reward": 0.64306784660767
      },
      {
        "metric": "swe-atlas-qna",
        "reward": 0.647849462365591
      },
      {
        "metric": "terminal-bench-v4",
        "reward": 0.575757575757576
      }
    ]
  },
  {
    "id": "33e80f383dc644527161435b076798ad",
    "label": "Claude Code - Sonnet 5.5 (high)",
    "index_score": 0.5501249049636148,
    "harness_versions": {
      "deep-swe-v1.1": {
        "min": {
          "version": "2.1.280",
          "dateReleased": null
        },
        "max": {
          "version": "2.1.280",
          "dateReleased": null
        }
      },
      "swe-atlas-qna": {
        "min": {
          "version": "2.1.280",
          "dateReleased": null
        },
        "max": {
          "version": "2.1.280",
          "dateReleased": null
        }
      },
      "terminal-bench-v4": {
        "min": {
          "version": "2.1.280",
          "dateReleased": null
        },
        "max": {
          "version": "2.1.280",
          "dateReleased": null
        }
      }
    },
    "host_model_slug": "private-anthropic_claude-beignet-eap-high",
    "mean": {
      "reward": 0.5501249049636148,
      "costUsd": 1.235497127172719,
      "agentWallTimeSec": 735.1809878987912,
      "steps": 36.772277227722775,
      "inputTokens": 2610884.190319032,
      "outputTokens": 41496.37513751376,
      "cacheTokens": 2507283.205720573,
      "cacheHitRate": 0.9477934370079523,
      "totalTokens": 2652380.565456549
    },
    "fallback_attempts": 33,
    "attempts": 909,
    "components": [
      {
        "metric": "deep-swe-v1.1",
        "reward": 0.666666666666667
      },
      {
        "metric": "swe-atlas-qna",
        "reward": 0.564516129032258
      },
      {
        "metric": "terminal-bench-v4",
        "reward": 0.419191919191919
      }
    ]
  },
  {
    "id": "4aee8aa2d4a0bad6068644f40edae235",
    "label": "Codex - GPT-6.1 Sol (high) ({'reasoning_effort': 'high'})",
    "index_score": 0.6014923716179784,
    "harness_versions": {
      "deep-swe-v1.1": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      },
      "swe-atlas-qna": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      },
      "terminal-bench-v4": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      }
    },
    "host_model_slug": "openai_quasar-alpha",
    "mean": {
      "reward": 0.6014923716179784,
      "costUsd": 0.8892413135295484,
      "agentWallTimeSec": 799.7652123212325,
      "steps": 38.97909790979098,
      "inputTokens": 2791163.2548648305,
      "outputTokens": 28748.735486663423,
      "cacheTokens": 2665415.3918375443,
      "cacheHitRate": 0.9278388366571144,
      "totalTokens": 2819911.9903514963
    },
    "fallback_attempts": 0,
    "attempts": 909,
    "components": [
      {
        "metric": "deep-swe-v1.1",
        "reward": 0.705014749262537
      },
      {
        "metric": "swe-atlas-qna",
        "reward": 0.599462365591398
      },
      {
        "metric": "terminal-bench-v4",
        "reward": 0.5
      }
    ]
  },
  {
    "id": "4d6aa9e9c8df931d54dba5956d94973c",
    "label": "Claude Code - Sonnet 5.5 (medium)",
    "index_score": 0.4588397535157457,
    "harness_versions": {
      "deep-swe-v1.1": {
        "min": {
          "version": "2.1.280",
          "dateReleased": null
        },
        "max": {
          "version": "2.1.280",
          "dateReleased": null
        }
      },
      "swe-atlas-qna": {
        "min": {
          "version": "2.1.280",
          "dateReleased": null
        },
        "max": {
          "version": "2.1.280",
          "dateReleased": null
        }
      },
      "terminal-bench-v4": {
        "min": {
          "version": "2.1.280",
          "dateReleased": null
        },
        "max": {
          "version": "2.1.280",
          "dateReleased": null
        }
      }
    },
    "host_model_slug": "private-anthropic_claude-beignet-eap-medium",
    "mean": {
      "reward": 0.4588397535157457,
      "costUsd": 0.6187507753025301,
      "agentWallTimeSec": 510.7721826182621,
      "steps": 22.20902090209021,
      "inputTokens": 1232830.7183718372,
      "outputTokens": 20520.774477447743,
      "cacheTokens": 1173805.1144114407,
      "cacheHitRate": 0.9341009883549257,
      "totalTokens": 1253351.4928492843
    },
    "fallback_attempts": 19,
    "attempts": 909,
    "components": [
      {
        "metric": "deep-swe-v1.1",
        "reward": 0.654867256637168
      },
      {
        "metric": "swe-atlas-qna",
        "reward": 0.448924731182796
      },
      {
        "metric": "terminal-bench-v4",
        "reward": 0.272727272727273
      }
    ]
  },
  {
    "id": "666605c66ad569c503ce53c01682e285",
    "label": "Claude Code - Sonnet 5.5 (xhigh)",
    "index_score": 0.628713868151493,
    "harness_versions": {
      "deep-swe-v1.1": {
        "min": {
          "version": "2.1.280",
          "dateReleased": null
        },
        "max": {
          "version": "2.1.280",
          "dateReleased": null
        }
      },
      "swe-atlas-qna": {
        "min": {
          "version": "2.1.280",
          "dateReleased": null
        },
        "max": {
          "version": "2.1.280",
          "dateReleased": null
        }
      },
      "terminal-bench-v4": {
        "min": {
          "version": "2.1.280",
          "dateReleased": null
        },
        "max": {
          "version": "2.1.280",
          "dateReleased": null
        }
      }
    },
    "host_model_slug": "private-anthropic_claude-beignet-eap-xhigh",
    "mean": {
      "reward": 0.628713868151493,
      "costUsd": 3.3333865348734886,
      "agentWallTimeSec": 1619.0134708470846,
      "steps": 78.31683168316832,
      "inputTokens": 6961982.357535753,
      "outputTokens": 125696.03850385037,
      "cacheTokens": 6693763.679867986,
      "cacheHitRate": 0.9574714756466342,
      "totalTokens": 7087678.3960396005
    },
    "fallback_attempts": 34,
    "attempts": 909,
    "components": [
      {
        "metric": "deep-swe-v1.1",
        "reward": 0.684365781710914
      },
      {
        "metric": "swe-atlas-qna",
        "reward": 0.620967741935484
      },
      {
        "metric": "terminal-bench-v4",
        "reward": 0.580808080808081
      }
    ]
  },
  {
    "id": "7a0dfd2a7a9d93056f860d4675586d66",
    "label": "Codex - GPT-6 Sol (max)",
    "index_score": 0.5666259127577997,
    "harness_versions": {
      "deep-swe-v1.1": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      },
      "swe-atlas-qna": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      },
      "terminal-bench-v4": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      }
    },
    "host_model_slug": "openai_gpt-6-sol",
    "mean": {
      "reward": 0.5666259127577997,
      "costUsd": 2.989729004180416,
      "agentWallTimeSec": 1338.0527051705185,
      "steps": 86.55005500550054,
      "inputTokens": 9778602.140814083,
      "outputTokens": 62175.62266226623,
      "cacheTokens": 9598928.887788778,
      "cacheHitRate": 0.9640119921972632,
      "totalTokens": 9840777.763476355
    },
    "fallback_attempts": 0,
    "attempts": 907,
    "components": [
      {
        "metric": "deep-swe-v1.1",
        "reward": 0.690265486725664
      },
      {
        "metric": "swe-atlas-qna",
        "reward": 0.575268817204301
      },
      {
        "metric": "terminal-bench-v4",
        "reward": 0.434343434343434
      }
    ]
  },
  {
    "id": "93acd59f31f8e5478686d650aa650f79",
    "label": "Codex - GPT-6.1 Sol (medium) ({'reasoning_effort': 'medium'})",
    "index_score": 0.6141474695571183,
    "harness_versions": {
      "deep-swe-v1.1": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      },
      "swe-atlas-qna": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      },
      "terminal-bench-v4": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      }
    },
    "host_model_slug": "openai_quasar-alpha",
    "mean": {
      "reward": 0.6141474695571183,
      "costUsd": 0.7046746644898908,
      "agentWallTimeSec": 651.3228668866898,
      "steps": 37.98899889988999,
      "inputTokens": 2275866.5144006205,
      "outputTokens": 22412.23102310231,
      "cacheTokens": 2176436.9277894986,
      "cacheHitRate": 0.9239660644436648,
      "totalTokens": 2298278.745423723
    },
    "fallback_attempts": 0,
    "attempts": 909,
    "components": [
      {
        "metric": "deep-swe-v1.1",
        "reward": 0.71976401179941
      },
      {
        "metric": "swe-atlas-qna",
        "reward": 0.60752688172043
      },
      {
        "metric": "terminal-bench-v4",
        "reward": 0.515151515151515
      }
    ]
  },
  {
    "id": "9e810ff3c6cce5ebaaff7c744c687a47",
    "label": "Codex - GPT-6 Astra (max)",
    "index_score": 0.6164504498789397,
    "harness_versions": {
      "deep-swe-v1.1": {
        "min": {
          "version": "0.153.4",
          "dateReleased": "2026-09-04"
        },
        "max": {
          "version": "0.153.4",
          "dateReleased": "2026-09-04"
        }
      },
      "swe-atlas-qna": {
        "min": {
          "version": "0.151.0",
          "dateReleased": "2026-08-29"
        },
        "max": {
          "version": "0.151.0",
          "dateReleased": "2026-08-29"
        }
      },
      "terminal-bench-v4": {
        "min": {
          "version": "0.153.4",
          "dateReleased": "2026-09-04"
        },
        "max": {
          "version": "0.153.4",
          "dateReleased": "2026-09-04"
        }
      }
    },
    "host_model_slug": "openai_vega-alpha",
    "mean": {
      "reward": 0.6164504498789397,
      "costUsd": 7.46709735138597,
      "agentWallTimeSec": 1762.0330264026397,
      "steps": 39.378437843784376,
      "inputTokens": 3299549.915435806,
      "outputTokens": 47061.09677361179,
      "cacheTokens": 3141588.9050659174,
      "cacheHitRate": 0.9335591564904552,
      "totalTokens": 3346611.012209417
    },
    "fallback_attempts": 0,
    "attempts": 909,
    "components": [
      {
        "metric": "deep-swe-v1.1",
        "reward": 0.675516224188791
      },
      {
        "metric": "swe-atlas-qna",
        "reward": 0.618279569892473
      },
      {
        "metric": "terminal-bench-v4",
        "reward": 0.555555555555555
      }
    ]
  },
  {
    "id": "a2c87c062f3cef73d8525e7578f14742",
    "label": "Claude Code - Sonnet 5.5 (max)",
    "index_score": 0.6835783373750831,
    "harness_versions": {
      "deep-swe-v1.1": {
        "min": {
          "version": "2.1.280",
          "dateReleased": null
        },
        "max": {
          "version": "2.1.280",
          "dateReleased": null
        }
      },
      "swe-atlas-qna": {
        "min": {
          "version": "2.1.280",
          "dateReleased": null
        },
        "max": {
          "version": "2.1.280",
          "dateReleased": null
        }
      },
      "terminal-bench-v4": {
        "min": {
          "version": "2.1.280",
          "dateReleased": null
        },
        "max": {
          "version": "2.1.280",
          "dateReleased": null
        }
      }
    },
    "host_model_slug": "private-anthropic_claude-beignet-eap-max",
    "mean": {
      "reward": 0.6835783373750831,
      "costUsd": 14.190800796369649,
      "agentWallTimeSec": 5244.965372937293,
      "steps": 265.93069306930687,
      "inputTokens": 27141368.959295932,
      "outputTokens": 600968.2744774478,
      "cacheTokens": 26006027.66886688,
      "cacheHitRate": 0.9506214943317096,
      "totalTokens": 27742337.233773403
    },
    "fallback_attempts": 45,
    "attempts": 909,
    "components": [
      {
        "metric": "deep-swe-v1.1",
        "reward": 0.71976401179941
      },
      {
        "metric": "swe-atlas-qna",
        "reward": 0.669354838709677
      },
      {
        "metric": "terminal-bench-v4",
        "reward": 0.661616161616162
      }
    ]
  },
  {
    "id": "ab5074e404be48bbd4ac3c488ad9a3e2",
    "label": "Claude Code - Opus 5.5 (max)",
    "index_score": 0.659885802549234,
    "harness_versions": {
      "deep-swe-v1.1": {
        "min": {
          "version": "2.1.280",
          "dateReleased": null
        },
        "max": {
          "version": "2.1.280",
          "dateReleased": null
        }
      },
      "swe-atlas-qna": {
        "min": {
          "version": "2.1.280",
          "dateReleased": null
        },
        "max": {
          "version": "2.1.280",
          "dateReleased": null
        }
      },
      "terminal-bench-v4": {
        "min": {
          "version": "2.1.280",
          "dateReleased": null
        },
        "max": {
          "version": "2.1.280",
          "dateReleased": null
        }
      }
    },
    "host_model_slug": "anthropic_claude-opus-5-5",
    "mean": {
      "reward": 0.659885802549234,
      "costUsd": 13.036383418426828,
      "agentWallTimeSec": 3867.0923410341024,
      "steps": 155.0957095709571,
      "inputTokens": 15218206.418041805,
      "outputTokens": 333176.801980198,
      "cacheTokens": 14568923.550055025,
      "cacheHitRate": 0.9500106599458848,
      "totalTokens": 15551383.220022004
    },
    "fallback_attempts": 78,
    "attempts": 909,
    "components": [
      {
        "metric": "deep-swe-v1.1",
        "reward": 0.684365781710915
      },
      {
        "metric": "swe-atlas-qna",
        "reward": 0.663978494623656
      },
      {
        "metric": "terminal-bench-v4",
        "reward": 0.631313131313131
      }
    ]
  },
  {
    "id": "ac8276cd5d15ac2daacb9956ad44c132",
    "label": "Codex - GPT-6 Luna (max) ({'reasoning_effort': 'max'})",
    "index_score": 0.4107438934016153,
    "harness_versions": {
      "deep-swe-v1.1": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      },
      "swe-atlas-qna": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      },
      "terminal-bench-v4": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      }
    },
    "host_model_slug": "openai_gpt-6-luna",
    "mean": {
      "reward": 0.4107438934016153,
      "costUsd": 0.17591172304455457,
      "agentWallTimeSec": 1285.9268514851478,
      "steps": 88.64301430143014,
      "inputTokens": 10119727.864136413,
      "outputTokens": 98341.38393839383,
      "cacheTokens": 9897312.69911991,
      "cacheHitRate": 0.9606438143896164,
      "totalTokens": 10218069.248074815
    },
    "fallback_attempts": 0,
    "attempts": 908,
    "components": [
      {
        "metric": "deep-swe-v1.1",
        "reward": 0.63716814159292
      },
      {
        "metric": "swe-atlas-qna",
        "reward": 0.443548387096774
      },
      {
        "metric": "terminal-bench-v4",
        "reward": 0.151515151515152
      }
    ]
  },
  {
    "id": "d1003684d8f29b95a3b14d2a7722221a",
    "label": "Antigravity SDK - Gemini 3.8 Flash (high)",
    "index_score": 0.4186315529449987,
    "harness_versions": {
      "deep-swe-v1.1": {
        "min": {
          "version": "0.1.16",
          "dateReleased": "2026-09-02"
        },
        "max": {
          "version": "0.1.16",
          "dateReleased": "2026-09-02"
        }
      },
      "swe-atlas-qna": {
        "min": {
          "version": "0.1.12",
          "dateReleased": "2026-08-13"
        },
        "max": {
          "version": "0.1.12",
          "dateReleased": "2026-08-13"
        }
      },
      "terminal-bench-v4": {
        "min": {
          "version": "0.1.16",
          "dateReleased": "2026-09-02"
        },
        "max": {
          "version": "0.1.16",
          "dateReleased": "2026-09-02"
        }
      }
    },
    "host_model_slug": "google_skimaki_ai-studio",
    "mean": {
      "reward": 0.4186315529449987,
      "costUsd": 2.4673922004125437,
      "agentWallTimeSec": 702.7110253025301,
      "steps": 188.4884488448845,
      "inputTokens": 13395246.748624865,
      "outputTokens": 115826.40649064905,
      "cacheTokens": 11871691.682068229,
      "cacheHitRate": 0.8697428234581179,
      "totalTokens": 13511073.155115515
    },
    "fallback_attempts": 0,
    "attempts": 909,
    "components": [
      {
        "metric": "deep-swe-v1.1",
        "reward": 0.657817109144543
      },
      {
        "metric": "swe-atlas-qna",
        "reward": 0.451612903225807
      },
      {
        "metric": "terminal-bench-v4",
        "reward": 0.146464646464646
      }
    ]
  },
  {
    "id": "e9f9d3ae5b03aa060114dd4b3275d67d",
    "label": "Claude Code - Sonnet 5.5 (low)",
    "index_score": 0.420593075103495,
    "harness_versions": {
      "deep-swe-v1.1": {
        "min": {
          "version": "2.1.280",
          "dateReleased": null
        },
        "max": {
          "version": "2.1.280",
          "dateReleased": null
        }
      },
      "swe-atlas-qna": {
        "min": {
          "version": "2.1.280",
          "dateReleased": null
        },
        "max": {
          "version": "2.1.280",
          "dateReleased": null
        }
      },
      "terminal-bench-v4": {
        "min": {
          "version": "2.1.280",
          "dateReleased": null
        },
        "max": {
          "version": "2.1.280",
          "dateReleased": null
        }
      }
    },
    "host_model_slug": "private-anthropic_claude-beignet-eap-low",
    "mean": {
      "reward": 0.420593075103495,
      "costUsd": 0.48311011474147336,
      "agentWallTimeSec": 380.85617931793183,
      "steps": 18.578657865786578,
      "inputTokens": 961481.2431243125,
      "outputTokens": 15803.526952695267,
      "cacheTokens": 917182.5137513754,
      "cacheHitRate": 0.9307814644010347,
      "totalTokens": 977284.7700770069
    },
    "fallback_attempts": 18,
    "attempts": 909,
    "components": [
      {
        "metric": "deep-swe-v1.1",
        "reward": 0.619469026548673
      },
      {
        "metric": "swe-atlas-qna",
        "reward": 0.389784946236559
      },
      {
        "metric": "terminal-bench-v4",
        "reward": 0.252525252525253
      }
    ]
  },
  {
    "id": "ea296ce785422d3b97cfff4427f2e0d2",
    "label": "Codex - GPT-6.1 Sol (max)",
    "index_score": 0.6014750704302517,
    "harness_versions": {
      "deep-swe-v1.1": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      },
      "swe-atlas-qna": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      },
      "terminal-bench-v4": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      }
    },
    "host_model_slug": "openai_quasar-alpha",
    "mean": {
      "reward": 0.6014750704302517,
      "costUsd": 1.554119528390054,
      "agentWallTimeSec": 1463.0464224422437,
      "steps": 40.45654565456545,
      "inputTokens": 3937843.0564663024,
      "outputTokens": 61658.55432920342,
      "cacheTokens": 3740552.1873072563,
      "cacheHitRate": 0.9281954641245814,
      "totalTokens": 3999501.6107955067
    },
    "fallback_attempts": 0,
    "attempts": 909,
    "components": [
      {
        "metric": "deep-swe-v1.1",
        "reward": 0.696165191740413
      },
      {
        "metric": "swe-atlas-qna",
        "reward": 0.577956989247312
      },
      {
        "metric": "terminal-bench-v4",
        "reward": 0.53030303030303
      }
    ]
  },
  {
    "id": "ed8c1d0b0a2fe4560bbfa0a9f47e8941",
    "label": "Codex - GPT-6.1 Sol (xhigh) ({'reasoning_effort': 'xhigh'})",
    "index_score": 0.6290776736822983,
    "harness_versions": {
      "deep-swe-v1.1": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      },
      "swe-atlas-qna": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      },
      "terminal-bench-v4": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      }
    },
    "host_model_slug": "openai_quasar-alpha",
    "mean": {
      "reward": 0.6290776736822983,
      "costUsd": 1.0397081014409635,
      "agentWallTimeSec": 931.3823091309132,
      "steps": 40.05500550055005,
      "inputTokens": 3177125.49569875,
      "outputTokens": 35402.61375317859,
      "cacheTokens": 3036065.628757956,
      "cacheHitRate": 0.9327572280083256,
      "totalTokens": 3212528.1094519314
    },
    "fallback_attempts": 0,
    "attempts": 909,
    "components": [
      {
        "metric": "deep-swe-v1.1",
        "reward": 0.731563421828909
      },
      {
        "metric": "swe-atlas-qna",
        "reward": 0.610215053763441
      },
      {
        "metric": "terminal-bench-v4",
        "reward": 0.545454545454545
      }
    ]
  },
  {
    "id": "f2de07373c3626f1b9269a6ecad7804c",
    "label": "Muse Code - Muse Spark 1.3 (max)",
    "index_score": 0.5430273329930766,
    "harness_versions": {
      "deep-swe-v1.1": {
        "min": {
          "version": "1.0.3-R2198.1",
          "dateReleased": "2026-09-04"
        },
        "max": {
          "version": "1.0.3-R2198.1",
          "dateReleased": "2026-09-04"
        }
      },
      "swe-atlas-qna": {
        "min": {
          "version": "1.0.3-R2198.1",
          "dateReleased": "2026-09-04"
        },
        "max": {
          "version": "1.0.3-R2198.1",
          "dateReleased": "2026-09-04"
        }
      },
      "terminal-bench-v4": {
        "min": {
          "version": "1.0.3-R2198.1",
          "dateReleased": "2026-09-04"
        },
        "max": {
          "version": "1.0.3-R2198.1",
          "dateReleased": "2026-09-04"
        }
      }
    },
    "host_model_slug": "meta_aa_glacier135",
    "mean": {
      "reward": 0.5430273329930766,
      "costUsd": 3.977504096947204,
      "agentWallTimeSec": 1106.8533531353135,
      "steps": 185.23817381738175,
      "inputTokens": 16445971.068206819,
      "outputTokens": 263473.5055005501,
      "cacheTokens": 16090656.487898799,
      "cacheHitRate": 0.9712826195905685,
      "totalTokens": 16709444.57370737
    },
    "fallback_attempts": 0,
    "attempts": 909,
    "components": [
      {
        "metric": "deep-swe-v1.1",
        "reward": 0.716814159292036
      },
      {
        "metric": "swe-atlas-qna",
        "reward": 0.594086021505376
      },
      {
        "metric": "terminal-bench-v4",
        "reward": 0.318181818181818
      }
    ]
  },
  {
    "id": "fa589dad83ffcdd6cce72a467e552063",
    "label": "Codex - GPT-6.1 Sol (low) ({'reasoning_effort': 'low'})",
    "index_score": 0.5721634943016617,
    "harness_versions": {
      "deep-swe-v1.1": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      },
      "swe-atlas-qna": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      },
      "terminal-bench-v4": {
        "min": {
          "version": "0.154.0",
          "dateReleased": null
        },
        "max": {
          "version": "0.154.0",
          "dateReleased": null
        }
      }
    },
    "host_model_slug": "openai_quasar-alpha",
    "mean": {
      "reward": 0.5721634943016617,
      "costUsd": 0.49880364510270725,
      "agentWallTimeSec": 517.8138943894386,
      "steps": 35.74037403740374,
      "inputTokens": 1699123.8601417518,
      "outputTokens": 14695.05805334632,
      "cacheTokens": 1627295.6717524226,
      "cacheHitRate": 0.9135762473536158,
      "totalTokens": 1713818.9181950986
    },
    "fallback_attempts": 0,
    "attempts": 909,
    "components": [
      {
        "metric": "deep-swe-v1.1",
        "reward": 0.675516224188791
      },
      {
        "metric": "swe-atlas-qna",
        "reward": 0.551075268817204
      },
      {
        "metric": "terminal-bench-v4",
        "reward": 0.48989898989899
      }
    ]
  }
]
```

**Validation.** Checked exact model/effort coverage against the pinned declarations (local exception above), unique row keys, required strict-schema fields, kind/billing/latency enums, finite 0–100 quality, nonnegative costs and uncertainty, Muse max-only, date/source presence, all 13 task-class frontier tables, tier recomputation and the I/C/O pricing equation. No metered configurations. Primary measured AA uncertainty remains unknown; local quality remains absent. Validation is structural/arithmetic; no harness availability or benchmark reruns are claimed.
