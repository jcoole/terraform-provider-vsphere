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
	"context"
	"fmt"
	"log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/hostsystem"
	"github.com/vmware/govmomi/object"
)

func dataSourceVSphereHostAdvancedSetting() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereHostAdvancedSettingRead,
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
		}
	}
}

// 'd' is a reference to the resource data which is what is returned.
// 'meta' is the TF specific connection
func dataSourceVSphereHostAdvancedSettingRead(d *schema.ResourceData, meta interface{}) error {
	
	/* this is the code from govmomi around optionmanager
	func NewOptionManager(c *vim25.Client, ref types.ManagedObjectReference) *OptionManager {
		return &OptionManager{
			Common: NewCommon(c, ref),
		}
	}
	// 'Query' is not a standalone function.
	// The parenthesis before it indicates Query is a method of OptionManager types.
	func (m OptionManager) Query(ctx context.Context, name string) ([]types.BaseOptionValue, error) {
		req := types.QueryOptions{
			This: m.Reference(),
			Name: name,
		}

		res, err := methods.QueryOptions(ctx, m.Client(), &req)
		if err != nil {
			return nil, err
		}

		return res.Returnval, nil
	}

	func (m OptionManager) Update(ctx context.Context, value []types.BaseOptionValue) error {
		req := types.UpdateOptions{
			This:         m.Reference(),
			ChangedValue: value,
		}

		_, err := methods.UpdateOptions(ctx, m.Client(), &req)
		return err
	}
*/

/* Example MorefID lookup
// FromID locates a HostSystem by its managed object reference ID.
func FromID(client *govmomi.Client, id string) (*object.HostSystem, error) {
	log.Printf("[DEBUG] Locating host system ID %s", id)
	
	// create a new finder
	finder := find.NewFinder(client.Client, false)
	// construct the MOREF
	ref := types.ManagedObjectReference{
		Type:  "HostSystem",
		Value: id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	// Reference searched and found.
	hs, err := finder.ObjectReference(ctx, ref)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Host system found: %s", hs.Reference().Value)
	// Unsure why hs.HostSystem is returned as a pointer.
	return hs.(*object.HostSystem), nil
}
*/
	
	
	// Get the VIM client from meta to do querying
	log.Printf("[DEBUG] dataSourceVSphereHostAdvancedSettingRead :: Beginning Advanced Setting lookup on host [%s] with key [%s]", d.Get("host_id").(string), d.Get("name").(string))
	client := meta.(*Client).vimClient

	// Begin Context. Consider moving this to a helper
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()

	// Get the Host object using the lookup helper function
	log.Printf("[DEBUG] dataSourceVSphereHostAdvancedSettingRead :: Beginning HostSystem lookup with host ID [%s]", d.Get("host_id").(string))
	host, err := hostsystem.FromID(client, d.Get("host_id").(string))
	log.Printf("[DEBUG] dataSourceVSphereHostAdvancedSettingRead :: HostSystem lookup with host ID [%s] successful, name is [%s]", d.Get("host_id").(string),host.Name)
	
	// Get the Host properties
	log.Printf("[DEBUG] dataSourceVSphereHostAdvancedSettingRead :: Beginning HostSystem Property lookup with host [%s]", host.Name)
	hostObject, err := Properties(host)
	log.Printf("[DEBUG] dataSourceVSphereHostAdvancedSettingRead :: Completed properties lookup on host [%s], power state property is [%s]", host.Name, hostObject.Runtime.PowerState)
	
	// Get Host Config Manager and subsequently Option Manager
	log.Printf("[DEBUG] dataSourceVSphereHostAdvancedSettingRead :: Beginning HostSystem OptionManagerlookup with host ID [%s]", d.Get("host_id").(string))
	om := host.NewHostConfigManager(client.Client).OptionManager
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] dataSourceVSphereHostAdvancedSettingRead :: Completed OptionManager lookup - got moref [%s]", om.Reference())
	// Get the key and its value
	// note: 'Query' is a function defined in go as a wrapper around 'queryOptions' in the api
	// note: the method to update is 'Update' which is wrapped around 'updateOptions' in the api
	log.Printf("[DEBUG] dataSourceVSphereHostAdvancedSettingRead :: Begin Key lookup of [%s]", d.Get("name").(string))
	qv, err := om.Query(ctx, d.Get("name").(string))
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] dataSourceVSphereHostAdvancedSettingRead :: Finished Key lookup of [%s], value is [%s]", d.Get("name").(string), qv.Value)
	
	// When you execute the query, it returns an array of {key, value} for the options.
	d.SetId(fmt.Sprint(qv.Key))
	// Set the value that currently exists.
	_ = d.Set("value", qv.Value)

	// TODO: Parse methods to determine the type? Assume string is default?
	return nil
}