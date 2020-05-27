package sonarqube

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sonargo "github.com/labd/sonargo/sonar"
)

func resourceSettingsValue() *schema.Resource {
	return &schema.Resource{
		Create: resourceSettingsValueCreate,
		Read:   resourceSettingsValueRead,
		Update: resourceSettingsValueUpdate,
		Delete: resourceSettingsValueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSettingsValueCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sonargo.Client)

	_, err := client.Settings.Set(&sonargo.SettingsSetOption{
		Key:   d.Get("key").(string),
		Value: d.Get("value").(string),
	})
	if err != nil {
		return err
	}

	d.SetId(d.Get("key").(string))
	return nil
}

func resourceSettingsValueRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*sonargo.Client)

	results, _, err := client.Settings.Values(&sonargo.SettingsValuesOption{
		Keys: d.Id(),
	})
	if err != nil {
		return err
	}

	if len(results.Settings) == 0 {
		return fmt.Errorf("No project found with key %s", d.Id())
	}

	d.Set("key", results.Settings[0].Key)
	d.Set("value", results.Settings[0].Value)
	return nil
}

func resourceSettingsValueUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sonargo.Client)

	if d.HasChange("key") {
		oldKey, newKey := d.GetChange("key")

		_, err := client.Settings.Reset(&sonargo.SettingsResetOption{
			Keys: oldKey.(string),
		})
		if err != nil {
			return err
		}
		d.SetId(newKey.(string))
	}

	_, err := client.Settings.Set(&sonargo.SettingsSetOption{
		Key:   d.Get("key").(string),
		Value: d.Get("value").(string),
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceSettingsValueDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*sonargo.Client)
	_, err := client.Settings.Reset(&sonargo.SettingsResetOption{
		Keys: d.Id(),
	})
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}
