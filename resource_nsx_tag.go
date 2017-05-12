package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceNSXTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceNSXTagCreate,
		Read: resourceNSXTagRead,
		Update: resourceNSXTagUpdate,
		Delete: resourceNSXTagDelete,

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


func resourceNSXTagCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("%s", "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! IN CREATE")
	return nil
}

func resourceNSXTagRead(d *schema.ResourceData, meta interface{}) error {
	// resp, err := client.Get(c.NSXManager + "/2.0/services/securitytags/tag/securitytag-9")

	// if err != nil {
	//	return err
	// }
	log.Printf("%s", "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! IN READ")
	log.Printf("!!! %+v", meta)

	return nil
}

func resourceNSXTagUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceNSXTagDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("%s", "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! IN DELETE")
	return nil
}