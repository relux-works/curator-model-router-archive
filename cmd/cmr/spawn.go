package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/relux-works/curator-model-router/pkg/recommend"
)

func runSpawn(args []string, stdout, stderr io.Writer) int {
	var wrapper recommendOptions
	f := recommendFlagSet("spawn", true, &wrapper)
	boundary, scanErr := optionBoundary(f, args)
	if scanErr != nil {
		return commandError(scanErr, hasJSON(args), stdout, stderr)
	}
	if boundary < 0 {
		return writeRefusal(&Refusal{Code: "cmr_invalid_arguments", Message: "spawn requires -- before task-board arguments"}, hasJSON(args), stdout, stderr)
	}
	o, err := recommendFlags("spawn", args[:boundary], true)
	if err != nil {
		return commandError(err, hasJSON(args[:boundary]), stdout, stderr)
	}
	if o.mode != "" && o.mode != "select" && o.mode != "recommend" && o.mode != "shadow" {
		return writeRefusal(&Refusal{Code: "cmr_invalid_arguments", Message: "mode must be select, recommend or shadow"}, o.json, stdout, stderr)
	}
	forwarded := args[boundary+1:]
	parsed, err := parseSpawnArgs(forwarded)
	if err != nil {
		return commandError(err, o.json, stdout, stderr)
	}
	locks := parsed.locks
	// Caller flags are locks, including flags after positional spawn arguments.
	for _, x := range []struct {
		name   string
		target *string
	}{{"agent", &o.agent}, {"model", &o.model}, {"reasoning-effort", &o.effort}, {"role", &o.role}} {
		if v, ok := locks[x.name]; ok {
			if *x.target != "" && *x.target != v && !(x.name == "model" && canonicalModel(*x.target) == canonicalModel(v)) {
				return writeRefusal(&Refusal{Code: "cmr_invalid_arguments", Message: "conflicting --" + x.name}, o.json, stdout, stderr)
			}
			*x.target = v
		}
	}
	o.admission = &recommend.AdmissionContext{BoardFlags: parsed.boardFlags}
	record, err := recommendation(o)
	if err != nil {
		return commandError(err, o.json, stdout, stderr)
	}
	mode := o.mode
	if mode == "" {
		mode = string(record.Inputs.Policy.Mode)
	}
	if mode != "select" && mode != "recommend" && mode != "shadow" {
		return writeRefusal(&Refusal{Code: "cmr_invalid_arguments", Message: "mode must be select, recommend or shadow"}, o.json, stdout, stderr)
	}
	if record.Recommendation.Refusal != nil {
		if err = outputRecommendation(record, o.json, nil, mode, stdout); err != nil {
			return writeRefusal(&Refusal{Code: "cmr_output_failed", Message: "could not write recommendation"}, o.json, stdout, stderr)
		}
		return 2
	}
	commands := spawnCommands(record, forwarded, parsed, mode)
	if mode == "recommend" {
		if err = outputRecommendation(record, o.json, commands, mode, stdout); err != nil {
			return writeRefusal(&Refusal{Code: "cmr_output_failed", Message: "could not write recommendation"}, o.json, stdout, stderr)
		}
		return 0
	}
	if _, err = fmt.Fprintf(stderr, "cmr:%s %s admission=%s\n", record.DecisionID, mode, record.Inputs.AdmissionSource); err != nil {
		return writeRefusal(&Refusal{Code: "cmr_output_failed", Message: "could not log decision"}, o.json, stdout, stderr)
	}
	binary, err := exec.LookPath("task-board")
	if err != nil {
		return writeRefusal(&Refusal{Code: "cmr_task_board_unavailable", Message: "task-board is not on PATH"}, o.json, stdout, stderr)
	}
	// Run sequentially so fan-out has one real spawn gate per configuration and
	// its status/output reaches the caller. The stored decision precedes launch.
	for _, command := range commands {
		cmd := exec.Command(binary, command[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err = cmd.Run(); err != nil {
			var exit *exec.ExitError
			if errors.As(err, &exit) {
				return exit.ExitCode()
			}
			return writeRefusal(&Refusal{Code: "cmr_spawn_failed", Message: "could not execute task-board"}, o.json, stdout, stderr)
		}
	}
	return 0
}
func spawnCommands(record recommend.DecisionRecord, args []string, parsed spawnArguments, mode string) [][]string {
	if mode == "shadow" {
		return [][]string{append([]string{"task-board", "spawn"}, args...)}
	}
	selected := []recommend.Candidate{*record.Recommendation.Selected}
	if record.Inputs.Task.Pipeline == "fanout" {
		selected = record.Recommendation.FanOut
	}
	commands := [][]string{}
	for _, c := range selected {
		added := []string{}
		// The launch must use the role whose admission and rules we evaluated.
		if _, ok := parsed.locks["role"]; !ok {
			added = append(added, "--role", record.Inputs.Task.Role)
		}
		// Fan-out launches several runs of one element; task-board returns the
		// existing live run unless parallel runs are requested explicitly.
		if len(selected) > 1 && !parsed.allowParallel {
			added = append(added, "--allow-parallel")
		}
		flags := spawnFlags(c)
		for i := 0; i < len(flags); i += 2 {
			if _, ok := parsed.locks[strings.TrimPrefix(flags[i], "--")]; !ok {
				added = append(added, flags[i], flags[i+1])
			}
		}
		if record.Inputs.Admission != nil && record.Inputs.Admission.RationaleRequired[c.Runtime] && !parsed.rationalePresent {
			added = append(added, "--selection-rationale", "cmr:"+record.DecisionID+" task="+record.Inputs.Task.TaskClass+" admission="+record.Inputs.AdmissionSource)
		}
		// Inject before task-board's own -- terminator, preserving all caller bytes.
		boundary := parsed.injectionIndex
		command := append([]string{"task-board", "spawn"}, args[:boundary]...)
		command = append(command, added...)
		command = append(command, args[boundary:]...)
		commands = append(commands, command)
	}
	return commands
}
