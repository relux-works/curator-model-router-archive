package evidence

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/relux-works/curator-model-router/pkg/canonical"
	"github.com/relux-works/skill-agents-management/pkg/vendorplugin"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var updateEvidence = flag.Bool("update-evidence", false, "rewrite pinned evidence goldens from source data")

// The module has no full data-only bench accessor. This independent static
// exporter reads only pkg/vendorplugin source literals, never loads vendor exec
// packages, and fails closed on new unsupported declarations. It is test-only;
// production imports the versioned JSON export. No rows are copied by hand.
type sourceExporter struct {
	funcs  map[string]*ast.FuncDecl
	values map[string]ast.Expr
}

func functionName(e ast.Expr) string {
	switch x := e.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.SelectorExpr:
		return x.Sel.Name
	}
	return ""
}
func literalString(e ast.Expr) (string, bool) {
	x, ok := e.(*ast.BasicLit)
	if !ok || x.Kind != token.STRING {
		return "", false
	}
	s, err := strconv.Unquote(x.Value)
	return s, err == nil
}
func literalNumber(e ast.Expr) (float64, bool) {
	x, ok := e.(*ast.BasicLit)
	if !ok {
		return 0, false
	}
	n, err := strconv.ParseFloat(strings.ReplaceAll(x.Value, "_", ""), 64)
	return n, err == nil
}
func (s sourceExporter) expand(e ast.Expr, env map[string]ast.Expr) ast.Expr {
	seen := map[string]bool{}
	for i := 0; i < 30; i++ {
		id, ok := e.(*ast.Ident)
		if !ok || seen[id.Name] {
			return e
		}
		seen[id.Name] = true
		if v, ok := env[id.Name]; ok {
			e = v
			continue
		}
		if v, ok := s.values[id.Name]; ok {
			e = v
			continue
		}
		return e
	}
	return e
}
func cloneEnv(env map[string]ast.Expr) map[string]ast.Expr {
	x := map[string]ast.Expr{}
	for k, v := range env {
		x[k] = v
	}
	return x
}
func (s sourceExporter) body(call *ast.CallExpr, env map[string]ast.Expr) (*ast.BlockStmt, map[string]ast.Expr) {
	fn := s.funcs[functionName(call.Fun)]
	if fn == nil {
		return nil, env
	}
	next := cloneEnv(env)
	arg := 0
	for _, field := range fn.Type.Params.List {
		for _, name := range field.Names {
			if arg < len(call.Args) {
				next[name.Name] = s.expand(call.Args[arg], env)
				arg++
			}
		}
	}
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		a, ok := n.(*ast.AssignStmt)
		if ok {
			for i, name := range a.Lhs {
				if id, ok := name.(*ast.Ident); ok && i < len(a.Rhs) {
					next[id.Name] = a.Rhs[i]
				}
			}
		}
		return true
	})
	return fn.Body, next
}
func (s sourceExporter) text(e ast.Expr, env map[string]ast.Expr) (string, bool) {
	e = s.expand(e, env)
	if x, ok := literalString(e); ok {
		return x, true
	}
	if functionName(e) == "BugHuntBenchFloor" {
		return vendorplugin.BugHuntBenchFloor, true
	}
	return "", false
}
func (s sourceExporter) score(e ast.Expr, env map[string]ast.Expr, depth int) (float64, bool) {
	if depth > 20 {
		return 0, false
	}
	e = s.expand(e, env)
	if c, ok := e.(*ast.CallExpr); ok {
		name := functionName(c.Fun)
		if (name == "rank" || name == "declaredRank") && len(c.Args) > 0 {
			return literalNumber(s.expand(c.Args[0], env))
		}
		if name == "declaredRankFrom" && len(c.Args) > 1 {
			return literalNumber(s.expand(c.Args[1], env))
		}
		body, next := s.body(c, env)
		if body != nil {
			var score float64
			found := false
			ast.Inspect(body, func(n ast.Node) bool {
				if r, ok := n.(*ast.ReturnStmt); ok {
					for _, e := range r.Results {
						if v, ok := s.score(e, next, depth+1); ok {
							score = v
							found = true
						}
					}
					return false
				}
				return true
			})
			return score, found
		}
	}
	if c, ok := e.(*ast.CompositeLit); ok {
		for _, field := range c.Elts {
			if kv, ok := field.(*ast.KeyValueExpr); ok && functionName(kv.Key) == "Score" {
				return literalNumber(s.expand(kv.Value, env))
			}
		}
	}
	return 0, false
}
func (s sourceExporter) claims(e ast.Expr, env map[string]ast.Expr, depth int) []vendorplugin.RankEvidence {
	if depth > 20 {
		return nil
	}
	e = s.expand(e, env)
	if call, ok := e.(*ast.CallExpr); ok {
		name := functionName(call.Fun)
		text := func(i int) string {
			if i >= len(call.Args) {
				panic("unsupported bench arity")
			}
			v, ok := s.text(call.Args[i], env)
			if !ok {
				panic("unsupported bench string")
			}
			return v
		}
		num := func(i int) int {
			v, ok := literalNumber(s.expand(call.Args[i], env))
			if !ok {
				panic("unsupported bench number")
			}
			return int(v)
		}
		switch name {
		case "BughuntMeasured":
			return []vendorplugin.RankEvidence{vendorplugin.BughuntMeasured(text(0), num(1))}
		case "BughuntMeasuredAs":
			return []vendorplugin.RankEvidence{vendorplugin.BughuntMeasuredAs(text(0), text(1), num(2))}
		case "BughuntInterpolated":
			return []vendorplugin.RankEvidence{vendorplugin.BughuntInterpolated(text(0), text(1), text(2))}
		case "BughuntCost":
			return []vendorplugin.RankEvidence{vendorplugin.BughuntCost(text(0))}
		}
		// A rank/declaredRank passes all of its explicit bench claim arguments.
		if name == "rank" || name == "declaredRank" || name == "declaredRankFrom" || name == "withCost" {
			out := []vendorplugin.RankEvidence{}
			for _, arg := range call.Args {
				out = append(out, s.claims(arg, env, depth+1)...)
			}
			return out
		}
		body, next := s.body(call, env)
		if body != nil {
			out := []vendorplugin.RankEvidence{}
			ast.Inspect(body, func(n ast.Node) bool {
				if ret, ok := n.(*ast.ReturnStmt); ok {
					for _, x := range ret.Results {
						out = append(out, s.claims(x, next, depth+1)...)
					}
					return false
				}
				return true
			})
			return out
		}
	}
	if c, ok := e.(*ast.CompositeLit); ok {
		out := []vendorplugin.RankEvidence{}
		for _, x := range c.Elts {
			if kv, ok := x.(*ast.KeyValueExpr); ok {
				x = kv.Value
			}
			out = append(out, s.claims(x, env, depth+1)...)
		}
		return out
	}
	return nil
}
func (s sourceExporter) efforts(e ast.Expr, env map[string]ast.Expr, depth int) []string {
	if depth > 20 {
		return nil
	}
	e = s.expand(e, env)
	if c, ok := e.(*ast.CallExpr); ok {
		name := functionName(c.Fun)
		if name == "effortNone" {
			return []string{"none"}
		}
		if name == "effortRequired" && len(c.Args) == 2 {
			return s.efforts(c.Args[1], env, depth+1)
		}
		body, next := s.body(c, env)
		// muse's local closure is a data literal, not an executable plugin.
		if body == nil {
			if id, ok := c.Fun.(*ast.Ident); ok {
				if x, ok := env[id.Name].(*ast.FuncLit); ok {
					body = x.Body
				}
			}
		}
		if body != nil {
			var out []string
			ast.Inspect(body, func(n ast.Node) bool {
				if ret, ok := n.(*ast.ReturnStmt); ok {
					for _, x := range ret.Results {
						if v := s.efforts(x, next, depth+1); v != nil {
							out = v
						}
					}
					return false
				}
				return true
			})
			return out
		}
	}
	if c, ok := e.(*ast.CompositeLit); ok {
		if _, ok := c.Type.(*ast.ArrayType); ok {
			out := []string{}
			for _, x := range c.Elts {
				v, ok := s.text(x, env)
				if !ok {
					return nil
				}
				out = append(out, v)
			}
			return out
		}
		none := false
		for _, x := range c.Elts {
			if kv, ok := x.(*ast.KeyValueExpr); ok {
				key := functionName(kv.Key)
				if key == "Vocabulary" {
					return s.efforts(kv.Value, env, depth+1)
				}
				if key == "Support" && functionName(kv.Value) == "EffortSupportNone" {
					none = true
				}
			}
		}
		if none {
			return []string{"none"}
		}
	}
	return nil
}
func moduleSource(t *testing.T) string {
	t.Helper()
	// Refuse a go.mod upgrade that leaves the importer/export pinned to an older
	// cached source tree; otherwise a newly added module row could be missed.
	manifest, err := os.ReadFile(filepath.Join("..", "..", "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	fields := strings.Fields(string(manifest))
	actual := ""
	for i, field := range fields {
		if field == "github.com/relux-works/skill-agents-management" && i+1 < len(fields) {
			actual = fields[i+1]
			break
		}
	}
	if actual != ModuleVersion {
		t.Fatalf("compiled dependency must match importer pin: go.mod=%s importer=%s", actual, ModuleVersion)
	}
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/relux-works/skill-agents-management")
	cmd.Env = append(os.Environ(), "GOFLAGS=-mod=mod", "GOPROXY=off")
	raw, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("offline pinned module lookup failed: %v: %s", err, raw)
	}
	dir := strings.TrimSpace(string(raw))
	if dir == "" {
		t.Fatal("offline pinned module lookup returned no directory")
	}
	path := filepath.Join(dir, "pkg", "vendorplugin")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("pinned module source absent: %v", err)
	}
	return path
}
func exportPinnedSource(t *testing.T) BugHuntExport {
	t.Helper()
	out := BugHuntExport{SchemaVersion: SchemaVersion, ModuleVersion: ModuleVersion, MappingVersion: "bughunt-model-map-v1", ImportedAt: "2026-10-01T00:00:00Z", RetrievedAt: "2026-10-01T00:00:00Z", Models: []RegistryModel{}, Mapping: []ModelMapping{}, Rows: []BugHuntRow{}}
	root := moduleSource(t)
	paths, err := sourcePaths(root)
	if err != nil {
		t.Fatal(err)
	}
	costPattern := regexp.MustCompile(`at ([a-z]+) this row fixed [0-9]+/105 for \$([0-9.]+) list cost`)
	for _, path := range paths {
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		modelRows := modelLiterals(file)
		s := sourceExporter{map[string]*ast.FuncDecl{}, map[string]ast.Expr{}}
		for _, d := range file.Decls {
			if fn, ok := d.(*ast.FuncDecl); ok {
				s.funcs[fn.Name.Name] = fn
			}
		}
		ast.Inspect(file, func(n ast.Node) bool {
			if v, ok := n.(*ast.ValueSpec); ok {
				for i, name := range v.Names {
					if i < len(v.Values) {
						s.values[name.Name] = v.Values[i]
					}
				}
			}
			if a, ok := n.(*ast.AssignStmt); ok {
				for i, name := range a.Lhs {
					if id, ok := name.(*ast.Ident); ok && i < len(a.Rhs) {
						s.values[id.Name] = a.Rhs[i]
					}
				}
			}
			return true
		})
		ast.Inspect(file, func(n ast.Node) bool {
			lit, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}
			fields := map[string]ast.Expr{}
			for _, x := range lit.Elts {
				if kv, ok := x.(*ast.KeyValueExpr); ok {
					fields[functionName(kv.Key)] = kv.Value
				}
			}
			if !modelRows[lit] && !looksLikeModel(lit, fields) {
				return true
			}
			id, skip, err := modelRow(s, lit, fields, path, root)
			if err != nil {
				t.Fatal(err)
			}
			if skip {
				return true
			}
			claims := s.claims(fields["Rank"], s.values, 0)
			score, ok := s.score(fields["Rank"], s.values, 0)
			if !ok {
				t.Fatalf("cannot statically export bench score for %s", id)
			}
			efforts := s.efforts(fields["Effort"], s.values, 0)
			if efforts == nil {
				t.Fatalf("cannot statically export vocabulary for %s", id)
			}
			canonical.SortByKey(efforts, func(s string) string { return s })
			out.Models = append(out.Models, RegistryModel{id, efforts})
			out.Mapping = append(out.Mapping, ModelMapping{id, id})
			for index, claim := range claims {
				parsed, err := vendorplugin.ParseBughuntEvidence(claim)
				if err != nil {
					t.Fatal(err)
				}
				row := BugHuntRow{Key: id + ":" + strconv.Itoa(index), ModelName: id, Kind: string(parsed.Kind), Value: score, OriginRef: "bug-hunt:" + id}
				switch parsed.Kind {
				case vendorplugin.BughuntMeasuredClaim:
					row.Value = float64(parsed.Fixed)
					row.Effort = parsed.Effort
					benchID := parsed.BenchID
					if benchID == "" {
						benchID = id
					}
					row.OriginRef = "bug-hunt:" + benchID + ":" + parsed.Effort
				case vendorplugin.BughuntInterpolatedClaim: // Effort not stated: remains unknown.
				case vendorplugin.BughuntCostClaim:
					m := costPattern.FindStringSubmatch(claim.Observation)
					if m == nil {
						t.Fatalf("new cost grammar requires exporter update: %s", id)
					}
					usd, err := strconv.ParseFloat(m[2], 64)
					if err != nil {
						t.Fatal(err)
					}
					row.Effort = m[1]
					row.Value = usd
					row.Cost = &Cost{USD: &usd}
					row.OriginRef = "bug-hunt:" + id + ":" + row.Effort + ":cost"
				default:
					t.Fatalf("unhandled claim kind %q", parsed.Kind)
				}
				out.Rows = append(out.Rows, row)
			}
			return false
		})
	}
	canonical.SortByKey(out.Models, func(m RegistryModel) string { return m.ModelID })
	canonical.SortByKey(out.Mapping, func(m ModelMapping) string { return m.ExternalName })
	canonical.SortByKey(out.Rows, func(r BugHuntRow) string { return r.Key })
	if len(out.Rows) == 0 {
		t.Fatal("empty module export")
	}
	return out
}
func TestAcceptance2BugHuntParity(t *testing.T) {
	export := exportPinnedSource(t)
	checkConstructorParity(t, moduleSource(t), export)
	raw, err := canonical.Marshal(export)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join("data", "bughunt-export.json")
	if *updateEvidence {
		if err = os.WriteFile(path, raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	golden, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != string(golden) {
		t.Fatal("pinned module rows changed or export dropped a row; review and regenerate export")
	}
	doc, _, err := ImportBugHunt(export, "2026-10-01T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Observations) != len(export.Rows) {
		t.Fatalf("got %d observations for %d rows", len(doc.Observations), len(export.Rows))
	}
	for _, row := range export.Rows {
		t.Run(row.Key, func(t *testing.T) {
			found := false
			for _, o := range doc.Observations {
				if o.Provenance.RawArtifactRef != "module:"+ModuleVersion+"/bench/"+row.Key {
					continue
				}
				found = true
				kind := row.Kind
				metric := "fixed"
				if kind == "cost" {
					kind = "measured"
					metric = "list_cost"
				}
				if o.Subject.ModelID != row.ModelName || o.Subject.Effort != row.Effort || o.EvidenceKind != kind || o.Metric != metric || o.Value != row.Value {
					t.Fatalf("row drift: %+v", o)
				}
				a, _ := json.Marshal(o.Cost)
				b, _ := json.Marshal(row.Cost)
				if string(a) != string(b) {
					t.Fatalf("cost drift: %s != %s", a, b)
				}
			}
			if !found {
				t.Fatal("module row dropped")
			}
		})
	}
}

func sourcePaths(root string) ([]string, error) {
	paths := []string{}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			paths = append(paths, path)
		}
		return nil
	})
	return paths, err
}

