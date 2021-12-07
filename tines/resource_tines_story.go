package tines

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tuckner/go-tines/tines"
)

func resourceTinesStory() *schema.Resource {
	return &schema.Resource{
		Create: resourceTinesStoryCreate,
		Read:   resourceTinesStoryRead,
		Update: resourceTinesStoryUpdate,
		Delete: resourceTinesStoryDelete,

		Schema: map[string]*schema.Schema{
			"story_id": {
				Type:     schema.TypeInt,
				Computed: true,
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
			"story_to_story_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"story_to_story_access": {
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
			"published_state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTinesStoryCreate(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	keepEventsFor := d.Get("keep_events_for").(int)
	teamID := d.Get("team_id").(int)
	disabled := d.Get("disabled").(bool)
	priority := d.Get("priority").(bool)
	folderID := d.Get("folder_id").(int)
	publishedState := d.Get("published_state").(string)

	s := tines.Story{
		Name:           name,
		Description:    description,
		KeepEventsFor:  keepEventsFor,
		TeamID:         teamID,
		Disabled:       &disabled,
		Priority:       &priority,
		FolderID:       folderID,
		PublishedState: publishedState,
	}

	story, _, err := tinesClient.Story.Create(&s)
	if err != nil {
		return err
	}

	ssid := strconv.Itoa(story.ID)

	d.SetId(ssid)

	return resourceTinesStoryRead(d, meta)
}

func resourceTinesStoryRead(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	sid, _ := strconv.ParseInt(d.Id(), 10, 32)
	story, _, err := tinesClient.Story.Get(int(sid))
	if err != nil {
		return err
	}

	ssid := strconv.Itoa(story.ID)

	d.SetId(ssid)
	d.Set("user_id", story.UserID)
	d.Set("story_id", story.ID)
	d.Set("name", story.Name)
	d.Set("description", story.Description)
	d.Set("send_to_story_enabled", story.SendToStoryEnabled)
	d.Set("send_to_story_access", story.SendToStoryAccess)
	d.Set("entry_agent_id", story.EntryAgentID)
	d.Set("exit_agent_ids", story.ExitAgents)
	d.Set("disabled", story.Disabled)
	d.Set("keep_events_for", story.KeepEventsFor)
	d.Set("priority", story.Priority)
	d.Set("team_id", story.TeamID)
	d.Set("folder_id", story.FolderID)
	d.Set("slug", story.Slug)
	d.Set("published_state", story.PublishedState)

	return nil
}

func resourceTinesStoryDelete(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)
	sid, _ := strconv.ParseInt(d.Id(), 10, 32)
	_, err := tinesClient.Story.Delete(int(sid))
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

func resourceTinesStoryUpdate(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	keepEventsFor := d.Get("keep_events_for").(int)
	teamID := d.Get("team_id").(int)
	disabled := d.Get("disabled").(bool)
	priority := d.Get("priority").(bool)
	folderID := d.Get("folder_id").(int)
	publishedState := d.Get("published_state").(string)
	sid, _ := strconv.ParseInt(d.Id(), 10, 32)

	s := tines.Story{
		Name:           name,
		Description:    description,
		KeepEventsFor:  keepEventsFor,
		TeamID:         teamID,
		Disabled:       &disabled,
		Priority:       &priority,
		FolderID:       folderID,
		PublishedState: publishedState,
	}

	story, _, err := tinesClient.Story.Update(int(sid), &s)
	if err != nil {
		return err
	}

	ssid := strconv.Itoa(story.ID)

	d.SetId(ssid)

	return resourceTinesStoryRead(d, meta)
}
