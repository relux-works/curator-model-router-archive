package catalog_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/relux-works/curator-model-router/pkg/recommend"
)

func TestNotesDescribeShippedTiers(t *testing.T) {
	data, err := os.ReadFile("NOTES.md")
	if err != nil {
		t.Fatal(err)
	}
	notes := string(data)
	if strings.Contains(notes, ".temp/rc/") {
		t.Fatal("worktree-specific path in notes")
	}
	catalog, err := recommend.LoadCatalog(recommend.DefaultCatalog())
	if err != nil {
		t.Fatal(err)
	}
	policy := recommend.DefaultPolicy()
	for _, edge := range []float64{policy.TierEdges.S, policy.TierEdges.A, policy.TierEdges.B} {
		if !strings.Contains(notes, fmt.Sprintf("≥%g", edge)) {
			t.Fatal("missing shipped edge", edge)
		}
	}
	counts := map[recommend.Tier]int{}
	for _, row := range catalog.Rows {
		quality := row.Quality.Coding
		if quality == nil {
			quality = row.Quality.Overall
		}
		tier := recommend.TierFor(quality, policy)
		counts[tier]++
		prefix := "| `" + row.Runtime + "/" + row.Model + "/" + row.Effort + "` |"
		matches := 0
		for _, line := range strings.Split(notes, "\n") {
			if strings.HasPrefix(line, prefix) && strings.HasSuffix(line, " | "+string(tier)+" |") {
				matches++
			}
		}
		if matches != 1 {
			t.Fatal("missing or stale tier row", prefix, tier)
		}
	}
	for tier, count := range counts {
		if !strings.Contains(notes, fmt.Sprintf("%s=%d", tier, count)) {
			t.Fatal("stale tier count", tier, count)
		}
	}
}
