package sonarqube

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sonargo "github.com/labd/sonargo/sonar"
)

func resourceProjectUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectUserCreate,
		Read:   resourceProjectUserRead,
		Delete: resourceProjectUserDelete,
		Exists: resourceProjectUserExists,
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
			"permission": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceProjectUserExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*sonargo.Client)

	_, _, err := client.Projects.Search(&sonargo.ProjectsSearchOption{
		Projects: "versie",
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func resourceProjectUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sonargo.Client)

	_, err := client.Permissions.AddUser(&sonargo.PermissionsAddUserOption{
		Login:      d.Get("login").(string),
		ProjectKey: d.Get("project_key").(string),
		Permission: d.Get("permission").(string),
	})
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s__%s", d.Get("project_key").(string), d.Get("login")))
	return nil
}

func resourceProjectUserRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceProjectUserDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*sonargo.Client)

	_, err := client.Permissions.RemoveUser(&sonargo.PermissionsRemoveUserOption{
		Login:      d.Get("login").(string),
		ProjectKey: d.Get("project_key").(string),
		Permission: d.Get("permission").(string),
	})
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
