package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/prasmussen/gandi-api/client"
	"github.com/prasmussen/gandi-api/domain/zone/record"
)

func resourceRecord() *schema.Resource {
	return &schema.Resource{
		Create: CreateRecord,
		Update: UpdateRecord,
		Read:   ReadRecord,
		Delete: DeleteRecord,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": &schema.Schema{
				Type:     schema.TypeString, // needs to be string cause API uses int64
				Required: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// getRecordClient wraps Gandi Client in Record Resource Methods
func getRecordClient(meta interface{}) *record.Record {
	return record.New(meta.(*client.Client))
}

// ZoneRecord
type ZoneRecord struct {
	record.RecordInfo
	Id      int64
	Zone    int64
	Version int64
}

func (zr *ZoneRecord) Parse(d *schema.ResourceData) error {
	//TODO: do not ignore parsing errors
	zr.Zone, _ = strconv.ParseInt(d.Get("zone_id").(string), 10, 64)
	zr.Ttl, _ = strconv.ParseInt(d.Get("ttl").(string), 10, 64)
	zr.Version, _ = strconv.ParseInt(d.Get("version").(string), 10, 64)
	zr.Id, _ = strconv.ParseInt(d.Id(), 10, 64)

	zr.Name = d.Get("name").(string)
	zr.Value = d.Get("value").(string)
	zr.Type = d.Get("type").(string)

	return nil
}

func (zr *ZoneRecord) toRecordAdd() record.RecordAdd {
	return record.RecordAdd{
		Zone:    zr.Zone,
		Ttl:     zr.Ttl,
		Version: zr.Version,
		Name:    zr.Name,
		Value:   zr.Value,
		Type:    zr.Type,
	}
}

func (zr *ZoneRecord) toRecordUpdate() record.RecordUpdate {
	return record.RecordUpdate{
		Zone:    zr.Zone,
		Ttl:     zr.Ttl,
		Version: zr.Version,
		Name:    zr.Name,
		Value:   zr.Value,
		Type:    zr.Type,
		Id:      zr.Id,
	}
}

// CreateRecord creates new record
func CreateRecord(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Entering CreateRecord")
	var err error
	var activeVersion int64
	client := getRecordClient(meta)

	var zr ZoneRecord
	zr.Parse(d)
	if zr.Version == 0 {
		log.Printf("[DEBUG] Looking for active zone version")
		version := getZoneVersionClient(meta)
		_, activeVersion, err = getActiveZoneVersion(meta, zr.Zone)
		newZoneVersion, err := createZoneVersion(version, zr.Zone, activeVersion, 0)
		if err != nil {
			return fmt.Errorf("Could not create new version for record: %v", err)
		}
		_, zr.Version = resourceIDSplit(newZoneVersion, "_")
	}

	log.Printf("[DEBUG] Creating new record from spec: %+v", zr)

	newRecord, err := client.Add(zr.toRecordAdd())
	if err != nil {
		return fmt.Errorf("Could not create new record: %v", err)
	}

	// Success
	d.SetId(strconv.FormatInt(newRecord.Id, 10))
	log.Printf("[INFO] Successfully created record: %v", d.Id())
	log.Printf("[INFO] Active zone version: %v New Zone Version: %v", activeVersion, zr.Version)
	if activeVersion != 0 {
		setActiveZoneVersion(meta, zr.Zone, zr.Version)
	}

	return ReadRecord(d, meta)
}

func getActiveZoneVersion(meta interface{}, zoneID int64) (string, int64, error) {
	zone := getZoneClient(meta)
	zoneInfo, err := zone.Info(zoneID)
	if err != nil {
		return "", 0, fmt.Errorf("Cannot get zone info for zone id: %d, %v", zoneID, err)
	}
	zoneVersion := strconv.FormatInt(zoneInfo.Version, 10)
	return zoneVersion, zoneInfo.Version, nil
}

func setActiveZoneVersion(meta interface{}, zoneID int64, version int64) error {
	client := getZoneVersionClient(meta)
	_, err := client.Set(zoneID, version)
	return err
}

// GetRecord returns record if exist in specified zone/version
func GetRecord(client *record.Record, zoneID interface{}, zoneVersion interface{}, recordID interface{}) (*record.RecordInfo, error) {
	var zid, zv, rid int64
	zid, _ = strconv.ParseInt(zoneID.(string), 10, 64)
	zv, _ = strconv.ParseInt(zoneVersion.(string), 10, 64)
	rid, _ = strconv.ParseInt(recordID.(string), 10, 64)

	records, err := client.List(zid, zv)

	if err != nil {
		return nil, fmt.Errorf("Cannot read record: %v", rid)
	}

	// TODO: need to implement this to be sorted to improve speed
	for _, r := range records {
		if r.Id == rid {
			log.Printf("[DEBUG] Record found: %v", rid)
			return r, nil
		}
	}

	// not found
	return nil, fmt.Errorf("Record not found")
}

// CheckRecord returns boolean value for record existence
func CheckRecord(client *record.Record, zoneID interface{}, zoneVersion interface{}, recordID interface{}) (bool, error) {
	record, err := GetRecord(client, zoneID, zoneVersion, recordID)
	if err != nil {
		return false, err
	}

	if record != nil {
		return true, nil
	}

	return false, nil
}

// ReadRecord fetches configuration
func ReadRecord(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Entering ReadRecord")
	var err error
	client := getRecordClient(meta)

	zoneID := d.Get("zone_id")
	zoneVersion := d.Get("version")
	// if the zoneVersion is nil, get the active version for the zone
	if zoneVersion == "" {
		log.Printf("[DEBUG] Looking for active version of zone %v", zoneID)
		zid, _ := strconv.ParseInt(zoneID.(string), 10, 64)
		zoneVersion, _, err = getActiveZoneVersion(meta, zid)
		log.Printf("[DEBUG] Found active version of zone %v: %#v", zoneID, zoneVersion)
	}
	recordID := d.Id()

	log.Printf("[DEBUG] Reading records from zone: %v version: %v", zoneID, zoneVersion)

	record, err := GetRecord(client, zoneID, zoneVersion, recordID)
	log.Printf("[DEBUG] %#v", record)

	if err != nil {
		if strings.Contains(err.Error(), "Record not found") {
			// not found
			log.Printf("[DEBUG] Deleting record from tfstate: %v", recordID)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find record: %s", err)
	}

	if record != nil {
		// XXX: Gandi quotes values for SRV and TXT records. They need to be unquoted for comparision
		value, err := strconv.Unquote(record.Value)
		// Cannot unquote, no quotes use as is
		if err != nil {
			value = record.Value
		}
		d.Set("value", value)
		d.Set("name", record.Name)
		d.Set("ttl", strconv.FormatInt(record.Ttl, 10))
		d.Set("type", record.Type)
	}

	return nil
}

// UpdateRecord updates record in zone/version according to the new spec
func UpdateRecord(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Entering UpdateRecord")
	var err error
	var activeVersion int64
	client := getRecordClient(meta)

	var zr ZoneRecord
	zr.Parse(d)
	log.Printf("[DEBUG] FINDME ZoneRecord: %#v", zr)
	if zr.Version == 0 {
		log.Printf("[DEBUG] Looking for active zone version")
		version := getZoneVersionClient(meta)
		_, activeVersion, err = getActiveZoneVersion(meta, zr.Zone)
		newZoneVersion, err := createZoneVersion(version, zr.Zone, activeVersion, 0)
		if err != nil {
			return fmt.Errorf("Could not create new version for record: %v", err)
		}
		_, zr.Version = resourceIDSplit(newZoneVersion, "_")
		// Find old record in active version by name, type, value
		oldRecords, err := client.List(zr.Zone, activeVersion)
		var oldRecord *record.RecordInfo
		for _, r := range oldRecords {
			if r.Id == zr.Id {
				oldRecord = r
				break
			}
		}
		// FIXME check for nil oldRecord
		newRecords, err := client.List(zr.Zone, zr.Version)
		for _, r := range newRecords {
			// Find ID in new version by name, type, value
			if r.Name == oldRecord.Name && r.Type == oldRecord.Type && r.Value == oldRecord.Value {
				// Fix ID of current zr record
				zr.Id = r.Id
				break
			}
		}
	}

	log.Printf("[DEBUG] Updating record: %v", zr.Id)
	//TODO: it returns []*record.RecordInfo. Does the driver update more than 1 record at the time?
	_, err = client.Update(zr.toRecordUpdate())
	if err != nil {
		return fmt.Errorf("Cannot update record: %v", err)
	}

	// Success
	log.Printf("[DEBUG] Updated record: %v", zr.Id)
	if activeVersion != 0 {
		setActiveZoneVersion(meta, zr.Zone, zr.Version)
	}
	return nil
}

//DeleteRecord deletes records from zone version by id
func DeleteRecord(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Entering DeleteRecord")
	var err error
	var activeVersion int64
	client := getRecordClient(meta)

	var zr ZoneRecord
	zr.Parse(d)
	if zr.Version == 0 {
		log.Printf("[DEBUG] Looking for active zone version")
		version := getZoneVersionClient(meta)
		_, activeVersion, err = getActiveZoneVersion(meta, zr.Zone)
		newZoneVersion, err := createZoneVersion(version, zr.Zone, activeVersion, 0)
		if err != nil {
			return fmt.Errorf("Could not create new version for record: %v", err)
		}
		_, zr.Version = resourceIDSplit(newZoneVersion, "_")
		// ID needs to also be updated.
		records, err := client.List(zr.Zone, zr.Version)
		for _, r := range records {
			if r.Name == zr.Name && r.Type == zr.Type && r.Value == zr.Value {
				zr.Id = r.Id
				break
			}
		}
	}

	log.Printf("[DEBUG] Deleting record: %v", zr.Id)
	success, err := client.Delete(zr.Zone, zr.Version, zr.Id)
	if err != nil {
		return fmt.Errorf("Cannot delete record: %v", err)
	}

	if success {
		log.Printf("[DEBUG] Deleted record: %v", zr.Id)
		d.SetId("")
	} else {
		log.Printf("[DEBUG] Failure Deleting record: %v %#v", zr.Id, err)
	}
	log.Printf("[INFO] Active zone version: %v New Zone Version: %v", activeVersion, zr.Version)
	if activeVersion != 0 {
		setActiveZoneVersion(meta, zr.Zone, zr.Version)
	}

	return err
}
