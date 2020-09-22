package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	print("==========provider=============")
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
				// DefaultFunc: schema.EnvDefaultFunc("METABASE_HOST", ""),
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				// DefaultFunc: schema.EnvDefaultFunc("METABASE_USERNAME", ""),
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
				// DefaultFunc: schema.EnvDefaultFunc("METABASE_PASSWORD", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"metabase_card": resourceCard(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// username := d.Get("username").(string)
	// password := d.Get("password").(string)

	// var host *string

	// hVal, ok := d.GetOk("host")
	return nil, nil
}
