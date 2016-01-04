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
				Type:     schema.TypeString,
				Optional: true,
				Default:  "_default",
			},
			// "run_list": &schema.Schema{
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Elem: &schema.Schema{
			// 		Type:      schema.TypeString,
			// 		StateFunc: runListEntryStateFunc,
			// 	},
			// },
		},
	}
}

// getZoneClient wraps Gandi Client in Zone Resource Methods
func getZoneClient(meta interface{}) *zone.Zone {
	return zone.New(meta.(*client.Client))
}

// TODO: need a function that would return the newest possible version based on the domain name

// TODO: get a function that would return a domain name the zone is associated to
func getAssociatedDomainID(zoneName string) int64 {
	return 6334583 // returns bemehow.com domain id
}

// UpdateZone changes zone properties
func UpdateZone(d *schema.ResourceData, meta interface{}) error { return nil }

// CreateZone creates new zone
func CreateZone(d *schema.ResourceData, meta interface{}) error {
	client := getZoneClient(meta)

	zone, err := client.Create(d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("Cannot create: %s", err)
	}

	// Assign the zone Id to a string repr of the zone.Id
	d.SetId(strconv.FormatInt(zone.Id, 10))
	log.Printf("[INFO] Zone ID: %s", zone.Id)

	// TODO: make the association happen here (under create)

	return ReadZone(d, meta)
}

// ReadZone fetches configuration
func ReadZone(d *schema.ResourceData, meta interface{}) error {
	client := getZoneClient(meta)

	//Id is a name after the resource "type" "name"
	log.Printf("[INFO] Reading Zone ID: %v", d.Id())

	// HACK: for now this is a hacks
	// need to mixin a client here to so it can perform a list and then fetch stuff
	// func getZoneID(zoneName string) int64 {
	// 	return 1762856 //bemehow.com zone id
	// }
	// zoneAPI := zone.New(c)
	// var zones []*zone.ZoneInfoBase
	// zones, _ = zoneAPI.List()
	// for _, z := range zones {
	// 	fmt.Println(z)
	// 	zinfo, _ := zoneAPI.Info(z.Id)
	// 	fmt.Println(zinfo.Domains, zinfo.Id, zinfo.Versions, zinfo.Name)
	// 	zzinfo := reflect.ValueOf(zinfo).Elem()
	// 	zType := zzinfo.Type()
	// 	for i := 0; i < zzinfo.NumField(); i++ {
	// 		fmt.Printf("%v\n", zType.Field(i))
	// 	}
	// }
	// enumerate all zones and build look-up table
	// zoneNames := make(map[string]int64)
	// zones, err := client.List()
	// if err != nil {
	// 	return fmt.Errorf("Unable to lookup zone Id: %s. Cannot continue", err)
	// 	// not updating status, it does not mean the zone does not exist out there
	// }
	// for _, z := range zones {
	// 	zoneNames[z.Name] = z.Id
	// }
	// // TODO: make it print nicer and skip map[string]int64
	// log.Printf("[DEBUG] Found Zones: %#v", zoneNames)

	//get zone.ZoneInfoBase
	// Id is stored as string in tfstate, API expects a int64
	ID, _ := strconv.ParseInt(d.Id(), 10, 64)
	zone, err := client.Info(ID)
	if err != nil {
		//TODO: add deletion for the zone that does not existCreateZone
		// set the name to ""
		d.SetId("")
		return fmt.Errorf("Unable to read zone: %s", err)
	}

	d.Set("name", zone.Name)
	d.Set("domain_id", getAssociatedDomainID(zone.Name))

	return nil
}

// DeleteZone deletes configuration
func DeleteZone(d *schema.ResourceData, meta interface{}) error {
	client := getZoneClient(meta)

	log.Printf("[DEBUG] Deleting Zone: %v", d.Id())

	ID, _ := strconv.ParseInt(d.Id(), 10, 64)
	success, err := client.Delete(ID)
	if err != nil {
		return fmt.Errorf("Cannot delete: %s", err)
	}

	if success {
		log.Printf("[DEBUG] Deleted Zone: %v", d.Id())
		d.SetId("")
	}

	return err
}
