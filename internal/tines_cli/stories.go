package tines_cli

import (
	"encoding/json"
	"fmt"
)

type StoryImportRequest struct {
	NewName  string                 `json:"new_name"`
	Data     map[string]interface{} `json:"data"`
	TeamID   int64                  `json:"team_id"`
	FolderID int64                  `json:"folder_id,omitempty"`
	Mode     string                 `json:"mode"`
}

type Story struct {
	ID                   int64    `json:"id"`
	Name                 string   `json:"name"`
	UserID               int64    `json:"user_id"`
	Description          string   `json:"description"`
	KeepEventsFor        int64    `json:"keep_events_for"`
	Disabled             bool     `json:"disabled"`
	Priority             bool     `json:"priority"`
	STSEnabled           bool     `json:"send_to_story_enabled"`
	STSAccessSource      string   `json:"send_to_story_access_source"`
	STSAccess            string   `json:"send_to_story_access"`
	STSSkillConfirmation bool     `json:"send_to_story_skill_use_requires_confirmation"`
	SharedTeamSlugs      []string `json:"shared_team_slugs,omitempty"`
	EntryAgentID         int64    `json:"entry_agent_id,omitempty"`
	ExitAgents           []int64  `json:"exit_agents,omitempty"`
	TeamID               int64    `json:"team_id"`
	Tags                 []string `json:"tags,omitempty"`
	Guid                 string   `json:"guid"`
	Slug                 string   `json:"slug"`
	CreatedAt            string   `json:"created_at"`
	UpdatedAt            string   `json:"updated_at"`
	EditedAt             string   `json:"edited_at"`
	Mode                 string   `json:"mode"`
	FolderID             int64    `json:"folder_id"`
	Published            bool     `json:"published"`
	ChangeControlEnabled bool     `json:"change_control_enabled"`
	Locked               bool     `json:"locked"`
	Owners               []int64  `json:"owners"`
}

// Import a new story, or update an existing one.
func (c *Client) ImportStory(story *StoryImportRequest) (*Story, error) {
	newStory := Story{}

	req, err := json.Marshal(&story)
	if err != nil {
		return &newStory, err
	}

	_, body, err := c.doRequest("POST", "/api/v1/stories/import", req)
	if err != nil {
		return &newStory, err
	}

	err = json.Unmarshal(body, &newStory)
	if err != nil {
		return &newStory, err
	}

	return &newStory, nil
}

// Get current state for a story.
func (c *Client) GetStory(id int64) (status int, story *Story, e error) {
	resource := fmt.Sprintf("/api/v1/stories/%d", id)

	status, body, err := c.doRequest("GET", resource, nil)
	if err != nil || status == 404 {
		return status, nil, err
	}

	err = json.Unmarshal(body, &story)
	if err != nil {
		return status, nil, err
	}

	return status, story, err
}

// Delete a story.
func (c *Client) DeleteStory(id int64) error {
	resource := fmt.Sprintf("/api/v1/stories/%d", id)

	_, _, err := c.doRequest("DELETE", resource, nil)

	return err
}
