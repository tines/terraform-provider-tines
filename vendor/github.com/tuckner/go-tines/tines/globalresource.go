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

// GlobalResourceService handles fields for the Tines instance / API.
type GlobalResourceService struct {
	client *Client
}

// GlobalResource structure
type GlobalResource struct {
	ID        int         `json:"id" structs:"id,omitempty"`
	UserID    int         `json:"user_id" structs:"user_id,omitempty"`
	Name      string      `json:"name" structs:"name,omitempty"`
	ValueType string      `json:"value_type" structs:"value_type,omitempty"`
	Value     string      `json:"value" structs:"value,omitempty"`
	CreatedAt time.Time   `json:"created_at" structs:"created_at,omitempty"`
	UpdatedAt time.Time   `json:"updated_at" structs:"updated_at,omitempty"`
	Slug      string      `json:"slug" structs:"slug,omitempty"`
	TeamID    int         `json:"team_id" structs:"team_id,omitempty"`
	FolderID  int `json:"folder_id" structs:"folder_id,omitempty"`
	Unknowns  tcontainer.MarshalMap
}

// MarshalJSON is a custom JSON marshal function for the GlobalResource* structs.
// It handles Tines options and maps those from / to "Unknowns" key.
func (i *GlobalResource) MarshalJSON() ([]byte, error) {
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

// UnmarshalJSON is a custom JSON marshal function for the GlobalResource structs.
// It handles Tines custom options and maps those from / to "Unknowns" key.
func (i *GlobalResource) UnmarshalJSON(data []byte) error {

	// Do the normal unmarshalling first
	// Details for this way: http://choly.ca/post/go-json-marshalling/
	type Alias GlobalResource
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
	i = (*GlobalResource)(aux.Alias)
	// all the tags found in the struct were removed. Whatever is left are unknowns to struct
	i.Unknowns = totalMap
	return nil

}

// GetWithContext returns an agent for the given global resource id.
func (s *GlobalResourceService) GetWithContext(ctx context.Context, globalResourceID int) (*GlobalResource, *Response, error) {
	apiEndpoint := fmt.Sprintf("global_resources/%v", globalResourceID)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	globalresource := new(GlobalResource)
	resp, err := s.client.Do(req, globalresource)
	if err != nil {
		return nil, resp, err
	}

	return globalresource, resp, nil
}

// Get wraps GetWithContext using the background context.
func (s *GlobalResourceService) Get(globalResourceID int) (*GlobalResource, *Response, error) {
	return s.GetWithContext(context.Background(), globalResourceID)
}

// DeleteWithContext deletes a global resource.
func (s *GlobalResourceService) DeleteWithContext(ctx context.Context, globalResourceID int) (*Response, error) {
	apiEndpoint := fmt.Sprintf("global_resources/%v", globalResourceID)
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
func (s *GlobalResourceService) Delete(globalResourceID int) (*Response, error) {
	return s.DeleteWithContext(context.Background(), globalResourceID)
}

// CreateWithContext creates a global resource.
func (s *GlobalResourceService) CreateWithContext(ctx context.Context, globalresource *GlobalResource) (*GlobalResource, *Response, error) {
	apiEndpoint := fmt.Sprintf("global_resources")
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, globalresource)
	if err != nil {
		return nil, nil, err
	}

	globalresourceresp := new(GlobalResource)
	resp, err := s.client.Do(req, globalresourceresp)
	if err != nil {
		return nil, resp, err
	}

	return globalresourceresp, resp, err
}

// Create wraps CreateWithContext using the background context.
func (s *GlobalResourceService) Create(globalresource *GlobalResource) (*GlobalResource, *Response, error) {
	return s.CreateWithContext(context.Background(), globalresource)
}

// UpdateWithContext Updates a global resource.
func (s *GlobalResourceService) UpdateWithContext(ctx context.Context, globalResourceID int, globalresource *GlobalResource) (*GlobalResource, *Response, error) {
	apiEndpoint := fmt.Sprintf("global_resources/%v", globalResourceID)
	req, err := s.client.NewRequestWithContext(ctx, "PUT", apiEndpoint, globalresource)
	if err != nil {
		return nil, nil, err
	}

	globalresourceresp := new(GlobalResource)
	resp, err := s.client.Do(req, globalresourceresp)
	if err != nil {
		return nil, resp, err
	}

	return globalresourceresp, resp, err
}

// Update wraps UpdateWithContext using the background context.
func (s *GlobalResourceService) Update(globalResouceID int, globalresource *GlobalResource) (*GlobalResource, *Response, error) {
	return s.UpdateWithContext(context.Background(), globalResouceID, globalresource)
}
