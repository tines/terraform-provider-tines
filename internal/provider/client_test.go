package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestBuildData(t *testing.T) {
	t.Parallel()
	client, err := NewClient()
	assert.Nil(t, err)

	sirB, err := client.buildStoryData(strconv.Quote(string(([]byte(`{}`)))), int64(2), int64(1))
	assert.Nil(t, err)
	var data StoryImportRequest
	err = json.Unmarshal(sirB, &data)
	assert.Nil(t, err)

	assert.Equal(t, data.TeamID, int64(2))
	assert.Equal(t, data.FolderID, int64(1))
	assert.Equal(t, data.Mode, "versionReplace")
}

func TestBuildDataReturnsError(t *testing.T) {
	t.Parallel()
	client, err := NewClient()
	assert.Nil(t, err)

	_, err = client.buildStoryData(string([]byte(`{}`)), int64(2), int64(1))
	expectedError := fmt.Errorf("invalid syntax")
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
}

func TestImportStory(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("testdata/imported_story.json")
		assert.Nil(t, err)

		defer f.Close()
		_, cErr := io.Copy(w, f)
		assert.Nil(t, cErr)

		assert.Equal(t, "application/json", r.Header.Get("content-type"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "dummyToken", r.Header.Get("x-user-token"))
		assert.Equal(t, "tines-terraform-client", r.Header.Get("User-Agent"))

		var data StoryImportRequest
		body, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		err = json.Unmarshal(body, &data)
		assert.Nil(t, err)

		assert.Equal(t, data.TeamID, int64(2))
		assert.Equal(t, data.FolderID, int64(1))
		assert.Equal(t, data.Mode, "versionReplace")
	}))
	defer ts.Close()
	client, err := NewClient()
	client.tenantUrl = ts.URL
	client.apiToken = "dummyToken"
	assert.Nil(t, err)

	ctx := context.Background()

	want := StoryResponse{
		ID:            2,
		Name:          "Simple story",
		UserID:        1,
		CreatedAt:     "2021-03-235T20:06:14.825Z",
		UpdatedAt:     "2021-03-23T20:06:14.825Z",
		Description:   "In the simple story we will create a fictional situation where a detection system is configured to send alerts to our Tines tenant",
		GUID:          "b7c81e0cb416ae8f4c00874ca7b1cdf8",
		KeepEventsFor: 604800,
		Disabled:      false,
		Priority:      false,
		TeamID:        2,
		FolderID:      1,
	}

	sirB, err := client.buildStoryData(strconv.Quote(string(([]byte(`{}`)))), int64(2), int64(1))
	assert.Nil(t, err)

	got, err := client.ImportStory(ctx, sirB)
	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestFailedImportStory(t *testing.T) {
	t.Parallel()
	errorData := []byte(`{"error": "Team Not found"}`)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, wErr := w.Write(errorData)
		assert.Nil(t, wErr)

		assert.Equal(t, "application/json", r.Header.Get("content-type"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "dummyToken", r.Header.Get("x-user-token"))
		assert.Equal(t, "tines-terraform-client", r.Header.Get("User-Agent"))
	}))
	defer ts.Close()

	client, err := NewClient()
	assert.Nil(t, err)
	client.tenantUrl = ts.URL
	client.apiToken = "dummyToken"

	sir := StoryImportRequest{}
	sirB, err := json.Marshal(sir)
	assert.Nil(t, err)
	ctx := context.Background()
	_, newErr := client.ImportStory(ctx, sirB)
	expectedError := fmt.Errorf("request to import story failed with status: %d and message: %s", 404, errorData)
	if assert.Error(t, newErr) {
		assert.Equal(t, expectedError, newErr)
	}
}

func TestDeleteStory(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)

		assert.Equal(t, "application/json", r.Header.Get("content-type"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "dummyToken", r.Header.Get("x-user-token"))
		assert.Equal(t, "tines-terraform-client", r.Header.Get("User-Agent"))
	}))
	defer ts.Close()
	client, err := NewClient()
	assert.Nil(t, err)
	client.tenantUrl = ts.URL
	client.apiToken = "dummyToken"

	ctx := context.Background()

	status, e := client.DeleteStory(ctx, types.Int64Value(2))
	assert.Equal(t, 204, status)
	assert.Equal(t, nil, e)
}

func TestFailedDeleteStory(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()
	client, err := NewClient()
	assert.Nil(t, err)
	client.tenantUrl = ts.URL
	client.apiToken = "dummyToken"

	ctx := context.Background()

	_, e := client.DeleteStory(ctx, types.Int64Value(2))
	expectedError := fmt.Errorf("request to delete story failed with status: %d", 404)
	if assert.Error(t, e) {
		assert.Equal(t, expectedError, e)
	}
}
