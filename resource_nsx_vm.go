package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"regexp"
	"log"
)

func resourceNSXVm() *schema.Resource {
	return &schema.Resource{
		Create: resourceNSXVmCreate,
		Read: resourceNSXVmRead,
		Update: resourceNSXVmUpdate,
		Delete: resourceNSXVmDelete,

		Schema: map[string]*schema.Schema{
			"vm_id": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				Description: "VM managed object ID or VM instance UUID.",
			},
			"security_tags": &schema.Schema{
				Type: schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{Type: schema.TypeString},
				Description: "List of security tag ids or names.",
			},
		},
	}
}

func resourceNSXVmCreate(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("security_tags"); ok {
		log.Printf("--- TAGS --- %+v ", v)
	}

	return nil
}

func resourceNSXVmRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceNSXVmUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceNSXVmDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func isSecurityTagId (value string) bool {
	matched, err := regexp.MatchString(`^securitytag-\d+$`, value)
	return err == nil && matched == true
}