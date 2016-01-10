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
	// Updates to the zone versions are theoretically possible but they involve
	// change to the ID since the version information is not available
	return nil
}

// func to create version with the specified client
func _createZoneVersion(client *zoneVersion.Version, zoneID int64, baseVersion int64, zoneVersion int64) (string, error) {
	zoneExist, err := _checkZoneVersion(client, zoneID, zoneVersion)
	if err != nil {
		return "", fmt.Errorf("Cannot check zone versions: %v", err)
	}

	if zoneExist {
		return "", fmt.Errorf("Zone version: %v already exist", zoneVersion)
	}

	// Create new version of the zone
	newZoneVersion, err := client.New(zoneID, baseVersion)

	if err != nil {
		return "", fmt.Errorf("Cannot create zone version: %s", err)
	}

	// Id is stored as compound string "zoneID_zoneVersion"
	ID := strconv.FormatInt(int64(zoneID), 10) + "_" + strconv.FormatInt(int64(newZoneVersion), 10)

	return ID, nil
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

	ID, err := _createZoneVersion(client, zoneID, baseVersion, zoneVersion)
	if err != nil {
		return nil
	}

	d.SetId(ID)
	log.Printf("[INFO] Created new zone version with ID: %v", ID)

	return ReadZoneVersion(d, meta)
}

// decode zoneID and versions from the custom resource ID
func _extractIDs(id string, separator string) (int64, int64) {
	zoneID, _ := strconv.ParseInt(strings.Split(id, separator)[0], 10, 64)
	zoneVersion, _ := strconv.ParseInt(strings.Split(id, separator)[1], 10, 64)
	// baseVersion, _ := strconv.ParseInt(strings.Split(id, separator)[2], 10, 64)

	// return zoneID, zoneVersion, baseVersion
	return zoneID, zoneVersion
}

// looks up version in the zone
func _checkZoneVersion(client *zoneVersion.Version, zoneID int64, zoneVersionNumber int64) (bool, error) {
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

// ReadZoneVersion fetches configuration
func ReadZoneVersion(d *schema.ResourceData, meta interface{}) error {
	client := getZoneVersionClient(meta)

	// Parse out version numbers from the resource ID
	zoneID, zoneVersion := _extractIDs(d.Id(), "_")

	zoneExist, err := _checkZoneVersion(client, zoneID, zoneVersion)
	if err != nil {
		return fmt.Errorf("Cannot verify if zone exist: %v", err)
	}

	if !zoneExist {
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
	zoneID, zoneVersion := _extractIDs(d.Id(), "_")

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
