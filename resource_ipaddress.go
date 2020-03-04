package main

import (
	"fmt"
	"net/http"
        "strconv"
	"strings"
	"sync"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/go-resty/resty/v2"
)

var resourceIPAddressCreateAvailableMutex = &sync.Mutex{}

func resourceIPAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPAddressCreate,
		Read:   resourceIPAddressRead,
		Update: resourceIPAddressUpdate,
		Delete: resourceIPAddressDelete,

		Schema: map[string]*schema.Schema{
			"prefix_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"address_cidr": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"address": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
                        "status":  &schema.Schema{
                                Type: schema.TypeString,
                                Optional: true,
                                Computed: true,
                        },
			"dns_name":  &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
			},
                        "description":  &schema.Schema{
                                Type: schema.TypeString,
                                Optional: true,
                        },
		},
	}
}

type IPAddressRequest struct {
       Id int `json:"id"`
       Address string `json:"address"`
       DNSName string `json:"dns_name"`
       Description string `json:"description"`
       Status string `json:"status"`
}

type IPAddressResponse struct {
       Id int `json:"id"`
       Address string `json:"address"`
       DNSName string `json:"dns_name"`
       Description string `json:"description"`
       Status Status
}

type Status struct {
       Value int `json:"value"`
       Label string `json:"label"`
}


func resourceIPAddressCreate(d *schema.ResourceData, m interface{}) error {
	var body IPAddressRequest
	var url string
        ipaddress_status_choise := map[string]string{
            "active":     "1",
            "reserved":   "2",
            "deprecated": "3",
            "dhcp":       "5",
        }

	client := m.(*resty.Client)
	if ip, ok := d.GetOk("address_cidr"); ok {
                body = IPAddressRequest{Address: ip.(string), DNSName: d.Get("dns_name").(string), Description: d.Get("description").(string), Status: ipaddress_status_choise[d.Get("status").(string)]}
		url = "/ipam/ip-addresses/"
	} else {
		if id, ok := d.GetOk("prefix_id"); !ok {
			return fmt.Errorf("prefix_id is required when no address is set")
		} else {
			url = fmt.Sprintf("/ipam/prefixes/%s/available-ips/", id.(string))
		}
		resourceIPAddressCreateAvailableMutex.Lock()
		defer resourceIPAddressCreateAvailableMutex.Unlock()
		body = IPAddressRequest{DNSName: d.Get("dns_name").(string), Description: d.Get("description").(string), Status: ipaddress_status_choise[d.Get("status").(string)]}
	}
	resp, err := client.R().
		SetBody(body).
		SetResult(IPAddressResponse{}).
		Post(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("POST: Unexpected HTTP status: %s", resp.Status())
	}
	ipaddress := resp.Result().(*IPAddressResponse)
	d.SetId(fmt.Sprintf("%d", ipaddress.Id))
	return resourceIPAddressRead(d, m)
}

func resourceIPAddressRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
        ipaddress_status_choise := map[string]string{
            "1": "active",
            "2": "reserved",
            "3": "deprecated",
            "5": "dhcp",
        }

	resp, err := client.R().
		SetResult(IPAddressResponse{}).
		Get(fmt.Sprintf("/ipam/ip-addresses/%s/", d.Id()))
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("GET: Unexpected HTTP status: %s", resp.Status())
	}
	ipaddress := resp.Result().(*IPAddressResponse)
	d.Set("address_cidr", ipaddress.Address)
	address, _ := SplitAddressMask(ipaddress.Address)
	d.Set("address", address)
        d.Set("status", ipaddress_status_choise[strconv.Itoa(ipaddress.Status.Value)])
        d.Set("description", ipaddress.Description)
	return nil
}

func resourceIPAddressUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
        ipaddress_status_choise := map[string]string{
            "active":     "1",
            "reserved":   "2",
            "deprecated": "3",
            "dhcp":       "5",
        }
	resp, err := client.R().
		SetBody(IPAddressRequest{Address: d.Get("address_cidr").(string), DNSName: d.Get("dns_name").(string), Description: d.Get("description").(string), Status: ipaddress_status_choise[d.Get("status").(string)]}).
                Put(fmt.Sprintf("/ipam/ip-addresses/%s/", d.Id()))
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("PUT: Unexpected HTTP status: %s", resp.Status())
	}
	return resourceIPAddressRead(d, m)
}

func resourceIPAddressDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	resp, err := client.R().
                Delete(fmt.Sprintf("/ipam/ip-addresses/%s/", d.Id()))
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("DELETE: Unexpected HTTP status: %s", resp.Status())
	}
	return nil
}

func SplitAddressMask(address string) (string, string) {
	x := strings.SplitN(address, "/", 2)
	return x[0], x[1]
}
