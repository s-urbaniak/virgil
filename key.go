package main

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	virgil "gopkg.in/virgil.v4"
	"gopkg.in/virgil.v4/virgilcrypto"
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
	cmd.AddCommand(keyEncryptCmd(cfg))

	return cmd
}

func keyCreateCmd(cfg *config) *cobra.Command {
	var kc keyConfig

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a public/private key pair, writing the public key to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			a := newApp(cfg)

			key, err := a.api.Keys.Generate()
			if err != nil {
				a.exit(errors.Wrap(err, "error creating key"))
			}

			pub, err := key.ExportPublicKey()
			if err != nil {
				a.exit(errors.Wrap(err, "error exporting public key"))
			}

			if err := ioutil.WriteFile(kc.identity+".pub", pub, 0644); err != nil {
				a.exit(errors.Wrap(err, "error writing public key"))
			}

			if err := key.Save(kc.identity, kc.password); err != nil {
				a.exit(errors.Wrap(err, "error saving key"))
			}
		},
	}

	cmd.Flags().StringVar(&kc.identity, "identity", "", "your identity")
	cmd.Flags().StringVar(&kc.password, "password", "", "your key password")

	markRequired(cmd.Flags(), "identity", "password")

	return cmd
}

func keyEncryptCmd(cfg *config) *cobra.Command {
	var identities []string

	cmd := &cobra.Command{
		Use:   "encrypt",
		Short: "encrypt data from stdin to stdout using a local public key, loaded from <identity>.pub",
		Run: func(cmd *cobra.Command, args []string) {
			a := newApp(cfg)

			pubkeys := make([]virgilcrypto.PublicKey, len(identities))
			for i, _ := range identities {
				pubbuf, err := ioutil.ReadFile(identities[i] + ".pub")
				if err != nil {
					a.exit(errors.Wrap(err, "error loading public key"))
				}

				pubkey, err := virgilcrypto.DecodePublicKey(pubbuf)
				if err != nil {
					a.exit(errors.Wrap(err, "error decoding public key"))
				}

				pubkeys[i] = pubkey
			}

			if err := virgil.Crypto().EncryptStream(os.Stdin, os.Stdout, pubkeys...); err != nil {
				a.exit(errors.Wrap(err, "error encrypting data"))
			}
		},
	}

	cmd.Flags().StringArrayVar(&identities, "identity", nil, "the identity")

	markRequired(cmd.Flags(), "identity")

	return cmd
}
