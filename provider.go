package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("NSX_USER", nil),
				Description: "The user name for NSX API operations.",
			},
			"password": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("NSX_PASSWORD", nil),
				Description: "The user password for NSX API operations.",
			},
			"nsx_manager": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("NSX_MANAGER", nil),
			},
			"allow_unverified_ssl": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSX_ALLOW_UNVERIFIED_SSL", false),
				Description: "If set, VMware vSphere client will permit unverifiable SSL certificates.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"nsx_tag": resourceNSXTag(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		User: d.Get("user").(string),
		Password: d.Get("password").(string),
		NSXManager: "https://" + d.Get("nsx_manager").(string) + "/api",
		InsecureFlag: d.Get("allow_unverified_ssl").(bool),
	}

	return config.Client()
}