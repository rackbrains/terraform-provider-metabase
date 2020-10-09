package provider

import (
	"testing"

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

// func TestAccItem_Basic(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		// PreCheck:     func() { testAccPreCheck(t) },
// 		Providers: metabaseAccProviders,
// 		// CheckDestroy: testAccCheckItemDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccCheckItemBasic(),
// 				Check: resource.ComposeTestCheckFunc(
// 					// testAccCheckExampleItemExists("example_item.test_item"),
// 					resource.TestCheckResourceAttr(
// 						"example_item.test_item", "name", "test"),
// 					resource.TestCheckResourceAttr(
// 						"example_item.test_item", "description", "hello"),
// 					resource.TestCheckResourceAttr(
// 						"example_item.test_item", "tags.#", "2"),
// 					resource.TestCheckResourceAttr("example_item.test_item", "tags.1931743815", "tag1"),
// 					resource.TestCheckResourceAttr("example_item.test_item", "tags.1477001604", "tag2"),
// 				),
// 			},
// 		},
// 	})
// }

// func testAccCheckItemBasic() string {
// 	return fmt.Sprintf(`
// resource "metabase_card" "test_item" {
//   name        = "test"
// 	description = "hello"
// 	query ="SELECT * from nico"
// 	connection_id = 23
// }
// `)
// }
