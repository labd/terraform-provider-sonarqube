package sonarqube

import (
	"log"

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

	_, _, err := client.Projects.Search(&sonargo.ProjectsSearchOption{
		Projects: "versie",
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sonargo.Client)

	visibility := "private"
	if d.Get("public").(bool) {
		visibility = "public"
	}

	project, resp, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
		Branch:     "",
		Name:       d.Get("name").(string),
		Project:    d.Get("key").(string),
		Visibility: visibility,
	})
	if err != nil {
		return err
	}
	log.Print(project)
	log.Print(resp)
	log.Print(err)
	d.SetId(project.Project.Key)
	d.Set("key", project.Project.Key)
	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*sonargo.Client)

	project, resp, err := client.Projects.Search(&sonargo.ProjectsSearchOption{
		Projects: d.Id(),
	})
	log.Print(project)
	log.Print(resp)
	log.Print(err)
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*sonargo.Client)
	_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
		Project: d.Id(),
	})
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}
