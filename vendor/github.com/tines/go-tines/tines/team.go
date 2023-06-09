package tines

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"

	"fmt"

	"github.com/fatih/structs"
	"github.com/trivago/tgo/tcontainer"
)

// TeamService handles fields for the Tines instance / API.
type TeamService struct {
	client *Client
}

// Team structure
type Team struct {
	ID       int    `json:"id" structs:"id,omitempty"`
	Name     string `json:"name" structs:"name,omitempty"`
	Unknowns tcontainer.MarshalMap
}

// MarshalJSON is a custom JSON marshal function for the Team* structs.
// It handles Tines options and maps those from / to "Unknowns" key.
func (i *Team) MarshalJSON() ([]byte, error) {
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

// UnmarshalJSON is a custom JSON marshal function for the Team structs.
// It handles Tines custom options and maps those from / to "Unknowns" key.
func (i *Team) UnmarshalJSON(data []byte) error {

	// Do the normal unmarshalling first
	// Details for this way: http://choly.ca/post/go-json-marshalling/
	type Alias Team
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
	i = (*Team)(aux.Alias)
	// all the tags found in the struct were removed. Whatever is left are unknowns to struct
	i.Unknowns = totalMap
	return nil

}

// GetWithContext returns an team for the given team key.
func (s *TeamService) GetWithContext(ctx context.Context, teamID int) (*Team, *Response, error) {
	apiEndpoint := fmt.Sprintf("teams/%v", teamID)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	team := new(Team)
	resp, err := s.client.Do(req, team)
	if err != nil {
		return nil, resp, err
	}

	return team, resp, nil
}

// Get wraps GetWithContext using the background context.
func (s *TeamService) Get(teamID int) (*Team, *Response, error) {
	return s.GetWithContext(context.Background(), teamID)
}

// DeleteWithContext deletes an team.
func (s *TeamService) DeleteWithContext(ctx context.Context, teamID int) (*Response, error) {
	apiEndpoint := fmt.Sprintf("teams/%v", teamID)
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
func (s *TeamService) Delete(teamID int) (*Response, error) {
	return s.DeleteWithContext(context.Background(), teamID)
}

// CreateWithContext creates an team.
func (s *TeamService) CreateWithContext(ctx context.Context, team *Team) (*Team, *Response, error) {
	apiEndpoint := fmt.Sprintf("teams")
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, team)
	if err != nil {
		return nil, nil, err
	}

	teamresp := new(Team)
	resp, err := s.client.Do(req, teamresp)
	if err != nil {
		return nil, resp, err
	}

	return teamresp, resp, err
}

// Create wraps CreateWithContext using the background context.
func (s *TeamService) Create(team *Team) (*Team, *Response, error) {
	return s.CreateWithContext(context.Background(), team)
}

// UpdateWithContext updates an team for the given team key.
func (s *TeamService) UpdateWithContext(ctx context.Context, teamID int, team *Team) (*Team, *Response, error) {
	apiEndpoint := fmt.Sprintf("teams/%v", teamID)
	req, err := s.client.NewRequestWithContext(ctx, "PUT", apiEndpoint, team)
	if err != nil {
		return nil, nil, err
	}

	teamresp := new(Team)
	resp, err := s.client.Do(req, teamresp)
	if err != nil {
		return nil, resp, err
	}

	return teamresp, resp, err
}

// Update wraps UpdateWithContext using the background context.
func (s *TeamService) Update(teamID int, team *Team) (*Team, *Response, error) {
	return s.UpdateWithContext(context.Background(), teamID, team)
}
