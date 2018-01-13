package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func main() {
	root := rootCmd()

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "error executing command"))
		os.Exit(1)
	}
}
