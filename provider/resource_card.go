package provider

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
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
			"description": {
				Type:        schema.TypeString,
				Required:    false,
				Description: "value may be nil, or if non-nil, value must be a non-blank string.",
			},
			"visualization_settings": {
				Type:        schema.TypeSet,
				Required:    false,
				Description: "An optional list of tags, represented as a key, value pair",
				// Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"collection_id": {
				Type:        schema.TypeInt,
				Required:    false,
				Description: "value may be nil, or if non-nil, value must be an integer greater than zero.",
			},
			"collection_position": {
				Type:        schema.TypeInt,
				Required:    false,
				Description: "value may be nil, or if non-nil, value must be an integer greater than zero.",
				// Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"result_metadata": {
				Type:        schema.TypeList,
				Required:    false,
				Description: "value may be nil, or if non-nil, value must be an array of valid results column metadata maps.",
				// Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"metadata_checksum": {
				Type:        schema.TypeString,
				Required:    false,
				Description: "value may be nil, or if non-nil, value must be a non-blank string.",
			},
			"dataset_query": {
				Required: true,
			},
			"display": {
				Required: true,
				Type:     schema.TypeString,
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		Create: resourceCreateCard,
		Read:   resourceReadCard,
		Update: resourceUpdateCard,
		Delete: resourceDeleteCard,
		Exists: resourceExistsCard,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

type query struct {
	type_    string `json:"type"`
	database int
}

type card struct {
	name                   string
	display                string
	visualization_settings map[string]string
	dataset_query          query
}

func resourceCreateCard(d *schema.ResourceData, m interface{}) error {
	query := &card{
		name: d.Get("name").(string),
	}
	queryJson, _ := json.Marshal(query)
	http.Post("https://metabase.perxtech.io/api/card", "application/json", bytes.NewBuffer(queryJson))
	return nil
}
func resourceReadCard(d *schema.ResourceData, m interface{}) error {
	return nil
}
func resourceUpdateCard(d *schema.ResourceData, m interface{}) error {
	return nil
}
func resourceDeleteCard(d *schema.ResourceData, m interface{}) error {
	return nil
}
func resourceExistsCard(d *schema.ResourceData, m interface{}) (bool, error) {
	return false, nil
}
