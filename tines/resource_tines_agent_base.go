package tines

import (
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tines/go-tines/tines"
	"github.com/trivago/tgo/tcontainer"
)

func resourceTinesAgentBase() *schema.Resource {
	return &schema.Resource{
		Create: resourceTinesAgentCreate,
		Read:   resourceTinesAgentRead,
		Update: resourceTinesAgentUpdate,
		Delete: resourceTinesAgentDelete,

		Schema: map[string]*schema.Schema{
			"agent_id": {
				Type:     schema.TypeInt,
				Computed: true,
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

func resourceTinesAgentBaseCreate(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	name := d.Get("name").(string)
	agentType := d.Get("agent_type").(string)
	storyID := d.Get("story_id").(int)
	keepEventsFor := d.Get("keep_events_for").(int)
	options := d.Get("agent_options").(string)
	receiveID := make([]int, 0)
	sourceID := make([]int, 0)
	custom := tcontainer.NewMarshalMap()
	custom["options"] = options

	a := tines.Agent{
		Name:          name,
		Type:          agentType,
		StoryID:       storyID,
		KeepEventsFor: keepEventsFor,
		SourceIds:     sourceID,
		ReceiverIds:   receiveID,
		Unknowns:      custom,
	}

	agent, _, err := tinesClient.Agent.Create(&a)
	if err != nil {
		return err
	}

	said := strconv.Itoa(agent.ID)

	d.SetId(said)

	return resourceTinesAgentRead(d, meta)
}

func resourceTinesAgentBaseRead(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	aid, _ := strconv.ParseInt(d.Id(), 10, 32)
	agent, _, err := tinesClient.Agent.Get(int(aid))
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
	d.Set("name", agent.Name)
	d.Set("guid", agent.GUID)
	d.Set("agent_id", agent.ID)
	d.Set("story_id", agent.StoryID)
	d.Set("user_id", agent.UserID)
	d.Set("agent_type", d.Get("agent_type").(string))

	return nil
}

func resourceTinesAgentBaseDelete(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)
	aid, _ := strconv.ParseInt(d.Id(), 10, 32)
	_, err := tinesClient.Agent.Delete(int(aid))
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

func resourceTinesAgentBaseUpdate(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	aid, _ := strconv.ParseInt(d.Id(), 10, 32)
	name := d.Get("name").(string)
	agentType := d.Get("agent_type").(string)
	storyID := d.Get("story_id").(int)
	keepEventsFor := d.Get("keep_events_for").(int)
	options := d.Get("agent_options").(string)
	receiveID := make([]int, 0)
	sourceID := make([]int, 0)
	custom := tcontainer.NewMarshalMap()
	custom["options"] = options

	a := tines.Agent{
		Name:          name,
		Type:          agentType,
		StoryID:       storyID,
		KeepEventsFor: keepEventsFor,
		SourceIds:     sourceID,
		ReceiverIds:   receiveID,
		Unknowns:      custom,
	}

	agent, _, err := tinesClient.Agent.Update(int(aid), &a)
	if err != nil {
		return err
	}

	said := strconv.Itoa(agent.ID)

	d.SetId(said)

	return resourceTinesAgentRead(d, meta)
}
