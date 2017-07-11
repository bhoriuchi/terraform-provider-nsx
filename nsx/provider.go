package nsx

import (
	"fmt"

	"github.com/coreos/go-semver/semver"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSX_USER", nil),
				Description: "The user name for NSX API operations.",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSX_PASSWORD", nil),
				Description: "The user password for NSX API operations.",
			},
			"nsx_manager": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSX_MANAGER", nil),
			},
			"nsx_version": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSX_VERSION", "6.3"),
			},
			"allow_unverified_ssl": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSX_ALLOW_UNVERIFIED_SSL", false),
				Description: "If set, VMware vSphere client will permit unverifiable SSL certificates.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"nsx_security_tag": dataSourceNSXSecurityTag(),
			"nsx_vm":           dataSourceNSXVm(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"nsx_security_tag": resourceNSXSecurityTag(),
			"nsx_vm":           resourceNSXVm(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	manager := "https://" + d.Get("nsx_manager").(string) + "/api"
	verString := fmt.Sprintf("%s.0", d.Get("nsx_version").(string))
	version, err := semver.NewVersion(verString)
	minVersion := semver.New("6.2.0")

	if err != nil {
		return nil, err
	}

	if version.LessThan(*minVersion) {
		return nil, fmt.Errorf("Unsupported NSX version %s. NSX 6.2.0 and higher is required", verString)
	}

	config := Config{
		User:         d.Get("user").(string),
		Password:     d.Get("password").(string),
		NSXManager:   manager,
		NSXVersion:   *version,
		TagEndpoint:  manager + "/2.0/services/securitytags",
		InsecureFlag: d.Get("allow_unverified_ssl").(bool),
	}

	return config.ClientInit()
}
