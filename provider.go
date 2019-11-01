package main

import (
	"fmt"
	"strings"
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETBOX_TOKEN", nil),
				Description: "The token used to connect to NetBox.",
			},
			"base_url": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("NETBOX_BASE_URL", ""),
				Description:  "The NetBox Base API URL",
				ValidateFunc: validateApiURL,
			},
			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A file containing the ca certificate to use in case ssl certificate is not from a standard chain",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable SSL verification of API calls",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"netbox_prefix": resourcePrefix(),
			"netbox_ipaddress": resourceIPAddress(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func validateApiURL(value interface{}, key string) (ws []string, es []error) {
	v := value.(string)
	if !strings.HasSuffix(v, "/api") {
		es = append(es, fmt.Errorf("base_url should end with /api"))
	}
	return
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: d.Get("insecure").(bool)})
	client.SetHostURL(d.Get("base_url").(string))
	client.SetHeader("Authorization", fmt.Sprintf("Token %s", d.Get("token").(string)))
	return client, nil
}
