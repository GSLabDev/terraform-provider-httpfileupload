package main

import (
	"github.com/GSLabDev/terraform-provider-httpfileupload/httpfileupload"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return httpfileupload.Provider()
		},
	})
}
