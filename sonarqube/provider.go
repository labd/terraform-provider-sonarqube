package sonarqube

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	sonargo "github.com/labd/sonargo/sonar"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SONAR_URL", nil),
				Sensitive:   true,
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SONAR_TOKEN", nil),
				Sensitive:   true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sonarqube_project":      resourceProject(),
			"sonarqube_project_user": resourceProjectUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sonarqube_user": dataSourceUser(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	token := d.Get("token").(string)
	url := d.Get("url").(string)

	client, err := sonargo.NewClientByToken(url+"/api", token)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return client, nil
}

// This is a global MutexKV for use within this plugin.
var ctMutexKV = mutexkv.NewMutexKV()
