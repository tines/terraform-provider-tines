package provider

import (
	"regexp"
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
				Config: providerConfig + testAccCreateImportStoryResourceNoFolder(),
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
			{
				RefreshState: true,
			},
			{
				Config: providerConfig + testAccUpdateImportStoryResourceNoFolder(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("tines_story.test_create_from_export_no_folder", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"tines_story.test_create_from_export_no_folder",
						tfjsonpath.New("team_id"),
						knownvalue.Int64Exact(32987),
					),
				},
			},
		},
	})
}

func TestAccTinesStory_fromExportWithFolder(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccCreateImportStoryResourceWithFolder(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectUnknownValue(
							"tines_story.test_create_from_export_with_folder",
							tfjsonpath.New("id"),
						),
						plancheck.ExpectKnownValue(
							"tines_story.test_create_from_export_with_folder",
							tfjsonpath.New("team_id"),
							knownvalue.Int64Exact(30906),
						),
						plancheck.ExpectKnownValue(
							"tines_story.test_create_from_export_with_folder",
							tfjsonpath.New("folder_id"),
							knownvalue.Int64Exact(7993),
						),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"tines_story.test_create_from_export_with_folder",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"tines_story.test_create_from_export_with_folder",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Test Story"),
					),
					statecheck.ExpectKnownValue(
						"tines_story.test_create_from_export_with_folder",
						tfjsonpath.New("folder_id"),
						knownvalue.Int64Exact(7993),
					),
				},
			},
		},
	})
}

func TestAccTinesStory_fromConfigNoFolder(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + testAccCreateConfigStoryResourceBadConfig(),
				ExpectError: regexp.MustCompile("Attribute \"data\" cannot be specified when \"name\" is specified"),
			},
			{
				Config: providerConfig + testAccCreateConfigStoryResourceOneStep(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
			},
			{
				Config: providerConfig + testAccCreateConfigStoryResourceMultiStep(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"tines_story.test_create_multi_step",
						tfjsonpath.New("change_control_enabled"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

func testAccCreateImportStoryResourceNoFolder() string {
	return `
resource "tines_story" "test_create_from_export_no_folder" {
	data = file("${path.module}/testdata/test-story.json")
	team_id = 30906
}
	`
}

func testAccUpdateImportStoryResourceNoFolder() string {
	return `
resource "tines_story" "test_create_from_export_no_folder" {
	data = file("${path.module}/testdata/test-story.json")
	team_id = 32987
}
	`
}

func testAccCreateImportStoryResourceWithFolder() string {
	return `
resource "tines_story" "test_create_from_export_with_folder" {
	data = file("${path.module}/testdata/test-story.json")
	team_id = 30906
	folder_id = 7993
}
	`
}

func testAccCreateConfigStoryResourceBadConfig() string {
	return `
resource "tines_story" "test_create_invalid" {
	data = file("${path.module}/testdata/test-story.json")
	team_id = 30906
	name = "Example Bad Config"
}
	`
}

func testAccCreateConfigStoryResourceOneStep() string {
	return `
resource "tines_story" "test_create_one_step" {
	team_id = 30906
	name = "Example One Step"
}
	`
}

func testAccCreateConfigStoryResourceMultiStep() string {
	return `
resource "tines_story" "test_create_multi_step" {
	team_id = 30906
	name = "Example Multi Step"
	change_control_enabled = true
}
	`
}
