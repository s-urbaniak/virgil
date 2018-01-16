package main

import (
	"fmt"
	"os"
	"testing"

	virgil "gopkg.in/virgil.v4"
	"gopkg.in/virgil.v4/virgilapi"
	"gopkg.in/virgil.v4/virgilcrypto"
)

var (
	api  *virgilapi.Api
	keys []*virgilapi.Key
)

var text = []byte(`Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. 

Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi. Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat. 

Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat. Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi.`)

func TestMain(m *testing.M) {
	var err error
	api, err = virgilapi.New("")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func BenchmarkGenerateKey(b *testing.B) {
	keys = make([]*virgilapi.Key, b.N)
	for i := 0; i < b.N; i++ {
		key, err := api.Keys.Generate()
		if err != nil {
			b.Fatal(err)
		}
		keys[i] = key
	}
}

func BenchmarkEncrypt(b *testing.B) {
	rn := []int{
		1,
		10,
		100,
		1000,
		10000,
	}

	for _, n := range rn {
		recipients := make([]virgilcrypto.PublicKey, n)

		for i := 0; i < n; i++ {
			key, err := api.Keys.Generate()
			if err != nil {
				b.Fatal(err)
			}
			buf, err := key.ExportPublicKey()
			if err != nil {
				b.Fatal(err)
			}

			pubkey, err := virgilcrypto.DecodePublicKey(buf)
			if err != nil {
				b.Fatal(err)
			}

			recipients[i] = pubkey
		}

		b.Run(fmt.Sprintf("%d recipients", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				virgil.Crypto().Encrypt(text, recipients...)
			}
		})
	}
}
