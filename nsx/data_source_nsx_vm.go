package nsx

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNSXVm() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNSXVmRead,

		Schema: map[string]*schema.Schema{
			"vm_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VM id",
			},
			"security_tag_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Security tag name",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"security_tag_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Security tag name",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceNSXVmRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vmId := d.Get("vm_id").(string)
	tagIds := make([]string, 0)
	tagNames := make([]string, 0)

	if vm, err := getNSXVmById(config, vmId); err != nil {
		return err
	} else {
		d.SetId(vm.ObjectId)

		for _, t := range vm.SecurityTags {
			tagIds = append(tagIds, t.ObjectId)
			tagNames = append(tagNames, t.Name)
		}
		d.Set("security_tag_ids", tagIds)
		d.Set("security_tag_names", tagNames)
	}

	return nil
}
