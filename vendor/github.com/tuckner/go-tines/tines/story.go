package tines

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"

	"fmt"
	"time"

	"github.com/fatih/structs"

	"github.com/trivago/tgo/tcontainer"
)

// StoryService handles fields for the Tines instance / API.
type StoryService struct {
	client *Client
}

// Story structure
type Story struct {
	ID            int         `json:"id" structs:"id,omitempty"`
	UserID        int         `json:"user_id" structs:"user_id,omitempty"`
	Name          string      `json:"name" structs:"name,omitempty"`
	CreatedAt     time.Time   `json:"created_at" structs:"created_at,omitempty"`
	UpdatedAt     time.Time   `json:"updated_at" structs:"updated_at,omitempty"`
	Description   string      `json:"description" structs:"description,omitempty"`
	GUID          string      `json:"guid" structs:"guid,omitempty"`
	SendToStory   bool        `json:"send_to_story_enabled" structs:"send_to_story_enabled,omitempty"`
	EntryAgentID  int         `json:"entry_agent_id" structs:"exit_agent_id,omitempty"`
	DiagramLayout interface{} `json:"diagram_layout" structs:"diagram_layout,omitempty"`
	Disabled      bool        `json:"disabled" structs:"disabled,omitempty"`
	KeepEventsFor int         `json:"keep_events_for" structs:"keep_events_for,omitempty"`
	Priority      bool        `json:"priority" structs:"priority,omitempty"`
	TeamID        int         `json:"team_id" structs:"team_id,omitempty"`
	FolderID      string      `json:"folder_id" structs:"folder_id,omitempty"`
	Unknowns      tcontainer.MarshalMap
}

// MarshalJSON is a custom JSON marshal function for the Story* structs.
// It handles Tines options and maps those from / to "Unknowns" key.
func (i *Story) MarshalJSON() ([]byte, error) {
	m := structs.Map(i)
	unknowns, okay := m["Unknowns"]
	if okay {
		// if unknowns present, shift all key value from unknown to a level up
		for key, value := range unknowns.(tcontainer.MarshalMap) {
			m[key] = value
		}
		delete(m, "Unknowns")
	}
	return json.Marshal(m)
}

// UnmarshalJSON is a custom JSON marshal function for the Story structs.
// It handles Tines custom options and maps those from / to "Unknowns" key.
func (i *Story) UnmarshalJSON(data []byte) error {

	// Do the normal unmarshalling first
	// Details for this way: http://choly.ca/post/go-json-marshalling/
	type Alias Story
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(i),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	totalMap := tcontainer.NewMarshalMap()
	err := json.Unmarshal(data, &totalMap)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(*i)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagDetail := field.Tag.Get("json")
		if tagDetail == "" {
			// ignore if there are no tags
			continue
		}
		options := strings.Split(tagDetail, ",")

		if len(options) == 0 {
			return fmt.Errorf("no tags options found for %s", field.Name)
		}
		// the first one is the json tag
		key := options[0]
		if _, okay := totalMap.Value(key); okay {
			delete(totalMap, key)
		}

	}
	i = (*Story)(aux.Alias)
	// all the tags found in the struct were removed. Whatever is left are unknowns to struct
	i.Unknowns = totalMap
	return nil

}

// GetWithContext returns an agent for the given global resource id.
func (s *StoryService) GetWithContext(ctx context.Context, storyID int) (*Story, *Response, error) {
	apiEndpoint := fmt.Sprintf("stories/%v", storyID)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	story := new(Story)
	resp, err := s.client.Do(req, story)
	if err != nil {
		return nil, resp, err
	}

	return story, resp, nil
}

// Get wraps GetWithContext using the background context.
func (s *StoryService) Get(storyID int) (*Story, *Response, error) {
	return s.GetWithContext(context.Background(), storyID)
}

// StoryWithContext deletes a story.
func (s *StoryService) DeleteWithContext(ctx context.Context, storyID int) (*Response, error) {
	apiEndpoint := fmt.Sprintf("stories/%v", storyID)
	req, err := s.client.NewRequestWithContext(ctx, "DELETE", apiEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Delete wraps DeleteWithContext using the background context.
func (s *StoryService) Delete(storyID int) (*Response, error) {
	return s.DeleteWithContext(context.Background(), storyID)
}

// // CreateWithContext creates a story.
func (s *StoryService) CreateWithContext(ctx context.Context, story *Story) (*Story, *Response, error) {
	apiEndpoint := fmt.Sprintf("stories")
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, story)
	if err != nil {
		return nil, nil, err
	}

	storyresp := new(Story)
	resp, err := s.client.Do(req, storyresp)
	if err != nil {
		return nil, resp, err
	}

	return storyresp, resp, err
}

// Create wraps CreateWithContext using the background context.
func (s *StoryService) Create(story *Story) (*Story, *Response, error) {
	return s.CreateWithContext(context.Background(), story)
}

// UpdateWithContext Updates a story.
func (s *StoryService) UpdateWithContext(ctx context.Context, storyID int, story *Story) (*Story, *Response, error) {
	apiEndpoint := fmt.Sprintf("stories/%v", storyID)
	req, err := s.client.NewRequestWithContext(ctx, "PUT", apiEndpoint, story)
	if err != nil {
		return nil, nil, err
	}

	storyresp := new(Story)
	resp, err := s.client.Do(req, storyresp)
	if err != nil {
		return nil, resp, err
	}

	return storyresp, resp, err
}

// Update wraps UpdateWithContext using the background context.
func (s *StoryService) Update(storyID int, story *Story) (*Story, *Response, error) {
	return s.UpdateWithContext(context.Background(), storyID, story)
}
