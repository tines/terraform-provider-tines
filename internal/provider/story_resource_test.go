package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccTinesStory_fromExportNoFolder(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccCreateStoryResourceNoFolder(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectUnknownValue(
							"tines_story.test_create_from_export_no_folder",
							tfjsonpath.New("id"),
						),
						plancheck.ExpectUnknownValue(
							"tines_story.test_create_from_export_no_folder",
							tfjsonpath.New("folder_id"),
						),
						plancheck.ExpectKnownValue(
							"tines_story.test_create_from_export_no_folder",
							tfjsonpath.New("team_id"),
							knownvalue.Int64Exact(30906),
						),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"tines_story.test_create_from_export_no_folder",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"tines_story.test_create_from_export_no_folder",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Test Story"),
					),
				},
			},
		},
	})
}

func testAccCreateStoryResourceNoFolder() string {
	return `
resource "tines_story" "test_create_from_export_no_folder" {
	data = file("${path.module}/testdata/test-story.json")
	team_id = 30906
}
	`
}
