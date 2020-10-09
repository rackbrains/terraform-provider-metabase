package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, resource.Get("name").(string), card.Name)
	assert.Equal(t, resource.Get("description").(string), card.Description)
	assert.Equal(t, resource.Get("query").(string), card.DatasetQuery.Native.Query)
}

func TestMetabaseCard_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		// PreCheck:     func() { testAccPreCheck(t) },
		Providers: metabaseAccProviders,
		// CheckDestroy: testAccCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: metabaseAccCheckCardBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"metabase_card.test_item", "name", "test"),
					resource.TestCheckResourceAttr(
						"metabase_card.test_item", "description", "hello"),
				),
			},
		},
	})
}

func metabaseAccCheckCardBasic() string {
	return fmt.Sprintf(`
resource "metabase_card" "test_item" {
  name        = "test"
	description = "hello"
	query ="SELECT * from nico"
	connection_id = 23
}
`)
}
