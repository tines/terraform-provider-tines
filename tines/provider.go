package tines

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/tuckner/go-tines/tines"
)

func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Optional:    false,
				DefaultFunc: schema.EnvDefaultFunc("TINES_TOKEN", nil),
				Description: descriptions["token"],
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    false,
				DefaultFunc: schema.EnvDefaultFunc("TINES_EMAIL", nil),
				Description: descriptions["email"],
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    false,
				DefaultFunc: schema.EnvDefaultFunc("TINES_URL", nil),
				Description: descriptions["base_url"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			// "tines_global_resource": resourceTinesGlobalResource(),
			// "tines_agent":           resourceTinesAgent(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			// "tines_global_resource": dataSourceTinesGlobalResource(),
			"tines_agent": dataSourceTinesAgent(),
		},
	}
	p.ConfigureFunc = providerConfigure(p)

	return p
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"token": "The API token used to connect to Tines. ",

		"base_url": "The Tines Base URL",

		"email": "The Tines user email",
	}
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {

		BaseURL := d.Get("base_url")
		Email := d.Get("email")
		Token := d.Get("token")

		client, err := tines.NewClient(nil, BaseURL, Email, Token)
		if err != nil {
			return nil, err
		}

		return client, err
	}
}
