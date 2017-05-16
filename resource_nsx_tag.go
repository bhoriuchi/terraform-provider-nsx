package main

import (
	"fmt"
	"log"
	"errors"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"gopkg.in/resty.v0"
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
	config 		:= meta.(*Config)

	if v, ok := d.GetOk("tag_id"); ok {
		if tag, err := getNSXTagById(config, v.(string)); err != nil {
			return err
		} else {
			d.Set("tag_description", tag.Description)
			setNSXTag(d, &tag)
		}
	} else if v, ok := d.GetOk("tag_name"); ok {
		if foundTag, foundErr := getNSXTagByName(config, v.(string)); foundErr == nil {
			d.Set("tag_description", foundTag.Description)
			setNSXTag(d, &foundTag)
		} else {
			if d.Get("create_if_missing").(bool) == true {
				if newTag, createErr := createNSXTag(d, meta); createErr != nil {
					return createErr
				} else {
					d.Set("tag_description", newTag.Description)
					setNSXTag(d, &newTag)
				}
			} else {
				return fmt.Errorf("%q not found and \"create_if_missing\" is false", v.(string))
			}
		}
	} else {
		return fmt.Errorf("must provide either %q or %q", "tag_id", "tag_name")
	}

	return nil
}

func resourceNSXTagRead(d *schema.ResourceData, meta interface{}) error {
	config 		:= meta.(*Config)

	if tag, err := getNSXTagById(config, d.Id()); err != nil {
		d.SetId("")
		return err
	} else {
		setNSXTag(d, &tag)
	}

	return nil
}

func resourceNSXTagUpdate(d *schema.ResourceData, meta interface{}) error {
	config 		:= meta.(*Config)
	endpoint 	:= fmt.Sprintf("%s/tag/%s", config.TagEndpoint, d.Id())

	if d.HasChange("tag_name") || d.HasChange("description") {
		tag := NSXTag{
			ObjectId: d.Id(),
			ObjectTypeName: "SecurityTag",
			Type: NSXTagType{ TypeName: "SecurityTag" },
			Name: d.Get("tag_name").(string),
			IsUniversal: d.Get("is_universal").(bool),
			Description: d.Get("description").(string),
		}

		if resp, err := resty.R().SetBody(tag).Put(endpoint); err != nil {
			return err
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("%s", resp.String())
		}
	}

	d.Set("tag_id", d.Id())

	return nil
}

func resourceNSXTagDelete(d *schema.ResourceData, meta interface{}) error {
	config 		:= meta.(*Config)
	endpoint 	:= fmt.Sprintf("%s/tag/%s", config.TagEndpoint, d.Id())

	if d.Get("persistent").(bool) != true {
		if resp, err := resty.R().Delete(endpoint); err != nil {
			return err
		} else if resp.StatusCode() == http.StatusNotFound {
			log.Printf("[NSXLOG] could not tag %q during destroy, considering it destroyed", d.Id())
		}
	} else {
		log.Printf("[NSXLOG] %q is persistent and will NOT be destroyed", d.Id())
	}

	return nil
}

func getNSXTagByName (config *Config, tagName string) (NSXTag, error) {
	tagList 	:= NSXTagList{}
	tagValue 	:= NSXTag{}
	endpoint 	:= fmt.Sprintf("%s/tag", config.TagEndpoint)

	if getErr := getRequest(endpoint, &tagList); getErr != nil {
		return tagValue, getErr
	}

	for _, tag := range tagList.SecurityTags {
		if tag.Name == tagName {
			return tag, nil

		}
	}

	return tagValue, fmt.Errorf("tag %q not found", tagName)
}

func getNSXTagById (config *Config, tagId string) (NSXTag, error) {
	tagValue 	:= NSXTag{}
	endpoint 	:= fmt.Sprintf("%s/tag/%s", config.TagEndpoint, tagId)

	if getErr := getRequest(endpoint, &tagValue); getErr != nil {
		return tagValue, getErr
	}

	return tagValue, nil
}

func createNSXTag(d *schema.ResourceData, meta interface{}) (NSXTag, error) {
	config 		:= meta.(*Config)
	endpoint 	:= fmt.Sprintf("%s/tag", config.TagEndpoint)

	newTag := NSXTag{
		ObjectTypeName: "SecurityTag",
		Type: NSXTagType{ TypeName: "SecurityTag" },
		Name: d.Get("tag_name").(string),
		IsUniversal: d.Get("is_universal").(bool),
		Description: d.Get("description").(string),
	}

	if resp, err := resty.R().SetBody(newTag).Post(endpoint); err != nil {
		return newTag, err
	} else if resp.StatusCode() != http.StatusCreated {
		log.Printf("[NSXLOG] %+v", resp)
		return newTag, errors.New(resp.String())
	} else {
		newTag.ObjectId = resp.String()
		return newTag, nil
	}
}

func setNSXTag (d *schema.ResourceData, tag *NSXTag) {
	d.SetId(tag.ObjectId)
	d.Set("tag_id", tag.ObjectId)
	d.Set("description", tag.Description)
	d.Set("tag_name", tag.Name)
	d.Set("is_universal", tag.IsUniversal)
}
