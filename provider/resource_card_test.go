package provider

import (
	"testing"
)

func TestUpdateResourceFromCard(t *testing.T) {
	card := CardResponse{
		Name:        "42",
		Description: "I am good",
		DatasetQuery: Query{
			Native: NativeQuery{
				Query: "SELECT * FROM myTable",
			},
		},
	}
	resource := resourceCard().Data(nil)

	updateResourceFromCard(card, resource)
	if resource.Get("name").(string) != "42" {
		t.Errorf("Name not properly mapped, got %s instead of %s", resource.Get("name"), card.Name)
	}
	if resource.Get("description").(string) != card.Description {
		t.Errorf("Description not properly mapped, got %s instead of %s", resource.Get("description"), card.Description)
	}
	if resource.Get("query").(string) != card.DatasetQuery.Native.Query {
		t.Errorf("Query not properly mapped, got '%s' instead of '%s'", resource.Get("query"), card.DatasetQuery.Native.Query)
	}
}
