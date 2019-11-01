# terraform-provider-netbox
Manages [NetBox](https://github.com/netbox-community/netbox) IP address management (IPAM) and data center infrastructure management (DCIM) tool. 

# Example
```hcl
resource "netbox_prefix" "example_prefix" {
        prefix = "192.168.0.0/24"
}

// register next available IP
resource "netbox_ipaddress" "example_ipaddress1" {
        prefix_id = netbox_prefix.example_prefix.id
        dns_name = "example-host1"
}

// register static IP
resource "netbox_ipaddress" "example_ipaddress2" {
        dns_name = "example-host2"
        address_cidr = "10.0.0.2/16"
}

// Terraform Mikrotik Provider - https://github.com/ddelnano/terraform-provider-mikrotik
resource "mikrotik_dhcp_lease" "example_lease" {
  address    = netbox_ipaddress.example_ipaddress1.address
  macaddress = "B2:8E:98:58:B3:74 "
  comment    = "Example DHCP Lease"
}
```
