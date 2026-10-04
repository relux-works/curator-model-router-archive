package main

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/relux-works/curator-model-router/pkg/recommend"
)

func TestReviewCLIShippedGoldenPicks(t *testing.T) {
	root := cliEnvironment(t)
	file := candidateFile(t, root, false)
	defaults := root + "/default-policy.json"
	if err := os.WriteFile(defaults, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
	var scenarios []struct {
		Difficulty string               `json:"difficulty"`
		Mini       bool                 `json:"mini"`
		Selected   *recommend.Candidate `json:"selected"`
		Tier       recommend.Tier       `json:"tier"`
		Index      string               `json:"index"`
		Refusal    string               `json:"refusal"`
	}
	b, err := os.ReadFile("../../pkg/recommend/testdata/review-picks.golden.json")
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(b, &scenarios); err != nil {
		t.Fatal(err)
	}
	for _, s := range scenarios {
		policy := defaults
		if s.Mini {
			policy = "../../docs/standing-rules.toml"
		}
		args := []string{"recommend", "--role", "reviewer", "--task-class", "review.code", "--difficulty", s.Difficulty, "--candidates", file, "--catalog", "../../catalog/catalog.json", "--policy", policy, "--host", "<your-mini-host>", "--json"}
		var out, stderr bytes.Buffer
		code := run(args, &out, &stderr)
		wantCode := 0
		if s.Refusal != "" {
			wantCode = 2
		}
		if code != wantCode {
			t.Fatal(code, out.String(), stderr.String())
		}
		var result recommendOutput
		if err = json.Unmarshal(out.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		if s.Refusal != "" {
			if result.Refusal == nil || result.Refusal.Code != s.Refusal || result.Selected != nil {
				t.Fatal(result)
			}
			continue
		}
		if result.Selected == nil || *result.Selected != *s.Selected {
			t.Fatal(s, result)
		}
		for _, x := range result.Explanation {
			if x.Selected && (x.Tier != s.Tier || x.QualityIndex != s.Index) {
				t.Fatal(x)
			}
		}
	}
}
