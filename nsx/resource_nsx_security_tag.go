package nsx

import (
	"fmt"
	"log"
	"net/http"

	"errors"
	"github.com/go-resty/resty"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNSXSecurityTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceNSXSecurityTagCreate,
		Read:   resourceNSXSecurityTagRead,
		Update: resourceNSXSecurityTagUpdate,
		Delete: resourceNSXSecurityTagDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Security tag name.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Security tag description.",
			},
			"is_universal": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Creates the security tag as a universal tag (NSX 6.3 and higher).",
			},
			"persistent": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, prevents removal of the security tag during a destroy operation.",
			},
			"safe_destroy": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "When true, prevents removal of the security tag if one or more virtual machines are attached to it.",
			},
		},
	}
}

func resourceNSXSecurityTagCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	endpoint := fmt.Sprintf("%s/tag", config.TagEndpoint)

	newTag := NSXTag{
		ObjectTypeName: "SecurityTag",
		Type:           NSXTagType{TypeName: "SecurityTag"},
		Name:           d.Get("tag_name").(string),
		IsUniversal:    d.Get("is_universal").(bool),
		Description:    d.Get("description").(string),
	}

	if resp, err := resty.R().SetBody(newTag).Post(endpoint); err != nil {
		return err
	} else if resp.StatusCode() != http.StatusCreated {
		log.Printf("[NSXLOG] %+v", resp)
		return errors.New(resp.String())
	} else {
		d.SetId(resp.String())
	}

	return nil
}

func resourceNSXSecurityTagRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if tag, err := getNSXTagById(config, d.Id()); err != nil {
		d.SetId("")
		return err
	} else {
		setNSXTag(d, &tag)
	}

	return nil
}

func resourceNSXSecurityTagUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	endpoint := fmt.Sprintf("%s/tag/%s", config.TagEndpoint, d.Id())

	if d.HasChange("tag_name") || d.HasChange("description") {
		tag := NSXTag{
			ObjectId:       d.Id(),
			ObjectTypeName: "SecurityTag",
			Type:           NSXTagType{TypeName: "SecurityTag"},
			Name:           d.Get("tag_name").(string),
			IsUniversal:    d.Get("is_universal").(bool),
			Description:    d.Get("description").(string),
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

func resourceNSXSecurityTagDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	endpoint := fmt.Sprintf("%s/tag/%s", config.TagEndpoint, d.Id())
	tag := NSXTag{}

	if d.Get("persistent").(bool) != true {
		if d.Get("safe_destroy").(bool) == true {
			if err := getRequest(endpoint, &tag); err != nil {
				return err
			} else if tag.VmCount > 0 {
				log.Printf("[NSXLOG] cannot safely destroy %q, %d vms are still attached", d.Id(), tag.VmCount)
				return nil
			}
		}

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
