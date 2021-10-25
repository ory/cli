package pkg

import (
	"github.com/pkg/errors"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/ory/x/stringsx"
)

type (
	Config struct {
		Project           oryProject     `yaml:"project"`
		PreReleaseHooks   []string       `yaml:"pre_release_hooks"`
		IgnoreTagPatterns []string       `yaml:"ignore_tags"`
		IgnoreTags        *regexp.Regexp `yaml:"-"`
	}
	oryProject string
)

func ReadConfig() (*Config, error) {
	raw, err := os.ReadFile(".orycli.yml")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.Errorf("The file .orycli.yml does not exist in this repository. It is required to be able to create a release. Did you create the file?")
		}
		return nil, err
	}

	var c Config
	if err := yaml.Unmarshal(raw, &c); err != nil {
		return nil, err
	}

	c.IgnoreTags, err = regexp.Compile(strings.Join(c.IgnoreTagPatterns, "|"))
	if err != nil {
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
