package sonarqube

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sonargo "github.com/labd/sonargo/sonar"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,

		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Computed values.
			"login": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sonargo.Client)

	result, _, err := client.Users.Search(&sonargo.UsersSearchOption{
		Q: d.Get("email").(string),
	})
	if err != nil {
		return err
	}

	if len(result.Users) < 1 {
		return fmt.Errorf("No user found with email address %s", d.Get("email").(string))
	}
	d.SetId(result.Users[0].Login)
	d.Set("login", result.Users[0].Login)
	return nil
}
