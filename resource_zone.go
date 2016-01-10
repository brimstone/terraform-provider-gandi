package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/prasmussen/gandi-api/client"
	"github.com/prasmussen/gandi-api/domain/zone"
)

func resourceZone() *schema.Resource {
	return &schema.Resource{
		Create: CreateZone,
		Update: UpdateZone,
		Read:   ReadZone,
		Delete: DeleteZone,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

// getZoneClient wraps Gandi Client in Zone Resource Methods
func getZoneClient(meta interface{}) *zone.Zone {
	return zone.New(meta.(*client.Client))
}

// TODO: get a function that would return a domain name the zone is associated to
// func getAssociatedDomainID(zoneName string) int64 {
// 	return 6334583 // returns bemehow.com domain id
// }

// UpdateZone changes zone properties
func UpdateZone(d *schema.ResourceData, meta interface{}) error { return nil }

// CreateZone creates new zone
func CreateZone(d *schema.ResourceData, meta interface{}) error {
	client := getZoneClient(meta)

	zone, err := client.Create(d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("Cannot create zone: %s", err)
	}

	// Assign the zone Id to a string repr of the zone.Id
	d.SetId(strconv.FormatInt(zone.Id, 10))
	log.Printf("[INFO] Created zone with ID: %v", zone.Id)

	// TODO: make the association happen here (under create) so it can be read later with ReadZone

	return ReadZone(d, meta)
}

// ReadZone fetches configuration
func ReadZone(d *schema.ResourceData, meta interface{}) error {
	client := getZoneClient(meta)

	//Id is a name after the resource "type" "name"
	log.Printf("[DEBUG] Reading zone: %v", d.Id())

	// Id is stored as string in tfstate, API expects a int64
	ID, _ := strconv.ParseInt(d.Id(), 10, 64)
	// Read info about the zone
	zone, err := client.Info(ID)
	if err != nil {
		// set the name to ""
		log.Printf("[DEBUG] Unable to read zone: %s. Cleaning resource reference", err)
		d.SetId("")
		return nil
	}

	d.Set("name", zone.Name)
	//TODO: figure out how to fetch the information which domain the zone is attached to
	// d.Set("domain_id", getAssociatedDomainID(zone.Name))

	return nil
}

// DeleteZone deletes configuration
func DeleteZone(d *schema.ResourceData, meta interface{}) error {
	client := getZoneClient(meta)

	log.Printf("[DEBUG] Deleting zone: %v", d.Id())

	ID, _ := strconv.ParseInt(d.Id(), 10, 64)
	success, err := client.Delete(ID)
	if err != nil {
		return fmt.Errorf("Cannot delete zone: %s", err)
	}

	if success {
		log.Printf("[DEBUG] Deleted zone: %v", d.Id())
		d.SetId("")
	}

	return err
}
