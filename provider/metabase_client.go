package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type MetabaseClient struct {
	host   string
	id     string
	client *http.Client
}

type CardResponse struct {
	Archived        bool   `json:"archived"`
	EnableEmbedding bool   `json:"enable_embedding"`
	Name            string `json:"name"`
	Id              int    `json:"id"`
	Display         string `json:"display"`
	Description     string `json:"description"`
	DatasetQuery    Query  `json:"dataset_query"`
	CollectionId    int    `json:"collection_id,omitempty"`
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

type putQuery struct {
	Name                  string            `json:"name,omitempty"`
	Display               string            `json:"display,omitempty"`
	VisualizationSettings map[string]string `json:"visualization_settings,omitempty"`
	DatasetQuery          *Query            `json:"dataset_query,omitempty"`
	Description           string            `json:"description,omitempty"`
	CollectionId          int               `json:"collection_id,omitempty"`
	EnableEmbedding       bool              `json:"enable_embedding,omitempty"`
	EmbeddingParams       map[string]string `json:"embedding_params,omitempty"`
	Archived              bool              `json:"archived,omitempty"`
}

func (c MetabaseClient) updateCard(id string, query putQuery) (*CardResponse, error) {
	queryJson, err := json.Marshal(query)
	if err != nil {
		print("json creation failed\n")
		return nil, err
	}
	print(string(queryJson), "\n")

	print("init http client\n")
	url := c.host + "/api/card/" + id
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(queryJson))
	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("X-Metabase-Session", c.id)
	resp, err := c.httpClient().Do(req)
	print("performed request\n")
	if err != nil {
		print("request failed\n")
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		print("request failed with status", resp.StatusCode, "\n")
		return nil, errors.New("Update Request failed " + string(body))
	}
	print("request succeeded\n")
	print(string(body), "\n")
	res := CardResponse{}
	json.Unmarshal(body, &res)
	return &res, nil
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type authResponse struct {
	Id string `json:"id"`
}

func GetMetabaseClient(host string, username string, password string) (*MetabaseClient, error) {
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
		return nil, err
	}
	req.Header.Add("Content-Type", `application/json`)
	resp, err := httpClient.Do(req)
	if err != nil {
		print("request failed\n")
		return nil, err
	}
	if resp.StatusCode >= 400 {
		print("request failed with status", resp.StatusCode, "\n")
		return nil, errors.New("Could not initialize session with metabase")
	}
	print("request succeeded\n")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print("response reading failed\n")
		return nil, err
	}
	response := authResponse{}
	json.Unmarshal(body, &response)
	res := MetabaseClient{
		host:   host,
		id:     response.Id,
		client: httpClient,
	}

	return &res, nil
}

func (c MetabaseClient) deleteCard(id string) error {
	url := c.host + "/api/card/" + id
	print("Deleting card @ ", url, "\n")
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("X-Metabase-Session", c.id)
	_, err := c.httpClient().Do(req)
	print("performed request\n")
	if err != nil {
		print("request failed\n")
		return err
	}
	return nil
}

func (c MetabaseClient) getCard(id string) (*CardResponse, error) {
	url := c.host + "/api/card/" + id
	print("Getting card @ ", url, "\n")
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("X-Metabase-Session", c.id)
	resp, err := c.httpClient().Do(req)
	print("performed request\n")
	if err != nil {
		print("request failed\n")
		return nil, err
	}
	print("request succeeded\n")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	print(string(body), "\n")
	res := CardResponse{}
	json.Unmarshal(body, &res)
	return &res, nil
}

type postQuery struct {
	Name                  string            `json:"name"`
	Display               string            `json:"display"`
	VisualizationSettings map[string]string `json:"visualization_settings"`
	DatasetQuery          Query             `json:"dataset_query"`
	Description           string            `json:"description,omitempty"`
	CollectionId          int               `json:"collection_id,omitempty"`
}

func (c MetabaseClient) postCard(query postQuery) (*CardResponse, error) {
	queryJson, err := json.Marshal(query)
	if err != nil {
		print("json creation failed\n")
		return nil, err
	}
	print(string(queryJson), "\n")

	print("init http client\n")
	req, _ := http.NewRequest("POST", c.host+"/api/card", bytes.NewBuffer(queryJson))
	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("X-Metabase-Session", c.id)
	resp, err := c.httpClient().Do(req)
	print("performed request\n")
	if err != nil {
		print("request failed")
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		print("request failed with status", resp.StatusCode, "\n")
		return nil, errors.New("Request failed: " + string(body))
	}
	print("request succeeded\n")
	print(string(body), "\n")
	res := CardResponse{}
	json.Unmarshal(body, &res)
	return &res, nil
}

func (c MetabaseClient) httpClient() *http.Client {
	if c.client == nil {
		c.client = new(http.Client)
	}
	return c.client
}
