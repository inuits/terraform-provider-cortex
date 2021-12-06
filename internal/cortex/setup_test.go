package cortex

import (
    "os"
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
    testAccProviderFactories map[string]func() (*schema.Provider, error)
)

func cortexProvider() (*schema.Provider, error) {
    err := os.Setenv(EnvCortexAddress,  "http://localhost:8080")
    if err != nil {
        return nil, err
    }

    return Provider(), nil
}

func TestMain(m *testing.M) {
    testAccProviderFactories = map[string]func() (*schema.Provider, error){
        "cortex": cortexProvider,
    }

    os.Exit(m.Run())
}
