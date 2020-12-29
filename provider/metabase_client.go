package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type MetabaseClientInterface interface {
	updateCard(id string, query putQuery) (*CardResponse, error)
	postCard(query postQuery) (*CardResponse, error)
	getCard(id string) (*CardResponse, error)
	deleteCard(id string) error
}

type MetabaseClient struct {
	host   string
	id     string
	client *http.Client
}

type CardResponse struct {
	Archived        bool              `json:"archived"`
	EnableEmbedding bool              `json:"enable_embedding"`
	Name            string            `json:"name"`
	Id              int               `json:"id"`
	Display         string            `json:"display"`
	Description     string            `json:"description"`
	DatasetQuery    Query             `json:"dataset_query"`
	CollectionId    int               `json:"collection_id,omitempty"`
	EmbeddingParams map[string]string `json:"embedding_params,omitempty"`
}

type TemplateTag struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	DisplayName string `json:"display-name"`
	Required    bool   `json:"required"`
	Default     string `json:"default,omitempty"`
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
		log.Printf("json creation failed\n")
		return nil, err
	}
	log.Printf("updateCard JSON: %s", string(queryJson))

	log.Printf("init http client\n")
	url := c.host + "/api/card/" + id
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(queryJson))
	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("X-Metabase-Session", c.id)
	resp, err := c.httpClient().Do(req)
	log.Printf("performed request: updateCard \n")
	if err != nil {
		log.Printf("request failed\n")
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		log.Printf("request failed with status %d \n", resp.StatusCode)
		return nil, errors.New("Update Request failed " + string(body))
	}
	log.Printf("updateCard request succeeded\n")
	log.Printf("updateCard Response: %s", string(body))
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
	log.Printf("Getting session @ %s", url)
	log.Printf("Getting session payload %s", string(queryJson))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(queryJson))
	if err != nil {
		log.Printf("request creation failed\n")
		return nil, err
	}
	req.Header.Add("Content-Type", `application/json`)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("request failed\n")
		return nil, err
	}
	if resp.StatusCode >= 400 {
		log.Printf("request failed with status %d", resp.StatusCode)
		return nil, errors.New("Could not initialize session with metabase")
	}
	log.Printf("GetMetabaseClient request succeeded\n")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("response reading failed\n")
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
	log.Printf("Deleting card @ %s \n", url)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("X-Metabase-Session", c.id)
	_, err := c.httpClient().Do(req)
	log.Printf("performed request: deleteCard\n")
	if err != nil {
		log.Printf("request failed\n")
		return err
	}
	return nil
}

func (c MetabaseClient) getCard(id string) (*CardResponse, error) {
	url := c.host + "/api/card/" + id
	log.Printf("Getting card @ %s \n", url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("X-Metabase-Session", c.id)
	resp, err := c.httpClient().Do(req)
	log.Printf("performed request getCard\n")
	if err != nil {
		log.Printf("request failed\n")
		return nil, err
	}
	log.Printf("getCard request succeeded\n")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Printf("getCard Response: %s", string(body))
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
		log.Printf("json creation failed\n")
		return nil, err
	}
	log.Printf("Post Card: %s", string(queryJson))

	log.Printf("init http client\n")
	req, _ := http.NewRequest("POST", c.host+"/api/card", bytes.NewBuffer(queryJson))
	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("X-Metabase-Session", c.id)
	resp, err := c.httpClient().Do(req)
	log.Printf("performed request postCard\n")
	if err != nil {
		log.Printf("request failed")
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		log.Printf("request failed with status %d \n", resp.StatusCode)
		return nil, errors.New("Request failed: " + string(body))
	}
	log.Printf("postCard request succeeded\n")
	log.Printf("postCard Response: %s", string(body))
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
