package office365

import (
	"context"
	"terraform-provider-office365/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLicense() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

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
			Skuid:         singleLicenseSchema["skuid"].(string),
			DisabledPlans: disabledPlanesData,
		}
		assignedLicenseArray = append(assignedLicenseArray, oneAssignedLicense)

	}
	main_license := client.License{
		AddLicenses: assignedLicenseArray,
	}

	err := c.CreateLicense(userPrincipalName, main_license)
	if err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(err)
	}
	var licenses []interface{}
	for _, v := range res.Values {
		License := make(map[string]string)
		License["skuid"] = v.Skuid
		licenses = append(licenses, License)
	}

	d.Set("license_collection", licenses)
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
