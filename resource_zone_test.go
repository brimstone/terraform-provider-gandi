package main

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/prasmussen/gandi-api/domain/zone"
)

func TestAccGandiZone(t *testing.T) {
	// var zone zone.ZoneInfo

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGandiZoneDestroy,
		Steps:        []resource.TestStep{},
	})
}

func TestAccGandiZoneCreate(t *testing.T) {
	var zone zone.ZoneInfo
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckZone(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGandiRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testGandiZoneConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGandiZoneExists("gandi_zone.test", &zone),
					resource.TestCheckResourceAttr(
						"gandi_zone.test", "name", "testing_zone"),
				),
			},
		},
	})
}

func testAccCheckGandiZoneExists(n string, z *zone.ZoneInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Zone ID is set")
		}

		client := getZoneClient(testAccProvider.Meta())

		ID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Invalid Zone ID")
		}

		foundZone, err := client.Info(ID)

		if err != nil {
			log.Printf("%+v", foundZone)
			return err
		}

		//TODO: make method ZoneInfo struct to make the string conversions easy
		if strconv.FormatInt(foundZone.Id, 10) != rs.Primary.ID {
			return fmt.Errorf("Zone not found")
		}

		*z = *foundZone

		return nil
	}
}

func testAccCheckGandiZoneDestroy(s *terraform.State) error {
	client := getZoneClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gandi_zone" {
			continue
		}

		ID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Invalid Zone ID")
		}

		_, err = client.Info(ID)

		if err == nil {
			return fmt.Errorf("Zone still exists")
		}
	}

	return nil
}

const testGandiZoneConfig = `
resource "gandi_zone" "test" {
  name = "testing_zone"
	}`
