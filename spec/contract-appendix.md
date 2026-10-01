# DIGEST contract (W0)

This appendix is normative for the frozen W0 contracts. The task-board boundary,
member-wise spawn validation, and delivery protocol portion belongs to W4.

1. Canonical JSON is UTF-8 without a BOM, insignificant whitespace, or a trailing
   newline. The top level MUST be an object with a non-empty string
   `schema_version`. Every nested `null` MUST be refused; omit optional unknown
   fields. Required unknown string/enum values use the literal `"unknown"`.
   Contract `Validate` refuses empty required strings with
   `contract_empty_field`, naming the JSON path and directing the caller to
   `"unknown"`. Required collections MUST be non-nil; `contract_missing_field`
   names the JSON path. Known-empty is `[]`. These checks recurse through all
   nested contract members, including map keys and values in `Facets` and
   `Features`: empty strings are `contract_empty_field`, and `"unknown"` is
   permitted for unknown map strings. Closed enums still require declared values;
   digest fields follow the shape and unknown exceptions in rule 6. Numeric zero
   is a value, never an unknown marker. Unknown `used_bp`, fitness,
   coverage, and optional timestamps are represented by omitted pointer fields.
2. Object keys MUST be sorted by unsigned UTF-16 code units, recursively, as in
   [RFC 8785](https://www.rfc-editor.org/rfc/rfc8785). Strings MUST use JCS escaping:
   quote/backslash escaped, the five short control escapes, other U+0000–U+001F
   as lowercase `\u00xx`, and all remaining Unicode verbatim. Do not escape slash,
   HTML characters, U+2028, or U+2029. Do not normalize Unicode. Invalid UTF-8,
   lone surrogates, duplicate object keys, malformed JSON, and trailing data MUST
   be refused. Nesting MUST NOT exceed 1000 simultaneously open object/array
   containers, counting the root object; excess depth is `canonical_invalid_json`.
3. Integer tokens represent signed or unsigned 64-bit integers and MUST serialize
   exactly, including above 2^53. Other tokens are IEEE-754 binary64: ECMAScript
   shortest round-trip digits, ties to even, decimal notation for absolute values
   in [10^-6, 10^21), scientific notation otherwise, lowercase `e`, explicit `+`
   for positive exponents, and no exponent leading zeros. Negative zero becomes
   `0`. NaN and either infinity MUST be refused, including overflow on parsing.
   An integral token beyond both 64-bit domains is interpreted as binary64.
   This is an explicit extension to JCS for exact contract integers.
4. Arrays preserve their declared semantic order. Collections declared as keyed
   MUST already be strictly ordered without duplicate keys; hashing MUST NOT
   silently repair them. The operator-facing `LoadPolicy` normalizes only the
   declared policy sets before validation, as rule 5 specifies. `SortByKey` is an explicit caller normalization helper;
   `CheckOrdered` refuses disorder and duplicates, even non-adjacent duplicates.
   Declared collection keys compare lexically by UTF-8 bytes: usage facts and
   in-flight counts by key, windows by id, candidates by id, scope entries by
   (runtime, scope), and scope model ids by id. Evaluation assessments are by id,
   estimates by candidate id, evaluator versions by name, and record addresses
   lexically. Policy masks, protected roles, window-id lists, and reason sets are
   lexical. Alternatives are by increasing base position. Task criteria/context/
   tools and declared group/pair lists retain explicit policy/task order.
5. Load policy defaults before hashing. Programmatic policies MUST start from
   `DefaultPolicy` and pass `Validate`; `LoadPolicy` materializes defaults while
   preserving explicit zero, false, and empty values. Every key MUST exactly
   match its declared JSON tag, case-sensitively, including nested group pairs;
   unknown or differently cased keys are `contract_unknown_field`.
   Strict decoding of snapshots/envelopes at the task-board boundary is W4's.
   `LoadPolicy` sorts and de-duplicates `free_field_mask`, `protected_roles`,
   and window-id lists. Programmatic `Validate` continues refusing unordered
   or duplicate sets. Other arrays retain their declared semantic order.
   `strategy` names ONLY the base: `config-order`, `quality-first`, or
   `cost-with-quality-floor`. `headroom.enabled` enables the headroom wrapper.
   `strategy: headroom-aware` MUST be refused as `routing_unknown_enum`.
   Defaults are schema version supplied by the caller, strategy `config-order`,
   mode `off`, empty free-field mask; headroom `enabled: false`, equivalence
   `same-pair-any-home`, groups `[]`, protected roles `[orchestrator, reviewer]`,
   reserve 2000 bp, band width 1000 bp, expiring band 3, conserve band -2,
   windows `all`, on-all-reserved `rank`.
   Non-empty groups with another equivalence are `headroom_groups_unused`;
   pair intersections are checked ONLY for `declared-groups`. Omitted runtime
   is a wildcard. Duplicate/intersecting pairs within one group are
   `headroom_group_duplicate_pair`; overlap BETWEEN groups is
   `headroom_groups_overlap`. Group and pair order is preserved.
6. A digest is `sha256:` followed by 64 lowercase hexadecimal digits of SHA-256
   over canonical bytes. `engine_profile.weight_digest` and
   `context_profile.lock_sha256` MUST be either that shape or literal `"unknown"`.
   `usage_facts[].record_digest` is optional in the wire shape: it MUST be omitted
   for state `absent` (a present value is `contract_field_forbidden`), and MUST
   be present and strictly shaped for every other state (`contract_missing_field`
   if omitted). It MUST NOT use `"unknown"` or a fabricated placeholder digest.
   `decision_id` and every other present decision/evaluation/snapshot/policy
   digest MUST retain the strict shape, with no `"unknown"` exception.
   Malformed digests are `contract_invalid_digest`.
   Use contract `Digest` methods to enforce presence and collection order.
   `RoutingDecision.ContentID` hashes the decision with `decision_id` omitted;
   `WithContentID` returns a copy carrying that identity. `ContentID` ignores an
   old id and can run before assignment; `VerifyContentID` MUST compare the
   stored id with the recomputed id and refuse missing, stale or substituted
   ids with `contract_decision_id_mismatch`.
7. `CandidateSnapshot.Validate` checks TTL-independent window facts.
   `fresh`/`stale` MUST have valid ranges, present `used_bp` and `observed_at`,
   observation no later than `as_of`, and absent reset or reset after `as_of`.
   `expired` MUST have valid ranges and observation plus reset at or before
   `as_of`. Contradictions are `contract_freshness_mismatch`. TTL is absent from
   the projection, so validation does not distinguish fresh from stale by age.
   `invalid` remains a permissible label for unknown or malformed facts.
   State `absent` MUST carry no windows: its required window list is `[]`;
   any window yields `contract_field_forbidden`. States `not_supported` and
   `unavailable` MAY retain windows for `explain`, subject to the label checks
   above, but NO window of either state is ever known, whatever its label
   (headroom.md §4). W1 MUST enforce this when deriving ordering quantities.
8. `selected` MUST carry `selected_candidate`; every other outcome MUST omit
   it. `baseline_after_abstain` MUST have outcome `selected`. Present expiry
   MUST be later than preparation; both use the supported timestamp range.
   Violations are `routing_invalid_decision`. Outcome tokens live in
   `OutcomeKind`; only spec-named reasons live in `ReasonCode` (including the
   explicitly named `no_eligible_candidates`). `bound_on_this_machine` is a
   catalog field, not a reason.
9. `FitnessEstimator` MUST be pure and deterministic over immutable in-memory
   evidence. It MUST refuse a reference whose `EvidenceSnapshotRef.Digest`
   differs from the evidence it holds with a typed `routing.Error`, code
   `estimator_snapshot_mismatch` (constructor
   `NewEstimatorSnapshotMismatchError`). No estimator implementation is part
   of W0. Fitness and coverage are 0..10000 inclusive; absent differs from zero.
10. Refusals are typed, with stable canonical codes: `canonical_null`,
    `canonical_nonfinite`, `canonical_missing_schema_version`,
    `canonical_unordered`, `canonical_duplicate_key`, `canonical_invalid_json`.
    Routing codes include the `contract_*`, `estimator_snapshot_mismatch`, and
    group codes specified above; `routing_unknown_enum`, `routing_invalid_policy`,
    `routing_invalid_as_of`, `routing_invalid_count`, `routing_missing_usage_key`,
    `routing_invalid_evidence_ref`, `routing_invalid_estimate`,
    `routing_invalid_selection`, `routing_invalid_decision`,
    `routing_invalid_percent`, `headroom_invalid_band_width`,
    `headroom_invalid_reserve`, and `headroom_invalid_bands` retain their meanings.
11. CLI refusals with `--json` MUST put the JSON error object on stdout and
    retain a nonzero exit code; plain refusals use stderr. Errors while writing
    output attempt the same channel rule and retain the refusal exit code.
    CLI codes are `cmr_unknown_subcommand`, `cmr_invalid_arguments`, and
    `cmr_output_failed`.

Vectors live in `pkg/canonical/testdata/vectors/`: `.input.json`, exact
`.canonical.json` bytes (no newline), `.digest`, or `.error` for a refusal.
`.ordered` declares `items` ordered by id for the collection-validation vectors.
`go test ./pkg/canonical -args -update` is the only golden rewrite mode; normal
runs MUST only compare. Expected bytes were authored independently of the encoder.
The three depth vectors below are generated in `TestDepthVectors` rather than
stored as large JSON files; both array and object nesting boundaries are checked.

| Vector | Digest or refusal |
| --- | --- |
| `array-top` | `canonical_missing_schema_version` |
| `duplicate-collection` | `canonical_duplicate_key` |
| `duplicate-object` | `canonical_duplicate_key` |
| `empty-schema` | `canonical_missing_schema_version` |
| `escapes` | `sha256:897cc57297e5765678b44dcf19cb8f01937bce2513c2c78bfe5e7f5362909d7f` |
| `float-boundaries` | `sha256:1c148f29706f1c4440bef011057bbea24dacc2457e244c0a64146638f28d7e24` |
| `floats` | `sha256:eeeedcfdb5051fafc8ed78707d4ecf2fa740f8e668e18a93c4a6956ec5ebb2f8` |
| `integers` | `sha256:0aa457ac99ba61dd40a007a71658af244866e50f6a8c0508e1390159a2b91fcb` |
| `invalid-utf8` | `canonical_invalid_json` |
| `invalid-json` | `canonical_invalid_json` |
| `keys` | `sha256:a87019fab9530d31e08504c96df4b3131ee722df5cd966654bf40e0df166be6a` |
| `lone-surrogate` | `canonical_invalid_json` |
| `missing-schema` | `canonical_missing_schema_version` |
| `multiple-values` | `canonical_invalid_json` |
| `negative-zero` | `sha256:6878ce440730e009d7baa2d5ee7c2aca368111e3b2d73b83933bbb93feb20f6a` |
| `nested` | `sha256:085e7e7b8e6e3b5ed4b300912bf610ad27bc64904fcf9e4db8419a12fce5e3c1` |
| `nonfinite` | `canonical_nonfinite` |
| `null` | `canonical_null` |
| `ordered` | `sha256:b3ae6b8681bb296c4db379608e92aa715d203d13a38543ee2ae792b9abd4be1b` |
| `surrogate-pair` | `sha256:218fd3c7f39f2de83f4cb32ada4d8266c887aad4da698ee46ed8b0b49096ca1f` |
| `unicode-keys` | `sha256:9a1ef3fd8204f252d62290c95d431ea8072ad41db6510b3b7425496eeabc7d87` |
| `unknown` | `sha256:b46f13e0660b01a272b240bc5e19c9b3b34ed1dc2b535644c0432a8120fa2f5e` |
| `unordered` | `canonical_unordered` |
| `wrong-schema` | `canonical_missing_schema_version` |
| `depth-999` (generated) | `sha256:8a3c9d6d78c6bb6b0f32aaa2b87e94244a6179ae703ec718be052759dd90c8e7` |
| `depth-1000` (generated) | `sha256:d314e865c86cb2ecfffec871138a28151ca9dd7de3a4aadc9792f5d828c2c7ec` |
| `depth-1001` (generated) | `canonical_invalid_json` |

Routing goldens live in `pkg/routing/testdata/`. The minimal-input default policy
pins both materialized canonical bytes and digest. The rich candidates and
non-default policy also pin canonical bytes; existing fixtures pin digests.
All newly added expected bytes/digests were independently authored from literal
JSON with sorted ASCII keys and SHA-256, without using the production defaults
or encoder. Normal runs only compare.

| Routing fixture | Digest (decision uses content id) |
| --- | --- |
| `envelope` | `sha256:e619e4332988da914296dce3089a38ef796ed56f3ad0a3e9b8c311d658eaf7df` |
| `candidates` | `sha256:22cd73316d314fdc6ba3f90ef725d5ed5a72193e03c9317768eff85b41305939` |
| `candidates-rich` | `sha256:ca0c49f546064904fdfda2aaa5f73606d5d6558aaacfc0d9a4b9253a57292c62` |
| `evaluation` | `sha256:8425c15c52a90e36df68b2357973e1edef45c68e990bec640b7a4850f9a6c68c` |
| `policy` | `sha256:4901016a8b88cd9979f8f3bfadba6f81a421a11c917d1ab738c28a0f587b53ef` |
| `policy-default` | `sha256:cc2682ca409195e06263f4a12627eda2014bee142f3ec34ac7b5d26dfd8b0117` |
| `policy-nondefault` | `sha256:c7c16d82cd085867dbb82c367b7437ca06356bc356c52c7b28a0ccaa212ed4c4` |
| `decision` | `sha256:9a859577ef79f81203fade534679ed5fd3856c1809a7e1febf3ee859ffc6bad4` |
