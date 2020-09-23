package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

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
			"metabase_card": resourceCard(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type authResponse struct {
	Id string `json:"id"`
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	print("provider initialization\n")
	var diags diag.Diagnostics
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	host := d.Get("host").(string)

	query := authRequest{
		Username: username,
		Password: password,
	}
	queryJson, err := json.Marshal(query)

	httpClient := &http.Client{}
	url := host + "/api/session"
	print("Getting session @ ", url, "\n")
	print("Getting session payload ", string(queryJson), "\n")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(queryJson))
	if err != nil {
		print("request creation failed\n")
		return nil, diag.FromErr(err)
	}
	req.Header.Add("Content-Type", `application/json`)
	resp, err := httpClient.Do(req)
	if err != nil {
		print("request failed\n")
		return nil, diag.FromErr(err)
	}
	if resp.StatusCode >= 400 {
		print("request failed with status", resp.StatusCode, "\n")
		return nil, diag.Errorf("Could not initialize session with metabase")
	}
	print("request succeeded\n")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print("response reading failed\n")
		return nil, diag.FromErr(err)
	}
	response := authResponse{}
	json.Unmarshal(body, &response)
	res := MetabaseClient{
		host: host,
		id:   response.Id,
	}

	return res, diags
}
