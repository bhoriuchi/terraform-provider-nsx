package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"gopkg.in/resty.v0"
	"log"
	"encoding/xml"
	"fmt"
)

type NSXTagType struct {
	TypeName string `xml:"typeName"`
}

type NSXTag struct {
	ObjectId string `xml:"objectId"`
	ObjectTypeName string `xml:"objectTypeName"`
	VsmUuid string `xml:"vsmUuid"`
	NodeId string `xml:"nodeId"`
	Revision int `xml:"revision"`
	Type NSXTagType `xml:"type"`
	Name string `xml:"name"`
	Description string `xml:"description"`
	IsUniversal bool `xml:"isUniversal"`
	UniversalRevision int `xml:"universalRevision"`
	SystemResource bool `xml:"systemResource"`
	VmCount int `xml:"vmCount"`
}

type NSXTagPost struct {
	XMLName xml.Name `xml:"securityTag"`
	ObjectTypeName string `xml:"objectTypeName"`
	Type NSXTagType `xml:"type"`
	Name string `xml:"name"`
	Description string `xml:"description"`
	IsUniversal bool `xml:"isUniversal"`
}

type NSXTagList struct {
	SecurityTags []NSXTag `xml:"securityTag"`
}

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
			"description": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
			},
			"create_if_missing": &schema.Schema{
				Type: schema.TypeBool,
				Optional: true,
				Default: true,
			},
			"is_universal": &schema.Schema{
				Type: schema.TypeBool,
				Optional: true,
				Default: false,
			},
		},
	}
}

func resourceNSXTagCreate(d *schema.ResourceData, meta interface{}) error {
	tagId := d.Get("tag_id").(string)
	tagName := d.Get("tag_name").(string)
	createMissing := d.Get("create_if_missing").(bool)
	config := meta.(*Config)
	tagValue := NSXTag{}

	if tagId != "" {
		resp, getErr := getRequest(config, "/2.0/services/securitytags/tag/" + tagId)
		if getErr != nil {
			return getErr
		}
		if resp.StatusCode() != 200 {
			return fmt.Errorf("tag %q not found", tagId)
		}

		parseErr := xml.Unmarshal(resp.Body(), &tagValue)
		if parseErr != nil {
			return parseErr
		}
		log.Printf("!!! %+v", tagValue)
		d.Set("tag_name", tagValue.Name)
		d.SetId(tagValue.ObjectId)

	} else if tagName != "" {
		found := false
		tagList := NSXTagList{}
		resp, getErr := getRequest(config, "/2.0/services/securitytags/tag")
		if getErr != nil {
			return getErr
		}
		if resp.StatusCode() != 200 {
			return fmt.Errorf("unable to query NSX tags")
		}
		parseErr := xml.Unmarshal(resp.Body(), &tagList)
		if parseErr != nil {
			return parseErr
		}

		for i := range tagList.SecurityTags {
			if tagList.SecurityTags[i].Name == tagName {
				found = true
				tagValue = tagList.SecurityTags[i]
				break
			}
		}
		if found != true {
			if createMissing == true {
				newTag, createErr := createNSXTag(d, meta)
				if createErr != nil {
					return createErr
				}
				log.Printf("=================CREATED TAG %+v", newTag)
				tagValue = newTag
			} else {
				return fmt.Errorf("Security Tag %q not found and %q is set to false", tagName, "create_if_missing")
			}
		}
	} else {
		return fmt.Errorf("must provide either %q or %q", "tag_id", "tag_name")
	}

	log.Printf("!!! %+v", tagValue)
	d.SetId(tagValue.ObjectId)
	d.Set("tag_id", tagValue.ObjectId)
	d.Set("tag_name", tagValue.Name)

	return nil
}

func createNSXTag(d *schema.ResourceData, meta interface{}) (NSXTag, error) {
	config := meta.(*Config)
	body := NSXTagPost{
		ObjectTypeName: "SecurityTag",
		Type: NSXTagType{ TypeName: "SecurityTag" },
		Name: d.Get("tag_name").(string),
		IsUniversal: d.Get("is_universal").(bool),
		Description: d.Get("description").(string),
	}

	resp, err := resty.R().SetBody(body).Post(config.NSXManager + "/2.0/services/securitytags/tag")

	if err != nil {
		return NSXTag{}, err
	}

	log.Printf("^^^^^^^^^^^^^^^^^^^^^ %+v, %d, %+v", resp, resp.StatusCode(), resp.Error())

	tagValue := NSXTag{}
	parseErr := xml.Unmarshal(resp.Body(), &tagValue)

	if parseErr != nil {
		return NSXTag{}, parseErr
	}

	return tagValue, nil
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
