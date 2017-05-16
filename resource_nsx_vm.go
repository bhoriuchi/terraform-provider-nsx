package main

import (
	"regexp"
	"fmt"
	"errors"
	"net/http"

	"gopkg.in/resty.v0"
	"github.com/hashicorp/terraform/helper/schema"
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
				Computed: true,
			},
		},
	}
}

func resourceNSXVmCreate(d *schema.ResourceData, meta interface{}) error {
	config 		:= meta.(*Config)
	vmId 		:= d.Get("vm_id").(string)

	if v, ok := d.GetOk("security_tags"); ok {
		tagList 	:= NSXTagList{}
		tagIds 		:= []string{}
		endpoint 	:= fmt.Sprintf("%s/tag", config.TagEndpoint)

		if err := getRequest(endpoint, &tagList); err != nil {
			return err
		}
		for _, ti := range v.([]interface{}) {
			if tagId, lookupErr := lookupTagId(tagList.SecurityTags, ti.(string)); lookupErr != nil {
				return lookupErr
			} else if attachErr := attachTag(config, vmId, tagId); attachErr != nil {
				return attachErr
			} else {
				tagIds = append(tagIds, tagId)
			}
		}

		d.Set("security_tags", tagIds)
		d.SetId(vmId)
	}

	return nil
}

func resourceNSXVmRead(d *schema.ResourceData, meta interface{}) error {
	config 		:= meta.(*Config)
	tagList 	:= NSXTagList{}
	tagIds 		:= []string{}
	endpoint 	:= fmt.Sprintf("%s/vm/%s", config.TagEndpoint, d.Id())

	if err := getRequest(endpoint, &tagList); err != nil {
		d.SetId("")
		return err
	}

	for _, tag := range tagList.SecurityTags {
		tagIds = append(tagIds, tag.ObjectId)
	}

	d.Set("security_tags", tagIds)
	d.SetId(d.Id())

	return nil
}

func resourceNSXVmUpdate(d *schema.ResourceData, meta interface{}) error {
	config 		:= meta.(*Config)
	tagList 	:= NSXTagList{}
	allTags		:= NSXTagList{}
	endpoint 	:= fmt.Sprintf("%s/vm/%s", config.TagEndpoint, d.Id())
	tagsEndpoint	:= fmt.Sprintf("%s/tag", config.TagEndpoint)

	if getErr := getRequest(endpoint, &tagList); getErr != nil {
		d.SetId("")
		return getErr
	}

	if v, ok := d.GetOk("security_tags"); ok {
		if err := getRequest(tagsEndpoint, &allTags); err != nil {
			return err
		}

		tagIds 	:= mapTagIds(tagList.SecurityTags)
		tfIds 	:= mapTfIds(allTags.SecurityTags, v)

		// remove tags
		for _, tagId := range tagIds {
			if !stringListContains(tfIds, tagId) {
				if err := detachTag(config, d.Id(), tagId); err != nil {
					return err
				}
			}
		}

		// add tags
		for _, tagId := range tfIds {
			if !stringListContains(tagIds, tagId) {
				if err := attachTag(config, d.Id(), tagId); err != nil {
					return err
				}
			}
		}
	}
	return resourceNSXVmRead(d, meta)
}

func resourceNSXVmDelete(d *schema.ResourceData, meta interface{}) error {
	config 		:= meta.(*Config)
	tagList		:= NSXTagList{}
	endpoint	:= fmt.Sprintf("%s/vm/%s", config.TagEndpoint, d.Id())


	if err := getRequest(endpoint, &tagList); err != nil {
		return err
	}

	for _, tag := range tagList.SecurityTags {
		detach := fmt.Sprintf("%s/tag/%s/vm/%s", config.TagEndpoint, tag.ObjectId, d.Id())
		resp, detachErr := resty.R().Delete(detach)

		if detachErr != nil {
			return detachErr
		} else if resp.StatusCode() != http.StatusNotFound && resp.StatusCode() != http.StatusOK {
			return errors.New(resp.String())
		}
	}

	return nil
}

func attachTag (config *Config, vmId string, tagId string) error {
	endpoint 	:= fmt.Sprintf("%s/tag/%s/vm/%s", config.TagEndpoint, tagId, vmId)

	if resp, err := resty.R().Put(endpoint); err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New(resp.String())
	}
	return nil
}

func detachTag (config *Config, vmId string, tagId string) error {
	endpoint 	:= fmt.Sprintf("%s/tag/%s/vm/%s", config.TagEndpoint, tagId, vmId)

	if resp, err := resty.R().Delete(endpoint); err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New(resp.String())
	}
	return nil
}

func mapTfIds (tagList []NSXTag, list interface{}) []string {
	mapped := []string{}

	for _, item := range list.([]interface{}) {
		if tagId, err := lookupTagId(tagList, item.(string)); err != nil {
			return err
		} else {
			mapped = append(mapped, tagId)
		}
	}
	return mapped
}

func mapTagIds (list []NSXTag) []string {
	mapped := []string{}

	for _, tag := range list {
		mapped = append(mapped, tag.ObjectId)
	}
	return mapped
}

func stringListContains (list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

func lookupTagId (tagList []NSXTag, tag string) (string, error) {
	if isSecurityTagId(tag) {
		return tag, nil
	}

	if tagId, err := lookupTagIdByName(tagList, tag); err != nil {
		return "", err
	} else {
		return tagId, nil
	}
}

func lookupTagIdByName (tagList []NSXTag, tagName string) (string, error) {
	for _, tag := range tagList {
		if tag.Name == tagName {
			return tag.ObjectId, nil
		}
	}
	return "", fmt.Errorf("Security tag %q not found", tagName)
}

func isSecurityTagId (value string) bool {
	matched, err := regexp.MatchString(`^securitytag-\d+$`, value)
	return err == nil && matched == true
}