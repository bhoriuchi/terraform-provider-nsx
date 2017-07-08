package main

import (
	"github.com/bhoriuchi/terraform-provider-nsx/nsx"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: nsx.Provider,
	})
}
