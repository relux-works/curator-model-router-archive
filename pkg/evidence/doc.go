// Package evidence implements append-only evidence documents and replayable
// suitability views. All times are explicit inputs; importers do no network I/O.
//
// TODO(decision): the CLI default is XDG_DATA_HOME/curator/model-router/evidence,
// with the local share directory as fallback. This keeps derived routing data
// separate from credentials and launch configuration. Store always requires an
// explicit root, and tests always pass a temporary directory.
//
// TODO(decision): require manual stale-lock recovery after confirming no writer.
// Writers use current.lock with O_EXCL and never steal a stale lock. Following a
// crash an operator must confirm that no writer is running, then remove the lock.
// Readers observe immutable files and an atomically replaced current pointer.
// Parent-directory fsync follows current's atomic rename, best effort where the
// platform or filesystem cannot sync a directory.
//
// Imports retain source resolution; status and effective model resolution are
// derived per snapshot using its bound registry. Rebind creates a new snapshot
// without rewriting imports or earlier snapshots. A source-declared unresolved
// observation remains unresolved. Notes created after a view's evaluated_at are
// excluded entirely, so a frozen view never uses future judgements.
package evidence
