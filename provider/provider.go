package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("METABASE_HOST", ""),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("METABASE_USERNAME", ""),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("METABASE_PASSWORD", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"metabase_card": resourceCard(),
			// "metabase_collection": resourceCollection(),
			// "metabase_connection": resourceConnection(),
		},
		// ConfigureFunc: providerConfigure,
	}
}

// func providerConfigure(d *schema.ResourceData) (interface{}, error) {
// 	address := d.Get("host").(string)
// 	port := d.Get("username").(string)
// 	token := d.Get("password").(string)
// 	return client.NewClient(address, port, token), nil
// }
