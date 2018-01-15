package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type cryptConfig struct {
	identity, password string
}

func decryptCmd(cfg *config) *cobra.Command {
	var dc cryptConfig

	cmd := &cobra.Command{
		Use:   "decrypt",
		Short: "decrypt data from stdin using a virgil private key",
		Run: func(cmd *cobra.Command, args []string) {
			a := newApp(cfg)

			key, err := a.api.Keys.Load(dc.identity, dc.password)
			if err != nil {
				a.exit(errors.Wrap(err, "error loading key"))
			}

			if err := key.DecryptStream(os.Stdin, os.Stdout); err != nil {
				a.exit(errors.Wrap(err, "errot decrypting data"))
			}
		},
	}

	cmd.PersistentFlags().StringVar(&dc.identity, "identity", "", "your identity")
	cmd.PersistentFlags().StringVar(&dc.password, "password", "", "your key password")

	markRequired(cmd.PersistentFlags(), "identity", "password")

	return cmd
}
