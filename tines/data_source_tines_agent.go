package tines

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tuckner/go-tines/tines"
)

func dataSourceTinesAgent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTinesAgentRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"story_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTinesAgentRead(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)
	aid := d.Get("id").(int)
	agent, _, err := tinesClient.Agent.Get(aid)
	if err != nil {
		return err
	}

	said := strconv.Itoa(agent.ID)

	d.SetId(said)
	d.Set("guid", agent.GUID)
	d.Set("agent_id", agent.ID)
	d.Set("name", agent.Name)
	d.Set("story_id", agent.StoryID)

	return nil
}
