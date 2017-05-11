package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNSXTag() *schema.Resource {
	return &schema.Resource{
		Read: resourceNSXTagRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceNSXTagRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}