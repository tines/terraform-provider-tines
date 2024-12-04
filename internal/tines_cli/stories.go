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
	ID                   int64    `json:"id,omitempty"`
	Name                 string   `json:"name,omitempty"`
	UserID               int64    `json:"user_id,omitempty"`
	Description          string   `json:"description,omitempty"`
	KeepEventsFor        int64    `json:"keep_events_for,omitempty"`
	Disabled             bool     `json:"disabled,omitempty"`
	Priority             bool     `json:"priority,omitempty"`
	STSEnabled           bool     `json:"send_to_story_enabled,omitempty"`
	STSAccessSource      string   `json:"send_to_story_access_source,omitempty"`
	STSAccess            string   `json:"send_to_story_access,omitempty"`
	STSSkillConfirmation bool     `json:"send_to_story_skill_use_requires_confirmation,omitempty"`
	SharedTeamSlugs      []string `json:"shared_team_slugs,omitempty"`
	EntryAgentID         int64    `json:"entry_agent_id,omitempty"`
	ExitAgents           []int64  `json:"exit_agents,omitempty"`
	TeamID               int64    `json:"team_id,omitempty"`
	Tags                 []string `json:"tags,omitempty"`
	Guid                 string   `json:"guid,omitempty"`
	Slug                 string   `json:"slug,omitempty"`
	CreatedAt            string   `json:"created_at,omitempty"`
	UpdatedAt            string   `json:"updated_at,omitempty"`
	EditedAt             string   `json:"edited_at,omitempty"`
	Mode                 string   `json:"mode,omitempty"`
	FolderID             int64    `json:"folder_id,omitempty"`
	Published            bool     `json:"published,omitempty"`
	ChangeControlEnabled bool     `json:"change_control_enabled,omitempty"`
	Locked               bool     `json:"locked,omitempty"`
	Owners               []int64  `json:"owners,omitempty"`
}

// Create a new story.
func (c *Client) CreateStory(s *Story) (*Story, error) {
	newStory := Story{}

	req, err := json.Marshal(&s)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest("POST", "/api/v1/stories", req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &newStory)
	if err != nil {
		return nil, err
	}

	return &newStory, nil
}

// Delete a story.
func (c *Client) DeleteStory(id int64) error {
	resource := fmt.Sprintf("/api/v1/stories/%d", id)

	_, err := c.doRequest("DELETE", resource, nil)

	return err
}

// Get current state for a story.
func (c *Client) GetStory(id int64) (story *Story, e error) {
	resource := fmt.Sprintf("/api/v1/stories/%d", id)

	body, err := c.doRequest("GET", resource, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &story)
	if err != nil {
		return nil, err
	}

	return story, nil
}

// Import a new story, or override an existing one.
func (c *Client) ImportStory(story *StoryImportRequest) (*Story, error) {
	newStory := Story{}

	req, err := json.Marshal(&story)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest("POST", "/api/v1/stories/import", req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &newStory)
	if err != nil {
		return nil, err
	}

	return &newStory, nil
}

// Update a story.
func (c *Client) UpdateStory(id int64, values *Story) (*Story, error) {
	updatedStory := Story{}
	resource := fmt.Sprintf("/api/v1/stories/%d", id)

	req, err := json.Marshal(&values)
	if err != nil {
		return &updatedStory, err
	}

	body, err := c.doRequest("PUT", resource, req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &updatedStory)
	if err != nil {
		return nil, err
	}

	return &updatedStory, nil
}
