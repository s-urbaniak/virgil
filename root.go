package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/virgil.v5/virgilapi"
)

type config struct {
	accessToken        string
	appID              string
	privateKeyFile     string
	privateKeyPassword string
	verbose            bool
}

type app struct {
	api  *virgilapi.Api
	exit func(error)
}

func rootCmd() *cobra.Command {
	var (
		cfg config
	)

	root := &cobra.Command{
		Use: "virgil",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(2)
		},
	}

	root.PersistentFlags().StringVar(&cfg.accessToken, "virgil-access-token", "", "virgil API access token")
	root.PersistentFlags().StringVar(&cfg.appID, "virgil-app-id", "", "virgil application ID")
	root.PersistentFlags().StringVar(&cfg.privateKeyFile, "virgil-private-key-file", "", "file location of the virgil private key")
	root.PersistentFlags().StringVar(&cfg.privateKeyPassword, "virgil-private-key-password", "", "password of the virgil private key")
	root.PersistentFlags().BoolVar(&cfg.verbose, "verbose", false, "verbose output of messages")

	rootFlags := []string{
		"virgil-access-token",
		"virgil-app-id",
		"virgil-private-key-file",
		"virgil-private-key-password",
	}

	readFromEnv(root.PersistentFlags(), rootFlags...)

	root.AddCommand(keyCmd(&cfg))
	root.AddCommand(cardCmd(&cfg))
	root.AddCommand(decryptCmd(&cfg))
	root.AddCommand(encryptCmd(&cfg))

	return root
}

func newApp(c *config) *app {
	var (
		api *virgilapi.Api
		err error

		a = &app{
			exit: newExitFunc(c),
		}
	)

	switch {
	case c.accessToken != "" && c.appID != "" && c.privateKeyFile != "" && c.privateKeyPassword != "":
		privateKey, err := ioutil.ReadFile(c.privateKeyFile)
		if err != nil {
			a.exit(errors.Wrap(err, "error opening private key file"))
		}

		api, err = virgilapi.NewWithConfig(virgilapi.Config{
			Token: c.accessToken,
			Credentials: &virgilapi.AppCredentials{
				AppId:              c.appID,
				PrivateKey:         privateKey,
				PrivateKeyPassword: c.privateKeyPassword,
			},
		})

	case c.accessToken != "":
		api, err = virgilapi.NewWithConfig(virgilapi.Config{
			Token: c.accessToken,
		})

	default:
		api, err = virgilapi.New("")
	}

	if err != nil {
		a.exit(errors.Wrap(err, "virgil api initialization failed"))
	}

	a.api = api
	return a
}

func newApi(c *config) (*virgilapi.Api, error) {
	privateKey, err := ioutil.ReadFile(c.privateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "error opening private key file")
	}

	return virgilapi.NewWithConfig(virgilapi.Config{
		Token: c.accessToken,
		Credentials: &virgilapi.AppCredentials{
			AppId:              c.appID,
			PrivateKey:         privateKey,
			PrivateKeyPassword: c.privateKeyPassword,
		},
	})
}

func newExitFunc(c *config) func(error) {
	return func(err error) {
		vf := ""
		if c.verbose {
			vf = "+"
		}
		fmt.Fprintf(os.Stderr, "%"+vf+"v\n", err) // %+v causes to print the stack trace
		os.Exit(1)
	}
}
