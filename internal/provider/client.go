package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Client - Tines API client.
type Client struct {
	apiToken  string
	tenantUrl string
	client    *http.Client
}

// StoryResponse - Response from Tines API when importing a story.
type StoryResponse struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	UserID        int64  `json:"user_id"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	Description   string `json:"description"`
	GUID          string `json:"guid"`
	KeepEventsFor int64  `json:"keep_events_for"`
	Disabled      bool   `json:"disabled"`
	Priority      bool   `json:"priority"`
	TeamID        int64  `json:"team_id"`
	FolderID      int64  `json:"folder_id"`
}

// NewClient - Creates a new Tines API client.
func NewClient() (*Client, error) {
	c := &Client{
		client: &http.Client{},
	}
	return c, nil
}

func (c *Client) buildStoryData(storyData string, teamId int64, folderId int64) (sirB []byte, e error) {
	sData, err := strconv.Unquote(storyData)
	if err != nil {
		return sirB, err
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(sData), &data)
	if err != nil {
		return sirB, err

	}

	name, _ := data["name"].(string)

	sir := StoryImportRequest{
		NewName:  name,
		TeamID:   teamId,
		FolderID: folderId,
		Mode:     "versionReplace",
		Data:     data,
	}

	sirB, err = json.Marshal(sir)

	if err != nil {
		return sirB, err
	}

	return sirB, nil
}

// ImportStory - Imports a story.
func (c *Client) ImportStory(ctx context.Context, storyData []byte) (sr StoryResponse, e error) {
	parsedTenantUrl, err := url.Parse(c.tenantUrl)
	if err != nil {
		return sr, err
	}

	fullPath := parsedTenantUrl.JoinPath("/api/v1/stories/import").String()
	req, err := http.NewRequest("POST", fullPath, bytes.NewBuffer(storyData))
	tflog.Debug(ctx, fmt.Sprintf("Sending %v request to %v", "POST", fullPath))
	if err != nil {
		return sr, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Set("x-user-token", c.apiToken)

	res, err := c.client.Do(req)
	if err != nil {
		return sr, err
	}
	defer res.Body.Close()
	if err != nil {
		return sr, err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return sr, err
	}

	if res.StatusCode != http.StatusOK {
		return sr, fmt.Errorf("request to import story failed with status: %d and message: %s", res.StatusCode, data)
	}

	jErr := json.Unmarshal(data, &sr)
	if jErr != nil {
		return sr, jErr
	}

	return sr, nil
}

func (c *Client) DeleteStory(ctx context.Context, storyID basetypes.Int64Value) (status int, e error) {
	parsedTenantUrl, err := url.Parse(c.tenantUrl)
	if err != nil {
		return status, err
	}

	fullPath := parsedTenantUrl.JoinPath("/api/v1/stories/" + storyID.String()).String()
	tflog.Debug(ctx, fmt.Sprintf("Sending %v request to %v", "DELETE", fullPath))
	req, err := http.NewRequest("DELETE", fullPath, nil)
	if err != nil {
		return -1, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Set("x-user-token", c.apiToken)

	res, err := c.client.Do(req)
	if err != nil {
		return -1, err
	}
	defer res.Body.Close()
	if err != nil {
		return -1, err
	}

	if res.StatusCode != 204 {
		return -1, fmt.Errorf("request to delete story failed with status: %d", res.StatusCode)
	}
	return res.StatusCode, nil
}
