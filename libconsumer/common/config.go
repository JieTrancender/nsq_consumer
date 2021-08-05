package common

import (
	"errors"

	"github.com/elastic/go-ucfg"
)

// Config object to store hierarchical configurations into.
// See https://godoc.org/github.com/elastic/go-ucfg#Config
type Config ucfg.Config

// ConfigNamespace storing at most one configuration section by name and sub-section.
type ConfigNamespace struct {
	name   string
	config *Config
}

var configOpts = []ucfg.Option{
	ucfg.PathSep("."),
	ucfg.ResolveEnv,
	ucfg.VarExp,
}

func NewConfig() *Config {
	return fromConfig(ucfg.New())
}

// NewConfigFrom creates a new Config object from the given input.
// From can be any kind of structured data(struct, map, array, slice).
func NewConfigFrom(from interface{}) (*Config, error) {
	c, err := ucfg.NewFrom(from, []ucfg.Option{}...)
	return fromConfig(c), err
}

func fromConfig(in *ucfg.Config) *Config {
	return (*Config)(in)
}

func (c *Config) Unpack(to interface{}) error {
	return c.access().Unpack(to, configOpts...)
}

func (c *Config) Child(name string, idx int) (*Config, error) {
	sub, err := c.access().Child(name, idx, configOpts...)
	return fromConfig(sub), err
}

func (c *Config) HasField(name string) bool {
	return c.access().HasField(name)
}

// Enabled return the configured enabled value or true by default
func (c *Config) Enabled() bool {
	testEnabled := struct {
		Enabled bool `config:"enabled"`
	}{true}

	if c == nil {
		return false
	}

	if err := c.Unpack(&testEnabled); err != nil {
		// if unpacking fails, expect 'enabled' being set to default value
		return true
	}
	return testEnabled.Enabled
}

func (c *Config) access() *ucfg.Config {
	return (*ucfg.Config)(c)
}

func (c *Config) GetFields() []string {
	return c.access().GetFields()
}

// Unpack unpacks a configuration with at most one sub object. Ab sub object is
// ignored if it is disabled by setting `enabled: false`. If the configuration
// passed contains multiple active sub objects, Unpack will return an error.
func (ns *ConfigNamespace) Unpack(cfg *Config) error {
	fields := cfg.GetFields()
	if len(fields) == 0 {
		return nil
	}

	var (
		err   error
		found bool
	)

	for _, name := range fields {
		var sub *Config
		sub, err = cfg.Child(name, -1)
		if err != nil {
			// element is no configuration object -> continue so a namespace
			// Config unpacked as a namespace can have other configuration
			// values as well
			continue
		}

		if !sub.Enabled() {
			continue
		}

		if ns.name != "" {
			return errors.New("more than one namespace configured")
		}

		ns.name = name
		ns.config = sub
		found = true
	}

	if !found {
		return err
	}
	return nil
}

func (ns *ConfigNamespace) Name() string {
	return ns.name
}

func (ns *ConfigNamespace) Config() *Config {
	return ns.config
}

func (ns *ConfigNamespace) IsSet() bool {
	return ns.config != nil
}

func (c *Config) String(name string, idx int) (string, error) {
	return c.access().String(name, idx, configOpts...)
}
