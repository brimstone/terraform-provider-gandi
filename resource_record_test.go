package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/prasmussen/gandi-api/domain/zone/record"
)

func TestAccGandiRecordA(t *testing.T) {
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
						"gandi_record.test", "zone_id", zoneID),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "name", "testa"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "type", "A"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "value", "1.1.1.1"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "ttl", "2000"),
				),
			},
		},
	})
}

func TestAccGandiRecordCNAME(t *testing.T) {
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
				Config: fmt.Sprintf(testGandiRecordConfigCNAME, zoneID, zoneVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGandiRecordExists("gandi_record.test", &record),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "zone_id", zoneID),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "name", "testcname"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "type", "CNAME"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "value", "foo"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "ttl", "2000"),
				),
			},
		},
	})
}

func TestAccGandiRecordMX(t *testing.T) {
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
				Config: fmt.Sprintf(testGandiRecordConfigMX, zoneID, zoneVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGandiRecordExists("gandi_record.test", &record),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "zone_id", zoneID),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "name", "testmx"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "type", "MX"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "value", "10 relay.mail.mx."),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "ttl", "2000"),
				),
			},
		},
	})
}

func TestAccGandiRecordTXT(t *testing.T) {
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
				Config: fmt.Sprintf(testGandiRecordConfigTXT, zoneID, zoneVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGandiRecordExists("gandi_record.test", &record),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "zone_id", zoneID),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "name", "testtxt"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "type", "TXT"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "value", "foo"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "ttl", "2000"),
				),
			},
		},
	})
}

func TestAccGandiRecordSPF(t *testing.T) {
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
				Config: fmt.Sprintf(testGandiRecordConfigSPF, zoneID, zoneVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGandiRecordExists("gandi_record.test", &record),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "zone_id", zoneID),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "name", "testspf"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "type", "SPF"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "value", "foo"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "ttl", "2000"),
				),
			},
		},
	})
}

func TestAccGandiRecordNS(t *testing.T) {
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
				Config: fmt.Sprintf(testGandiRecordConfigNS, zoneID, zoneVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGandiRecordExists("gandi_record.test", &record),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "zone_id", zoneID),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "name", "testns"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "type", "NS"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "value", "foo"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "ttl", "2000"),
				),
			},
		},
	})
}

func TestAccGandiRecordAAAA(t *testing.T) {
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
				Config: fmt.Sprintf(testGandiRecordConfigAAAA, zoneID, zoneVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGandiRecordExists("gandi_record.test", &record),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "zone_id", zoneID),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "name", "testaaaa"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "type", "AAAA"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "value", "fe80::202:b3ff:fe1e:8329"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "ttl", "2000"),
				),
			},
		},
	})
}

func TestAccGandiRecordSRV(t *testing.T) {
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
				Config: fmt.Sprintf(testGandiRecordConfigSRV, zoneID, zoneVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGandiRecordExists("gandi_record.test", &record),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "zone_id", zoneID),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "name", "_testsrv._tcp"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "type", "SRV"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "value", "10 20 5060 old-slow-sip-box.example.com."),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "ttl", "2000"),
				),
			},
		},
	})
}

func TestAccGandiRecordModifyAintoCNAME(t *testing.T) {
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
						"gandi_record.test", "zone_id", zoneID),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "name", "testa"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "type", "A"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "value", "1.1.1.1"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "ttl", "2000"),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(testGandiRecordConfigCNAME, zoneID, zoneVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGandiRecordExists("gandi_record.test", &record),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "zone_id", zoneID),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "name", "testcname"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "type", "CNAME"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "value", "foo"),
					resource.TestCheckResourceAttr(
						"gandi_record.test", "ttl", "2000"),
				),
			},
		},
	})
}

const testGandiRecordConfigA = `
resource "gandi_record" "test" {
  zone_id = "%s"
	version = "%s"
  name = "testa"
  type = "A"
  value = "1.1.1.1"
  ttl = "2000"
}`

const testGandiRecordConfigCNAME = `
resource "gandi_record" "test" {
  zone_id = "%s"
	version = "%s"
  name = "testcname"
  type = "CNAME"
  value = "foo"
  ttl = 2000
}`

const testGandiRecordConfigMX = `
resource "gandi_record" "test" {
  zone_id = "%s"
	version = "%s"
  name = "testmx"
  type = "MX"
  value = "10 relay.mail.mx."
  ttl = 2000
}`

const testGandiRecordConfigTXT = `
resource "gandi_record" "test" {
  zone_id = "%s"
	version = "%s"
  name = "testtxt"
  type = "TXT"
  value = "foo"
  ttl = 2000
}`

const testGandiRecordConfigSPF = `
resource "gandi_record" "test" {
  zone_id = "%s"
	version = "%s"
  name = "testspf"
  type = "SPF"
  value = "foo"
  ttl = 2000
}`

const testGandiRecordConfigPTR = `
resource "gandi_record" "test" {
  zone_id = "%s"
	version = "%s"
  name = "testptr"
  type = "PTR"
  value = "foo"
  ttl = 2000
}`

const testGandiRecordConfigNS = `
resource "gandi_record" "test" {
  zone_id = "%s"
	version = "%s"
  name = "testns"
  type = "NS"
  value = "foo"
  ttl = 2000
}`

// TODO: looks like the IPv6 needs to be small letters only
const testGandiRecordConfigAAAA = `
resource "gandi_record" "test" {
  zone_id = "%s"
	version = "%s"
  name = "testaaaa"
  type = "AAAA"
  value = "fe80::202:b3ff:fe1e:8329"
  ttl = 2000
}`

const testGandiRecordConfigSRV = `
resource "gandi_record" "test" {
  zone_id = "%s"
	version = "%s"
  name = "_testsrv._tcp"
  type = "SRV"
  value = "10 20 5060 old-slow-sip-box.example.com."
  ttl = 2000
}`

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
			log.Printf("%+v", foundRecord)
			return err
		}

		//TODO: make method Record struct to make the string conversions easy
		if strconv.FormatInt(foundRecord.Id, 10) != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*r = *foundRecord

		return nil
	}
}
