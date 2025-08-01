package provider

import (
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccTinesResource_String(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create the Tines Resource.
				Config: providerConfig + testAccCreateTinesResourceStringVal(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectUnknownValue(
							"tines_resource.test_example_string",
							tfjsonpath.New("id"),
						),
						plancheck.ExpectKnownValue(
							"tines_resource.test_example_string",
							tfjsonpath.New("team_id"),
							knownvalue.Int64Exact(30906),
						),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"tines_resource.test_example_string",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"tines_resource.test_example_string",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Terraform Test String Resource"),
					),
				},
			},
			{
				// Update the Tines Resource value.
				Config: providerConfig + testAccUpdateTinesResourceStringVal(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("tines_resource.test_example_string", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("tines_resource.test_example_string", tfjsonpath.New("id"), knownvalue.NotNull()),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"tines_resource.test_example_string",
						tfjsonpath.New("value"),
						knownvalue.StringExact("new string"),
					),
				},
			},
			{
				// Import the existing Tines Resource.
				ResourceName:      "tines_resource.test_example_string",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTinesResource_Array(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create the Tines Resource with array value.
				Config: providerConfig + testAccCreateTinesResourceArrayVal(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectUnknownValue(
							"tines_resource.test_example_array",
							tfjsonpath.New("id"),
						),
						plancheck.ExpectKnownValue(
							"tines_resource.test_example_array",
							tfjsonpath.New("team_id"),
							knownvalue.Int64Exact(30906),
						),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"tines_resource.test_example_array",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Terraform Test Array Resource"),
					),
					statecheck.ExpectKnownValue(
						"tines_resource.test_example_array",
						tfjsonpath.New("value"),
						knownvalue.TupleExact([]knownvalue.Check{
							knownvalue.NumberExact(big.NewFloat(1)),
							knownvalue.NumberExact(big.NewFloat(2)),
							knownvalue.NumberExact(big.NewFloat(3)),
						}),
					),
				},
			},
			{
				// Update the Tines Resource array value.
				Config: providerConfig + testAccUpdateTinesResourceArrayVal(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("tines_resource.test_example_array", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("tines_resource.test_example_array", tfjsonpath.New("id"), knownvalue.NotNull()),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"tines_resource.test_example_array",
						tfjsonpath.New("value"),
						knownvalue.TupleExact([]knownvalue.Check{
							knownvalue.NumberExact(big.NewFloat(1)),
							knownvalue.NumberExact(big.NewFloat(2)),
							knownvalue.NumberExact(big.NewFloat(3)),
							knownvalue.NumberExact(big.NewFloat(4)),
						}),
					),
				},
			},
			{
				// Import the existing Tines Resource.
				ResourceName:      "tines_resource.test_example_array",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCreateTinesResourceStringVal() string {
	return `
resource "tines_resource" "test_example_string" {
	team_id = 30906
	name = "Terraform Test String Resource"
	value = "example string"
}
	`
}

func testAccUpdateTinesResourceStringVal() string {
	return `
resource "tines_resource" "test_example_string" {
	team_id = 30906
	name = "Terraform Test String Resource"
	value = "new string"
}
	`
}

func testAccCreateTinesResourceArrayVal() string {
	return `
resource "tines_resource" "test_example_array" {
	team_id = 30906
	name = "Terraform Test Array Resource"
	value = [1, 2, 3]
}
	`
}

func testAccUpdateTinesResourceArrayVal() string {
	return `
resource "tines_resource" "test_example_array" {
	team_id = 30906
	name = "Terraform Test Array Resource"
	value = [1, 2, 3, 4]
}
	`
}
