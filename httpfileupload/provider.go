package httpfileupload

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"httpfileupload_file": resourceUpload(),
		},
	}
}
