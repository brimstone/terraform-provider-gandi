package main

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/prasmussen/gandi-api/domain/zone/record"
)

func TestAccDMERecordA(t *testing.T) {
	var record record.RecordInfo
	// zone id to perform tests with
	zoneID := os.Getenv("GANDI_ZONE_ID")
	zoneVersion := os.Getenv("GANDI_ZONE_VERSION")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGandiRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testGandiRecordConfigA, zoneID, zoneVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGandiRecordExists("gandi_record.test", &record),
					resource.TestCheckResourceAttr(
						"gandi_recor.test", "zone_id", zoneID),
					resource.TestCheckResourceAttr(
						"gandi_recor.test", "name", "testa"),
					resource.TestCheckResourceAttr(
						"gandi_recor.test", "type", "A"),
					resource.TestCheckResourceAttr(
						"gandi_recor.test", "value", "1.1.1.1"),
					resource.TestCheckResourceAttr(
						"gandi_recor.test", "ttl", "2000"),
				),
			},
		},
	})
}

func testAccCheckGandiRecordDestroy(s *terraform.State) error {
	client := getRecordClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gandi_record" {
			continue
		}

		_, err := GetRecord(client, rs.Primary.Attributes["zone_id"], rs.Primary.Attributes["version"], rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Record still exists")
		}
	}

	return nil
}

func testAccCheckGandiRecordExists(n string, r *record.RecordInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := getRecordClient(testAccProvider.Meta())

		foundRecord, err := GetRecord(client, rs.Primary.Attributes["zone_id"], rs.Primary.Attributes["version"], rs.Primary.ID)

		if err != nil {
			return err
		}

		//TODO: there should be a method on Record struct to make the string conversions easy
		if strconv.FormatInt(foundRecord.Id, 10) != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*r = *foundRecord

		return nil
	}
}

const testGandiRecordConfigA = `
resource "gandi_record" "test" {
  zone_id = "%s"
	version = "%s"
  name = "testa"
  type = "A"
  value = "1.1.1.1"
  ttl = 2000
}`

const testGandiRecordConfigCNAME = `
resource "gandi_record" "test" {
  zone_id = "%s"
  name = "testcname"
  type = "CNAME"
  value = "foo"
  ttl = 2000
}`

const testGandiRecordConfigMX = `
resource "gandi_record" "test" {
  zone_id = "%s"
  name = "testmx"
  type = "MX"
  value = "10 relay.mail.mx."
  ttl = 2000
}`

const testGandiRecordConfigTXT = `
resource "gandi_record" "test" {
  zone_id = "%s"
  name = "testtxt"
  type = "TXT"
  value = "foo"
  ttl = 2000
}`

const testGandiRecordConfigSPF = `
resource "gandi_record" "test" {
  zone_id = "%s"
  name = "testspf"
  type = "SPF"
  value = "foo"
  ttl = 2000
}`

const testGandiRecordConfigPTR = `
resource "gandi_record" "test" {
  zone_id = "%s"
  name = "testptr"
  type = "PTR"
  value = "foo"
  ttl = 2000
}`

const testGandiRecordConfigNS = `
resource "gandi_record" "test" {
  zone_id = "%s"
  name = "testns"
  type = "NS"
  value = "foo"
  ttl = 2000
}`

const testGandiRecordConfigAAAA = `
resource "gandi_record" "test" {
  zone_id = "%s"
  name = "testaaaa"
  type = "AAAA"
  value = "FE80::0202:B3FF:FE1E:8329"
  ttl = 2000
}`

const testGandiRecordConfigSRV = `
resource "gandi_record" "test" {
  zone_id = "%s"
  name = "_testsrv._tcp"
  type = "SRV"
  value = "10 20 5060 old-slow-sip-box.example.com."
  ttl = 2000
}`
