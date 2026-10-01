package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Refusal struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Refusal) Error() string { return e.Code + ": " + e.Message }
func main()                      { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }
func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return runHelp(nil, stdout, stderr)
	}
	switch args[0] {
	case "route":
		return runRoute(args[1:], stdout, stderr)
	case "explain":
		return runExplain(args[1:], stdout, stderr)
	case "replay":
		return runReplay(args[1:], stdout, stderr)
	case "headroom":
		return runHeadroom(args[1:], stdout, stderr)
	case "evidence":
		return runEvidence(args[1:], stdout, stderr)
	case "note":
		return runNote(args[1:], stdout, stderr)
	case "suitability":
		return runSuitability(args[1:], stdout, stderr)
	case "version":
		return runVersion(args[1:], stdout, stderr)
	case "help", "--help", "-h":
		return runHelp(args[1:], stdout, stderr)
	default:
		return writeRefusal(&Refusal{Code: "cmr_unknown_subcommand", Message: "unknown subcommand: " + strings.ReplaceAll(strings.ReplaceAll(args[0], "\n", " "), "\r", " ")}, hasJSON(args), stdout, stderr)
	}
}
func hasJSON(args []string) bool {
	for _, arg := range args {
		if arg == "--json" {
			return true
		}
	}
	return false
}
func writeRefusal(err *Refusal, asJSON bool, stdout, stderr io.Writer) int {
	if asJSON {
		_ = json.NewEncoder(stdout).Encode(struct {
			Error *Refusal `json:"error"`
		}{err})
	} else {
		_, _ = fmt.Fprintln(stderr, err.Error())
	}
	return 2
}
