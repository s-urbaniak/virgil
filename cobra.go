package main

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func markRequired(f *pflag.FlagSet, flags ...string) {
	for i := range flags {
		err := cobra.MarkFlagRequired(f, flags[i])
		if err != nil {
			panic(err)
		}
	}
}

func readFromEnv(f *pflag.FlagSet, flags ...string) {
	for i := range flags {
		v := os.Getenv(strings.Replace(strings.ToUpper(flags[i]), "-", "_", -1))

		if v == "" {
			continue
		}

		_ = f.Set(flags[i], v) // ignore err, causing the flag to be unset
	}
}
