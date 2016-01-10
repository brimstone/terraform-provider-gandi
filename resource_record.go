package main

import (
	"fmt"
	"log"
	"strconv"

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
				Type:     schema.TypeInt,
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
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

// getRecordClient wraps Gandi Client in Record Resource Methods
func getRecordClient(meta interface{}) *record.Record {
	return record.New(meta.(*client.Client))
}

func recordIDExist(records []*record.RecordInfo, recordID int64) int {
	var recordIDs sortutil.Int64Slice
	for _, r := range records {
		recordIDs = append(recordIDs, r.Id)
	}

	// sort the list before lookup
	recordIDs.Sort()
	i := sortutil.SearchInt64s(recordIDs, recordID)
	if i < len(recordIDs) && recordIDs[i] == recordID {
		log.Printf("[DEBUG] Record: %v found...", recordID)
		return i
	}
	log.Printf("[DEBUG] Record %v not found...", recordID)
	return -1
}

// CreateRecord creates new record
func CreateRecord(d *schema.ResourceData, meta interface{}) error {
	client := getRecordClient(meta)

	// zoneID is stored as string in tfstate, API expects an int64
	zoneID, _ := strconv.ParseInt(d.Get("zone_id").(string), 10, 64)

	newRecordSpec := record.RecordAdd{
		Name:    d.Get("name").(string),
		Value:   d.Get("value").(string),
		Ttl:     int64(d.Get("ttl").(int)),
		Type:    d.Get("type").(string),
		Zone:    zoneID,
		Version: int64(d.Get("version").(int)),
	}

	log.Printf("[DEBUG] Creating new record from spec: %+v", newRecordSpec)

	newRecord, err := client.Add(newRecordSpec)
	if err != nil {
		return fmt.Errorf("Could not create new record: %v", err)
	}

	// Success
	d.SetId(strconv.FormatInt(newRecord.Id, 10))
	log.Printf("[INFO] Successfully created record: %v", d.Id())

	return ReadRecord(d, meta)
}

// ReadRecord fetches configuration
func ReadRecord(d *schema.ResourceData, meta interface{}) error {
	client := getRecordClient(meta)

	// zoneID is stored as string in tfstate, API expects an int64
	zoneID, _ := strconv.ParseInt(d.Get("zone_id").(string), 10, 64)

	log.Printf("[DEBUG] Reading records from zone:%v, version:%v", zoneID, d.Get("version"))
	records, err := client.List(int64(zoneID), int64(d.Get("version").(int)))

	if err != nil {
		return fmt.Errorf("Cannot read record: %v", d.Id())
	}

	// ID is stored as string in tfstate, API expects an int64
	ID, _ := strconv.ParseInt(d.Id(), 10, 64)

	i := recordIDExist(records, ID) //index of the Record
	if i < 0 {
		log.Printf("[DEBUG] Record: %v not found", ID)
		d.SetId("")
		return nil
	}

	d.Set("name", records[i].Name)
	d.Set("value", records[i].Value)
	d.Set("ttl", records[i].Ttl)
	d.Set("type", records[i].Type)
	// zone_id and version properties are already set properly since
	// they were used for the lookup

	return nil
}

// UpdateRecord updates record in zone/version according to the new spec
func UpdateRecord(d *schema.ResourceData, meta interface{}) error {
	client := getRecordClient(meta)

	// zoneID is stored as string in tfstate, API expects an int64
	zoneID, _ := strconv.ParseInt(d.Get("zone_id").(string), 10, 64)
	ID, _ := strconv.ParseInt(d.Id(), 10, 64)

	updatedRecordSpec := record.RecordUpdate{
		Name:    d.Get("name").(string),
		Value:   d.Get("value").(string),
		Ttl:     int64(d.Get("ttl").(int)),
		Type:    d.Get("type").(string),
		Zone:    zoneID,
		Version: int64(d.Get("version").(int)),
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

	//Delete(zoneId, version, recordId int64)
	// zoneID is stored as string in tfstate, API expects an int64
	zoneID, _ := strconv.ParseInt(d.Get("zone_id").(string), 10, 64)
	ID, _ := strconv.ParseInt(d.Id(), 10, 64)

	log.Printf("[DEBUG] Deleting record: %v", d.Id())
	success, err := client.Delete(zoneID, int64(d.Get("version").(int)), ID)
	if err != nil {
		return fmt.Errorf("Cannot delete record: %v", err)
	}

	if success {
		log.Printf("[DEBUG] Deleted record: %v", d.Id())
		d.SetId("")
	}

	return err
}
