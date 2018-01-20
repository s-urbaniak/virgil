package main

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	virgil "gopkg.in/virgil.v5"
	"gopkg.in/virgil.v5/virgilcrypto"
)

func decryptCmd(cfg *config) *cobra.Command {
	var (
		identity, password, sender string
	)

	cmd := &cobra.Command{
		Use:   "decrypt",
		Short: "decrypt data from stdin using a virgil private key",
		Run: func(cmd *cobra.Command, args []string) {
			a := newApp(cfg)

			privbuf, err := ioutil.ReadFile(identity + ".cup")
			if err != nil {
				a.exit(errors.Wrap(err, "error loading private key"))
			}

			key, err := virgil.Crypto().ImportPrivateKey(privbuf, password)
			if err != nil {
				a.exit(errors.Wrap(err, "error importing private key"))
			}

			if sender != "" {
				pubbuf, err := ioutil.ReadFile(sender + ".pub")
				if err != nil {
					a.exit(errors.Wrap(err, "error loading public key"))
				}

				pubkey, err := virgilcrypto.DecodePublicKey(pubbuf)
				if err != nil {
					a.exit(errors.Wrap(err, "error decoding public key"))
				}

				buf, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					a.exit(errors.Wrap(err, "error reading payload from stdin"))
				}

				dec, err := virgil.Crypto().DecryptThenVerify(buf, key, pubkey)
				if err != nil {
					a.exit(errors.Wrap(err, "error decrypting payload"))
				}
				os.Stdout.Write(dec)
				return
			}
		},
	}

	cmd.PersistentFlags().StringVar(&identity, "identity", "", "your identity")
	cmd.PersistentFlags().StringVar(&password, "password", "", "your key password")
	cmd.PersistentFlags().StringVar(&sender, "sender", "", "the public key of the sender who signed the payload")

	markRequired(cmd.PersistentFlags(), "identity", "password")

	return cmd
}
