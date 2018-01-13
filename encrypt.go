package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/virgil.v4/virgilapi"
)

func encryptCmd(cfg *config) *cobra.Command {
	var id string

	cmd := &cobra.Command{
		Use:   "encrypt",
		Short: "encrypt data from stdin using a virgil",
		Run: func(cmd *cobra.Command, args []string) {
			a := newApp(cfg)

			card, err := a.api.Cards.Get(id)
			if err != nil {
				a.exit(errors.Wrap(err, "error finding card"))
			}

			buf, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				a.exit(errors.Wrap(err, "error reading input data"))
			}

			enc, err := card.Encrypt(virgilapi.Buffer(buf))
			if err != nil {
				a.exit(errors.Wrap(err, "error encrypting data"))
			}

			if _, err := io.Copy(os.Stdout, bytes.NewBuffer(enc)); err != nil {
				a.exit(errors.Wrap(err, "error writing encrypted data"))
			}
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "the id of the card")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}
