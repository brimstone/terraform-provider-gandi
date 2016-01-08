package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cznic/sortutil"
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
			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"base_version": &schema.Schema{
				Type:     schema.TypeString,
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

	baseVersionNumber, _ := strconv.ParseInt(d.Get("base_version").(string), 10, 64)
	zoneID, _ := strconv.ParseInt(d.Get("zone_id").(string), 10, 64)

	zoneVersionNumber, err := client.New(zoneID, baseVersionNumber)

	if err != nil {
		return fmt.Errorf("Cannot create zone version: %s", err)
	}

	// Id is stored as compound string "zoneID@baseVersionNumber@zoneVersion"
	ID := strconv.FormatInt(int64(zoneID), 10) + "|" + strconv.FormatInt(int64(baseVersionNumber), 10)
	ID = ID + "|" + strconv.FormatInt(int64(zoneVersionNumber), 10)
	d.SetId(ID)
	log.Printf("[INFO] Created new version: %v of zone: %v", zoneVersionNumber, zoneID)

	return ReadZoneVersion(d, meta)
}

// decode zoneID and versions from the custom resource ID
func extractIDs(id string, separator string) (int64, int64, int64) {
	zoneID, _ := strconv.ParseInt(strings.Split(id, separator)[0], 10, 64)
	zoneVersionNumber, _ := strconv.ParseInt(strings.Split(id, separator)[1], 10, 64)
	baseVersionNumber, _ := strconv.ParseInt(strings.Split(id, separator)[2], 10, 64)

	return zoneID, zoneVersionNumber, baseVersionNumber
}

// lookup version by Id in the list of configured versions
func zoneVersionNumberExist(versions []*zoneVersion.VersionInfo, zoneVersionNumber int64) int {
	var zoneVersionNumbers sortutil.Int64Slice
	for _, v := range versions {
		zoneVersionNumbers = append(zoneVersionNumbers, v.Id)
	}

	// sort the list before lookup
	zoneVersionNumbers.Sort()
	i := sortutil.SearchInt64s(zoneVersionNumbers, zoneVersionNumber)
	if i < len(zoneVersionNumbers) && zoneVersionNumbers[i] == zoneVersionNumber {
		log.Print("[DEBUG] Zone version found")
		return i
	}
	log.Print("[DEBUG] Zone version not found")
	return -1
}

// ReadZone fetches configuration
func ReadZoneVersion(d *schema.ResourceData, meta interface{}) error {
	client := getZoneVersionClient(meta)

	// Parse out version numbers from the resource ID
	zoneID, _, baseVersionNumber := extractIDs(d.Id(), "|")

	log.Printf("[DEBUG] Reading zone: %v versions", zoneID)
	versions, err := client.List(zoneID)

	if err != nil {
		return fmt.Errorf("Cannot read zone version: %v", d.Id())
	}

	i := zoneVersionNumberExist(versions, zoneID) //index of the Record
	if i < 0 {
		log.Printf("[DEBUG] Zone version with ID: %v not found. Cleaning local state reference", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("base_version", baseVersionNumber)
	d.Set("zone_id", zoneID)

	return nil
}

// DeleteZone deletes configuration
func DeleteZoneVersion(d *schema.ResourceData, meta interface{}) error {
	client := getZoneVersionClient(meta)

	log.Printf("[DEBUG] Deleting zone version: %v", d.Id())

	// Parse out version numbers from the resource ID
	zoneID, zoneVersionNumber, _ := extractIDs(d.Id(), "|")

	log.Printf("[DEBUG] Deleting zone version: %v", d.Id())
	success, err := client.Delete(zoneID, zoneVersionNumber)
	if err != nil {
		return fmt.Errorf("Cannot delete: %v", err)
	}

	if success {
		log.Printf("[DEBUG] Deleted zone version: %v", d.Id())
		d.SetId("")
	}

	return nil
}
