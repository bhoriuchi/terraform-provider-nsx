package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"encoding/xml"
)

func resourceNSXTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceNSXTagCreate,
		Read: resourceNSXTagRead,
		Update: resourceNSXTagUpdate,
		Delete: resourceNSXTagDelete,

		Schema: map[string]*schema.Schema{
			"tag_id": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
			},
			"tag_name": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
			},
		},
	}
}


func resourceNSXTagCreate(d *schema.ResourceData, meta interface{}) error {
	// tagId := d.Get("tag_id").(string)
	// tagName := d.Get("tag_name").(string)
	config := meta.(*Config)

	// if a taglist has not been requested or is empty, request it again
	if config.RequestedTagList == false && len(config.TagsList.SecurityTags) == 0 {
		config.RequestedTagList = true
		resp, getErr := getRequest(config, "/2.0/services/securitytags/tag")
		if getErr != nil {
			return getErr
		}
		parseErr := xml.Unmarshal(resp.Body(), &config.TagsList)
		if parseErr != nil {
			return parseErr
		}
	}

	// log.Printf("!!! %+v", config.TagsList)

	return nil
}
/*
func getTagById (config Config, tagId string) (NSXTag, error) {
	resp, err := getRequest(config, "/2.0/services/securitytags/tag/" + tagId)

	if err != nil {
		return nil, err
	}

	v := NSXTag{}
	err2 := xml.Unmarshal(resp.Body(), &v)

	if err2 != nil {
		return nil, err2
	}

}
*/

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
