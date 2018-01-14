package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type keyConfig struct {
	identity, password string
}

func keyCmd(cfg *config) *cobra.Command {
	var kc keyConfig

	cmd := &cobra.Command{
		Use:   "key",
		Short: "virgil key management",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(2)
		},
	}

	cmd.PersistentFlags().StringVar(&kc.identity, "identity", "", "your identity")
	cmd.PersistentFlags().StringVar(&kc.password, "password", "", "your key password")

	markRequired(cmd.PersistentFlags(), "identity", "password")

	cmd.AddCommand(keyCreateCmd(cfg, &kc))

	return cmd
}

func keyCreateCmd(cfg *config, c *keyConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a public/private key pair",
		Run: func(cmd *cobra.Command, args []string) {
			a := newApp(cfg)

			key, err := a.api.Keys.Generate()
			if err != nil {
				a.exit(errors.Wrap(err, "error creating key"))
			}

			if err := key.Save(c.identity, c.password); err != nil {
				a.exit(errors.Wrap(err, "error saving key"))
			}
		},
	}

	return cmd
}