func looksLikeModel(lit *ast.CompositeLit, fields map[string]ast.Expr) bool {
	if len(lit.Elts) == 0 {
		return false
	}
	if sel, ok := lit.Type.(*ast.SelectorExpr); ok && functionName(sel.X) == "agentic" && sel.Sel.Name == "Model" {
		return false
	}
	return functionName(lit.Type) == "Model" || fields["Rank"] != nil || fields["ID"] != nil && fields["Effort"] != nil || fields["ID"] != nil && (fields["Systems"] != nil || fields["ContextWindowTokens"] != nil)
}

func modelRow(s sourceExporter, lit *ast.CompositeLit, fields map[string]ast.Expr, path, root string) (string, bool, error) {
	for _, key := range []string{"ID", "Rank", "Effort"} {
		if fields[key] == nil {
			return "", false, fmt.Errorf("model row lacks %s in %s at offset %d", key, path, lit.Pos())
		}
	}
	claims := s.claims(fields["Rank"], s.values, 0)
	id, literal := literalString(fields["ID"])
	if len(claims) == 0 {
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return "", false, err
		}
		// Explicit zero-claim allow-list: the pin's sole local-models template is
		// operator-populated, has dynamic modelID and localCapabilityRank(), and
		// declares no benchmark. No other local row or changed template is exempt.
		allowed := map[string]string{"vendors/local-models/models.go:modelID": "localCapabilityRank"}
		key := filepath.ToSlash(relative) + ":" + functionName(fields["ID"])
		rank, call := fields["Rank"].(*ast.CallExpr)
		if !literal && allowed[key] != "" && call && functionName(rank.Fun) == allowed[key] && len(rank.Args) == 0 {
			return "", true, nil
		}
		return "", false, fmt.Errorf("bench row cannot be exported (zero claims): %s at offset %d", path, lit.Pos())
	}
	if !literal {
		return "", false, fmt.Errorf("model row ID cannot be exported in %s at offset %d", path, lit.Pos())
	}
	return id, false, nil
}

