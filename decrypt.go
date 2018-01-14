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

			buf, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				a.exit(errors.Wrap(err, "error reading input data"))
			}

			dec, err := key.Decrypt(virgilapi.Buffer(buf))
			if err != nil {
				a.exit(errors.Wrap(err, "error decrypting data"))
			}

			if _, err := io.Copy(os.Stdout, bytes.NewBuffer(dec)); err != nil {
				a.exit(errors.Wrap(err, "error writing decrypted data"))
			}
		},
	}

	cmd.PersistentFlags().StringVar(&dc.identity, "identity", "", "your identity")
	cmd.PersistentFlags().StringVar(&dc.password, "password", "", "your key password")

	markRequired(cmd.PersistentFlags(), "identity", "password")

	return cmd
}
