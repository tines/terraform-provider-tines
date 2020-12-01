package tines

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tuckner/go-tines/tines"
)

func dataSourceTinesStory() *schema.Resource {
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

func dataSourceTinesStoryRead(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)
	sid := d.Get("id").(int)
	story, _, err := tinesClient.Story.Get(sid)
	if err != nil {
		return err
	}

	ssid := strconv.Itoa(story.ID)

	d.SetId(ssid)
	d.Set("user_id", story.User_ID)
	d.Set("story_id", story.ID)
	d.Set("name", story.Name)
	d.Set("description", story.Description)
	d.Set("send_to_story", story.SendToStory)
	d.Set("entry_agent_id", story.EntryAgentID)
	d.Set("disabled", story.Disabled)
	d.Set("keep_events_for", story.KeepEvents_for)
	d.Set("priority", story.Priority)
	d.Set("team_id", story.TeamID)
	d.Set("folder_id", story.FolderID)

	return nil
}
