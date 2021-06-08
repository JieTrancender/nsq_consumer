package common

import (
	"github.com/elastic/go-ucfg"
)

// Config object to store hierarchical configurations into.
// See https://godoc.org/github.com/elastic/go-ucfg#Config
type Config ucfg.Config

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

func (c *Config) HasField(name string) bool {
	return c.access().HasField(name)
}

func (c *Config) access() *ucfg.Config {
	return (*ucfg.Config)(c)
}

func (c *Config) Child(name string, idx int) (*Config, error) {
	sub, err := c.access().Child(name, idx, configOpts...)
	return fromConfig(sub), err
}
