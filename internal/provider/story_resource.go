package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/tines/terraform-provider-tines/internal/tines_cli"
)

// storyResource is the resource implementation.
type storyResource struct {
	client *tines_cli.Client
}

// Prior schema data to enable upgrade of existing tfstate files to
// the current schema version.
type storyResourceModelV0 struct {
	Data          types.String `tfsdk:"data"`
	ID            types.Int64  `tfsdk:"id"`
	TinesApiToken types.String `tfsdk:"tines_api_token"`
	TenantUrl     types.String `tfsdk:"tenant_url"`
	TeamID        types.Int64  `tfsdk:"team_id"`
	FolderID      types.Int64  `tfsdk:"folder_id"`
}

type storyResourceModel struct {
	Data                 types.String `tfsdk:"data"`
	TinesApiToken        types.String `tfsdk:"tines_api_token"` // Deprecated
	TenantUrl            types.String `tfsdk:"tenant_url"`      // Deprecated
	ID                   types.Int64  `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	UserID               types.Int64  `tfsdk:"user_id"`
	Description          types.String `tfsdk:"description"`
	KeepEventsFor        types.Int64  `tfsdk:"keep_events_for"`
	Disabled             types.Bool   `tfsdk:"disabled"`
	Priority             types.Bool   `tfsdk:"priority"`
	STSEnabled           types.Bool   `tfsdk:"send_to_story_enabled"`
	STSAccessSource      types.String `tfsdk:"send_to_story_access_source"`
	STSAccess            types.String `tfsdk:"send_to_story_access"`
	STSSkillConfirmation types.Bool   `tfsdk:"send_to_story_skill_use_requires_confirmation"`
	SharedTeamSlugs      types.List   `tfsdk:"shared_team_slugs"`
	EntryAgentID         types.Int64  `tfsdk:"entry_agent_id"`
	ExitAgents           types.List   `tfsdk:"exit_agents"`
	TeamID               types.Int64  `tfsdk:"team_id"`
	Tags                 types.List   `tfsdk:"tags"`
	Guid                 types.String `tfsdk:"guid"`
	Slug                 types.String `tfsdk:"slug"`
	CreatedAt            types.String `tfsdk:"created_at"`
	EditedAt             types.String `tfsdk:"edited_at"`
	Mode                 types.String `tfsdk:"mode"`
	FolderID             types.Int64  `tfsdk:"folder_id"`
	Published            types.Bool   `tfsdk:"published"`
	ChangeControlEnabled types.Bool   `tfsdk:"change_control_enabled"`
	Locked               types.Bool   `tfsdk:"locked"`
	Owners               types.List   `tfsdk:"owners"`
	LastUpdated          types.String `tfsdk:"last_updated"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &storyResource{}
	_ resource.ResourceWithConfigure = &storyResource{}
)

// NewStoryResource is a helper function to simplify the provider implementation.
func NewStoryResource() resource.Resource {
	return &storyResource{}
}

// Metadata returns the resource type name.
func (r *storyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_story"
}

const STORY_RESOURCE_DESCRIPTION = `
A Tines Story resource can be managed either via a Story JSON export file, or by setting configuration values on the resource.
We recommend managing Stories via JSON export files if you rely on Terraform to manage change control for storyboard content. Otherwise, if
you use Tines' built-in Change Control feature, we recommend only enforcing configuration values via Terraform.
`

