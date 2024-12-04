package tines_cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportStory(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := []byte(`{
				"name": "Test Story",
				"user_id": 1,
				"description": null,
				"keep_events_for": 86400,
				"disabled": false,
				"priority": false,
				"send_to_story_enabled": false,
				"send_to_story_access_source": "OFF",
				"send_to_story_access": null,
				"shared_team_slugs": [],
				"entry_agent_id": null,
				"exit_agents": [],
				"send_to_story_skill_use_requires_confirmation": true,
				"team_id": 1,
				"tags": [],
				"guid": "3ef721e341e953727b057d4bd7bd65eb",
				"slug": "test_story",
				"created_at": "2024-11-04T23:04:43Z",
				"updated_at": "2024-11-04T23:06:51Z",
				"edited_at": "2024-11-05T01:45:48Z",
				"mode": "LIVE",
				"story_container": {
					"id": 1,
					"team_id": 1,
					"folder_id": null,
					"published": true,
					"locked": false,
					"change_control_enabled": false,
					"owners": []
				},
				"id": 1,
				"folder_id": null,
				"published": true,
				"locked": false
			}`)

		w.Write(res) //nolint:errcheck

	}))
	defer ts.Close()

	tenant := ts.URL
	apiKey := "foo"
	version := "test"

	c, err := NewClient(tenant, apiKey, version)

	assert.Nil(err, "should instantiate the Tines API client without errors")

	var data map[string]interface{}
	f, err := os.ReadFile("testdata/test-import-story.json")
	assert.Nil(err, "It should open and read the JSON test file")

	err = json.Unmarshal(f, &data)
	assert.Nil(err, "It should unmarshal the JSON without errors")

	name, ok := data["name"].(string)
	assert.True(ok, "It should assert that the 'name' field in the data interface is a string type")

	storyImportRequest := StoryImportRequest{
		NewName: name,
		Data:    data,
		TeamID:  1,
		Mode:    "versionReplace",
	}

	story, err := c.ImportStory(&storyImportRequest)

	assert.Nil(err, "It should import the story without API errors")
	assert.Equal(int64(0), story.FolderID, "The folder ID should be coerced to 0 if the story is not in a folder")
	assert.Equal(int64(1), story.TeamID, "The story should be imported to the correct team")
	assert.False(story.Disabled, "The story should not be disabled")
	assert.False(story.ChangeControlEnabled, "The story should not have change control enabled")
	assert.ElementsMatch([]int64{}, story.ExitAgents, "The story should have an empty list of exit agents")
}
