package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCollection() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "value must be a non-blank string.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "value may be nil, or if non-nil, value must be a non-blank string.",
				Default:      "Managed by Terraform.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"color": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "value must be a string that matches the regex ^#[0-9A-Fa-f]{6}$",
				Default:      "#509EE3",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`), "Not a valid color code, should be of the form #f77594"),
			},
			"parent_id": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Required:     false,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"namespace": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "value may be nil, or if non-nil, value must be a non-blank string.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		CreateContext: resourceCreateCollection,
		ReadContext:   resourceReadCollection,
		UpdateContext: resourceUpdateCollection,
		DeleteContext: resourceDeleteCollection,
		// Exists:        resourceExistsCollection,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateCollection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	print("resourceCreateCollection\n")
	// c := m.(MetabaseClientInterface)
	var diags diag.Diagnostics

	return diags
}

func resourceReadCollection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	print("resourceReadCollection\n")
	c := m.(MetabaseClientInterface)
	var diags diag.Diagnostics
	collection, err := c.getCollection(d.Id())
	if err != nil {
		print("request failed\n")
		return diag.FromErr(err)
	}
	updateResourceFromCollection(*collection, d)
	return diags
}

func resourceUpdateCollection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	print("resourceUpdateCollection\n")
	// c := m.(MetabaseClientInterface)
	print("got client\n")
	var diags diag.Diagnostics

	return diags
}

func resourceDeleteCollection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	print("resourceDeleteCollection\n")
	// c := m.(MetabaseClientInterface)
	var diags diag.Diagnostics

	d.SetId("")
	return diags
}

func updateResourceFromCollection(collec CollectionResponse, d *schema.ResourceData) {
	d.SetId(fmt.Sprint(collec.Id))
	d.Set("name", collec.Name)
	if collec.Description != "" {
		d.Set("description", collec.Description)
	}
	d.Set("color", collec.Color)
	print("parent_id", collec.ParentId, "\n")
	if collec.ParentId != 0 {
		d.Set("parent_id", collec.ParentId)
	}
}
