package gandi

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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

func testAccCheckGandiZoneDestroy(s *terraform.State) error { return nil }
