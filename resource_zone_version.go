package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/prasmussen/gandi-api/client"
	zoneVersion "github.com/prasmussen/gandi-api/domain/zone/version"
)

func resourceZoneVersion() *schema.Resource {
	return &schema.Resource{
		Create: CreateZoneVersion,
		Update: UpdateZoneVersion,
		Read:   ReadZoneVersion,
		Delete: DeleteZoneVersion,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

// getZoneClient wraps Gandi Client in Zone Resource Methods
func getZoneVersionClient(meta interface{}) *zoneVersion.Version {
	return zoneVersion.New(meta.(*client.Client))
}

// UpdateZone changes zone properties
func UpdateZoneVersion(d *schema.ResourceData, meta interface{}) error { return nil }

// CreateZone creates new zone
func CreateZoneVersion(d *schema.ResourceData, meta interface{}) error {
	client := getZoneVersionClient(meta)

	zoneBase := d.Get("base_version").(int64)
	zoneID := d.Get("zone_id").(int64)

	newZoneVersion, err := client.New(zoneID, zoneBase)

	if err != nil {
		return fmt.Errorf("Cannot create version: %s", err)
	}

	// Id is a zoneID concatenated with newZoneVersion
	d.SetId(strconv.FormatInt(int64(zoneID), 10) + strconv.FormatInt(int64(newZoneVersion), 10))
	log.Printf("[INFO] Created new version: %v of: %v", newZoneVersion, zoneID)

	return ReadZoneVersion(d, meta)
}

// ReadZone fetches configuration
func ReadZoneVersion(d *schema.ResourceData, meta interface{}) error {
	// client := getZoneVersionClient(meta)

	//Id is a name after the resource "type" "name"
	log.Printf("[DEBUG] Reading zone: %v versions", d.Id())

	// // Id is stored as string in tfstate, API expects a int64
	// ID, _ := strconv.ParseInt(d.Id(), 10, 64)
	// // Read info about the zone
	// zone, err := client.Info(ID)
	// if err != nil {
	// 	// set the name to ""
	// 	log.Printf("[DEBUG] Unable to read zone: %s. Cleaning resource reference", err)
	// 	d.SetId("")
	// 	return nil
	// }
	//
	// d.Set("name", zone.Name)

	return nil
}

// DeleteZone deletes configuration
func DeleteZoneVersion(d *schema.ResourceData, meta interface{}) error {
	// client := getZoneVersionClient(meta)

	log.Printf("[DEBUG] Deleting zone version: %v", d.Id())

	// ID, _ := strconv.ParseInt(d.Id(), 10, 64)
	// success, err := client.Delete(ID)
	// if err != nil {
	// 	return fmt.Errorf("Cannot delete: %s", err)
	// }
	//
	// if success {
	// 	log.Printf("[DEBUG] Deleted Zone: %v", d.Id())
	// 	d.SetId("")
	// }

	// return err
	return nil
}
