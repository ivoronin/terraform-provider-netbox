package main

import (
	"fmt"
	"net/http"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/go-resty/resty/v2"
)

func resourcePrefix() *schema.Resource {
	return &schema.Resource{
		Create: resourcePrefixCreate,
		Read:   resourcePrefixRead,
		Update: resourcePrefixUpdate,
		Delete: resourcePrefixDelete,

		Schema: map[string]*schema.Schema{
			"prefix": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

type Prefix struct {
	Id int `json:"id"`
	Prefix string `json:"prefix"`
}

func resourcePrefixCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	resp, err := client.R().
		SetBody(Prefix{Prefix: d.Get("prefix").(string)}).
		SetResult(Prefix{}).
		Post("/ipam/prefixes/")
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("POST: Unexpected HTTP status: %s", resp.Status())
	}
	prefix := resp.Result().(*Prefix)
	d.SetId(fmt.Sprintf("%d", prefix.Id))
	return resourcePrefixRead(d, m)
}

func resourcePrefixRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	resp, err := client.R().
		SetResult(Prefix{}).
		Get(fmt.Sprintf("/ipam/prefixes/%s/", d.Id()))
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("GET: Unexpected HTTP status: %s", resp.Status())
	}
	prefix := resp.Result().(*Prefix)
	d.Set("prefix", prefix.Prefix)
	return nil
}

func resourcePrefixUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	resp, err := client.R().
		SetBody(Prefix{Prefix: d.Get("prefix").(string)}).
                Put(fmt.Sprintf("/ipam/prefixes/%s/", d.Id()))
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("PUT: Unexpected HTTP status: %s", resp.Status())
	}
	return resourcePrefixRead(d, m)
}

func resourcePrefixDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	resp, err := client.R().
                Delete(fmt.Sprintf("/ipam/prefixes/%s/", d.Id()))
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("DELETE: Unexpected HTTP status: %s", resp.Status())
	}
	return nil
}