// constructorCounts is independent of the model-row interpreter: it visits
// every call in every non-test Go file, including unsupported declaration shapes.
func constructorCounts(root string) (map[string]int, error) {
	paths, err := sourcePaths(root)
	if err != nil {
		return nil, err
	}
	counts := map[string]int{"measured": 0, "interpolated": 0, "cost": 0}
	for _, path := range paths {
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			return nil, err
		}
		ast.Inspect(file, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				switch functionName(call.Fun) {
				case "BughuntMeasured", "BughuntMeasuredAs":
					counts["measured"]++
				case "BughuntInterpolated":
					counts["interpolated"]++
				case "BughuntCost":
					counts["cost"]++
				}
			}
			return true
		})
	}
	return counts, nil
}

func modelLiterals(file *ast.File) map[*ast.CompositeLit]bool {
	rows := map[*ast.CompositeLit]bool{}
	ast.Inspect(file, func(n ast.Node) bool {
		if lit, ok := n.(*ast.CompositeLit); ok {
			if a, ok := lit.Type.(*ast.ArrayType); ok && functionName(a.Elt) == "Model" {
				for _, e := range lit.Elts {
					if row, ok := e.(*ast.CompositeLit); ok {
						rows[row] = true
					}
				}
			}
		}
		return true
	})
	return rows
}

