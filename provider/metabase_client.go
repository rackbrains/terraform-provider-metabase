package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

func NewClient(host string, username string, password string) (*MetabaseClient, error) {
	query := AuthRequest{
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
	log.Printf("NewClient request succeeded\n")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("response reading failed\n")
		return nil, err
	}
	response := AuthResponse{}
	json.Unmarshal(body, &response)
	res := MetabaseClient{
		host:   host,
		id:     response.Id,
		client: httpClient,
	}

	return &res, nil
}

func (c MetabaseClient) UpdateCard(id string, query UpdateCardQuery) (*CardResponse, error) {
	queryJson, err := json.Marshal(query)
	if err != nil {
		log.Printf("json creation failed\n")
		return nil, err
	}
	log.Printf("UpdateCard JSON: %s", string(queryJson))

	log.Printf("init http client\n")
	url := c.host + "/api/card/" + id
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(queryJson))
	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("X-Metabase-Session", c.id)
	resp, err := c.Client().Do(req)
	log.Printf("performed request: UpdateCard \n")
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
	log.Printf("UpdateCard request succeeded\n")
	log.Printf("UpdateCard Response: %s", string(body))
	res := CardResponse{}
	json.Unmarshal(body, &res)
	return &res, nil
}

func (c MetabaseClient) DeleteCard(id string) error {
	url := c.host + "/api/card/" + id
	log.Printf("Deleting card @ %s \n", url)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("X-Metabase-Session", c.id)
	_, err := c.Client().Do(req)
	log.Printf("performed request: deleteCard\n")
	if err != nil {
		log.Printf("request failed\n")
		return err
	}
	return nil
}

func (c MetabaseClient) GetCard(id string) (*CardResponse, error) {
	url := c.host + "/api/card/" + id
	log.Printf("Getting card @ %s \n", url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("X-Metabase-Session", c.id)
	resp, err := c.Client().Do(req)
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

func (c MetabaseClient) CreateCard(query CreateCardQuery) (*CardResponse, error) {
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
	resp, err := c.Client().Do(req)
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

func (c MetabaseClient) Client() *http.Client {
	if c.client == nil {
		c.client = new(http.Client)
	}
	return c.client
}
