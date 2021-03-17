package cortex

import (
	"context"
	"fmt"
	"strings"

	"github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"gopkg.in/yaml.v3"
)

func resourceRules() *schema.Resource {
	return &schema.Resource{
		Description:   "This resource enables you to manage rule groups inside a cortex cluster.",
		CreateContext: resourceRulesCreate,
		ReadContext:   resourceRulesRead,
		// Updates use the same API as create.
		UpdateContext: resourceRulesCreate,
		DeleteContext: resourceRulesDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the rule group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"tenant_id": &schema.Schema{
				Description: "Tenant ID, passed as X-Org-ScopeID HTTP header. If empty, the provider tenant ID is used.",
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeString,
			},
			"namespace": &schema.Schema{
				Description: "The namespace of the rule group.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"content": &schema.Schema{
				Description:      "Rule group content.",
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressYAMLDiff,
			},
		},
	}
}

func resourceRulesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var (
		c         = m.(cortexClientFunc)
		namespace = d.Get("namespace").(string)
		content   = d.Get("content").(string)
		tenantID  string
		diags     diag.Diagnostics
	)

	if data, ok := d.GetOk("tenant_id"); ok {
		tenantID = data.(string)
	}

	client, err := c(tenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	var rg rulefmt.RuleGroup
	err = yaml.Unmarshal([]byte(content), &rg)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleGroup := rwrulefmt.RuleGroup{
		RuleGroup: rg,
	}

	err = client.CreateRuleGroup(ctx, namespace, ruleGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", rg.Name, namespace, tenantID))

	return diags
}

func resourceRulesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var (
		c         = m.(cortexClientFunc)
		namespace = d.Get("namespace").(string)
		groupName = strings.Split(d.Id(), "/")[0]
		tenantID  string
		diags     diag.Diagnostics
	)

	if data, ok := d.GetOk("tenant_id"); ok {
		tenantID = data.(string)
	}

	client, err := c(tenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteRuleGroup(ctx, namespace, groupName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourceRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var (
		c         = m.(cortexClientFunc)
		namespace = d.Get("namespace").(string)
		groupName = strings.Split(d.Id(), "/")[0]
		tenantID  string
		diags     diag.Diagnostics
	)

	if data, ok := d.GetOk("tenant_id"); ok {
		tenantID = data.(string)
	}

	client, err := c(tenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	rg, err := client.GetRuleGroup(ctx, namespace, groupName)
	if err != nil {
		return diag.FromErr(err)
	}

	out, err := yaml.Marshal(rg)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("content", string(out))

	return diags
}
