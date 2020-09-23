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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "value must be a non-blank string.",
				// ForceNew:    true,
			},
			// "description": {
			// 	Type:        schema.TypeString,
			// 	Required:    false,
			// 	Description: "value may be nil, or if non-nil, value must be a non-blank string.",
			// },
			// "visualization_settings": {
			// 	Type:        schema.TypeSet,
			// 	Required:    false,
			// 	Description: "An optional list of tags, represented as a key, value pair",
			// 	// Elem:        &schema.Schema{Type: schema.TypeString},
			// },
			// "collection_id": {
			// 	Type:        schema.TypeInt,
			// 	Required:    false,
			// 	Description: "value may be nil, or if non-nil, value must be an integer greater than zero.",
			// },
			// "collection_position": {
			// 	Type:        schema.TypeInt,
			// 	Required:    false,
			// 	Description: "value may be nil, or if non-nil, value must be an integer greater than zero.",
			// 	// Elem:        &schema.Schema{Type: schema.TypeString},
			// },
			// "result_metadata": {
			// 	Type:        schema.TypeList,
			// 	Required:    false,
			// 	Description: "value may be nil, or if non-nil, value must be an array of valid results column metadata maps.",
			// 	// Elem:        &schema.Schema{Type: schema.TypeString},
			// },
			// "metadata_checksum": {
			// 	Type:        schema.TypeString,
			// 	Required:    false,
			// 	Description: "value may be nil, or if non-nil, value must be a non-blank string.",
			// },
			// "dataset_query": {
			// 	Required: true,
			// },
			// "display": {
			// 	Required: true,
			// 	Type:     schema.TypeString,
			// },
			// "id": {
			// 	Type:     schema.TypeInt,
			// 	Computed: true,
			// },
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

type query struct {
	Type     string `json:"type,omitempty"`
	Database int    `json:"database,omitempty"`
}

type CardResponse struct {
	Archived        bool   `json:"archived"`
	CanWrite        bool   `json:"can_write"`
	EnableEmbedding bool   `json:"enable_embedding"`
	Name            string `json:"name"`
	Id              int    `json:"id"`
	Display         string `json:"display"`
}

type postQuery struct {
	Name                  string            `json:"name,omitempty"`
	Display               string            `json:"display,omitempty"`
	VisualizationSettings map[string]string `json:"visualization_settings,omitempty"`
	DatasetQuery          query             `json:"dataset_query,omitempty"`
}

func resourceCreateCard(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	print("resourceCreateCard\n")
	c := m.(Client)
	print("got client\n")
	var diags diag.Diagnostics
	print("init diags\n")

	query := postQuery{
		Name:                  d.Get("name").(string),
		Display:               "table",
		VisualizationSettings: map[string]string{},
		DatasetQuery:          query{Type: "native", Database: 15},
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
	if resp.StatusCode >= 400 {
		print("request failed with status", resp.StatusCode, "\n")
		return diag.Errorf("Could not initialize session with metabase")
	}
	print("request succeeded\n")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	print(string(body))
	res := CardResponse{}
	json.Unmarshal(body, &res)
	updateResourceFromCard(res, d)
	return diags
}

func resourceReadCard(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(Client)

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
	print(string(body))
	res := CardResponse{}
	json.Unmarshal(body, &res)
	updateResourceFromCard(res, d)
	return diags
}

func resourceUpdateCard(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	print("resourceUpdateCard\n")
	c := m.(Client)
	print("got client\n")
	var diags diag.Diagnostics
	print("init diags\n")

	query := postQuery{}
	if d.HasChange("name") {
		query.Name = d.Get("name").(string)
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
	url := c.host + "/api/card/" + d.Id()
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(queryJson))
	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("X-Metabase-Session", c.id)
	resp, err := client.Do(req)
	print("performed request\n")
	if err != nil {
		print("request failed")
		return diag.FromErr(err)
	}
	if resp.StatusCode >= 400 {
		print("request failed with status", resp.StatusCode, "\n")
		return diag.Errorf("Update Request failed")
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

func resourceDeleteCard(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(Client)
	print("resourceReadCard\n")
	var diags diag.Diagnostics
	client := &http.Client{}
	url := c.host + "/api/card/" + d.Id()
	print("Deleting card @ ", url, "\n")
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Content-Type", `application/json`)
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
	print("updateResourceFromCard ", card.Id, card.Name, "\n")
	d.SetId(fmt.Sprint(card.Id))
	d.Set("name", card.Name)
}
