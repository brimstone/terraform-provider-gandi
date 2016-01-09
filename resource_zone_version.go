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
			"zone_version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

// getZoneVersionClient wraps Gandi Client in Zone Resource Methods
func getZoneVersionClient(meta interface{}) *zoneVersion.Version {
	return zoneVersion.New(meta.(*client.Client))
}

// UpdateZoneVersion changes zone properties
func UpdateZoneVersion(d *schema.ResourceData, meta interface{}) error {
	// when the base zone version is changed, the version is re-created
	// WARNING: all the records added to the version of the zone are going
	// to be in a loose state.
	// TODO: this needs to be handled in the records too!

	if d.HasChange("base_version") {
		// delete zone
		// create new zone based on the new base_version
		// change the id so the records referencing this zone version can also detect the change
		// zoneID, zoneVersion, baseVersionNumber := extractIDs(d.Id(), "_")

		// Open transaction
		log.Printf("[DEBUG] Updating zone (create step)")
		successCreate := CreateZoneVersion(d, meta)
		if successCreate != nil {
			return fmt.Errorf("Cannot update zone version (create step): %v")
		}

		log.Printf("[DEBUG] Updating zone (delete step)")
		successDelete := DeleteZoneVersion(d, meta)
		if successDelete != nil {
			return fmt.Errorf("Cannot update zone version (delete step): %v")
		}
	}
	return nil
}

// CreateZoneVersion creates new zone
func CreateZoneVersion(d *schema.ResourceData, meta interface{}) error {
	client := getZoneVersionClient(meta)

	baseVersion, _ := strconv.ParseInt(d.Get("base_version").(string), 10, 64)
	zoneID, _ := strconv.ParseInt(d.Get("zone_id").(string), 10, 64)
	zoneVersion, _ := strconv.ParseInt(d.Get("zone_version").(string), 10, 64)

	// Zone Version API does not specify desired version of the zone (only the base_version)
	// It needs to be checked for consistency to allow tight tracking of local+remote states
	// of zone versioning

	zoneExist, err := checkZoneVersion(client, zoneID, zoneVersion)
	if err != nil {
		return fmt.Errorf("Cannot check zone versions: %v", err)
	}

	if zoneExist {
		return fmt.Errorf("Zone version: %v already exist", zoneVersion)
	}

	// Create new version of the zone
	newZoneVersion, err := client.New(zoneID, baseVersion)

	if err != nil {
		return fmt.Errorf("Cannot create zone version: %s", err)
	}

	// Id is stored as compound string "zoneID_baseVersionNumber_zoneVersion"
	ID := strconv.FormatInt(int64(zoneID), 10) + "_" + strconv.FormatInt(int64(baseVersion), 10)
	ID = ID + "_" + strconv.FormatInt(int64(zoneVersion), 10)
	d.SetId(ID)
	log.Printf("[INFO] Created new version: %v of zone: %v", newZoneVersion, zoneID)

	return ReadZoneVersion(d, meta)
}

// decode zoneID and versions from the custom resource ID
func extractIDs(id string, separator string) (int64, int64, int64) {
	zoneID, _ := strconv.ParseInt(strings.Split(id, separator)[0], 10, 64)
	zoneVersion, _ := strconv.ParseInt(strings.Split(id, separator)[1], 10, 64)
	baseVersion, _ := strconv.ParseInt(strings.Split(id, separator)[2], 10, 64)

	return zoneID, zoneVersion, baseVersion
}

// looks up version in the zone
func checkZoneVersion(client *zoneVersion.Version, zoneID int64, zoneVersionNumber int64) (bool, error) {
	var zoneVersionNumbers sortutil.Int64Slice

	log.Printf("[DEBUG] Reading zone: %v versions", zoneID)
	versions, err := client.List(zoneID)

	if err != nil {
		return false, fmt.Errorf("Cannot read zone version: %v", zoneID)
	}

	for _, v := range versions {
		zoneVersionNumbers = append(zoneVersionNumbers, v.Id)
	}

	// sort the list before lookup
	zoneVersionNumbers.Sort()

	i := sortutil.SearchInt64s(zoneVersionNumbers, zoneVersionNumber)

	if i < len(zoneVersionNumbers) && zoneVersionNumbers[i] == zoneVersionNumber {
		log.Print("[DEBUG] Zone version found")
		return true, nil
	}
	log.Print("[DEBUG] Zone version not found")
	return false, nil
}

// ReadZone fetches configuration
func ReadZoneVersion(d *schema.ResourceData, meta interface{}) error {
	client := getZoneVersionClient(meta)

	// Parse out version numbers from the resource ID
	zoneID, zoneVersion, baseVersionNumber := extractIDs(d.Id(), "_")

	zoneExist, err := checkZoneVersion(client, zoneID, zoneVersion)
	if err != nil {
		return fmt.Errorf("Cannot check zone version: %v", err)
	}

	if zoneExist {
		d.Set("base_version", baseVersionNumber)
		d.Set("zone_id", zoneID)
	} else {
		log.Printf("[DEBUG] Zone version with ID: %v not found. Cleaning local state reference", d.Id())
		d.SetId("")
	}

	return nil
}

// DeleteZone deletes configuration
func DeleteZoneVersion(d *schema.ResourceData, meta interface{}) error {
	client := getZoneVersionClient(meta)

	log.Printf("[DEBUG] Deleting zone version: %v", d.Id())

	// Parse out version numbers from the resource ID
	zoneID, zoneVersion, baseVersionNumber := extractIDs(d.Id(), "_")

	log.Printf("[DEBUG] Deleting zone version: %v", d.Id())
	success, err := client.Delete(zoneID, zoneVersion)
	if err != nil {
		return fmt.Errorf("Cannot delete: %v", err)
	}

	if success {
		log.Printf("[DEBUG] Deleted zone version: %v", d.Id())
		d.SetId("")
	}

	return nil
}
