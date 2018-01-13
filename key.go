package main

import (
	"fmt"
	"os"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	v "gopkg.in/virgil.v4"
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
	cmd.AddCommand(keyRevokeCmd(cfg, &kc))

	return cmd
}

func keyCreateCmd(cfg *config, c *keyConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a public/private key pair and a virgil card",
		Run: func(cmd *cobra.Command, args []string) {
			a := newApp(cfg)

			key, err := a.api.Keys.Generate()
			if err != nil {
				a.exit(errors.Wrap(err, "error creating key"))
			}

			card, err := a.api.Cards.Create(c.identity, key, nil)
			if err != nil {
				a.exit(errors.Wrap(err, "error creating card"))
			}

			card, err = a.api.Cards.Publish(card)
			if err != nil {
				a.exit(errors.Wrap(err, "error publishing card"))
			}

			fmt.Fprintln(os.Stdout, "created card ID", card.ID, "using identity", c.identity)

			if err := key.Save(c.identity, c.password); err != nil {
				a.exit(errors.Wrap(err, "error saving key"))
			}
		},
	}

	return cmd
}

func keyRevokeCmd(cfg *config, c *keyConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke",
		Short: "revoke a virgil card",
		Run: func(cmd *cobra.Command, args []string) {
			a := newApp(cfg)

			_, err := a.api.Keys.Load(c.identity, c.password)
			if err != nil {
				a.exit(errors.Wrap(err, "error loading key"))
			}

			cs, err := a.api.Cards.Find(c.identity)
			if err != nil {
				a.exit(errors.Wrap(err, "error finding card"))
			}

			if len(cs) == 0 {
				a.exit(errors.New("no cards found"))
			}

			var errs error
			for i := range cs {
				if err := a.api.Cards.Revoke(cs[i], v.RevocationReason.Unspecified); err != nil {
					multierror.Append(errs, err)
				}
			}

			if errs != nil {
				a.exit(errors.Wrap(errs, "error revoking card"))
			}
		},
	}

	return cmd
}