func TestParityModelRowsFailClosed(t *testing.T) {
	for _, tc := range []struct {
		name, row, path string
		skip            bool
		fail            bool
	}{
		{"missing-ID", `Model{Rank: rank(), Effort: axis()}`, "models.go", false, true},
		{"missing-Rank", `Model{ID: "new", Effort: axis()}`, "models.go", false, true},
		{"missing-Effort", `Model{ID: "new", Rank: rank()}`, "models.go", false, true},
		{"aliased-package-missing-Rank", `registry.Model{ID: "new", Effort: axis()}`, "models.go", false, true},
		{"implicit-missing-fields", `[]Model{{ID: "new"}}`, "models.go", false, true},
		{"dynamic-ID", `Model{ID: modelID, Rank: rank(BughuntMeasured("max", 1)), Effort: axis()}`, "models.go", false, true},
		{"new-local-zero", `Model{ID: "new-local", Rank: localCapabilityRank(), Effort: axis()}`, "vendors/local-models/models.go", false, true},
		{"new-local-template", `Model{ID: newID, Rank: localCapabilityRank(), Effort: axis()}`, "vendors/local-models/models.go", false, true},
		{"allow-listed-template", `Model{ID: modelID, Rank: localCapabilityRank(), Effort: axis()}`, "vendors/local-models/models.go", true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			file, err := parser.ParseFile(token.NewFileSet(), "fixture.go", "package fixture; var row = "+tc.row, 0)
			if err != nil {
				t.Fatal(err)
			}
			marked := modelLiterals(file)
			found := false
			ast.Inspect(file, func(n ast.Node) bool {
				lit, ok := n.(*ast.CompositeLit)
				if !ok {
					return true
				}
				fields := map[string]ast.Expr{}
				for _, e := range lit.Elts {
					if kv, ok := e.(*ast.KeyValueExpr); ok {
						fields[functionName(kv.Key)] = kv.Value
					}
				}
				if !marked[lit] && !looksLikeModel(lit, fields) {
					return true
				}
				found = true
				_, skip, err := modelRow(sourceExporter{map[string]*ast.FuncDecl{}, map[string]ast.Expr{}}, lit, fields, tc.path, ".")
				if (err != nil) != tc.fail || skip != tc.skip {
					t.Fatalf("row skip=%v error=%v", skip, err)
				}
				return false
			})
			if !found {
				t.Fatal("model-shaped literal silently skipped")
			}
		})
	}
}
func TestParityConstructorWalkAllFiles(t *testing.T) {
	root := t.TempDir()
	for _, tc := range []struct{ path, code string }{
		{"extra.go", `package fixture; var a = BughuntMeasured("max", 1)`},
		{"nested/unusual.go", `package fixture; var a = vp.BughuntMeasuredAs("alias", "max", 1); var b = vp.BughuntInterpolated("a", "b", "why"); var c = BughuntCost("cost")`},
		{"excluded_test.go", `package fixture; var ignored = BughuntMeasured("max", 1)`},
	} {
		path := filepath.Join(root, tc.path)
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(tc.code), 0600); err != nil {
			t.Fatal(err)
		}
	}
	counts, err := constructorCounts(root)
	if err != nil {
		t.Fatal(err)
	}
	if counts["measured"] != 2 || counts["interpolated"] != 1 || counts["cost"] != 1 {
		t.Fatalf("constructor calls dropped: %v", counts)
	}
}

