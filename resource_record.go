package gandi

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cznic/sortutil"
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
				Required: true,
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

type ZoneRecord struct {
	record.RecordAdd
	Id int64
}

func (zr *ZoneRecord) Parse(d *schema.ResourceData) {
	zr.Zone, _ = strconv.ParseInt(d.Get("zone_id").(string), 10, 64)
	zr.Ttl, _ = strconv.ParseInt(d.Get("ttl").(string), 10, 64)
	zr.Version, _ = strconv.ParseInt(d.Get("version").(string), 10, 64)
	zr.Name = d.Get("name").(string)
	zr.Value = d.Get("value").(string)
	zr.Type = d.Get("type").(string)
	zr.Id, _ = strconv.ParseInt(d.Id(), 10, 64)
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

// CreateRecord creates new record
func CreateRecord(d *schema.ResourceData, meta interface{}) error {
	client := getRecordClient(meta)

	var zr ZoneRecord
	zr.Parse(d)

	// log.Printf("[DEBUG] Creating new record from spec: %+v", &ZoneRecord.New(d))

	newRecord, err := client.Add(zr.toRecordAdd())
	if err != nil {
		return fmt.Errorf("Could not create new record: %v", err)
	}

	// Success
	d.SetId(strconv.FormatInt(newRecord.Id, 10))
	log.Printf("[INFO] Successfully created record: %v", d.Id())

	return ReadRecord(d, meta)
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

	// need an int64 slice for sorting
	var recordIDs sortutil.Int64Slice
	for _, r := range records {
		recordIDs = append(recordIDs, r.Id)
	}

	recordIDs.Sort()
	i := sortutil.SearchInt64s(recordIDs, rid)
	if i < len(recordIDs) && recordIDs[i] == rid {
		log.Printf("[DEBUG] Record: %v found...", rid)
		return records[i], nil
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
	client := getRecordClient(meta)

	// zoneID is stored as string in tfstate, API expects an int64
	zoneID := d.Get("zone_id")
	zoneVersion := d.Get("version")
	recordID := d.Id()

	log.Printf("[DEBUG] Reading records from zone: %v version: %v", zoneID, zoneVersion)

	record, err := GetRecord(client, zoneID, zoneVersion, recordID)
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
		d.Set("name", record.Name)
		d.Set("value", record.Value)
		d.Set("ttl", strconv.FormatInt(record.Ttl, 10))
		d.Set("type", record.Type)
	}

	return nil
}

// UpdateRecord updates record in zone/version according to the new spec
func UpdateRecord(d *schema.ResourceData, meta interface{}) error {
	client := getRecordClient(meta)

	// zoneID is stored as string in tfstate, API expects an int64
	zoneID, _ := strconv.ParseInt(d.Get("zone_id").(string), 10, 64)
	ID, _ := strconv.ParseInt(d.Id(), 10, 64)
	ttl, _ := strconv.ParseInt(d.Get("ttl").(string), 10, 64)
	version, _ := strconv.ParseInt(d.Get("version").(string), 10, 64)

	updatedRecordSpec := record.RecordUpdate{
		Name:    d.Get("name").(string),
		Value:   d.Get("value").(string),
		Ttl:     ttl,
		Type:    d.Get("type").(string),
		Zone:    zoneID,
		Version: version,
		Id:      ID,
	}

	log.Printf("[DEBUG] Updating record: %v", d.Id())
	_, err := client.Update(updatedRecordSpec)
	if err != nil {
		return fmt.Errorf("Cannot update record: %v", err)
	}

	// Success
	log.Printf("[DEBUG] Updated record: %v", d.Id())
	return nil
}

//DeleteRecord deletes records from zone version by id
func DeleteRecord(d *schema.ResourceData, meta interface{}) error {
	client := getRecordClient(meta)

	zoneID, _ := strconv.ParseInt(d.Get("zone_id").(string), 10, 64)
	ID, _ := strconv.ParseInt(d.Id(), 10, 64)
	version, _ := strconv.ParseInt(d.Get("version").(string), 10, 64)

	log.Printf("[DEBUG] Deleting record: %v", d.Id())
	success, err := client.Delete(zoneID, version, ID)
	if err != nil {
		return fmt.Errorf("Cannot delete record: %v", err)
	}

	if success {
		log.Printf("[DEBUG] Deleted record: %v", d.Id())
		d.SetId("")
	}

	return err
}
