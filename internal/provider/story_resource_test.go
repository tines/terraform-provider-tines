package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/assert"
	"github.com/tines/terraform-provider-tines/internal/tines_cli"
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"tines": providerserver.NewProtocol6WithError(New("test")()),
	}
)

func TestStoryResourceCreate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			w.WriteHeader(204)
			return
		} else {
			f, err := os.Open("testdata/imported_story.json")
			assert.Nil(t, err)

			defer f.Close()
			_, cErr := io.Copy(w, f)
			assert.Nil(t, cErr)

			var data tines_cli.StoryImportRequest
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &data)
			assert.Nil(t, err)

			assert.Equal(t, data.NewName, "Simple story")
			assert.Equal(t, data.TeamID, int64(2))
			assert.Equal(t, data.FolderID, int64(1))
			assert.Equal(t, data.Mode, "versionReplace")
		}
	}))
	defer ts.Close()
	resourceConfig := `provider "tines" {}
		resource "tines_story" "test" {
			data 			= file("./testdata/test_story.json")
			tenant_url      = "%s"
			tines_api_token = "dummyToken"
			team_id 				= 2
			folder_id 			= 1
		}
	`
	resourceConfig = fmt.Sprintf(resourceConfig, ts.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tines_story.test", "id", "2"),
				),
			},
		},
	})
}

func TestStoryUpdateInPlace(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			w.WriteHeader(204)
			return
		} else {
			f, err := os.Open("testdata/imported_story.json")
			assert.Nil(t, err)

			defer f.Close()
			_, cErr := io.Copy(w, f)
			assert.Nil(t, cErr)

			var data tines_cli.StoryImportRequest
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &data)
			assert.Nil(t, err)

			assert.Equal(t, data.TeamID, int64(2))
			assert.Equal(t, data.FolderID, int64(1))
			assert.Equal(t, data.Mode, "versionReplace")
		}
	}))
	defer ts.Close()
	resourceConfig := `
		provider "tines" {}
		resource "tines_story" "test" {
			data 						= file("%s")
			tenant_url      = "%s"
			tines_api_token = "dummyToken"
			team_id 				= 2
			folder_id 			= 1
		}
	`

	resourceConfigOne := fmt.Sprintf(resourceConfig, "./testdata/test_story.json", ts.URL)
	resourceConfigTwo := fmt.Sprintf(resourceConfig, "./testdata/test_story_update_in_place.json", ts.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigOne,
			},
			{
				Config: resourceConfigTwo,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("tines_story.test", plancheck.ResourceActionUpdate),
					},
				},
			},
		},
	})
}
