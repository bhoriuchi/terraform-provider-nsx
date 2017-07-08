package nsx

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty"
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
	"regexp"
)

func getRequest(route string, obj interface{}) error {
	resp, reqErr := resty.R().
		SetResult(&obj).
		Get(route)
	if reqErr != nil {
		return reqErr
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New(resp.String())
	}

	return nil
}

func getNSXSecurityTagByNameRegEx(config *Config, regexStr string) (NSXTag, error) {
	tagList := NSXTagList{}
	tagValue := NSXTag{}
	endpoint := fmt.Sprintf("%s/tag", config.TagEndpoint)
	regex := regexp.MustCompile(regexStr)

	if getErr := getRequest(endpoint, &tagList); getErr != nil {
		return tagValue, getErr
	}

	for _, tag := range tagList.SecurityTags {
		if len(regex.FindAllStringSubmatch(tag.Name, -1)) > 0 {
			return tag, nil

		}
	}

	return tagValue, fmt.Errorf("tag not found using regex %q", regexStr)
}

func getNSXVmById(config *Config, vmId string) (NSXVm, error) {
	vm := NSXVm{}
	tagList := NSXTagList{}
	endpoint := fmt.Sprintf("%s/vm/%s", config.TagEndpoint, vmId)

	if err := getRequest(endpoint, &tagList); err != nil {
		return vm, err
	}

	vm.ObjectId = vmId
	vm.SecurityTags = tagList.SecurityTags

	return vm, nil
}

func getNSXTagByName(config *Config, tagName string) (NSXTag, error) {
	tagList := NSXTagList{}
	tagValue := NSXTag{}
	endpoint := fmt.Sprintf("%s/tag", config.TagEndpoint)

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

func getNSXTagById(config *Config, tagId string) (NSXTag, error) {
	tagValue := NSXTag{}
	endpoint := fmt.Sprintf("%s/tag/%s", config.TagEndpoint, tagId)

	if getErr := getRequest(endpoint, &tagValue); getErr != nil {
		return tagValue, getErr
	}

	return tagValue, nil
}

func setNSXTag(d *schema.ResourceData, tag *NSXTag) {
	d.SetId(tag.ObjectId)
	d.Set("tag_id", tag.ObjectId)
	d.Set("description", tag.Description)
	d.Set("tag_name", tag.Name)
	d.Set("is_universal", tag.IsUniversal)
}
