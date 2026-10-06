package recommend

import "testing"

// BenchmarkRecommendEmbeddedCatalog measures a complete in-process decision,
// including input freezing, ranking, explanation, and the canonical decision ID.
// Catalog decoding is initialization; file and usage-cache reads are excluded.
func BenchmarkRecommendEmbeddedCatalog(b *testing.B) {
	catalog, err := LoadCatalog(DefaultCatalog())
	if err != nil {
		b.Fatal(err)
	}
	candidates := make([]Candidate, 0, len(catalog.Rows))
	for _, row := range catalog.Rows {
		candidates = append(candidates, row.Candidate)
	}
	request := Request{
		Catalog:         catalog,
		Policy:          DefaultPolicy(),
		Task:            TaskProfile{Role: "developer", TaskClass: "code.implement", Pipeline: "single"},
		Candidates:      candidates,
		Usage:           UsageSnapshot{AsOf: 1791321600}, // Frozen UTC clock; usage is unknown.
		AdmissionSource: "benchmark-admitted",
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := Recommend(request)
		if err != nil {
			b.Fatal(err)
		}
		if result.Refusal != nil || result.Selected == nil || result.DecisionID == "" {
			b.Fatalf("incomplete decision: %+v", result)
		}
	}
}
