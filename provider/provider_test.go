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

func (c MetabaseClientMock) updateCard(id string, query putCardQuery) (*CardResponse, error) {
	res := cards[lastId]

	return res, nil
}
func (c MetabaseClientMock) postCard(query postCardQuery) (*CardResponse, error) {
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
func (c MetabaseClientMock) getCard(id string) (*CardResponse, error) {
	return cards[lastId], nil
}
func (c MetabaseClientMock) deleteCard(id string) error {
	return nil
}
