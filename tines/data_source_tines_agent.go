package tines

import (
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tines/go-tines/tines"
)

func dataSourceTinesAgent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTinesAgentRead,

		Schema: map[string]*schema.Schema{
			"agent_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"story_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"agent_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"agent_options": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"keep_events_for": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  604800,
			},
			"cron": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"source_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"receiver_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func dataSourceTinesAgentRead(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)
	aid := d.Get("id").(int)
	agent, _, err := tinesClient.Agent.Get(aid)
	if err != nil {
		log.Printf("[DEBUG] Error: %v", err)
		if strings.Contains(err.Error(), "404") {
			log.Printf("[INFO] Action %v no longer exists", d.Id())
			d.SetId("")
			return nil
		} else {
			return err
		}
	}

	said := strconv.Itoa(agent.ID)

	d.SetId(said)
	d.Set("guid", agent.GUID)
	d.Set("agent_id", agent.ID)
	d.Set("name", agent.Name)
	d.Set("story_id", agent.StoryID)

	return nil
}
