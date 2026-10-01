package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/relux-works/curator-model-router/pkg/canonical"
	"strings"
)

// RecordAddress hashes exactly the canonical record without id. W0's encoder
// requires schema_version on the root, so canonicalise a versioned envelope and
// extract its already canonical record bytes before hashing (no envelope bytes
// enter the address). TODO(contract): add canonical.RecordDigest for schema-less
// nested records; retain the same null, Unicode and numeric validation.
func RecordAddress(prefix string, record any) (string, error) {
	raw, err := canonical.Marshal(struct {
		SchemaVersion string `json:"schema_version"`
		Record        any    `json:"record"`
	}{SchemaVersion, record})
	if err != nil {
		return "", err
	}
	var original struct {
		Record map[string]json.RawMessage `json:"record"`
	}
	if err = json.Unmarshal(raw, &original); err != nil {
		return "", err
	}
	fields := original.Record
	delete(fields, "id")
	b, err := canonical.Marshal(struct {
		SchemaVersion string                     `json:"schema_version"`
		Record        map[string]json.RawMessage `json:"record"`
	}{SchemaVersion, fields})
	if err != nil {
		return "", err
	}
	var env struct {
		Record json.RawMessage `json:"record"`
	}
	if err = json.Unmarshal(b, &env); err != nil {
		return "", err
	}
	sum := sha256.Sum256(env.Record)
	return prefix + hex.EncodeToString(sum[:]), nil
}
func validDigest(s string) bool {
	if !strings.HasPrefix(s, "sha256:") || len(s) != 71 {
		return false
	}
	for _, c := range s[7:] {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}
func validAddress(s string) bool {
	for _, p := range []string{"obs:", "note:", "ret:"} {
		if strings.HasPrefix(s, p) {
			return validDigest("sha256:" + strings.TrimPrefix(s, p))
		}
	}
	return false
}
func (v Import) Digest() (string, error) {
	if err := v.Validate(); err != nil {
		return "", err
	}
	return canonical.Digest(v)
}
func (v Snapshot) Digest() (string, error) {
	if v.SchemaVersion != SchemaVersion || v.Kind != "evidence-snapshot" || v.TaxonomyVersion != "v1" || v.Imports == nil || !validDigest(v.RegistryDigest) || v.RegistryRef.Name == "" || v.RegistryRef.Version == "" {
		return "", refuse("evidence_invalid_snapshot", "invalid snapshot header")
	}
	if err := canonical.CheckOrdered(v.Imports, func(s string) string { return s }); err != nil {
		return "", err
	}
	for _, s := range v.Imports {
		if !validDigest(s) {
			return "", refuse("evidence_invalid_digest", "invalid import digest")
		}
	}
	return canonical.Digest(v)
}
