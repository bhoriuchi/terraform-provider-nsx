package main

import (
	"fmt"
	"log"
	"errors"
	"net/http"
	"encoding/xml"

	"github.com/hashicorp/terraform/helper/schema"
	"gopkg.in/resty.v0"
)

type NSXTagType struct {
	TypeName string `xml:"typeName"`
}

type NSXTag struct {
	XMLName xml.Name `xml:"securityTag"`
	ObjectId string `xml:"objectId"`
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
		setNSXTag(d, &tag)

	} else if v, ok := d.GetOk("tag_name"); ok {
		tagName := v.(string)
		foundTag, foundErr := getNSXTagByName(config, tagName)

		if foundErr == nil {
			d.Set("tag_description", foundTag.Description)
			setNSXTag(d, &foundTag)
		} else {
			if d.Get("create_if_missing").(bool) == true {
				newTag, createErr := createNSXTag(d, meta)
				if createErr != nil {
					return createErr
				}
				d.Set("tag_description", newTag.Description)
				setNSXTag(d, &newTag)

			} else {
				return fmt.Errorf("Security Tag %q not found and %q is set to false", tagName, "create_if_missing")
			}
		}
	} else {
		return fmt.Errorf("must provide either %q or %q", "tag_id", "tag_name")
	}

	return nil
}

func setNSXTag (d *schema.ResourceData, tag *NSXTag) {
	d.SetId(tag.ObjectId)
	d.Set("tag_id", tag.ObjectId)
	d.Set("description", tag.Description)
	d.Set("tag_name", tag.Name)
	d.Set("is_universal", tag.IsUniversal)
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

	newTag := NSXTag{
		ObjectTypeName: "SecurityTag",
		Type: NSXTagType{ TypeName: "SecurityTag" },
		Name: d.Get("tag_name").(string),
		IsUniversal: d.Get("is_universal").(bool),
		Description: d.Get("description").(string),
	}

	resp, err := resty.R().
		SetBody(newTag).
		Post(config.TagEndpoint)

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

	tag, err := getNSXTagById(config, d.Id())

	if err != nil {
		d.SetId("")
		return err
	}

	setNSXTag(d, &tag)

	return nil
}

func resourceNSXTagUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if d.HasChange("tag_name") || d.HasChange("description") {
		tag := NSXTag{
			ObjectId: d.Id(),
			ObjectTypeName: "SecurityTag",
			Type: NSXTagType{ TypeName: "SecurityTag" },
			Name: d.Get("tag_name").(string),
			IsUniversal: d.Get("is_universal").(bool),
			Description: d.Get("description").(string),
		}

		resp, err := resty.R().
			SetBody(tag).
			Put(fmt.Sprintf("%s/%s", config.TagEndpoint, d.Id()))

		if err != nil {
			return err
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("%s", resp.String())
		}
	}

	d.Set("tag_id", d.Id())

	return nil
}

func resourceNSXTagDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if d.Get("persistent").(bool) != true {
		resp, err := resty.R().
			Delete(fmt.Sprintf("%s/%s", config.TagEndpoint, d.Id()))

		if err != nil {
			return err
		}

		if resp.StatusCode() == http.StatusNotFound {
			log.Printf("[NSXLOG] could not find security tag %q during destroy operation, considering it destroyed", d.Id())
		}
	} else {
		log.Printf("[NSXLOG] %q is persistent and will NOT be destroyed", d.Id())
	}

	return nil
}
