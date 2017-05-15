package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"gopkg.in/resty.v0"
	"encoding/xml"
	"errors"
	"strings"
	"strconv"
	"fmt"
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
			"nsx_version": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc("NSX_VERSION", "6.3"),
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
	manager := "https://" + d.Get("nsx_manager").(string) + "/api"
	verString := d.Get("nsx_version").(string)
	ver := strings.Split(verString, ".")
	major, magErr := strconv.Atoi(ver[0])
	minor, minErr := strconv.Atoi(ver[1])

	if magErr != nil || minErr != nil || major < 6 || (major == 6 && minor < 2) {
		return nil, fmt.Errorf("Unsupported NSX version %s. NSX 6.2 and higher is required", verString)
	}

	config := Config{
		User: d.Get("user").(string),
		Password: d.Get("password").(string),
		NSXManager: manager,
		NSXVersion: Semver{Major: major, Minor: minor},
		TagEndpoint: manager + "/2.0/services/securitytags/tag",
		InsecureFlag: d.Get("allow_unverified_ssl").(bool),
	}

	return config.Client()
}

func getRequest (route string, obj interface{}) error {
	resp, reqErr := resty.R().Get(route)
	if reqErr != nil {
		return reqErr
	}

	if resp.StatusCode() >= 400 {
		return errors.New(resp.String())
	}

	parseErr := xml.Unmarshal(resp.Body(), &obj)
	if parseErr != nil {
		return parseErr
	}

	return nil
}