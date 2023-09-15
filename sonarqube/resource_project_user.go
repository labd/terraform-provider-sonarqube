package sonarqube

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	sonargo "github.com/labd/sonargo/sonar"
)

func resourceProjectUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectUserCreate,
		Read:   resourceProjectUserRead,
		Delete: resourceProjectUserDelete,
		Schema: map[string]*schema.Schema{
			"login": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permissions": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"admin",
						"codeviewer",
						"issueadmin",
						"securityhotspotadmin",
						"scan",
						"user",
					}, false),
				},
			},
		},
	}
}

func resourceProjectUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sonargo.Client)
	for _, permission := range d.Get("permissions").([]interface{}) {
		_, err := client.Permissions.AddUser(&sonargo.PermissionsAddUserOption{
			Login:      d.Get("login").(string),
			ProjectKey: d.Get("project_key").(string),
			Permission: permission.(string),
		})
		if err != nil {
			return err
		}
		d.SetId(fmt.Sprintf("%s:%s:%s", d.Get("project_key"), d.Get("login"), d.Get("permission")))
	}
	return nil
}

func resourceProjectUserRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceProjectUserDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*sonargo.Client)
	var email string = d.Get("login").(string)
	fmt.Print(email)
	result, _, err1 := client.Users.Search(&sonargo.UsersSearchOption{
		Q: d.Get("login").(string),
	})
	if err1 != nil {
		return fmt.Errorf("No user found with email address %s", d.Get("login").(string))
	}

	if len(result.Users) < 1 {
		return fmt.Errorf("No user found with email address %s", d.Get("login").(string))
	}
	for _, permission := range d.Get("permissions").([]interface{}) {
		_, err := client.Permissions.RemoveUser(&sonargo.PermissionsRemoveUserOption{
			Login:      d.Get("login").(string),
			Permission: permission.(string),
			ProjectKey: d.Get("project_key").(string),
		})
		if err != nil {
			return err
		}
	}
	d.SetId("")
	return nil
}