// Schema defines the schema for the resource.
func (r *storyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: STORY_RESOURCE_DESCRIPTION,
		Version:     1, // This needs to be incremented every time we change the schema, and accompanied by a schema migration.
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The Tines-generated identifier for this story.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"data": schema.StringAttribute{
				Description: "A local JSON file containing an exported Tines story. Setting this value can only be combined with the team_id and folder_id attributes.",
				Optional:    true,
			},
			"tines_api_token": schema.StringAttribute{
				Description:        "API token for Tines Tenant",
				Optional:           true,
				DeprecationMessage: "Value will be overridden by the value set in the provider credentials. This field will be removed in a future version.",
				Sensitive:          true,
			},
			"tenant_url": schema.StringAttribute{
				Description:        "Tines tenant URL",
				Optional:           true,
				DeprecationMessage: "Value will be overridden by the value set in the provider credentials. This field will be removed in a future version.",
			},
			"team_id": schema.Int64Attribute{
				Description: "The ID of the team that this story belongs to.",
				Required:    true,
			},
			"folder_id": schema.Int64Attribute{
				Description: "The ID of the folder where this story should be organized. The folder ID must belong to the associated team that owns this story.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the Tines story.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("data")),
				},
			},
			"user_id": schema.Int64Attribute{
				Description: "ID of the story creator.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A user-defined description of the story.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("data")),
				},
			},
			"keep_events_for": schema.Int64Attribute{
				Description: "Defined event retention period in seconds.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.MatchRoot("data")),
				},
			},
			"disabled": schema.BoolAttribute{
				Description: "Boolean flag indicating whether the story is disabled from running.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(path.MatchRoot("data")),
				},
			},
			"priority": schema.BoolAttribute{
				Description: "Boolean flag indicating whether story runs with high priority.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(path.MatchRoot("data")),
				},
			},
			"send_to_story_enabled": schema.BoolAttribute{
				Description:        "Boolean flag indicating if Send to Story is enabled. If enabling Send to Story, the entry_agent_id and exit_agent_ids attributes must also be specified.",
				Optional:           true,
				Computed:           true,
				DeprecationMessage: "This attribute will be removed in a future version. Set the `send_to_story_access_source` attribute instead.",
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(path.MatchRoot("data")),
					boolvalidator.AlsoRequires(path.MatchRoot("entry_agent_id"), path.MatchRoot("exit_agent_ids")),
				},
			},
			"send_to_story_access_source": schema.StringAttribute{
				Description: "Valid values are STS, STS_AND_WORKBENCH, WORKBENCH or OFF indicating where the Send to Story can be used.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("data")),
					stringvalidator.OneOf("STS", "STS_AND_WORKBENCH", "WORKBENCH", "OFF"),
				},
			},
			"send_to_story_access": schema.StringAttribute{
				Description: "Controls who is allowed to send to this story (TEAM, GLOBAL, SPECIFIC_TEAMS). default: TEAM.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("data")),
					stringvalidator.OneOf("TEAM", "GLOBAL", "SPECIFIC_TEAMS"),
				},
			},
			"send_to_story_skill_use_requires_confirmation": schema.BoolAttribute{
				Description: "Boolean flag indicating whether Workbench should ask for confirmation before running this story.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(path.MatchRoot("data")),
				},
			},
			"shared_team_slugs": schema.ListAttribute{
				Description: "Array of team slugs that can send to this story. Required to set send_to_story_access to SPECIFIC_TEAMS.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("data")),
					listvalidator.AlsoRequires(path.MatchRoot("send_to_story_access")),
				},
			},
			"entry_agent_id": schema.Int64Attribute{
				Description: "The ID of the entry action for this story (action must be of type Webhook).",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.MatchRoot("data")),
				},
			},
			"exit_agents": schema.ListAttribute{
				Description: "An Array of IDs describing exit actions for this story (actions must be message-only mode event transformation).",
				ElementType: types.Int64Type,
				Optional:    true,
				Computed:    true,
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("data")),
				},
			},
			"tags": schema.ListAttribute{
				Description: "An array of tag names to apply to the story.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("data")),
				},
			},
			"guid": schema.StringAttribute{
				Description: "The globally unique identifier of the story.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"slug": schema.StringAttribute{
				Description: "An underscored representation of the story name.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "ISO 8601 Timestamp representing date and time the story was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"edited_at": schema.StringAttribute{
				Description: "ISO 8601 Timestamp representing date and time the story was last logically updated.",
				Computed:    true,
			},
			"mode": schema.StringAttribute{
				Description: "The mode of the story (LIVE or TEST).",
				Computed:    true,
			},
			"published": schema.BoolAttribute{
				Description: "Boolean flag indicating whether the story is published.",
				Computed:    true,
			},
			"change_control_enabled": schema.BoolAttribute{
				Description: "Boolean flag indicating if change control is enabled.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(path.MatchRoot("data")),
				},
			},
			"locked": schema.BoolAttribute{
				Description: "Boolean flag indicating whether the story is locked, preventing edits.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(path.MatchRoot("data")),
				},
			},
			"owners": schema.ListAttribute{
				// This is currently a read-only API attribute, but may become read-write in the future.
				Description: "List of user IDs that are listed as owners on the story.",
				ElementType: types.Int64Type,
				Computed:    true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create imports the json data as Tines Story and sets the initial Terraform state.
func (r *storyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Story")

	var plan storyResourceModel
	var story *tines_cli.Story
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Data.IsNull() {
		tflog.Info(ctx, "Exported Story payload detected, using the Import strategy")
		story, diags = r.runImportStory(&plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		tflog.Info(ctx, "No exported Story payload detected, using the Create strategy")
		// Some fields cannot be set at creation time, and require a subsequent update to the Story
		// in order to be set properly.
		var requiresUpdate bool
		var updateStory tines_cli.Story
		var err error

		// Checking for both null and unknown ensures that we're only setting the value
		// in the model if it has been explicitly set in the resource configuration.
		if !plan.ChangeControlEnabled.IsNull() && !plan.ChangeControlEnabled.IsUnknown() {
			requiresUpdate = true
			updateStory.ChangeControlEnabled = plan.ChangeControlEnabled.ValueBool()
		}
		if !plan.STSAccess.IsNull() && !plan.STSAccess.IsUnknown() {
			requiresUpdate = true
			updateStory.STSAccess = plan.STSAccess.ValueString()
		}
		if !plan.STSAccessSource.IsNull() && !plan.STSAccessSource.IsUnknown() {
			requiresUpdate = true
			updateStory.STSAccessSource = plan.STSAccessSource.ValueString()
		}
		if !plan.STSEnabled.IsNull() && !plan.STSEnabled.IsUnknown() {
			requiresUpdate = true
			updateStory.STSEnabled = plan.STSEnabled.ValueBool()
		}
		if !plan.STSSkillConfirmation.IsNull() && !plan.STSSkillConfirmation.IsUnknown() {
			requiresUpdate = true
			updateStory.STSSkillConfirmation = plan.STSSkillConfirmation.ValueBool()
		}
		if !plan.SharedTeamSlugs.IsNull() && !plan.SharedTeamSlugs.IsUnknown() {
			requiresUpdate = true
			diags = plan.SharedTeamSlugs.ElementsAs(ctx, updateStory.SharedTeamSlugs, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
		if !plan.Locked.IsNull() && !plan.Locked.IsUnknown() {
			requiresUpdate = true
			updateStory.Locked = plan.Locked.ValueBool()
		}

		// Create the new story first, so we have something to update if needed.
		newStory := tines_cli.Story{
			TeamID: plan.TeamID.ValueInt64(),
		}

		// Iterate through our optional values to ensure we're only setting parameters
		// in the API request body if they are non-default, otherwise we could unintentionally
		// set a value to an unexpected default. For example, if the keep_events_for value
		// is not explicitly set, the ValueInt64() function will return 0, and we obviously
		// don't want to set the value to 0 by default.
		if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
			newStory.Name = plan.Name.ValueString()
		}

		if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
			newStory.Description = plan.Description.ValueString()
		}

		if !plan.KeepEventsFor.IsNull() && !plan.KeepEventsFor.IsUnknown() {
			newStory.KeepEventsFor = plan.KeepEventsFor.ValueInt64()
		}

		if !plan.FolderID.IsNull() && !plan.FolderID.IsUnknown() {
			newStory.FolderID = plan.FolderID.ValueInt64()
		}

		if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
			diags = plan.Tags.ElementsAs(ctx, newStory.Tags, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		if !plan.Disabled.IsNull() && !plan.Disabled.IsUnknown() {
			newStory.Disabled = plan.Disabled.ValueBool()
		}

		if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
			newStory.Priority = plan.Priority.ValueBool()
		}

		story, err = r.client.CreateStory(&newStory)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Creating Tines Story",
				"Could not create story, unexpected error: "+err.Error(),
			)
			return
		}

		// If the story requires an update to set all values, we set the new fields here
		// and then return the latest API response values for use in updating our Terraform plan.
		// We're not worried about overwriting the story ID value here because it won't change.
		if requiresUpdate {
			tflog.Info(ctx, "Some fields present require an additional update to the Story for the values to be set, running Story Update.")
			story, err = r.client.UpdateStory(story.ID, &updateStory)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Updating Tines Story",
					"Could not update story, unexpected error: "+err.Error(),
				)
				return
			}
		}
	}

	// Populate all the computed values in the plan.
	tflog.Info(ctx, "Populating new plan values")
	diags = r.convertStoryToPlan(ctx, &plan, story)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Retrieve the current infrastructure state.
