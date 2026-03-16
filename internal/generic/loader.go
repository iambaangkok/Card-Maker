package generic

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadProjectConfig loads a ProjectConfig from the given path.
// The file may be YAML or JSON; YAML is preferred.
func LoadProjectConfig(path string) (ProjectConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ProjectConfig{}, err
	}

	var cfg ProjectConfig
	switch ext := filepath.Ext(path); ext {
	case ".yaml", ".yml", "":
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return ProjectConfig{}, fmt.Errorf("unmarshal project config yaml: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, &cfg); err != nil {
			return ProjectConfig{}, fmt.Errorf("unmarshal project config json: %w", err)
		}
	default:
		return ProjectConfig{}, fmt.Errorf("unsupported project config extension %q", ext)
	}

	return cfg, nil
}

// LoadCards loads GenericCard instances for a given card type schema from the
// project's data directory. Each schema.DataFiles entry is resolved relative
// to dataDir and may be YAML or JSON.
func LoadCards(schema CardTypeSchema, dataDir string) ([]GenericCard, error) {
	var all []GenericCard

	for _, file := range schema.DataFiles {
		fullPath := filepath.Join(dataDir, file)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("read card data %q: %w", fullPath, err)
		}

		var cards []GenericCard
		switch ext := filepath.Ext(fullPath); ext {
		case ".yaml", ".yml", "":
			if err := yaml.Unmarshal(data, &cards); err != nil {
				return nil, fmt.Errorf("unmarshal cards yaml %q: %w", fullPath, err)
			}
		case ".json":
			if err := json.Unmarshal(data, &cards); err != nil {
				return nil, fmt.Errorf("unmarshal cards json %q: %w", fullPath, err)
			}
		default:
			return nil, fmt.Errorf("unsupported card data extension %q for %q", ext, fullPath)
		}

		all = append(all, cards...)
	}

	return all, nil
}

// LoadReferenceData loads reference data files declared in project.ReferenceData.
// Returns a map of key -> loaded data (typically []map[string]interface{} for YAML lists).
func LoadReferenceData(project ProjectConfig) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	if project.ReferenceData == nil {
		return out, nil
	}

	for key, file := range project.ReferenceData {
		fullPath := filepath.Join(project.DataDir, file)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("read reference data %q (%s): %w", key, fullPath, err)
		}

		var loaded interface{}
		switch ext := filepath.Ext(fullPath); ext {
		case ".yaml", ".yml", "":
			if err := yaml.Unmarshal(data, &loaded); err != nil {
				return nil, fmt.Errorf("unmarshal reference data %q: %w", key, err)
			}
		case ".json":
			if err := json.Unmarshal(data, &loaded); err != nil {
				return nil, fmt.Errorf("unmarshal reference data %q: %w", key, err)
			}
		default:
			return nil, fmt.Errorf("unsupported reference data extension %q for %s", ext, fullPath)
		}

		out[key] = loaded
	}

	return out, nil
}

