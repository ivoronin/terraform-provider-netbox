package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccIPAddressConfig(prefix string, dns_name1 string, dns_name2 string, ipaddress2 string) string {
	return fmt.Sprintf(`
resource "netbox_prefix" "prefix_test" {
	prefix = "%s"
}

resource "netbox_ipaddress" "ipaddress_test1" {
	prefix_id = netbox_prefix.prefix_test.id
	dns_name = "%s"
}

resource "netbox_ipaddress" "ipaddress_test2" {
	dns_name = "%s"
	address_cidr = "%s"
}
`, prefix, dns_name1, dns_name2, ipaddress2)
}

func TestAccNetboxIPAddress_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIPAddressConfig("192.168.100.0/24", "test-name1", "test-name2", "192.168.90.1/32"),
				Check: resource.ComposeTestCheckFunc(
					testAccNetboxIPAddressCheck("netbox_prefix.prefix_test"),
					testAccNetboxIPAddressCheck("netbox_ipaddress.ipaddress_test1"),
					testAccNetboxIPAddressCheck("netbox_ipaddress.ipaddress_test2"),
				),
			},
			{
				Config: testAccIPAddressConfig("192.168.100.0/24", "test-name3", "test-name2", "192.168.90.2/32"),
				Check: resource.ComposeTestCheckFunc(
					testAccNetboxIPAddressCheck("netbox_prefix.prefix_test"),
					testAccNetboxIPAddressCheck("netbox_ipaddress.ipaddress_test1"),
					testAccNetboxIPAddressCheck("netbox_ipaddress.ipaddress_test2"),
				),
			},
		},
	})
}

func testAccNetboxIPAddressCheck(id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		//address := rs.Primary.Attributes["address"]

		return nil
	}
}
