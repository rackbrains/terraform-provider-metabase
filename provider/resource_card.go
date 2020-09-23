package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
			// "collection_position": &schema.Schema{
			// 	Type:        schema.TypeInt,
			// Optional:    true,
			// 	Description: "value may be nil, or if non-nil, value must be an integer greater than zero.",
			// },
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
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
	}
}

type CardResponse struct {
	// Archived        bool   `json:"archived"`
	// CanWrite        bool   `json:"can_write"`
	EnableEmbedding bool   `json:"enable_embedding"`
	Name            string `json:"name"`
	Id              int    `json:"id"`
	Display         string `json:"display"`
	Description     string `json:"description"`
	DatasetQuery    Query  `json:"dataset_query"`
	CollectionId    int    `json:"collection_id,omitempty"`
}

type postQuery struct {
	Name                  string            `json:"name"`
	Display               string            `json:"display"`
	VisualizationSettings map[string]string `json:"visualization_settings"`
	DatasetQuery          Query             `json:"dataset_query"`
	Description           string            `json:"description,omitempty"`
	CollectionId          int               `json:"collection_id,omitempty"`
}

type TemplateTag struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	DisplayName string `json:"display_name"`
	Required    bool   `json:"required"`
	Default     string `json:"default"`
}

type Query struct {
	Type     string      `json:"type,omitempty"`
	Database int         `json:"database,omitempty"`
	Native   NativeQuery `json:"native,omitempty"`
}

type NativeQuery struct {
	Query        string                 `json:"query,omitempty"`
	TemplateTags map[string]TemplateTag `json:"template-tags,omitempty"`
}

func resourceCreateCard(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	print("resourceCreateCard\n")
	c := m.(MetabaseClient)
	print("got client\n")
	var diags diag.Diagnostics
	print("init diags\n")

	query := postQuery{
		Name:                  d.Get("name").(string),
		Display:               d.Get("display").(string),
		VisualizationSettings: map[string]string{},
		DatasetQuery: Query{
			Type:     d.Get("query_type").(string),
			Database: 15,
			Native: NativeQuery{
				Query:        d.Get("query").(string),
				TemplateTags: extractTags(d),
			},
		},
		Description:  d.Get("description").(string),
		CollectionId: d.Get("collection_id").(int),
	}
	print("built query\n")

	queryJson, err := json.Marshal(query)
	if err != nil {
		print("json creation failed\n")
		return diag.FromErr(err)
	}
	print(string(queryJson), "\n")

	client := &http.Client{}
	print("init http client\n")
	req, _ := http.NewRequest("POST", c.host+"/api/card", bytes.NewBuffer(queryJson))
	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("X-Metabase-Session", c.id)
	resp, err := client.Do(req)
	print("performed request\n")
	if err != nil {
		print("request failed")
		return diag.FromErr(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		print("request failed with status", resp.StatusCode, "\n")
		return diag.Errorf("Request failed: " + string(body))
	}
	print("request succeeded\n")
	print(string(body), "\n")
	res := CardResponse{}
	json.Unmarshal(body, &res)
	updateResourceFromCard(res, d)

	//todo update enable_embedding and embedding_params

	return diags
}

func resourceReadCard(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(MetabaseClient)

	print("resourceReadCard\n")
	var diags diag.Diagnostics
	client := &http.Client{}
	url := c.host + "/api/card/" + d.Id()
	print("Getting card @ ", url, "\n")
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("X-Metabase-Session", c.id)
	resp, err := client.Do(req)
	print("performed request\n")
	if err != nil {
		print("request failed\n")
		return diag.FromErr(err)
	}
	print("request succeeded\n")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	print(string(body), "\n")
	res := CardResponse{}
	json.Unmarshal(body, &res)
	updateResourceFromCard(res, d)
	return diags
}

func resourceUpdateCard(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	print("resourceUpdateCard\n")
	c := m.(MetabaseClient)
	print("got client\n")
	var diags diag.Diagnostics
	print("init diags\n")

	query := putQuery{}
	if d.HasChange("name") {
		query.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		query.Description = d.Get("description").(string)
	}
	if d.HasChange("query_type") || d.HasChange("query") || d.HasChange("variables") {
		query.DatasetQuery.Type = d.Get("query_type").(string)
		query.DatasetQuery.Native.Query = d.Get("query").(string)
		query.DatasetQuery.Database = 15
		query.DatasetQuery.Native.TemplateTags = extractTags(d)
	} else {
		query.DatasetQuery = nil
	}
	if d.HasChange("collection_id") {
		query.CollectionId = d.Get("collection_id").(int)
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
	c := m.(MetabaseClient)
	print("resourceReadCard\n")
	var diags diag.Diagnostics
	client := &http.Client{}
	url := c.host + "/api/card/" + d.Id()
	print("Deleting card @ ", url, "\n")
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("X-Metabase-Session", c.id)
	_, err := client.Do(req)
	print("performed request\n")
	if err != nil {
		print("request failed\n")
		return diag.FromErr(err)
	}
	print("request succeeded\n")

	d.SetId("")
	return diags
}

func updateResourceFromCard(card CardResponse, d *schema.ResourceData) {
	print("updateResourceFromCard ", card.Id, " ", card.Name, "\n")
	d.SetId(fmt.Sprint(card.Id))
	d.Set("name", card.Name)
	d.Set("description", card.Description)
	d.Set("display", card.Display)
	d.Set("query", card.DatasetQuery.Native.Query)
	d.Set("query_type", card.DatasetQuery.Type)
	d.Set("collection_id", card.CollectionId)
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
