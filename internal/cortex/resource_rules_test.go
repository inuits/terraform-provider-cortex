package cortex

import (
    "context"
    "fmt"
    "strings"
    "testing"

    "github.com/grafana/cortex-tools/pkg/client"
    "github.com/grafana/cortex-tools/pkg/rules"
    "github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
    "github.com/stretchr/testify/require"
    "gopkg.in/yaml.v3"
)

func TestAccRule_Basic(t *testing.T) {
    var (
    	ruleGroup rwrulefmt.RuleGroup
        resourceName = "watchdog"
        resourceId = "cortex_rules.watchdog"
    )

    ruleContent :=  `name: watchdog
rules:
- alert: watchdog
  expr: vector(1)
  for: 20m`
    ruleConfig := testRuleConfig(resourceName, ruleContent)

    resource.Test(t, resource.TestCase{
        ProviderFactories: testAccProviderFactories,
        PreCheck: func() { testAccPreCheck(t, "default", resourceName) },
        Steps: []resource.TestStep{
            {
                Config: ruleConfig,
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr(resourceId, "id", fmt.Sprintf("%s/default/", resourceName)),
                    resource.TestCheckResourceAttr(resourceId, "namespace", "default"),
                    testAccCheckResourceContent(resourceId, ruleContent),
                    testAccCheckRuleGroupExists(resourceId, &ruleGroup),
                    testAccCheckRuleGroupAttrContent(&ruleGroup, ruleContent),
                ),
            },
        },
    })
}

func TestAccRule_YAMLFormatCausesNoChanges(t *testing.T) {
    ruleContent := `name: watchdog
rules:
- alert: watchdog
  expr: vector(1)
  for: 20m
  labels:
    unquoted_bool: true
    single_quoted_bool: 'true'
    double_quoted_bool: "true"
    unquoted_string: string
    single_quoted_string: 'string'
    double_quoted_string: "string"
    multi_line_string: |-
      Some description
      Link https://github.com`
    ruleConfig := testRuleConfig("watchdog", ruleContent)

    resource.Test(t, resource.TestCase{
        ProviderFactories: testAccProviderFactories,
        PreCheck: func() { testAccPreCheck(t, "default", "watchdog") },
        Steps: []resource.TestStep{
            {
                Config: ruleConfig,
            },
            {
                Config: ruleConfig,
                PlanOnly: true,
            },
        },
    })
}

func testRuleConfig(name, content string) string {
    return fmt.Sprintf(`
resource "cortex_rules" "%s" {
  namespace = "default"
  content   = <<EOT
%s
EOT
}
`, name, content)
}

func testAccPreCheck(t *testing.T, namespace, name string) {
    cortexClient, err := client.New(client.Config{
        Address: "http://localhost:8080",
    })
    require.NoError(t, err)

    rg, err := cortexClient.GetRuleGroup(context.Background(), namespace, name)
    require.EqualError(t, err, "requested resource not found", fmt.Sprintf("test precondition failed: rule group '%s/%s' already exists", namespace, name))
    require.Nil(t, rg)
}

func testAccCheckResourceContent(resourceName, wantContent string) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        var (
            wantRg *rwrulefmt.RuleGroup
            gotRg *rwrulefmt.RuleGroup
        )

        err := yaml.Unmarshal([]byte(wantContent), &wantRg)
        if err != nil {
            return err
        }

        rs, ok := s.RootModule().Resources[resourceName]
        if !ok {
            return fmt.Errorf("not found: %s", resourceName)
        }
        gotContent := rs.Primary.Attributes["content"]
        err = yaml.Unmarshal([]byte(gotContent), &gotRg)
        if err != nil {
            return err
        }

        err = rules.CompareGroups(*gotRg, *wantRg)
        if err != nil {
            return fmt.Errorf("unexpected rule group content: %w", err)
        }
        return nil
    }
}

func testAccCheckRuleGroupExists(resourceName string, ruleGroup *rwrulefmt.RuleGroup) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        rs, ok := s.RootModule().Resources[resourceName]
        if !ok {
            return fmt.Errorf("not found: %s", resourceName)
        }

        if rs.Primary.ID == "" {
            return fmt.Errorf("resource ID is not set")
        }

        if rs.Primary.Attributes["namespace"] == "" {
            return fmt.Errorf("resource namespace is not set")
        }
        namespace := rs.Primary.Attributes["namespace"]

        split := strings.Split(rs.Primary.ID, "/")
        if len(split) != 3 {
            return fmt.Errorf("unexpected resource ID format")
        }
        groupName := split[0]

        cortexClient, err := client.New(client.Config{
            Address: "http://localhost:8080",
        })
        if err != nil {
            return err
        }

        rg, err := cortexClient.GetRuleGroup(context.Background(), namespace, groupName)
        if err != nil {
            return err
        }
        *ruleGroup = *rg

        return nil
    }
}

func testAccCheckRuleGroupAttrContent(gotRg *rwrulefmt.RuleGroup, wantContent string) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        var wantRg *rwrulefmt.RuleGroup
        err := yaml.Unmarshal([]byte(wantContent), &wantRg)
        if err != nil {
            return err
        }

        err = rules.CompareGroups(*gotRg, *wantRg)
        if err != nil {
            return fmt.Errorf("unexpected rule group content: %w", err)
        }

        return nil
    }
}
