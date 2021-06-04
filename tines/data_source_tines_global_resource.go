package tines

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tuckner/go-tines/tines"
)

func dataSourceTinesGlobalResource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTinesGlobalResourceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"read_access": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"team_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"folder_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"global_resource_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceTinesGlobalResourceRead(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	grid, _ := strconv.ParseInt(d.Id(), 10, 32)
	globalresource, _, err := tinesClient.GlobalResource.Get(int(grid))
	if err != nil {
		return err
	}

	sgrid := strconv.Itoa(globalresource.ID)

	d.SetId(sgrid)
	d.Set("name", globalresource.Name)
	d.Set("value", globalresource.Value)
	d.Set("team_id", globalresource.TeamID)
	d.Set("read_access", globalresource.ReadAccess)
	d.Set("folder_id", globalresource.FolderID)
	d.Set("global_resource_id", globalresource.ID)

	return nil
}
