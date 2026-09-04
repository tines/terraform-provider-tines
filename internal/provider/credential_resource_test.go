package provider

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// testAccTeamID returns the Tines Team ID to use in acceptance tests. It can be
// overridden with the TINES_TEST_TEAM_ID environment variable so the tests can
// run against any tenant without hardcoding a team, and defaults to the value
// used by the existing tines_resource acceptance tests.
func testAccTeamID() int64 {
	if v := os.Getenv("TINES_TEST_TEAM_ID"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			return id
		}
	}
	return 30906
}

func TestAccTinesCredential_Text(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create the Tines Credential.
				Config: providerConfig + testAccCreateTinesCredentialText(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectUnknownValue(
							"tines_credential.test_example_text",
							tfjsonpath.New("id"),
						),
						plancheck.ExpectKnownValue(
							"tines_credential.test_example_text",
							tfjsonpath.New("team_id"),
							knownvalue.Int64Exact(testAccTeamID()),
						),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"tines_credential.test_example_text",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"tines_credential.test_example_text",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Terraform Test Text Credential"),
					),
					statecheck.ExpectKnownValue(
						"tines_credential.test_example_text",
						tfjsonpath.New("mode"),
						knownvalue.StringExact("TEXT"),
					),
				},
			},
			{
				// Rotate the Tines Credential value.
				Config: providerConfig + testAccUpdateTinesCredentialText(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("tines_credential.test_example_text", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("tines_credential.test_example_text", tfjsonpath.New("id"), knownvalue.NotNull()),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"tines_credential.test_example_text",
						tfjsonpath.New("value"),
						knownvalue.StringExact("rotated_secret_value"),
					),
				},
			},
			{
				// Import the existing Tines Credential. The secret value is never
				// returned by the API, so it cannot be verified on import.
				ResourceName:            "tines_credential.test_example_text",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
		},
	})
}

func testAccCreateTinesCredentialText() string {
	return fmt.Sprintf(`
resource "tines_credential" "test_example_text" {
	team_id = %d
	name = "Terraform Test Text Credential"
	value = "initial_secret_value"
}
	`, testAccTeamID())
}

func testAccUpdateTinesCredentialText() string {
	return fmt.Sprintf(`
resource "tines_credential" "test_example_text" {
	team_id = %d
	name = "Terraform Test Text Credential"
	value = "rotated_secret_value"
}
	`, testAccTeamID())
}

func TestAccTinesCredential_WriteOnly(t *testing.T) {
	resource.Test(t, resource.TestCase{
		// Write-only attributes are only supported in Terraform 1.11 and later,
		// so skip this test on the older versions covered by CI.
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create the Tines Credential using a write-only secret value.
				Config: providerConfig + testAccCreateTinesCredentialWriteOnly(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectUnknownValue(
							"tines_credential.test_example_wo",
							tfjsonpath.New("id"),
						),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"tines_credential.test_example_wo",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"tines_credential.test_example_wo",
						tfjsonpath.New("mode"),
						knownvalue.StringExact("TEXT"),
					),
					// The write-only value must never be persisted to state.
					statecheck.ExpectKnownValue(
						"tines_credential.test_example_wo",
						tfjsonpath.New("value_wo"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"tines_credential.test_example_wo",
						tfjsonpath.New("value_wo_version"),
						knownvalue.Int64Exact(1),
					),
				},
			},
			{
				// Rotate the secret by bumping value_wo_version, which triggers
				// an in-place update.
				Config: providerConfig + testAccRotateTinesCredentialWriteOnly(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("tines_credential.test_example_wo", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"tines_credential.test_example_wo",
						tfjsonpath.New("value_wo_version"),
						knownvalue.Int64Exact(2),
					),
					statecheck.ExpectKnownValue(
						"tines_credential.test_example_wo",
						tfjsonpath.New("value_wo"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

func testAccCreateTinesCredentialWriteOnly() string {
	return fmt.Sprintf(`
resource "tines_credential" "test_example_wo" {
	team_id = %d
	name = "Terraform Test Write-Only Credential"
	value_wo = "initial_secret_value"
	value_wo_version = 1
}
	`, testAccTeamID())
}

func testAccRotateTinesCredentialWriteOnly() string {
	return fmt.Sprintf(`
resource "tines_credential" "test_example_wo" {
	team_id = %d
	name = "Terraform Test Write-Only Credential"
	value_wo = "rotated_secret_value"
	value_wo_version = 2
}
	`, testAccTeamID())
}
