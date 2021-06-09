package office365

import (
	"context"
	"log"
	"strings"
	"terraform-provider-office365/client"

	val "terraform-provider-office365/validate"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"account_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"force_change_password_nextsignin": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"mail_nick_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_principal_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: val.ValidateEmail,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"job_title": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mail": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"office_location": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mobile_phone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"preferred_language": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"postal_code": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"surname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"given_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"street_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"usage_location": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"city": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"department": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"country": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"object_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"proxy_addresses": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	Id := d.Get("user_principal_name").(string)
	password := d.Get("password").(string)
	if password == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Password can not be empty",
			Detail:   "password can not be emtpy",
		})
		return diags
	}
	mailNickName := d.Get("mail_nick_name").(string)
	if mailNickName == "" {
		mailNickName = strings.Split(Id, "@")[0]
	}
	mail := d.Get("mail").(string)
	if mail == "" {
		mail = d.Get("user_principal_name").(string)
	}
	req_json := client.CreatUser{
		AccountEnabled:    d.Get("account_enabled").(bool),
		DisplayName:       d.Get("display_name").(string),
		MailNickName:      mailNickName,
		UserPrincipalName: d.Get("user_principal_name").(string),
		PasswordProfile: client.PasswordProfileModel{
			ForceChangePasswordNextSignIn: d.Get("force_change_password_nextsignin").(bool),
			Password:                      d.Get("password").(string),
		},
		GivenName:         d.Get("given_name").(string),
		JobTitle:          d.Get("job_title").(string),
		OfficeLocation:    d.Get("office_location").(string),
		MobilePhone:       d.Get("mobile_phone").(string),
		PreferredLanguage: d.Get("preferred_language").(string),
		Surname:           d.Get("surname").(string),
		State:             d.Get("state").(string),
		StreetAddress:     d.Get("street_address").(string),
		UsageLocation:     d.Get("usage_location").(string),
		PostalCode:        d.Get("postal_code").(string),
		City:              d.Get("city").(string),
		Country:           d.Get("country").(string),
		Department:        d.Get("department").(string),
		Mail:              mail,
	}
	_, er := c.CreateUser(req_json)
	if er != nil {
		log.Println("[ERROR]: ", er)
		return diag.FromErr(er)
	}
	d.SetId(Id)
	resourceUserRead(ctx, d, m)
	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	UserInfo, err := c.GetUser(d.Id())
	if err != nil {
		log.Println("[ERROR]: ", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get User Information",
			Detail:   err.Error(),
		})
		return diags
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
	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	if d.HasChange("user_principal_name") {
		d.SetId("")
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't update Principal Name",
			Detail:   "Can't update Principal Name",
		})
		return diags
	}
	if d.HasChange("password") {
		d.SetId("")
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't update password",
			Detail:   "Can't update Password",
		})
		return diags
	}
	req_json := client.UpdateUser{
		MailNickName:      d.Get("mail_nick_name").(string),
		DisplayName:       d.Get("display_name").(string),
		GivenName:         d.Get("given_name").(string),
		JobTitle:          d.Get("job_title").(string),
		Mail:              d.Get("mail").(string),
		PreferredLanguage: d.Get("preferred_language").(string),
		Surname:           d.Get("surname").(string),
		MobilePhone:       d.Get("mobile_phone").(string),
		OfficeLocation:    d.Get("office_location").(string),
		AccountEnabled:    d.Get("account_enabled").(bool),
		StreetAddress:     d.Get("street_address").(string),
		PostalCode:        d.Get("postal_code").(string),
		City:              d.Get("city").(string),
		Department:        d.Get("department").(string),
		Country:           d.Get("country").(string),
		State:             d.Get("state").(string),
	}
	errr := c.UpdateUser(d.Id(), req_json)
	if errr != nil {
		log.Println("[ERROR]: ", errr)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "!Update failed",
			Detail:   errr.Error(),
		})
		return diags
	}
	d.Set("last_updated", time.Now().Format(time.RFC850))
	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	UserID := d.Id()
	err := c.DeleteUser(UserID)
	if err != nil {
		log.Println("[ERROR]: ", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Delete Failed",
			Detail:   err.Error(),
		})
		return diags
	}
	d.SetId("")
	return diags
}