// The pinned source has 53 constructor calls, expanded to 58 exported claims.
// Five aliases reuse constructors; astra has its own call and is not an expansion.
// Counting calls as if all 58 were distinct would reject the unchanged pin.
func checkConstructorParity(t *testing.T, root string, export BugHuntExport) {
	t.Helper()
	calls, err := constructorCounts(root)
	if err != nil {
		t.Fatal(err)
	}
	expanded := map[string]int{"measured": 0, "interpolated": 0, "cost": 0}
	// Explicit alias expansions observed in the pinned AST, not inferred from
	// equal ranks: shared sonnet55Rank, opus55Rank, sol61Rank, lunaRank and measured13.
	aliases := []struct{ name, target, kind string }{
		{"sonnet", "claude-sonnet-5-5", "interpolated"},
		{"opus", "claude-opus-5-5", "interpolated"},
		{"sol", "gpt-6.1-sol", "interpolated"},
		{"luna", "gpt-6-luna", "interpolated"},
		{"muse-spark", "muse-spark-1.3-contributor", "measured"},
	}
	paths, err := sourcePaths(root)
	if err != nil {
		t.Fatal(err)
	}
	declared := map[string]string{}
	for _, path := range paths {
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(file, func(n ast.Node) bool {
			if lit, ok := n.(*ast.CompositeLit); ok {
				var id, target string
				for _, e := range lit.Elts {
					if kv, ok := e.(*ast.KeyValueExpr); ok {
						if functionName(kv.Key) == "ID" {
							id, _ = literalString(kv.Value)
						}
						if functionName(kv.Key) == "AliasOf" {
							target, _ = literalString(kv.Value)
						}
					}
				}
				if id != "" && target != "" {
					declared[id] = target
				}
			}
			return true
		})
	}
	for _, alias := range aliases {
		if declared[alias.name] != alias.target {
			t.Fatalf("alias expansion no longer declared: %s -> %s", alias.name, alias.target)
		}
		var row, target *BugHuntRow
		for i := range export.Rows {
			x := &export.Rows[i]
			if x.ModelName == alias.name && x.Kind == alias.kind {
				row = x
			}
			if x.ModelName == alias.target && x.Kind == alias.kind {
				target = x
			}
		}
		if row == nil || target == nil || row.Value != target.Value || row.Effort != target.Effort {
			t.Fatalf("alias claims missing or unequal: %s", alias.name)
		}
		expanded[alias.kind]++
	}
	exported := map[string]int{"measured": 0, "interpolated": 0, "cost": 0}
	for _, row := range export.Rows {
		exported[row.Kind]++
	}
	for _, kind := range []string{"measured", "interpolated", "cost"} {
		if calls[kind]+expanded[kind] != exported[kind] {
			t.Fatalf("independent constructor parity failed for %s: %d calls + %d aliases != %d exported claims", kind, calls[kind], expanded[kind], exported[kind])
		}
	}
	total := calls["measured"] + calls["interpolated"] + calls["cost"]
	if total+len(aliases) != len(export.Rows) {
		t.Fatalf("constructor count %d + aliases %d != claim count %d", total, len(aliases), len(export.Rows))
	}
	t.Logf("independent AST: %d constructor calls (%d/%d/%d) + %d declared alias expansions = %d exported claims (%d/%d/%d)", total, calls["measured"], calls["interpolated"], calls["cost"], len(aliases), len(export.Rows), exported["measured"], exported["interpolated"], exported["cost"])
}
