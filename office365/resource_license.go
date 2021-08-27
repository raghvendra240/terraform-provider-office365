package office365

import (
	"context"
	"encoding/json"
	"terraform-provider-office365/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLicense() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{

			"user_principal_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"license_collection": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disabled_plans": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"skuid": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceUserImporter,
		},
		CreateContext: resourceLicenseCreate,
		ReadContext:   resourceLicenseRead,
		UpdateContext: resourceLicenseUpdate,
		DeleteContext: resourceLicenseDelete,
	}
}

func resourceLicenseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics

	userPrincipalName := d.Get("user_principal_name").(string)
	var assignedLicenseArray []client.AssignedLicenses

	v := d.Get("license_collection")
	licenseCollection := v.(*schema.Set).List()

	for _, v := range licenseCollection {
		singleLicenseSchema := v.(map[string]interface{})

		tfDisablePlanes := singleLicenseSchema["disabled_plans"].(*schema.Set).List()
		disabledPlanesData := make([]string, len(tfDisablePlanes))
		for i, data := range tfDisablePlanes {
			disabledPlanesData[i] = data.(string)
		}

		oneAssignedLicense := client.AssignedLicenses{
			Skid:          singleLicenseSchema["skuid"].(string),
			DisabledPlans: disabledPlanesData,
		}
		assignedLicenseArray = append(assignedLicenseArray, oneAssignedLicense)

	}
	main_license := client.License{
		AddLicenses: assignedLicenseArray,
	}

	err := c.CreateLicense(userPrincipalName, main_license)
	if err != nil {
		e, _ := json.Marshal(main_license)

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  string(e),
			Detail:   err.Error(),
		})
		return diags

	}
	Id := userPrincipalName
	d.SetId(Id)
	return diags
}

func resourceLicenseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)

	res, err := c.GetLicense(d.Id())

	if err != nil {
		return diag.FromErr(&json.UnsupportedTypeError{})
	}

	out, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  string(out),
		Detail:   err.Error(),
	})

	return diags
}

func resourceLicenseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	return diags
}

func resourceLicenseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)

	userPrincipalName := d.Id()
	var SkuidCollection []string

	v := d.Get("license_collection")
	licenseCollection := v.(*schema.Set).List()
	for _, v := range licenseCollection {
		singleLicenseSchema := v.(map[string]interface{})
		SkuidCollection = append(SkuidCollection, singleLicenseSchema["skuid"].(string))

	}
	main_license := client.License{
		RemoveLicenses: SkuidCollection,
	}
	err := c.CreateLicense(userPrincipalName, main_license)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceLicenseImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	return []*schema.ResourceData{d}, nil
}
