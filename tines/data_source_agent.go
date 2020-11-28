package tines

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tuckner/go-tines/tines"
)

func dataSourceTinesAgent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTinesAgentRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
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
	agentID := d.Get("id")
	log.Printf("[INFO] Reading AgentID: %s", agentID)

	agent, resp, err := tinesClient.Agent.Get(agentID)
	if err != nil {
		return err
	}

	d.SetId("guid", agent.guid)
	d.Set("agent_id", agent.id)
	d.set("name", agent.name)
	d.set("story_id", agent.story_id)

	return nil
}
