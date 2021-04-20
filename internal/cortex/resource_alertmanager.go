package cortex

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlertmanager() *schema.Resource {
	return &schema.Resource{
		Description:   "This alermanager resource enables you to manage Alertmanger configuration in Cortex.",
		CreateContext: resourceAlertsCreate,
		ReadContext:   resourceAlertsRead,
		// Updates use the same API as create.
		UpdateContext: resourceAlertsCreate,
		DeleteContext: resourceAlertsDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the alertmanager config resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"tenant_id": &schema.Schema{
				Description: "Tenant ID, passed as X-Org-ScopeID HTTP header. If empty, the provider tenant ID is used.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"alertmanager_config": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Content of the alertmanager configuration.",
				Required:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					olds, err := formatYAML(old)
					if err != nil {
						return false
					}
					news, err := formatYAML(new)
					if err != nil {
						return false
					}
					return olds == news
				},
			},
			"template_files": &schema.Schema{
				Description: "Alert templates.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAlertsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var (
		c         = m.(cortexClientFunc)
		config    = d.Get("alertmanager_config").(string)
		templates = make(map[string]string)
		tenantID  string
		diags     diag.Diagnostics
	)

	if data, ok := d.GetOk("tenant_id"); ok {
		tenantID = data.(string)
	}
	if data, ok := d.GetOk("template_files"); ok {
		tpls := data.(map[string]interface{})
		for k, v := range tpls {
			templates[k] = v.(string)
		}
	}

	client, err := c(tenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.CreateAlertmanagerConfig(ctx, config, templates)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("alertmanager%s", tenantID))

	return diags
}

func resourceAlertsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var (
		c        = m.(cortexClientFunc)
		tenantID string
		diags    diag.Diagnostics
	)

	if data, ok := d.GetOk("tenant_id"); ok {
		tenantID = data.(string)
	}

	client, err := c(tenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteAlermanagerConfig(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourceAlertsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var (
		c        = m.(cortexClientFunc)
		tenantID string
		diags    diag.Diagnostics
	)

	if data, ok := d.GetOk("tenant_id"); ok {
		tenantID = data.(string)
	}

	client, err := c(tenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	config, templates, err := client.GetAlertmanagerConfig(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("alertmanager_config", config)
	d.Set("template_files", templates)

	return diags
}