func (r *storyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var localState storyResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &localState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	remoteState, err := r.client.GetStory(localState.ID.ValueInt64())
	if err != nil {
		// Treat HTTP 404 Not Found status as a signal to recreate resource
		// and return early.
		if tinesErr, ok := err.(tines_cli.Error); ok {
			if tinesErr.StatusCode == 404 {
				resp.State.RemoveResource(ctx)
				return
			}
		}

		resp.Diagnostics.AddError(
			"Unable to Refresh Resource",
			"An unexpected error occurred while attempting to refresh resource state. "+
				"Please retry the operation or report this issue to the provider developers.\n\n"+
				"HTTP Error: "+err.Error(),
		)
		return
	}

	diags := r.convertStoryToPlan(ctx, &localState, remoteState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &localState)...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update performs an in-place update of the import Tines Story and sets the updated Terraform state on success.
func (r *storyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Story")

	var plan, state storyResourceModel
	var story *tines_cli.Story
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Data.IsNull() && !plan.Data.Equal(state.Data) {
		tflog.Info(ctx, "Exported Story payload detected, using the Import strategy")
		story, diags = r.runImportStory(&plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		var storyUpdate tines_cli.Story
		var err error
		diags = r.convertPlanToStory(ctx, &plan, &storyUpdate)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		story, err = r.client.UpdateStory(plan.ID.ValueInt64(), &storyUpdate)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting Tines Story",
				"Could not delete story, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Populate all the computed values in the plan.
	diags = r.convertStoryToPlan(ctx, &plan, story)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the Tines Story and removes the Terraform state on success.
func (r *storyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state storyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing story.
	err := r.client.DeleteStory(state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Tines Story",
			"Could not delete story, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *storyResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (existing state version) to 1 (current Schema.Version).
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "Tines Story identifier",
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"data": schema.StringAttribute{
						Description: "Tines Story export that gets read in from a JSON file",
						Required:    true,
					},
					"tines_api_token": schema.StringAttribute{
						Description: "API token for Tines Tenant",
						Required:    true,
						Sensitive:   true,
					},
					"tenant_url": schema.StringAttribute{
						Description: "Tines tenant URL",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"team_id": schema.Int64Attribute{
						Description: "Tines Team ID",
						Optional:    true,
					},
					"folder_id": schema.Int64Attribute{
						Description: "Tines Folder ID",
						Optional:    true,
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var priorStateData storyResourceModelV0

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)

				if resp.Diagnostics.HasError() {
					return
				}

				upgradedStateData := storyResourceModel{
					ID:            priorStateData.ID,
					Data:          priorStateData.Data,
					TinesApiToken: priorStateData.TinesApiToken,
					TenantUrl:     priorStateData.TenantUrl,
				}

				if !priorStateData.FolderID.IsNull() {
					upgradedStateData.FolderID = priorStateData.FolderID
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
			},
		},
	}

}

// Configure adds the provider configured client to the resource.
func (r *storyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*tines_cli.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Tines Client Configure Type",
			fmt.Sprintf("Expected *tines_cli.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *storyResource) convertPlanToStory(ctx context.Context, plan *storyResourceModel, story *tines_cli.Story) (diags diag.Diagnostics) {
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		story.Name = plan.Name.ValueString()
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		story.Description = plan.Description.ValueString()
	}

	if !plan.KeepEventsFor.IsNull() && !plan.KeepEventsFor.IsUnknown() {
		story.KeepEventsFor = plan.KeepEventsFor.ValueInt64()
	}

	if !plan.Disabled.IsNull() && !plan.Disabled.IsUnknown() {
		story.Disabled = plan.Disabled.ValueBool()
	}

	if !plan.Locked.IsNull() && !plan.Locked.IsUnknown() {
		story.Locked = plan.Locked.ValueBool()
	}

	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		story.Priority = plan.Priority.ValueBool()
	}

	if !plan.STSEnabled.IsNull() && !plan.Priority.IsUnknown() {
		story.STSEnabled = plan.STSEnabled.ValueBool()
	}

	if !plan.STSAccessSource.IsNull() && !plan.STSAccessSource.IsUnknown() {
		story.STSAccessSource = plan.STSAccessSource.ValueString()
	}

	if !plan.STSAccess.IsNull() && !plan.STSAccess.IsUnknown() {
		story.STSAccess = plan.STSAccess.ValueString()
	}

	if !plan.SharedTeamSlugs.IsNull() && !plan.SharedTeamSlugs.IsUnknown() {
		diags = plan.SharedTeamSlugs.ElementsAs(ctx, story.SharedTeamSlugs, false)
		if diags.HasError() {
			return
		}
	}

	if !plan.STSSkillConfirmation.IsNull() && !plan.STSSkillConfirmation.IsUnknown() {
		story.STSSkillConfirmation = plan.STSSkillConfirmation.ValueBool()
	}

	if !plan.EntryAgentID.IsNull() && !plan.EntryAgentID.IsUnknown() {
		story.EntryAgentID = plan.EntryAgentID.ValueInt64()
	}

	if !plan.ExitAgents.IsNull() && !plan.ExitAgents.IsUnknown() {
		diags = plan.ExitAgents.ElementsAs(ctx, story.ExitAgents, false)
		if diags.HasError() {
			return
		}
	}

	if !plan.TeamID.IsNull() && !plan.TeamID.IsUnknown() {
		story.TeamID = plan.TeamID.ValueInt64()
	}

	if !plan.FolderID.IsNull() && !plan.FolderID.IsUnknown() {
		story.FolderID = plan.FolderID.ValueInt64()
	}

	return diags
}

// This is reused in both the Create and Update methods.
func (r *storyResource) convertStoryToPlan(ctx context.Context, plan *storyResourceModel, story *tines_cli.Story) (diags diag.Diagnostics) {
	// Populate all the computed values in the plan.

	plan.ID = types.Int64Value(story.ID)
	plan.Name = types.StringValue(story.Name)
	plan.UserID = types.Int64Value(story.UserID)
	plan.Description = types.StringValue(story.Description)
	plan.KeepEventsFor = types.Int64Value(story.KeepEventsFor)
	plan.Disabled = types.BoolValue(story.Disabled)
	plan.Priority = types.BoolValue(story.Priority)
	plan.STSEnabled = types.BoolValue(story.STSEnabled)
	plan.STSAccessSource = types.StringValue(story.STSAccessSource)
	plan.STSAccess = types.StringValue(story.STSAccess)
	plan.STSSkillConfirmation = types.BoolValue(story.STSSkillConfirmation)
	plan.SharedTeamSlugs, diags = types.ListValueFrom(ctx, types.StringType, story.SharedTeamSlugs)
	if diags.HasError() {
		return diags
	}
	plan.EntryAgentID = types.Int64Value(story.EntryAgentID)
	plan.ExitAgents, diags = types.ListValueFrom(ctx, types.Int64Type, story.ExitAgents)
	if diags.HasError() {
		return diags
	}
	plan.TeamID = types.Int64Value(story.TeamID)
	plan.Tags, diags = types.ListValueFrom(ctx, types.StringType, story.Tags)
	if diags.HasError() {
		return diags
	}
	plan.Guid = types.StringValue(story.Guid)
	plan.Slug = types.StringValue(story.Slug)
	plan.CreatedAt = types.StringValue(story.CreatedAt)
	plan.EditedAt = types.StringValue(story.EditedAt)
	plan.Mode = types.StringValue(story.Mode)
	plan.FolderID = types.Int64Value(story.FolderID)
	plan.Published = types.BoolValue(story.Published)
	plan.ChangeControlEnabled = types.BoolValue(story.ChangeControlEnabled)
	plan.Locked = types.BoolValue(story.Locked)
	plan.Owners, diags = types.ListValueFrom(ctx, types.Int64Type, story.Owners)
	if diags.HasError() {
		return diags
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	return diags
}

// This is reused in both the Create and Update methods when the resource management strategy is set to use
// imported stories to override all values.
func (r *storyResource) runImportStory(plan *storyResourceModel) (story *tines_cli.Story, diags diag.Diagnostics) {
	var data map[string]interface{}

	encData := plan.Data.ValueString()

	err := json.Unmarshal([]byte(encData), &data)
	if err != nil {
		diags.AddError("Invalid JSON in file", err.Error())
		return
	}

	name, ok := data["name"].(string)
	if !ok {
		diags.AddError("Invalid string", "The 'name' field in the imported story must be a string")
		return
	}

	var importRequest = tines_cli.StoryImportRequest{
		NewName: name,
		Data:    data,
		TeamID:  plan.TeamID.ValueInt64(),
		Mode:    "versionReplace",
	}

	if !plan.FolderID.IsNull() && !plan.FolderID.IsUnknown() {
		importRequest.FolderID = plan.FolderID.ValueInt64()
	}

	story, err = r.client.ImportStory(&importRequest)
	if err != nil {
		diags.AddError(
			"Error Importing Tines Story",
			"Could not import story, unexpected error: "+err.Error(),
		)
		return
	}

	return story, diags
}
