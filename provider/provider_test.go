package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var metabaseAccProviders map[string]*schema.Provider
var metabaseAccProvider *schema.Provider

func init() {
	metabaseAccProvider = Provider()
	metabaseAccProvider.ConfigureContextFunc = nil
	metabaseAccProvider.ConfigureFunc = nil
	metabaseAccProvider.SetMeta(MetabaseClientMock{})
	metabaseAccProviders = map[string]*schema.Provider{
		"metabase": metabaseAccProvider,
	}
}

func TestProvider(t *testing.T) {
	provider := Provider()
	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

type MetabaseClientMock struct{}

var cards = map[int]*CardResponse{}
var lastId int = 0

func (c MetabaseClientMock) UpdateCard(id string, query UpdateCardQuery) (*CardResponse, error) {
	res := cards[lastId]

	return res, nil
}
func (c MetabaseClientMock) CreateCard(query CreateCardQuery) (*CardResponse, error) {
	lastId++
	res := &CardResponse{
		Name:         query.Name,
		Id:           lastId,
		Description:  query.Description,
		Display:      query.Display,
		DatasetQuery: query.DatasetQuery,
	}
	cards[lastId] = res
	return res, nil
}
func (c MetabaseClientMock) GetCard(id string) (*CardResponse, error) {
	return cards[lastId], nil
}
func (c MetabaseClientMock) DeleteCard(id string) error {
	return nil
}
