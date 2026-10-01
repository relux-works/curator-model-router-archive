# Evidence v1

`Import`, `Snapshot`, `Benchmark`, `Observation`, `Note` and `Retraction` follow
`spec/evidence-format.md`. `kind` discriminates import and snapshot documents;
`schema_version` is `evidence-v1`. Required collections are present, including
known-empty arrays. Optional unknown fields are omitted, required unknown strings
use `unknown`, and zero remains a value. Native decoding rejects null, duplicate
keys, extra keys, differently cased keys and missing required fields.

`Prepare` explicitly sorts declared sets: benchmarks by `id@version`, metrics by
name, benchmark categories by category, observations/notes/retractions by address,
category/effort/facet/supersession lists lexically, and note evidence references by
`kind:ref`. Duplicates are refused. `Import.Digest` validates this ordering without
repairing it. Record addresses hash exactly the nested canonical object without
`id`; the temporary encoder envelope does not enter that hash. Snapshot import
sets are lexical. Candidate sets use the lexical canonical full fingerprint as
key. Suitability results preserve that order and never sort by a module rank.

Unknown categories and unresolved models stay in imports and are reported.
Unrecognised efforts stay verbatim and carry `effort_unrecognised` in the report.
Only the frozen model's effort vocabulary can establish an exact effort match.
`Prepare` and `Native` take a frozen registry projection explicitly. They never
register a model. Imports retain source-declared `subject_resolution`; effective
resolution is derived against each snapshot's bound registry. Explicit source
`unresolved` stays unresolved even if the registry knows the model. The
`subject_resolution` stored in an import file is therefore only what the source
declared; `cmr evidence show` and `ls` print the derived value. Unknown
models receive no `effort_unrecognised` issue because their vocabulary is unknown.

The active set is derived from all imports. Withdrawn records have no outgoing
supersession edges; suppression follows paths in the remaining graph. Both active
replacements stay active in a conflict. Retractions cannot target retractions;
unknown future targets are allowed so arrival order does not determine status.
Cycle validation runs before supplied address verification and again over the
complete snapshot during publication. Status and conflicts are inspection output,
never imported fields.

## Store and locking

`NewStore` requires an explicit root. The CLI suggestion is
`XDG_DATA_HOME/curator/model-router/evidence`, with the local share directory as
fallback (`TODO(decision)`). Every command accepts `--store DIR`, and every test
passes `t.TempDir()`.

Imports, snapshots and views are immutable. The additional
`registries/<digest>.json` directory retains the exact vocabulary projection bound
by a snapshot. Snapshot identity includes both the public registry reference and
its digest. Store readers verify content addresses. Views verify the six identity
inputs, embedded derivation/requirements/candidates and contribution arithmetic.

Writers exclusively create `current.lock` with `O_EXCL`. Existing locks return
`evidence_locked`. Locks are never stolen based on an ambient clock or PID.
Following a crash, confirm no writer is running and manually remove the lock.
Immutable files are written and synced to temporary files, then atomically linked
into place. `current` is synced in a temporary file and published by atomic rename.
After renaming `current`, the parent directory is fsynced on a best-effort basis;
platforms or filesystems without directory fsync still retain atomic publication.
Readers need no lock. A failed publication may leave unreferenced immutable files;
it cannot make `current` refer to partly written bytes.

## Bug Hunt interchange

The dependency is pinned to `skill-agents-management v0.5.32`. Its complete vendor
bench table lacks an accessor that can be used without importing executable
vendor plugins. `BugHuntExport` is the supported JSON fallback:

- `schema_version`, `module_version`, `mapping_version`, `imported_at`, `retrieved_at`;
- `models`: `{model_id, efforts[]}` (no-effort axes explicitly use `none`);
- `mapping`: `{external_name, model_id}`, ordered by external name;
- `rows`: `{key, model_name, effort?, kind, value, cost?, origin_ref}`, ordered by key.

`kind` is `measured`, `interpolated` or `cost`. A cost claim is a `list_cost`
observation with `evidence_kind: measured` and `cost.usd`. It is never used as a
quality measurement. The pinned module's cost annotation supplies the launched
max effort and USD list cost; the remaining annotation stays in the referenced
module artifact. An interpolation has unknown effort because the module does
not state one. The observation time is `unknown`; the leaderboard update date is
a benchmark version, never an invented measurement time.

The export in `pkg/evidence/data/` records generation time
`2026-10-01T00:00:00Z` as fixed `imported_at` and `retrieved_at` metadata.
The importer uses these values by default; `--imported-at` explicitly overrides
only the import time. Re-importing a pin retains the import digest and all
observation addresses; the store reports `already present` for duplicates.

The checked-in JSON was generated statically from the module's public data
literals, never transcribed. The parity test independently walks every vendor model
literal and the frozen runtime model declarations across every non-test Go file.
A separate AST walk counts every bench constructor call. The actual pin has 53
calls (11 measured, 41 interpolated, 1 cost), plus five declared alias expansions
(one measured, four interpolated), yielding 58 exported claims (12/45/1). The
cross-check verifies both alias declarations and corresponding exported claims.
Missing fields fail; only the explicitly allow-listed, dynamic local-models
template may have zero claims. It fails on missing rows,
unknown declaration syntax or a changed cost grammar. Production reads only JSON
and the exported bench total; it never imports a vendor implementation, registers
plugins, calls availability, executes a harness or fetches data. Regenerate after
review with `go test ./pkg/evidence -run TestAcceptance2 -args -update-evidence`.

