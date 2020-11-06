package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("METABASE_HOST", ""),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("METABASE_USERNAME", ""),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("METABASE_PASSWORD", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"metabase_card":       resourceCard(),
			"metabase_collection": resourceCollection(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	print("provider initialization\n")
	var diags diag.Diagnostics
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	host := d.Get("host").(string)

	client, err := GetMetabaseClient(host, username, password)

	if err != nil {
		print("client initialization failed\n")
		return nil, diag.FromErr(err)
	}

	return *client, diags
}
