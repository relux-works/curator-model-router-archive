// Package cmrio owns the command layer's configuration and decision IO.
package cmrio

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/relux-works/curator-model-router/pkg/recommend"
)

func StateRoot() string {
	root := os.Getenv("XDG_STATE_HOME")
	if root == "" {
		root = filepath.Join(os.Getenv("HOME"), ".local", "state")
	}
	return filepath.Join(root, "curator", "model-router")
}
func LoadCatalog(path string) (recommend.Catalog, error) {
	if path == "" {
		root := os.Getenv("XDG_CONFIG_HOME")
		if root == "" {
			root = filepath.Join(os.Getenv("HOME"), ".config")
		}
		path = filepath.Join(root, "curator", "model-router", "catalog.json")
		b, err := os.ReadFile(path)
		if os.IsNotExist(err) {
			return recommend.LoadCatalog(recommend.DefaultCatalog())
		}
		if err != nil {
			return recommend.Catalog{}, &recommend.Refusal{Code: recommend.InvalidCatalog, Message: "cannot read catalog"}
		}
		return recommend.LoadCatalog(b)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return recommend.Catalog{}, &recommend.Refusal{Code: recommend.InvalidCatalog, Message: "cannot read catalog"}
	}
	return recommend.LoadCatalog(b)
}
func LoadPolicy(path string) (recommend.Policy, error) {
	explicit := path != ""
	if !explicit {
		path = filepath.Join(os.Getenv("HOME"), ".curator", "routing.toml")
	}
	b, err := os.ReadFile(path)
	if !explicit && os.IsNotExist(err) {
		return recommend.DefaultPolicy(), nil
	}
	if err != nil {
		return recommend.Policy{}, &recommend.Refusal{Code: recommend.InvalidPolicy, Message: "cannot read policy"}
	}
	return DecodePolicy(b, strings.EqualFold(filepath.Ext(path), ".json"))
}
func DecodePolicy(b []byte, asJSON bool) (recommend.Policy, error) {
	if !asJSON {
		// Convert TOML's wire object to JSON to use the library's strict overlay,
		// including partial difficulty defaults and custom headroom window filters.
		var wire map[string]any
		if err := toml.NewDecoder(bytes.NewReader(b)).Decode(&wire); err != nil {
			return recommend.Policy{}, &recommend.Refusal{Code: recommend.InvalidPolicy, Message: "invalid TOML policy"}
		}
		var err error
		b, err = json.Marshal(wire)
		if err != nil {
			return recommend.Policy{}, &recommend.Refusal{Code: recommend.InvalidPolicy, Message: "invalid TOML policy values"}
		}
	}
	return recommend.LoadPolicy(b)
}
func SaveDecision(root string, record recommend.DecisionRecord) error {
	b, err := record.JSON()
	if err != nil {
		return err
	}
	return WriteAtomic(filepath.Join(root, "decisions", strings.TrimPrefix(record.DecisionID, "sha256:")+".json"), append(b, '\n'))
}
func WriteAtomic(path string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".cmr-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err = f.Write(b); err != nil {
		f.Close()
		return err
	}
	if err = f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}