The module should supply data-only `BugHuntRows()` and `RegistryFacts()` accessors;
see `W2-REPORT.md`. Until then the JSON export is the explicit pinned registry
projection for the benchmark scope. `--registry FILE` supplies another frozen
projection for native imports and local models. `Add` refuses a different
projection with `evidence_registry_mismatch` and points to
`cmr evidence rebind --registry FILE`. Rebind creates a new snapshot over the
same import digests with the supplied registry; earlier snapshots and imports
remain immutable and replay with their original resolution.

## Derivation and inspection

`DefaultMapping` declares the four playbook roles and fix, implementation, ops,
research, review and routine workloads. Projects supply fully mapped
`RequirementsInput` documents, including their mapping table and facets. V1
weights and thresholds are explicit `TODO(decision)` choices, not admission
rules. Only declared exact benchmark/version/metric normalisations contribute.
A different benchmark or version has no implicit compatibility. An active
observation with no normalisation for its benchmark/metric/category is reported
as `observation_unmapped`, with its observation address and unmapped triple.

V1 normalises Bug Hunt fixed counts to `[0,1]`, discounts partial observations by
0.5, weighs notes by confidence and basis, and averages mapped contributions
within a category. Covered categories are combined by declared positive weights;
a missing category contributes `none`, never a fabricated zero measurement.
Signed contributions sum to relative utility on `[-1,1]`. Strong starts at 0.4,
adequate at 0.2, and weaker supported scores are weak. No supported evidence means
unknown. This is not a calibrated probability.

Known conflicting fingerprint members reject a measurement. Unknown runtime,
harness, engine or engine member, tool or execution profile makes it partial.
Reprints with the same original publication, benchmark version and metric count
once per candidate/category; a direct matching reprint takes precedence over a
partial one, then lexical record address breaks ties. Notes have their own scope:
known efforts, runtime and harness range are direct; unspecified axes are partial.
Harness ranges accept exact versions or inclusive numeric dotted `[a,b]`; other
range syntax remains stored and does not match (`TODO(decision)`). Language,
platform, context-size and horizon requirements refine matching; a scoped role
must match the requested role.

A date-only `review_by` includes that whole UTC date (`TODO(decision)`); RFC3339
expiry is an exact instant. Expired notes remain visible with age and zero weight.
Notes whose `created_at` is later than `evaluated_at` are excluded entirely
from the view, including contributions and coverage; the boundary is inclusive.
This preserves knowledge as of the frozen evaluation time.
Evaluation time is a required library input. Only the CLI may default it to now,
and the resulting view always prints the exact `inputs.evaluated_at` it used.
The view id binds snapshot, derivation, role requirements/mapping, weights/policy,
evaluation time and sorted candidate fingerprints. All inputs needed for replay
are retained in the immutable view.

Example commands, using a store inside the working directory:

```sh
cmr evidence import pkg/evidence/data/bughunt-export.json --importer bughunt --imported-at 2026-10-01T00:00:00Z --store .build/evidence --json
cmr evidence ls --store .build/evidence --json
cmr evidence show RECORD_ADDRESS --store .build/evidence --json
cmr evidence unresolved --store .build/evidence --json
cmr evidence rebind --registry REGISTRY_FILE --store .build/evidence --json
cmr note add --model gpt-6.1-sol --efforts high --categories review.code --statement 'Useful review prior' --author operator --confidence high --review-by 2026-11-01 --at 2026-10-01T00:00:00Z --store .build/evidence --json
cmr note retract RECORD_ADDRESS --reason withdrawn --author operator --at 2026-10-02T00:00:00Z --store .build/evidence --json
cmr suitability --role reviewer --project-reqs REQUIREMENTS_FILE --candidates CANDIDATES_FILE --evaluated-at 2026-10-01T00:00:00Z --store .build/evidence --json
```

`--candidates` accepts a versioned `{schema_version, candidates[]}` document of
public fingerprints. Without it, inspection derives advisory fingerprints from
resolved exact-effort observations and note scopes. A note without an effort
restriction covers the frozen vocabulary; an empty store enumerates registry
model/effort pairs with unknown grades (`TODO(decision)`). These are
never admission decisions or executable offers. `note add --file` accepts a
versioned `{schema_version, note}` document; otherwise flags form a one-note import.
CLI flags may appear before or after positional arguments. JSON refusals go to
stdout with a nonzero exit code; plain refusals go to stderr.

`--store` always overrides the XDG default; the default may move once the Curator
data layout is fixed (`TODO(decision)`). Every command that defaults a time
prints the value in plain and JSON output: import results carry `imported_at`,
Bug Hunt results also carry the pinned `retrieved_at`; note results also carry `created_at` (add) or `at` (retract), and suitability
prints `inputs.evaluated_at`.
