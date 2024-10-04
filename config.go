package main

import (
	"flag"
	"fmt"
	"strings"
)

type Config struct {
	Source             Path
	Destination        Path
	SSHAdd             StringsFlag
	RSyncOverrideFlags string
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := cfg.flags(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) flags() error {
	flag.Var(&c.Source, "src", "")
	flag.Var(&c.Destination, "dest", "")
	flag.Var(&c.SSHAdd, "ssh-add", "List of SSH keys to add to ssh-agent")

	flag.Parse()

	if c.Source == "" {
		return fmt.Errorf("missing required flag: -src")
	}
	if c.Destination == "" {
		return fmt.Errorf("missing required flag: -dest")
	}

	c.gatherRsyncFlagsIfExist()

	return nil
}

func (c *Config) gatherRsyncFlagsIfExist() {
	args := flag.Args()
	sepIndex := -1
	for i, arg := range args {
		if arg == "--" {
			sepIndex = i
			break
		}
	}

	var rsyncArgs []string
	if sepIndex != -1 {
		rsyncArgs = args[sepIndex+1:]
	}

	c.RSyncOverrideFlags = strings.Join(rsyncArgs, " ")
}

type Path string

func (p Path) IsRemoteSpecification() bool {
	return strings.Contains(string(p), ":")
}

func (p *Path) String() string {
	if p == nil {
		return ""
	}
	return string(*p)
}
func (p *Path) Set(s string) error {
	*p = Path(s)
	return nil
}

type StringsFlag []string

func (s *StringsFlag) String() string {
	if s == nil {
		return ""
	}
	return strings.Join(*s, ", ")
}

func (s *StringsFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}
