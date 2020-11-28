package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/tuckner/terraform-provider-tines/tines"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: tines.Provider})
}
