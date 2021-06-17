package tines

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tuckner/go-tines/tines"
)

func resourceTinesSendToStory() *schema.Resource {
	return &schema.Resource{
		Create: resourceTinesSendToStoryCreate,
		Read:   resourceTinesSendToStoryRead,
		Update: resourceTinesSendToStoryUpdate,
		Delete: resourceTinesSendToStoryDelete,

		Schema: map[string]*schema.Schema{
			"story_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"send_to_story_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"send_to_story_access": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"entry_agent_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"exit_agent_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"keep_events_for": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  604800,
			},
			"priority": {
				Type:     schema.TypeBool,
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
		},
	}
}

func resourceTinesSendToStoryCreate(d *schema.ResourceData, meta interface{}) error {

	return resourceTinesSendToStoryUpdate(d, meta)
}

func resourceTinesSendToStoryRead(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	sid, _ := strconv.ParseInt(d.Id(), 10, 32)
	story, _, err := tinesClient.Story.Get(int(sid))
	if err != nil {
		return err
	}

	ssid := strconv.Itoa(story.ID)

	d.SetId(ssid)
	d.Set("story_id", story.ID)
	d.Set("name", story.Name)
	d.Set("send_to_story_enabled", story.SendToStoryEnabled)
	d.Set("send_to_story_access", story.SendToStoryAccess)
	d.Set("entry_agent_id", story.EntryAgentID)
	d.Set("exit_agent_ids", story.ExitAgents)
	return nil
}

func resourceTinesSendToStoryDelete(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)
	sid := d.Get("story_id").(int)

	False := false

	s := tines.Story{
		SendToStoryEnabled: &False,
	}

	_, _, err := tinesClient.Story.Update(int(sid), &s)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceTinesSendToStoryUpdate(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)
	sid := d.Get("story_id").(int)
	sendToStoryAccess := d.Get("send_to_story_access").(string)
	entryAgentID := d.Get("entry_agent_id").(int)
	exitAgentsRaw := d.Get("exit_agent_ids").([]interface{})

	exitAgents := make([]int, len(exitAgentsRaw))
	for i, v := range exitAgentsRaw {
		exitAgents[i] = v.(int)
	}

	True := true

	s := tines.Story{
		SendToStoryEnabled: &True,
		SendToStoryAccess:  sendToStoryAccess,
		EntryAgentID:       entryAgentID,
		ExitAgents:         exitAgents,
	}

	story, _, err := tinesClient.Story.Update(int(sid), &s)
	if err != nil {
		return err
	}

	ssid := strconv.Itoa(story.ID)

	d.SetId(ssid)

	return resourceTinesSendToStoryRead(d, meta)
}
