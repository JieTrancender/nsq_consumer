package cli

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

func exitOnPanice() {
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "panic: %s\n", r)
		debug.PrintStack()
		os.Exit(1)
	}
}

func RunWith(
	fn func(cmd *cobra.Command, args []string) error,
) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		defer exitOnPanice()
		if err := fn(cmd, args); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}
