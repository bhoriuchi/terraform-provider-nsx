package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"gopkg.in/resty.v0"
	"log"
	"encoding/xml"
	"net/http"
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
				Description: "Security tag id (i.e. \"securitytag-1\").",
			},
			"tag_name": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Description: "Case sensetive security tag name to search for.",
			},
			"description": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Description: "Security tag description.",
			},
			"create_if_missing": &schema.Schema{
				Type: schema.TypeBool,
				Optional: true,
				Default: false,
				Description: "Creates the security tag if it is not found during create.",
			},
			"is_universal": &schema.Schema{
				Type: schema.TypeBool,
				Optional: true,
				Default: false,
				Description: "Creates the security tag as a universal tag (NSX 6.3 and higher).",
			},
			"persistent": &schema.Schema{
				Type: schema.TypeBool,
				Optional: true,
				Default: true,
				Description: "When true, prevents removal of the security tag during a destroy operation.",
			},
		},
	}
}

func resourceNSXTagCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if v, ok := d.GetOk("tag_id"); ok {
		tagId := v.(string)
		tag, err := getNSXTagById(config, tagId)

		if err != nil {
			return err
		}

		d.Set("tag_description", tag.Description)
		setNSXTagId(d, tag.ObjectId)

	} else if v, ok := d.GetOk("tag_name"); ok {
		tagName := v.(string)
		foundTag, foundErr := getNSXTagByName(config, tagName)

		if foundErr == nil {
			d.Set("tag_description", foundTag.Description)
			setNSXTagId(d, foundTag.ObjectId)
		} else {
			if d.Get("create_if_missing").(bool) == true {
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

	if resp.StatusCode() != http.StatusCreated {
		log.Printf("%+v", resp)
		return newTag, errors.New(resp.String())
	}

	newTag.ObjectId = resp.String()
	return newTag, nil
}

func resourceNSXTagRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if v, ok := d.GetOk("tag_id"); ok {
		tagId := v.(string)
		_, err := getNSXTagById(config, tagId)

		if err != nil {
			d.SetId("")
			return err
		}
	} else if v, ok := d.GetOk("tag_name"); ok {
		tagName := v.(string)
		_, err := getNSXTagByName(config, tagName)
		if err != nil {
			d.SetId("")
			return err
		}
	} else {
		return fmt.Errorf("no %q or %q specified", "tag_id", "tag_name")
	}

	return nil
}

func resourceNSXTagUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceNSXTagDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("%s", "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! IN DELETE")
	return nil
}
