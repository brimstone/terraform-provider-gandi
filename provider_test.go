package main

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Acceptance Tests for the Gandi Provider

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// To run these acceptance tests, you will need a Gandi Account
// There is a need to set-up the access credentials and enable API access
//
// With all of that done, you can run like this:
//    make testacc TEST=./builtin/providers/gandi

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"gandi": testAccProvider,
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("GANDI_KEY"); v == "" {
		t.Fatal("GANDI_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("GANDI_TESTING"); v == "" {
		t.Fatal("GANDI_TESTING must be set for acceptance tests")
	}
}
