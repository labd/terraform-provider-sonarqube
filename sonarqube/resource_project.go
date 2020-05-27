package sonarqube

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sonargo "github.com/labd/sonargo/sonar"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Exists: resourceProjectExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceProjectExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*sonargo.Client)

	results, _, err := client.Projects.Search(&sonargo.ProjectsSearchOption{
		Projects: d.Id(),
	})
	if err != nil {
		return false, err
	}
	return len(results.Components) == 1, nil
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sonargo.Client)

	visibility := "private"
	if d.Get("public").(bool) {
		visibility = "public"
	}

	result, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
		Branch:     "",
		Name:       d.Get("name").(string),
		Project:    d.Get("key").(string),
		Visibility: visibility,
	})
	if err != nil {
		return err
	}

	d.SetId(result.Project.Key)
	d.Set("name", result.Project.Name)
	d.Set("key", result.Project.Key)
	d.Set("public", result.Project.Visibility == "public")
	return nil
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*sonargo.Client)

	results, _, err := client.Projects.Search(&sonargo.ProjectsSearchOption{
		Projects: d.Id(),
	})
	if err != nil {
		return err
	}

	if len(results.Components) == 0 {
		return fmt.Errorf("No project found with key %s", d.Id())
	}

	d.Set("name", results.Components[0].Name)
	d.Set("key", results.Components[0].Key)
	d.Set("public", results.Components[0].Visibility == "public")
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sonargo.Client)

	d.Partial(true)

	if d.HasChange("key") {
		oldValue, newValue := d.GetChange("key")
		_, err := client.Projects.UpdateKey(&sonargo.ProjectsUpdateKeyOption{
			From: oldValue.(string),
			To:   newValue.(string),
		})

		if err != nil {
			return err
		}

		d.SetPartial("key")
	}

	if d.HasChange("public") {

		visibility := "private"
		if d.Get("public").(bool) {
			visibility = "public"
		}

		_, err := client.Projects.UpdateVisibility(&sonargo.ProjectsUpdateVisibilityOption{
			Project:    d.Get("key").(string),
			Visibility: visibility,
		})

		if err != nil {
			return err
		}

		d.SetPartial("public")
	}

	d.Partial(false)

	return nil
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*sonargo.Client)
	_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
		Project: d.Get("key").(string),
	})
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}
