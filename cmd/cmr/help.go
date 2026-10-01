package main

import (
	"encoding/json"
	"io"
)

func runHelp(args []string, stdout, stderr io.Writer) int {
	if len(args) > 1 || (len(args) == 1 && args[0] != "--json") {
		return writeRefusal(&Refusal{Code: "cmr_invalid_arguments", Message: "usage: cmr help [--json]"}, hasJSON(args), stdout, stderr)
	}
	var err error
	if hasJSON(args) {
		err = json.NewEncoder(stdout).Encode(struct {
			Commands []string `json:"commands"`
		}{[]string{"help", "version"}})
	} else {
		_, err = io.WriteString(stdout, "Usage: cmr <command> [--json]\n\nCommands:\n  help     Show help\n  version  Show version\n")
	}
	if err != nil {
		return writeRefusal(&Refusal{Code: "cmr_output_failed", Message: "could not write help"}, hasJSON(args), stdout, stderr)
	}
	return 0
}
