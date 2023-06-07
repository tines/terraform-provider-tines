package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/tines/terraform-provider-tines/tines"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
// go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: tines.Provider})
}
