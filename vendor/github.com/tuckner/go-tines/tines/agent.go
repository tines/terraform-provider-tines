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

// AgentService handles fields for the Tines instance / API.
type AgentService struct {
	client *Client
}

// Agent structure
type Agent struct {
	ID                 int         `json:"id" structs:"id,omitempty"`
	UserID             int         `json:"user_id" structs:"user_id,omitempty"`
	Name               string      `json:"name" structs:"name"`
	Schedule           interface{} `json:"schedule" structs:"schedule,omitempty"`
	EventsCount        int         `json:"events_count" structs:"events_count,omitempty"`
	LastCheckAt        interface{} `json:"last_check_at" structs:"last_check_at,omitempty"`
	LastReceiveAt      time.Time   `json:"last_receive_at" structs:"last_receive_at,omitempty"`
	LastCheckedEventID int         `json:"last_checked_event_id" structs:"last_checked_event_id,omitempty"`
	CreatedAt          time.Time   `json:"created_at" structs:"created_at,omitempty"`
	UpdatedAt          time.Time   `json:"updated_at" structs:"updated_at,omitempty"`
	LastWebRequestAt   interface{} `json:"last_web_request_at" structs:"last_web_request_at,omitempty"`
	KeepEventsFor      int         `json:"keep_events_for" structs:"keep_events_for,omitempty"`
	LastEventAt        time.Time   `json:"last_event_at" structs:"last_event_at,omitempty"`
	LastErrorLogAt     interface{} `json:"last_error_log_at" structs:"last_error_log_at,omitempty"`
	Disabled           bool        `json:"disabled" structs:"disabled,omitempty"`
	GUID               string      `json:"guid" structs:"guid,omitempty"`
	StoryID            int         `json:"story_id" structs:"story_id"`
	SourceIds          []int       `json:"source_ids" structs:"source_ids,omitempty"`
	ReceiverIds        []int       `json:"receiver_ids" structs:"receiver_ids,omitempty"`
	Type               string      `json:"type" structs:"type,omitempty"`
	Unknowns           tcontainer.MarshalMap
}

// MarshalJSON is a custom JSON marshal function for the Agent* structs.
// It handles Tines options and maps those from / to "Unknowns" key.
func (i *Agent) MarshalJSON() ([]byte, error) {
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

// UnmarshalJSON is a custom JSON marshal function for the Agent structs.
// It handles Tines custom options and maps those from / to "Unknowns" key.
func (i *Agent) UnmarshalJSON(data []byte) error {

	// Do the normal unmarshalling first
	// Details for this way: http://choly.ca/post/go-json-marshalling/
	type Alias Agent
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
	i = (*Agent)(aux.Alias)
	// all the tags found in the struct were removed. Whatever is left are unknowns to struct
	i.Unknowns = totalMap
	return nil

}

// GetWithContext returns an agent for the given agent key.
func (s *AgentService) GetWithContext(ctx context.Context, agentID int) (*Agent, *Response, error) {
	apiEndpoint := fmt.Sprintf("agents/%v", agentID)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	agent := new(Agent)
	resp, err := s.client.Do(req, agent)
	if err != nil {
		return nil, resp, err
	}

	return agent, resp, nil
}

// Get wraps GetWithContext using the background context.
func (s *AgentService) Get(agentID int) (*Agent, *Response, error) {
	return s.GetWithContext(context.Background(), agentID)
}

// DeleteWithContext deletes an agent.
func (s *AgentService) DeleteWithContext(ctx context.Context, agentID int) (*Response, error) {
	apiEndpoint := fmt.Sprintf("agents/%v", agentID)
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

// Delete wraps GetWithContext using the background context.
func (s *AgentService) Delete(agentID int) (*Response, error) {
	return s.DeleteWithContext(context.Background(), agentID)
}

// CreateWithContext creates an agent.
func (s *AgentService) CreateWithContext(ctx context.Context, agent *Agent) (*Agent, *Response, error) {
	apiEndpoint := fmt.Sprintf("agents")
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, agent)
	if err != nil {
		return nil, nil, err
	}

	agentresp := new(Agent)
	resp, err := s.client.Do(req, agentresp)
	if err != nil {
		return nil, resp, err
	}

	return agentresp, resp, err
}

// Create wraps CreateWithContext using the background context.
func (s *AgentService) Create(agent *Agent) (*Agent, *Response, error) {
	return s.CreateWithContext(context.Background(), agent)
}

// UpdateWithContext updates an agent for the given agent key.
func (s *AgentService) UpdateWithContext(ctx context.Context, agentID int, agent *Agent) (*Agent, *Response, error) {
	apiEndpoint := fmt.Sprintf("agents/%v", agentID)
	req, err := s.client.NewRequestWithContext(ctx, "PUT", apiEndpoint, agent)
	if err != nil {
		return nil, nil, err
	}

	agentresp := new(Agent)
	resp, err := s.client.Do(req, agentresp)
	if err != nil {
		return nil, resp, err
	}

	return agentresp, resp, err
}

// Update wraps UpdateWithContext using the background context.
func (s *AgentService) Update(agentID int, agent *Agent) (*Agent, *Response, error) {
	return s.UpdateWithContext(context.Background(), agentID, agent)
}
