package cortex

import (
	"os"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"cortex": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CORTEX_ADDRESS"); v == "" {
		t.Fatal("CORTEX_ADDRESS must be set for acceptance tests")
	}
	if v := os.Getenv("CORTEX_API_KEY"); v == "" {
		t.Fatal("CORTEX_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("CORTEX_TENANT_ID"); v == "" {
		t.Fatal("CORTEX_TENANT_ID must be set for acceptance tests")
	}
}
