package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCard() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "value must be a non-blank string.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "value may be nil, or if non-nil, value must be a non-blank string.",
				Default:     "Managed by Terraform.",
			},
			"collection_id": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "value may be nil, or if non-nil, value must be an integer greater than zero.",
			},
			"query": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"query_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "native",
			},
			"display": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "table",
			},
			"enable_embedding": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"connection_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"variables": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"display_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"required": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"default": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  nil,
						},
						"embedding_param": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "enabled",
						},
					},
				},
			},
		},
		CreateContext: resourceCreateCard,
		ReadContext:   resourceReadCard,
		UpdateContext: resourceUpdateCard,
		DeleteContext: resourceDeleteCard,
		// Exists:        resourceExistsCard,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateCard(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	print("resourceCreateCard\n")
	c := m.(MetabaseClientInterface)
	print("resourceCreateCard::got client\n")
	var diags diag.Diagnostics
	print("resourceCreateCard::init diags\n")

	query := postCardQuery{
		Name:                  d.Get("name").(string),
		Display:               d.Get("display").(string),
		VisualizationSettings: map[string]string{},
		DatasetQuery: Query{
			Type:     d.Get("query_type").(string),
			Database: d.Get("connection_id").(int),
			Native: NativeQuery{
				Query:        d.Get("query").(string),
				TemplateTags: extractTags(d),
			},
		},
		Description:  d.Get("description").(string),
		CollectionId: d.Get("collection_id").(int),
	}
	print("resourceCreateCard::built query\n")

	//build update query before overriding it with postCard results
	updateQuery := putCardQuery{
		EnableEmbedding: d.Get("enable_embedding").(bool),
		EmbeddingParams: extractEmbeddingParams(d),
	}

	card, err := c.postCard(query)
	if err != nil {
		print("card creation failed\n")
		return diag.FromErr(err)
	}
	updateResourceFromCard(*card, d)

	// update enable_embedding and embedding_params
	resUpdate, err := c.updateCard(d.Id(), updateQuery)
	if err != nil {
		print("card update failed")
		return diag.FromErr(err)
	}
	updateResourceFromCard(*resUpdate, d)

	return diags
}

func resourceReadCard(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(MetabaseClientInterface)

	print("resourceReadCard\n")
	var diags diag.Diagnostics
	card, err := c.getCard(d.Id())
	if err != nil {
		print("request failed\n")
		return diag.FromErr(err)
	}
	updateResourceFromCard(*card, d)
	return diags
}

func resourceUpdateCard(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	print("resourceUpdateCard\n")
	c := m.(MetabaseClientInterface)
	print("got client\n")
	var diags diag.Diagnostics
	print("init diags\n")

	query := putCardQuery{}
	if d.HasChange("name") {
		query.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		query.Description = d.Get("description").(string)
	}
	if d.HasChanges("query_type", "query", "variables", "connection_id") {
		query.DatasetQuery = &Query{
			Type:     d.Get("query_type").(string),
			Database: d.Get("connection_id").(int),
			Native: NativeQuery{
				Query:        d.Get("query").(string),
				TemplateTags: extractTags(d),
			},
		}

		query.EmbeddingParams = extractEmbeddingParams(d)
	} else {
		query.DatasetQuery = nil
		query.EmbeddingParams = nil
	}
	if d.HasChange("collection_id") {
		query.CollectionId = d.Get("collection_id").(int)
	}
	if d.HasChange("enable_embedding") {
		query.EnableEmbedding = d.Get("enable_embedding").(bool)
	}
	print("built query\n")

	res, err := c.updateCard(d.Id(), query)
	if err != nil {
		print("error while updating card")
		return diag.FromErr(err)
	}
	updateResourceFromCard(*res, d)
	return diags
}

func resourceDeleteCard(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(MetabaseClientInterface)
	print("resourceDeleteCard\n")
	var diags diag.Diagnostics
	err := c.deleteCard(d.Id())
	if err != nil {
		print("deletion failed\n")
		return diag.FromErr(err)
	}
	print("deletion succeeded\n")

	d.SetId("")
	return diags
}

func updateResourceFromCard(card CardResponse, d *schema.ResourceData) {
	print("updateResourceFromCard Id:", card.Id, ", name:", card.Name, "\n")
	d.SetId(fmt.Sprint(card.Id))
	d.Set("name", card.Name)
	d.Set("description", card.Description)
	d.Set("display", card.Display)
	d.Set("query", card.DatasetQuery.Native.Query)
	d.Set("query_type", card.DatasetQuery.Type)
	d.Set("collection_id", card.CollectionId)
	d.Set("enable_embedding", card.EnableEmbedding)

	updateVariablesFromTags(card.DatasetQuery.Native.TemplateTags, card.EmbeddingParams, d)
}

func updateVariablesFromTags(tags map[string]TemplateTag, embedding map[string]string, d *schema.ResourceData) error {
	var ois []interface{}
	for _, element := range tags {
		oi := make(map[string]interface{})
		oi["default"] = element.Default
		oi["display_name"] = element.DisplayName
		oi["id"] = element.Id
		oi["name"] = element.Name
		oi["required"] = element.Required
		oi["type"] = element.Type
		oi["embedding_param"] = embedding[element.Name]

		ois = append(ois, oi)
	}

	d.Set("variables", ois)

	return nil
}

func extractTags(d *schema.ResourceData) map[string]TemplateTag {
	variables := d.Get("variables").([]interface{})
	tags := make(map[string]TemplateTag)
	for _, variable := range variables {
		i := variable.(map[string]interface{})
		name := i["name"].(string)

		tag := TemplateTag{
			Id:          i["id"].(string),
			Name:        name,
			Type:        i["type"].(string),
			DisplayName: i["display_name"].(string),
			Required:    i["required"].(bool),
			Default:     i["default"].(string),
		}
		tags[name] = tag
	}
	return tags
}

func extractEmbeddingParams(d *schema.ResourceData) map[string]string {
	variables := d.Get("variables").([]interface{})
	embeddingParams := make(map[string]string)
	for _, variable := range variables {
		i := variable.(map[string]interface{})
		name := i["name"].(string)
		if val, ok := i["embedding_param"]; ok {
			embeddingParams[name] = val.(string)
		}
	}
	return embeddingParams
}
