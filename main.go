package main

import (
	// "github.com/PerxTech/terraform-provider-metabase/provider"
	"./provider/provider"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
