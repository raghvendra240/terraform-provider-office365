package office365

import (
	"context"
	"encoding/json"
	"strings"
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

						"single_license": &schema.Schema{
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
				},
			},
			"remove_licenses": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
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
		singleLicense := singleLicenseSchema["single_license"].(*schema.Set).List()

		for _, u := range singleLicense {

			oneLicense := u.(map[string]interface{})
			tfDisablePlanes := oneLicense["disabled_plans"].(*schema.Set).List()
			disabledPlanesData := make([]string, len(tfDisablePlanes))
			for i, data := range tfDisablePlanes {
				disabledPlanesData[i] = data.(string)
			}

			oneAssignedLicense := client.AssignedLicenses{
				Skid:          oneLicense["skuid"].(string),
				DisabledPlans: disabledPlanesData,
			}
			assignedLicenseArray = append(assignedLicenseArray, oneAssignedLicense)

		}

	}

	tfRemoveLicense := d.Get("remove_licenses").(*schema.Set).List()
	removeLicenseData := make([]string, len(tfRemoveLicense))
	for i, data := range tfRemoveLicense {
		removeLicenseData[i] = data.(string)
	}

	main_license := client.License{
		AddLicenses:    assignedLicenseArray,
		RemoveLicenses: removeLicenseData,
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
	return diags
}

func resourceLicenseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	return diags
}

func resourceLicenseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	partsOfId := ParseId(d.Id())
	userPrincipalName := partsOfId[0]
	SkUID := partsOfId[1]

	skUidArray := make([]string, 1)
	skUidArray[0] = SkUID

	licenseStruct := client.License{
		RemoveLicenses: skUidArray,
	}
	err := c.CreateLicense(userPrincipalName, licenseStruct)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceLicenseImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(*client.Client)
	UserInfo, err := c.GetUser(d.Id())
	if err != nil {
		return nil, err
	}
	d.Set("display_name", UserInfo.DisplayName)
	d.Set("job_title", UserInfo.JobTitle)
	d.Set("mail", UserInfo.Mail)
	d.Set("user_principal_name", UserInfo.UserPrincipalName)
	d.Set("office_location", UserInfo.OfficeLocation)
	d.Set("mobile_phone", UserInfo.MobilePhone)
	d.Set("preferred_language", UserInfo.PreferredLanguage)
	d.Set("surname", UserInfo.Surname)
	d.Set("object_id", UserInfo.ObjectId)
	d.Set("given_name", UserInfo.GivenName)
	return []*schema.ResourceData{d}, nil
}

func ParseId(id string) []string {
	parts := strings.Split(id, ":")
	return parts
}
