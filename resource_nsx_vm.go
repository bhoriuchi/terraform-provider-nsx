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
	config := meta.(*Config)
	vmId := d.Get("vm_id").(string)

	if v, ok := d.GetOk("security_tags"); ok {
		tagList := NSXTagList{}
		err := getRequest(fmt.Sprintf("%s/tag", config.TagEndpoint), &tagList)

		if err != nil {
			return err
		}

		tagIds := []string{}

		for _, ti := range v.([]interface{}) {
			tag := ti.(string)

			if isSecurityTagId(tag) {
				err := putTag(config, vmId, tag)
				if err != nil {
					return err
				}
				tagIds = append(tagIds, tag)
			} else {
				tagId, le := lookupTagIdByName(tagList, tag)

				if le != nil {
					return err
				}

				err := putTag(config, vmId, tagId)
				if err != nil {
					return err
				}

				tagIds = append(tagIds, tagId)
			}
		}

		d.Set("security_tags", tagIds)
		d.SetId(vmId)
	}

	return nil
}

func resourceNSXVmRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	tagList := NSXTagList{}

	err := getRequest(fmt.Sprintf("%s/vm/%s", config.TagEndpoint, d.Id()), &tagList)

	if err != nil {
		return err
	}

	tagIds := []string{}

	for _, tag := range tagList.SecurityTags {
		tagIds = append(tagIds, tag.ObjectId)
	}

	d.Set("security_tags", tagIds)
	d.SetId(d.Id())

	return nil
}

func resourceNSXVmUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceNSXVmDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func putTag (config *Config, vmId string, tagId string) error {
	resp, err := resty.R().
		Put(fmt.Sprintf("%s/tag/%s/vm/%s", config.TagEndpoint, tagId, vmId))

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New(resp.String())
	}

	return nil
}

func lookupTagIdByName (tagList NSXTagList, tagName string) (string, error) {
	for _, tag := range tagList.SecurityTags {
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