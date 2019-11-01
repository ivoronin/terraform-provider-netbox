package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccPrefixConfig(prefix string) string {
	return fmt.Sprintf(`
resource "netbox_prefix" "prefix_test" {
	prefix = "%s"
}
`, prefix)
}

func TestAccNetboxPrefix_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPrefixConfig("192.168.100.0/24"),
				Check: resource.ComposeTestCheckFunc(
					testAccNetboxPrefixCheck("netbox_prefix.prefix_test"),
				),
			},
			{
				Config: testAccPrefixConfig("192.168.200.0/24"),
				Check: resource.ComposeTestCheckFunc(
					testAccNetboxPrefixCheck("netbox_prefix.prefix_test"),
				),
			},
		},
	})
}

func testAccNetboxPrefixCheck(id string) resource.TestCheckFunc {
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
