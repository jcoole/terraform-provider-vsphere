// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vsphere

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/hostsystem"
	"github.com/vmware/govmomi/object"
)

// Resource CRUD and definition.
// Host Advanced Settings can only be Read/Updated
func resourceVSphereHostAdvancedSetting() *schema.Resource {
	return &schema.Resource{
		//Create: resourceVSphereCustomAttributeCreate,
		Read:   resourceVSphereHostAdvancedSettingRead,
		Update: resourceVSphereHostAdvancedSettingUpdate,
		//Delete: resourceVSphereCustomAttributeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereHostAdvancedSettingImport,
		},

		Schema: map[string]*schema.Schema{
			"host_id": {
				Type: schema.TypeString,
				Description: "The Managed Object ID of the Host System.",
				Required: true
			},
			"name": {
				Type: schema.TypeString,
				Description: "The name of the Advanced Setting.",
				Required: true
			},
			"value": {
				Type: schema.TypeString,
				Description: "The value of the Advanced Setting.",
				Required: true
			},
			"type": {
				Type: schema.TypeString,
				Description: "The value type. String, Long, Bool, Int, Choice."
			}
		},
	}
}

// Read function.
// Basically the same as the data source, aside from context related stuff
func resourceVSphereHostAdvancedSettingRead(d *schema.ResourceData, meta interface{}) error {
	// Get the VIM client from meta to do querying
	client := meta.(*Client).vimClient

	// Instantiate HostConfigManager for the host to query settings, error out otherwise
	om := object.NewHostConfigManager(client.Client)
	if err != nil {
		return err
	}
	
	// Context stuff.
	// Context is used for i/o with apis so that they can run concurrently, in this case background within the api timeout.
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	// Query settings for the name/key and its value
	qv, err := om.Query(d.Get("name").(string))
	if err != nil {
		return err
	}

	// When you execute the query, it returns an array of {key, value} for the options.
	d.SetId(fmt.Sprint(qv.Key))

	// Set the values that currently exists.
	_ = d.Set("key", qv.Key)
	_ = d.Set("value", qv.Value)
	// todo: type logic. make a helper module?
	return nil
}

// Update function
func resourceVSphereHostAdvancedSettingUpdate(d *schema.ResourceData, meta interface{}) error {
	// For update to function, probably a good idea to add some validation for the type.
	// NYI
	type, err := d.GetHostAdvancedSettingType(d.Get("type").(string))
	if err != nil {
		return err
	}
	// Get the VIM client from meta to do querying
	client := meta.(*Client).vimClient

	// Instantiate HostConfigManager for the host to update settings, error out otherwise
	om := object.NewHostConfigManager(client.Client)
	if err != nil {
		return err
	}

	// Context stuff.
	// Context is used for i/o with apis so that they can run concurrently, in this case background within the api timeout.
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	// Update the value.
	return om.Update(ctx, d.Get("value").(string))

	qv, err := om.Update(ctx, d.Get("value").(string))
	if err != nil {
		return err
	}

}

// Import function
func resourceVSphereHostAdvancedSettingImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	
}

/* OLD REF BELOW
func resourceVSphereCustomAttribute() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereCustomAttributeCreate,
		Read:   resourceVSphereCustomAttributeRead,
		Update: resourceVSphereCustomAttributeUpdate,
		Delete: resourceVSphereCustomAttributeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereCustomAttributeImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The display name of the custom attribute.",
				Required:    true,
			},
			"managed_object_type": {
				Type:        schema.TypeString,
				Description: "Object type for which the custom attribute is valid. If not specified, the attribute is valid for all managed object types.",
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceVSphereCustomAttributeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).vimClient
	if err := customattribute.VerifySupport(client); err != nil {
		return err
	}

	fm, err := object.GetCustomFieldsManager(client.Client)
	if err != nil {
		return err
	}
	key, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	fields, err := fm.Field(ctx)
	if err != nil {
		return err
	}
	field := fields.ByKey(int32(key))
	if field == nil {
		return fmt.Errorf("could not locate category with id '%d'", key)
	}
	_ = d.Set("name", field.Name)
	_ = d.Set("managed_object_type", field.ManagedObjectType)
	return nil
}

func resourceVSphereCustomAttributeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).vimClient
	if err := customattribute.VerifySupport(client); err != nil {
		return err
	}

	fm, err := object.GetCustomFieldsManager(client.Client)
	if err != nil {
		return err
	}
	key, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	return fm.Rename(ctx, int32(key), d.Get("name").(string))
}

func resourceVSphereCustomAttributeImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Client).vimClient
	if err := customattribute.VerifySupport(client); err != nil {
		return nil, err
	}

	fm, err := object.GetCustomFieldsManager(client.Client)
	if err != nil {
		return nil, err
	}

	field, err := customattribute.ByName(fm, d.Id())
	if err != nil {
		return nil, err
	}

	d.SetId(fmt.Sprint(field.Key))
	return []*schema.ResourceData{d}, nil
}
*/