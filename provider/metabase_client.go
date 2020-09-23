package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type putQuery struct {
	Name                  string            `json:"name,omitempty"`
	Display               string            `json:"display,omitempty"`
	VisualizationSettings map[string]string `json:"visualization_settings,omitempty"`
	DatasetQuery          *Query            `json:"dataset_query,omitempty"`
	Description           string            `json:"description,omitempty"`
	CollectionId          int               `json:"collection_id,omitempty"`
	EnableEmbedding       bool              `json:"enable_embedding,omitempty"`
	EmbeddingParams       map[string]string `json:"embedding_params,omitempty"`
}

type MetabaseClient struct {
	host string
	id   string
}

func (c MetabaseClient) updateCard(id string, query putQuery) (*CardResponse, error) {
	queryJson, err := json.Marshal(query)
	if err != nil {
		print("json creation failed\n")
		return nil, err
	}
	print(string(queryJson), "\n")

	client := &http.Client{}
	print("init http client\n")
	url := c.host + "/api/card/" + id
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(queryJson))
	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("X-Metabase-Session", c.id)
	resp, err := client.Do(req)
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
