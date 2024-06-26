// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/ghodss/yaml"

	"github.com/ory/x/osx"
	"github.com/ory/x/stringsx"
)

// ReadAndParseFiles reads and parses JSON/YAML files from the given sources.
func ReadAndParseFiles(files []string) ([]json.RawMessage, error) {
	var fileContents []json.RawMessage
	for _, source := range files {
		config, err := readAndParseFile(source)
		if err != nil {
			return nil, err
		}
		fileContents = append(fileContents, config)
	}
	return fileContents, nil
}

func readAndParseFile(source string) (json.RawMessage, error) {
	contents, err := osx.ReadFileFromAllSources(source, osx.WithEnabledBase64Loader(), osx.WithEnabledHTTPLoader(), osx.WithEnabledFileLoader())
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", source, err)
	}

	switch f := stringsx.SwitchExact(filepath.Ext(source)); {
	case f.AddCase(".yaml"), f.AddCase(".yml"):
		var config json.RawMessage
		if err := yaml.Unmarshal(contents, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML file %q: %w", source, err)
		}
		return config, nil
	case f.AddCase(".json"):
		var config json.RawMessage
		if err := json.NewDecoder(bytes.NewReader(contents)).Decode(&config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON file %q: %w", source, err)
		}
		return config, nil
	default:
		return nil, f.ToUnknownCaseErr()
	}
}
