package main

import (
	"log"

	// "github.com/gites/vault-auth-file/authfile"
	"vault-auth-file/authfile"
	"github.com/hashicorp/vault/sdk/plugin"
)

func main() {
	if err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: authfile.Factory,
		TLSProviderFunc:    plugin.DefaultTLSProvider(),
	}); err != nil {
		log.Fatal(err)
	}
}
