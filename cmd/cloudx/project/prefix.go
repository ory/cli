package project

import (
	"encoding/json"

	"github.com/tidwall/sjson"
)

func prefixConfig(prefix string, s []string) []string {
	for k := range s {
		s[k] = prefix + s[k]
	}
	return s
}

func prefixIdentityConfig(s []string) []string {
	return prefixConfig("/services/identity/config", s)
}

func prefixPermissionConfig(s []string) []string {
	return prefixConfig("/services/permission/config", s)
}

func prefixFileConfig(prefix string, configs []json.RawMessage) ([]json.RawMessage, error) {
	for k := range configs {
		raw, err := sjson.SetRawBytes(json.RawMessage("{}"), prefix, configs[k])
		if err != nil {
			return nil, err
		}
		configs[k] = raw
	}
	return configs, nil
}

func prefixFileIdentityConfig(configs []json.RawMessage) ([]json.RawMessage, error) {
	return prefixFileConfig("services.identity.config", configs)
}

func prefixFilePermissionConfig(configs []json.RawMessage) ([]json.RawMessage, error) {
	return prefixFileConfig("services.permission.config", configs)
}
