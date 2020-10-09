package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var metabaseAccProviders map[string]*schema.Provider
var metabaseAccProvider *schema.Provider

func init() {
	metabaseAccProvider = Provider()
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
