package tines

import (
	"context"

	"fmt"
)

// NoteService handles fields for the Tines instance / API.
type NoteService struct {
	client *Client
}

// Note structure
type Note struct {
	ID       int                    `json:"id" structs:"id,omitempty"`
	StoryID  int                    `json:"story_id" structs:"story_id"`
	Content  string                 `json:"content" structs:"content,omitempty"`
	Position map[string]interface{} `json:"position" structs:"position,omitempty"`
}

// GetWithContext returns an note for the given note key.
func (s *NoteService) GetWithContext(ctx context.Context, noteID int) (*Note, *Response, error) {
	apiEndpoint := fmt.Sprintf("diagram_notes/%v", noteID)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	note := new(Note)
	resp, err := s.client.Do(req, note)
	if err != nil {
		return nil, resp, err
	}

	return note, resp, nil
}

// Get wraps GetWithContext using the background context.
func (s *NoteService) Get(noteID int) (*Note, *Response, error) {
	return s.GetWithContext(context.Background(), noteID)
}

// DeleteWithContext deletes an note.
func (s *NoteService) DeleteWithContext(ctx context.Context, noteID int) (*Response, error) {
	apiEndpoint := fmt.Sprintf("diagram_notes/%v", noteID)
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
func (s *NoteService) Delete(noteID int) (*Response, error) {
	return s.DeleteWithContext(context.Background(), noteID)
}

// CreateWithContext creates an note.
func (s *NoteService) CreateWithContext(ctx context.Context, note *Note) (*Note, *Response, error) {
	apiEndpoint := fmt.Sprintf("diagram_notes")
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, note)
	if err != nil {
		return nil, nil, err
	}

	noteresp := new(Note)
	resp, err := s.client.Do(req, noteresp)
	if err != nil {
		return nil, resp, err
	}

	return noteresp, resp, err
}

// Create wraps CreateWithContext using the background context.
func (s *NoteService) Create(note *Note) (*Note, *Response, error) {
	return s.CreateWithContext(context.Background(), note)
}

// UpdateWithContext updates an note for the given note key.
func (s *NoteService) UpdateWithContext(ctx context.Context, noteID int, note *Note) (*Note, *Response, error) {
	apiEndpoint := fmt.Sprintf("diagram_notes/%v", noteID)
	req, err := s.client.NewRequestWithContext(ctx, "PUT", apiEndpoint, note)
	if err != nil {
		return nil, nil, err
	}

	noteresp := new(Note)
	resp, err := s.client.Do(req, noteresp)
	if err != nil {
		return nil, resp, err
	}

	return noteresp, resp, err
}

// Update wraps UpdateWithContext using the background context.
func (s *NoteService) Update(noteID int, note *Note) (*Note, *Response, error) {
	return s.UpdateWithContext(context.Background(), noteID, note)
}
