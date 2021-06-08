package instance

import "github.com/spf13/pflag"

type Settings struct {
	Name        string
	IndexPrefix string
	Version     string

	RunFlags *pflag.FlagSet
}
