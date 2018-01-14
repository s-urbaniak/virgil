package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	v "gopkg.in/virgil.v4"
)

func cardCmd(c *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "card",
		Short: "virgil card management",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(2)
		},
	}

	cmd.AddCommand(cardFindCmd(c))
	cmd.AddCommand(cardCreateCmd(c))
	cmd.AddCommand(cardRevokeCmd(c))
	cmd.AddCommand(cardExportCmd(c))

	return cmd
}

func cardCreateCmd(cfg *config) *cobra.Command {
	var kc keyConfig

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a public/private key pair and a virgil card",
		Run: func(cmd *cobra.Command, args []string) {
			a := newApp(cfg)

			key, err := a.api.Keys.Generate()
			if err != nil {
				a.exit(errors.Wrap(err, "error creating key"))
			}

			card, err := a.api.Cards.Create(kc.identity, key, nil)
			if err != nil {
				a.exit(errors.Wrap(err, "error creating card"))
			}

			card, err = a.api.Cards.Publish(card)
			if err != nil {
				a.exit(errors.Wrap(err, "error publishing card"))
			}

			fmt.Fprintln(os.Stdout, "created card ID", card.ID, "using identity", kc.identity)

			if err := key.Save(kc.identity, kc.password); err != nil {
				a.exit(errors.Wrap(err, "error saving key"))
			}
		},
	}

	cmd.PersistentFlags().StringVar(&kc.identity, "identity", "", "your identity")
	cmd.PersistentFlags().StringVar(&kc.password, "password", "", "your key password")

	markRequired(cmd.PersistentFlags(), "identity", "password")

	return cmd
}

func cardRevokeCmd(cfg *config) *cobra.Command {
	var kc keyConfig

	cmd := &cobra.Command{
		Use:   "revoke",
		Short: "revoke a virgil card",
		Run: func(cmd *cobra.Command, args []string) {
			a := newApp(cfg)

			_, err := a.api.Keys.Load(kc.identity, kc.password)
			if err != nil {
				a.exit(errors.Wrap(err, "error loading key"))
			}

			cs, err := a.api.Cards.Find(kc.identity)
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

	cmd.PersistentFlags().StringVar(&kc.identity, "identity", "", "your identity")
	cmd.PersistentFlags().StringVar(&kc.password, "password", "", "your key password")

	markRequired(cmd.PersistentFlags(), "identity", "password")

	return cmd
}

func cardFindCmd(cfg *config) *cobra.Command {
	var identities []string

	cmd := &cobra.Command{
		Use:   "find",
		Short: "find virgil cards",
		Run: func(cmd *cobra.Command, args []string) {
			a := newApp(cfg)

			cs, err := a.api.Cards.Find(identities...)
			if err != nil {
				a.exit(errors.Wrap(err, "error finding cards"))
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug)
			fmt.Fprintln(w, "ID\tIdentity\tType\tScope\tVersion\tCreated At")
			for i := range cs {
				fmt.Fprintf(
					w,
					"%s\t%s\t%s\t%s\t%s\t%s\n",
					cs[i].ID, cs[i].Identity, cs[i].IdentityType, cs[i].Scope, cs[i].CardVersion, cs[i].CreatedAt,
				)
			}
			w.Flush()
		},
	}

	cmd.Flags().StringArrayVar(&identities, "identity", nil, "the identities to find")
	_ = cmd.MarkFlagRequired("identity")

	return cmd
}

func cardExportCmd(cfg *config) *cobra.Command {
	var id string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "export a virgil card",
		Run: func(cmd *cobra.Command, args []string) {
			a := newApp(cfg)

			card, err := a.api.Cards.Get(id)
			if err != nil {
				a.exit(errors.Wrap(err, "error finding card"))
			}

			cs, err := card.Export()
			if err != nil {
				a.exit(errors.Wrap(err, "error exporting card"))
			}

			fmt.Fprintln(os.Stdout, cs)
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "the id of the card")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}
