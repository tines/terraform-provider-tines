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

// FolderService handles fields for the Tines instance / API.
type FolderService struct {
	client *Client
}

// Folder structure
type Folder struct {
	ID          int    `json:"id" structs:"id,omitempty"`
	Name        string `json:"name" structs:"name,omitempty"`
	ContentType string `json:"content_type" structs:"content_type,omitempty"`
	TeamID      int    `json:"team_id" structs:"team_id,omitempty"`
	Size        int    `json:"size" structs:"size,omitempty"`
	Unknowns    tcontainer.MarshalMap
}

// MarshalJSON is a custom JSON marshal function for the Folder* structs.
// It handles Tines options and maps those from / to "Unknowns" key.
func (i *Folder) MarshalJSON() ([]byte, error) {
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

// UnmarshalJSON is a custom JSON marshal function for the Folder structs.
// It handles Tines custom options and maps those from / to "Unknowns" key.
func (i *Folder) UnmarshalJSON(data []byte) error {

	// Do the normal unmarshalling first
	// Details for this way: http://choly.ca/post/go-json-marshalling/
	type Alias Folder
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
	i = (*Folder)(aux.Alias)
	// all the tags found in the struct were removed. Whatever is left are unknowns to struct
	i.Unknowns = totalMap
	return nil

}

// GetWithContext returns an folder for the given folder key.
func (s *FolderService) GetWithContext(ctx context.Context, folderID int) (*Folder, *Response, error) {
	apiEndpoint := fmt.Sprintf("folders/%v", folderID)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	folder := new(Folder)
	resp, err := s.client.Do(req, folder)
	if err != nil {
		return nil, resp, err
	}

	return folder, resp, nil
}

// Get wraps GetWithContext using the background context.
func (s *FolderService) Get(folderID int) (*Folder, *Response, error) {
	return s.GetWithContext(context.Background(), folderID)
}

// DeleteWithContext deletes an folder.
func (s *FolderService) DeleteWithContext(ctx context.Context, folderID int) (*Response, error) {
	apiEndpoint := fmt.Sprintf("folders/%v", folderID)
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
func (s *FolderService) Delete(folderID int) (*Response, error) {
	return s.DeleteWithContext(context.Background(), folderID)
}

// CreateWithContext creates an folder.
func (s *FolderService) CreateWithContext(ctx context.Context, folder *Folder) (*Folder, *Response, error) {
	apiEndpoint := fmt.Sprintf("folders")
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, folder)
	if err != nil {
		return nil, nil, err
	}

	folderresp := new(Folder)
	resp, err := s.client.Do(req, folderresp)
	if err != nil {
		return nil, resp, err
	}

	return folderresp, resp, err
}

// Create wraps CreateWithContext using the background context.
func (s *FolderService) Create(folder *Folder) (*Folder, *Response, error) {
	return s.CreateWithContext(context.Background(), folder)
}

// UpdateWithContext updates an folder for the given folder key.
func (s *FolderService) UpdateWithContext(ctx context.Context, folderID int, folder *Folder) (*Folder, *Response, error) {
	apiEndpoint := fmt.Sprintf("folders/%v", folderID)
	req, err := s.client.NewRequestWithContext(ctx, "PUT", apiEndpoint, folder)
	if err != nil {
		return nil, nil, err
	}

	folderresp := new(Folder)
	resp, err := s.client.Do(req, folderresp)
	if err != nil {
		return nil, resp, err
	}

	return folderresp, resp, err
}

// Update wraps UpdateWithContext using the background context.
func (s *FolderService) Update(folderID int, folder *Folder) (*Folder, *Response, error) {
	return s.UpdateWithContext(context.Background(), folderID, folder)
}
