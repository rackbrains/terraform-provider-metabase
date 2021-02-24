package provider

import "net/http"

type MetabaseClientInterface interface {
	updateCard(id string, query UpdateCardQuery) (*CardResponse, error)
	postCard(query CreateCardQuery) (*CardResponse, error)
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

type UpdateCardQuery struct {
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

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type AuthResponse struct {
	Id string `json:"id"`
}

type CreateCardQuery struct {
	Name                  string            `json:"name"`
	Display               string            `json:"display"`
	VisualizationSettings map[string]string `json:"visualization_settings"`
	DatasetQuery          Query             `json:"dataset_query"`
	Description           string            `json:"description,omitempty"`
	CollectionId          int               `json:"collection_id,omitempty"`
}
