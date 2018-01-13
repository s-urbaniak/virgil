package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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
	cmd.AddCommand(cardExportCmd(c))

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
