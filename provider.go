package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns Gandi Resoruce Provider...
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GANDI_KEY", nil),
				Description: "A Gandi XMLRPC API Key.",
			},
			"testing": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GANDI_TESTING", true),
				Description: "Set it to use the Test Environment.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"gandi_zone":         resourceZone(),
			"gandi_record":       resourceRecord(),
			"gandi_zone_version": resourceZoneVersion(),
		},

		ConfigureFunc: providerConfigure,
	}
}
