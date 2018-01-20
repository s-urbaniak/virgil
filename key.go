package main

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	virgil "gopkg.in/virgil.v5"
)

type keyConfig struct {
	identity, password string
}

func keyCmd(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key",
		Short: "virgil key management",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(2)
		},
	}

	cmd.AddCommand(keyCreateCmd(cfg))

	return cmd
}

func keyCreateCmd(cfg *config) *cobra.Command {
	var kc keyConfig

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a public/private key pair",
		Run: func(cmd *cobra.Command, args []string) {
			a := newApp(cfg)

			key, err := virgil.Crypto().GenerateKeypair()
			pub, err := key.PublicKey().Encode()
			if err != nil {
				a.exit(errors.Wrap(err, "error encoding public key"))
			}

			priv, err := key.PrivateKey().Encode([]byte(kc.password))
			if err != nil {
				a.exit(errors.Wrap(err, "error encoding private key"))
			}

			if err := ioutil.WriteFile(kc.identity+".pub", pub, 0644); err != nil {
				a.exit(errors.Wrap(err, "error writing public key"))
			}

			if err := ioutil.WriteFile(kc.identity+".cup", priv, 0644); err != nil {
				a.exit(errors.Wrap(err, "error writing CUP"))
			}
		},
	}

	cmd.Flags().StringVar(&kc.identity, "identity", "", "your identity")
	cmd.Flags().StringVar(&kc.password, "password", "", "your key password")

	markRequired(cmd.Flags(), "identity", "password")

	return cmd
}
