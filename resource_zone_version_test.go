package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccGandiZoneVersion(t *testing.T) {
	zoneID := os.Getenv("GANDI_ZONE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckZoneVersion(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGandiRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testGandiZoneVersionConfig, zoneID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGandiZoneVersionExists("gandi_zone_version.test"),
					resource.TestCheckResourceAttr(
						"gandi_zone_version.test", "base_version", "1"),
					resource.TestCheckResourceAttr(
						"gandi_zone_version.test", "zone_version", "2"),
					resource.TestCheckResourceAttr(
						"gandi_zone_version.test", "zone_id", zoneID),
				),
			},
		},
	})
}

func testAccCheckGandiZoneVersionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Zone Version ID is set")
		}

		client := getZoneVersionClient(testAccProvider.Meta())
		zoneID, zoneVersion := resourceIDSplit(rs.Primary.ID, "_")
		zoneExists, err := CheckZoneVersion(client, zoneID, zoneVersion)

		if err != nil {
			return err
		}

		if zoneExists {
			return nil
		}

		return fmt.Errorf("Zone Version not found")
	}
}

func testAccCheckGandiZoneVersionDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gandi_zone_version" {
			continue
		}

		client := getZoneVersionClient(testAccProvider.Meta())
		zoneID, zoneVersion := resourceIDSplit(rs.Primary.ID, "_")
		zoneExists, _ := CheckZoneVersion(client, zoneID, zoneVersion)

		if zoneExists {
			return fmt.Errorf("Zone Version still exists")
		}
	}

	return nil
}

const testGandiZoneVersionConfig = `
resource "gandi_zone_version" "test" {
	base_version = 1
	zone_version = 2
	zone_id = "%s"
}`
