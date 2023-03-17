// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vsphere

/*
Notes:
	Path to settings at host level -- HostSystem.ConfigManager.AdvancedOption (moref is OptionManager-EsxHostAdvSettings-xxxx)
	Look for 'SupportedOption' to get the list of settings.
		Filter out the few that have the 'ValueIsReadOnly' flag set in their respective 'OptionType'
		Add checks for the 'Choice' option maybe, or just return the error.
		To see what the 'valid' values are on choice types, get the .OptionType.ChoiceInfo.Key list.
*/
import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/hostsystem"
	"github.com/vmware/govmomi/object"
)

func dataSourceVSphereHostAdvancedSetting() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereHostAdvancedSettingRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
				Description: "The name of the Advanced Setting.",
				Required: true,
			},
			// TBD: Should value be a jsonencode due to option choices, or use Go to infer type based on value?
			"value": {
				Type: schema.TypeString,
				Description: "The value of the Advanced Setting.",
				Required: true
			}
		}
	}
}

// 'd' is a reference to the resource data which is what is returned.
// 'meta' is the TF specific connection
func dataSourceVSphereHostAdvancedSettingRead(d *schema.ResourceData, meta interface{}) error {
	// Get the VIM client from meta to do querying
	client := meta.(*Client).vimClient

	// Instantiate HostConfigManager for the host, then get optionManager
	// note: 'om' is a go wrapper around the api, so methods use slightly different names
	om := object.NewHostConfigManager(client.Client)
	if err != nil {
		return err
	}
	// Get the key and its value
	// note: 'Query' is a function defined in go as a wrapper around 'queryOptions' in the api
	// note: the method to update is 'Update' which is wrapped around 'updateOptions' in the api
	qv, err := om.Query(d.Get("name").(string))
	if err != nil {
		return err
	}
	// When you execute the query, it returns an array of {key, value} for the options.
	d.SetId(fmt.Sprint(qv.Key))
	// Set the value that currently exists.
	_ = d.Set("value", qv.Value)
	return nil
}