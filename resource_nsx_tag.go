package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"gopkg.in/resty.v0"
	"log"
	"encoding/xml"
	"fmt"
	"errors"
)

type NSXTagType struct {
	TypeName string `xml:"typeName"`
}

type NSXTag struct {
	ObjectId string `xml:"objectId"`
	ObjectTypeName string `xml:"objectTypeName"`
	Type NSXTagType `xml:"type"`
	Name string `xml:"name"`
	Description string `xml:"description"`
}

type NSXTagList struct {
	SecurityTags []NSXTag `xml:"securityTag"`
}

// support multiple api versions
type NSXTagPost_6_2 struct {
	XMLName xml.Name `xml:"securityTag"`
	ObjectTypeName string `xml:"objectTypeName"`
	Type NSXTagType `xml:"type"`
	Name string `xml:"name"`
	Description string `xml:"description"`
}

type NSXTagPost_6_3 struct {
	XMLName xml.Name `xml:"securityTag"`
	ObjectTypeName string `xml:"objectTypeName"`
	Type NSXTagType `xml:"type"`
	Name string `xml:"name"`
	Description string `xml:"description"`
	IsUniversal bool `xml:"isUniversal"`
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
				Default: false,
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

	if tagId != "" {
		tag, err := getNSXTagById(config, tagId)

		if err != nil {
			return err
		}

		d.Set("tag_description", tag.Description)
		setNSXTagId(d, tag.ObjectId)

	} else if tagName != "" {
		foundTag, foundErr := getNSXTagByName(config, tagName)

		if foundErr == nil {
			d.Set("tag_description", foundTag.Description)
			setNSXTagId(d, foundTag.ObjectId)
		} else {
			if createMissing == true {
				newTag, createErr := createNSXTag(d, meta)
				if createErr != nil {
					return createErr
				}
				d.Set("tag_description", newTag.Description)
				setNSXTagId(d, newTag.ObjectId)

			} else {
				return fmt.Errorf("Security Tag %q not found and %q is set to false", tagName, "create_if_missing")
			}
		}
	} else {
		return fmt.Errorf("must provide either %q or %q", "tag_id", "tag_name")
	}

	return nil
}

func setNSXTagId (d *schema.ResourceData, tagId string) {
	d.SetId(tagId)
	d.Set("tag_id", tagId)
}

func getNSXTagByName (config *Config, tagName string) (NSXTag, error) {
	found := false
	tagList := NSXTagList{}
	tagValue := NSXTag{}
	getErr := getRequest(config.TagEndpoint, &tagList)

	if getErr != nil {
		return tagValue, getErr
	}

	for i := range tagList.SecurityTags {
		if tagList.SecurityTags[i].Name == tagName {
			found = true
			tagValue = tagList.SecurityTags[i]
			break
		}
	}

	if found == false {
		return tagValue, errors.New("tag not found")
	}

	return tagValue, nil
}

func getNSXTagById (config *Config, tagId string) (NSXTag, error) {
	tagValue := NSXTag{}
	getErr := getRequest(fmt.Sprintf("%s/%s", config.TagEndpoint, tagId), &tagValue)

	if getErr != nil {
		return tagValue, getErr
	}

	return tagValue, nil
}

func createNSXTag(d *schema.ResourceData, meta interface{}) (NSXTag, error) {
	config := meta.(*Config)
	var body interface{}

	// use the correct post body
	if config.NSXVersion.Major == 6 && config.NSXVersion.Minor == 2 {
		body = NSXTagPost_6_2{
			ObjectTypeName: "SecurityTag",
			Type: NSXTagType{ TypeName: "SecurityTag" },
			Name: d.Get("tag_name").(string),
			Description: d.Get("description").(string),
		}
	} else {
		body = NSXTagPost_6_3{
			ObjectTypeName: "SecurityTag",
			Type: NSXTagType{ TypeName: "SecurityTag" },
			Name: d.Get("tag_name").(string),
			IsUniversal: d.Get("is_universal").(bool),
			Description: d.Get("description").(string),
		}
	}

	newTag := NSXTag{
		ObjectTypeName: "SecurityTag",
		Type: NSXTagType{ TypeName: "SecurityTag" },
		Name: d.Get("tag_name").(string),
		Description: d.Get("description").(string),
	}

	resp, err := resty.R().SetBody(body).Post(config.TagEndpoint)

	if err != nil {
		return newTag, err
	}

	if resp.StatusCode() != 201 {
		log.Printf("%+v", resp)
		return newTag, errors.New(resp.String())
	}

	newTag.ObjectId = resp.String()
	return newTag, nil
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
