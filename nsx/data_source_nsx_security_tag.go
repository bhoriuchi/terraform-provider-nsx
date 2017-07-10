package nsx

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNSXSecurityTag() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNSXSecurityTagRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Regular expression to use when searching for a tag by name",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security tag name",
			},
			"is_universal": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is a universal security tag",
			},
			"type_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Object type name",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of security tag",
			},
			"vm_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of VMs attached to the security group",
			},
			"revision": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Revision number of the security tag",
			},
		},
	}
}

func dataSourceNSXSecurityTagRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	rxStr := d.Get("name_regex").(string)

	if tag, err := getNSXSecurityTagByNameRegEx(config, rxStr); err != nil {
		return err
	} else {
		d.SetId(tag.ObjectId)
		d.Set("name", tag.Name)
		d.Set("is_universal", tag.IsUniversal)
		d.Set("description", tag.Description)
		d.Set("type_name", tag.ObjectTypeName)
		d.Set("vm_count", tag.VmCount)
		d.Set("revision", tag.Revision)
	}

	return nil
}
