package main

import (
	"os"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()//.(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"netbox": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("NETBOX_TOKEN"); v == "" {
		t.Fatal("NETBOX_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("NETBOX_BASE_URL"); v == "" {
		t.Fatal("NETBOX_BASE_URL must be set for acceptance tests")
	}
}
