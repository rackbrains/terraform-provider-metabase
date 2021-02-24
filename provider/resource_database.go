package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// details

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Database name",
			},
			"engine": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Database Engine",
			},

			"is_full_sync": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_on_demand": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_run_queries": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"user": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"db": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		CreateContext: resourceCreateDatabase,
		ReadContext:   resourceReadDatabase,
		UpdateContext: resourceUpdateDatabase,
		DeleteContext: resourceDeleteDatabase,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateDatabase(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func resourceUpdateDatabase(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func resourceReadDatabase(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func resourceDeleteDatabase(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}
