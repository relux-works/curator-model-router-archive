package cmrio

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/relux-works/curator-model-router/pkg/recommend"
)

func TestTOMLPolicy(t *testing.T) {
	raw, err := os.ReadFile("../../docs/standing-rules.toml")
	if err != nil {
		t.Fatal(err)
	}
	p, err := DecodePolicy(raw, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Rules) == 0 || p.Rules[0].Source == "" || p.TierEdges.S != 62 {
		t.Fatalf("policy: %+v", p)
	}
	p, err = DecodePolicy([]byte("allow_metered = false\nfanout_k = 2\n[difficulties.standard]\nminimum_tier = 'A'\n[headroom]\nwindows = ['weekly']\n"), false)
	if err != nil {
		t.Fatal(err)
	}
	if p.Difficulties["standard"].Objective != "standard" || p.Headroom.Windows.IDs[0] != "weekly" {
		t.Fatalf("overlay: %+v", p)
	}
	for _, raw := range []string{"unknown = true", "fanout_k = 0", "allow_metered = true\nallow_metered = false", "[tier_edges]\ns = 1"} {
		if _, err := DecodePolicy([]byte(raw), false); err == nil {
			t.Fatalf("accepted %s", raw)
		}
	}
}
func TestConfigPrecedence(t *testing.T) {
	root := t.TempDir()
	t.Setenv("HOME", root)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))
	t.Setenv("XDG_STATE_HOME", filepath.Join(root, "state"))
	catalog, err := LoadCatalog("")
	if err != nil || len(catalog.Rows) == 0 {
		t.Fatalf("embedded %v", err)
	}
	p, err := LoadPolicy("")
	if err != nil || p.Mode != "select" {
		t.Fatalf("defaults %v", err)
	}
	if _, err = LoadPolicy(filepath.Join(root, "absent.toml")); err == nil {
		t.Fatal("explicit missing policy")
	}
	if _, err = LoadCatalog(filepath.Join(root, "absent.json")); err == nil {
		t.Fatal("explicit missing catalog")
	}
	path := filepath.Join(root, "config", "curator", "model-router", "catalog.json")
	if err = WriteAtomic(path, []byte(`{"schema_version":"catalog-v1","rows":[]}`)); err != nil {
		t.Fatal(err)
	}
	catalog, err = LoadCatalog("")
	if err != nil || len(catalog.Rows) != 0 {
		t.Fatal(catalog, err)
	}
	explicit := filepath.Join(root, "catalog.json")
	if err = WriteAtomic(explicit, recommend.DefaultCatalog()); err != nil {
		t.Fatal(err)
	}
	catalog, err = LoadCatalog(explicit)
	if err != nil || len(catalog.Rows) == 0 {
		t.Fatal(catalog, err)
	}
}
