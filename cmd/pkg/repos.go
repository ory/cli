package pkg

import (
	"gopkg.in/yaml.v2"
	"os"

	"github.com/ory/x/stringsx"
)

type (
	Config struct {
		Project         oryProject `yaml:"project"`
		PreReleaseHooks []string   `yaml:"pre_release_hooks"`
	}
	oryProject string
)

func ReadConfig() (*Config, error) {
	raw, err := os.ReadFile(".orycli.yml")
	if err != nil {
		return nil, err
	}

	var c Config
	if err := yaml.Unmarshal(raw, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (p *oryProject) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw string
	if err := unmarshal(&raw); err != nil {
		return err
	}

	switch r := stringsx.SwitchExact(raw); {
	case r.AddCase("hydra"),
		r.AddCase("kratos"),
		r.AddCase("keto"),
		r.AddCase("oathkeeper"),
		r.AddCase("cli"):
		*p = oryProject(raw)
	default:
		return r.ToUnknownCaseErr()
	}
	return nil
}
