package tines

import (
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tines/go-tines/tines"
)

func dataSourceTinesStory() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTinesStoryRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Optional: true,
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
			"story_to_story": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"entry_agent_id": {
				Type:     schema.TypeInt,
				Optional: true,
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
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceTinesStoryRead(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)
	sid := d.Get("id").(int)
	story, _, err := tinesClient.Story.Get(sid)
	if err != nil {
		log.Printf("[DEBUG] Error: %v", err)
		if strings.Contains(err.Error(), "404") {
			log.Printf("[INFO] Story %v no longer exists", d.Id())
			d.SetId("")
			return nil
		} else {
			return err
		}
	}
	ssid := strconv.Itoa(story.ID)

	d.SetId(ssid)
	d.Set("user_id", story.UserID)
	d.Set("sid", story.ID)
	d.Set("name", story.Name)
	d.Set("description", story.Description)
	d.Set("send_to_story", story.SendToStoryEnabled)
	d.Set("entry_agent_id", story.EntryAgentID)
	d.Set("disabled", story.Disabled)
	d.Set("keep_events_for", story.KeepEventsFor)
	d.Set("priority", story.Priority)
	d.Set("team_id", story.TeamID)
	d.Set("folder_id", story.FolderID)

	return nil
}
