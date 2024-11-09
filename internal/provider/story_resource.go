package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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

// Schema defines the schema for the resource.
func (r *storyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a Tines Story",
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
				Required:    true,
			},
			"tines_api_token": schema.StringAttribute{
				Description:        "API token for Tines Tenant",
				Optional:           true,
				DeprecationMessage: "Value will be overridden by the value set in the provider credentials. This field will be removed in a future version.",
				Sensitive:          true,
			},
			"tenant_url": schema.StringAttribute{
				Description:        "[DEPRECATED] Tines tenant URL",
				Optional:           true,
				DeprecationMessage: "Value will be overridden by the value set in the provider credentials. This field will be removed in a future version.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"team_id": schema.Int64Attribute{
				Description: "The ID of the team that this story belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"folder_id": schema.Int64Attribute{
				Description: "The ID of the folder where this story should be organized. The folder ID must belong to the associated team that owns this story.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
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
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"keep_events_for": schema.Int64Attribute{
				Computed: true,
			},
			"disabled": schema.BoolAttribute{
				Computed: true,
			},
			"priority": schema.BoolAttribute{
				Computed: true,
			},
			"send_to_story_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"send_to_story_access_source": schema.StringAttribute{
				Computed: true,
			},
			"send_to_story_access": schema.StringAttribute{
				Computed: true,
			},
			"send_to_story_skill_use_requires_confirmation": schema.BoolAttribute{
				Computed: true,
			},
			"shared_team_slugs": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"entry_agent_id": schema.Int64Attribute{
				Computed: true,
			},
			"exit_agents": schema.ListAttribute{
				ElementType: types.Int64Type,
				Computed:    true,
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"guid": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"slug": schema.StringAttribute{
				Computed: true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"edited_at": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mode": schema.StringAttribute{
				Computed: true,
			},
			"published": schema.BoolAttribute{
				Computed: true,
			},
			"change_control_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"locked": schema.BoolAttribute{
				Computed: true,
			},
			"owners": schema.ListAttribute{
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
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Data.IsNull() {
		tflog.Info(ctx, "Exported Story payload detected, using the Import strategy")
	}

	var data map[string]interface{}

	encData := plan.Data.ValueString()

	err := json.Unmarshal([]byte(encData), &data)
	if err != nil {
		resp.Diagnostics.AddError("Invalid JSON in file", err.Error())
		return
	}

	name, ok := data["name"].(string)
	if !ok {
		resp.Diagnostics.AddError("Invalid string", "The 'name' field in the imported story must be a string")
		return
	}

	var importRequest = tines_cli.StoryImportRequest{
		NewName:  name,
		Data:     data,
		TeamID:   plan.TeamID.ValueInt64(),
		FolderID: plan.FolderID.ValueInt64(),
		Mode:     "versionReplace",
	}

	story, err := r.client.ImportStory(&importRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Tines Story",
			"Could not import story, unexpected error: "+err.Error(),
		)
		return
	}

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
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.EntryAgentID = types.Int64Value(story.EntryAgentID)
	plan.ExitAgents, diags = types.ListValueFrom(ctx, types.Int64Type, story.ExitAgents)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.TeamID = types.Int64Value(story.TeamID)
	plan.Tags, diags = types.ListValueFrom(ctx, types.StringType, story.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
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
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

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

	status, remoteState, err := r.client.GetStory(localState.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Refresh Resource",
			"An unexpected error occurred while attempting to refresh resource state. "+
				"Please retry the operation or report this issue to the provider developers.\n\n"+
				"HTTP Error: "+err.Error(),
		)
		return
	}

	// Treat HTTP 404 Not Found status as a signal to recreate resource
	// and return early.
	if status == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	localState.TeamID = types.Int64Value(remoteState.TeamID)
	localState.FolderID = types.Int64Value(remoteState.FolderID)

	// Set refreshed state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &localState)...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update performs an in-place update of the import Tines Story and sets the updated Terraform state on success.
func (r *storyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Story")

	var plan storyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Data.IsNull() {
		tflog.Info(ctx, "Exported Story payload detected, using the Import strategy")
	}

	var data map[string]interface{}

	encData := plan.Data.ValueString()

	err := json.Unmarshal([]byte(encData), &data)
	if err != nil {
		resp.Diagnostics.AddError("Invalid JSON in file", err.Error())
		return
	}

	name, ok := data["name"].(string)
	if !ok {
		resp.Diagnostics.AddError("Invalid string", "The 'name' field in the imported story must be a string")
		return
	}

	var importRequest = tines_cli.StoryImportRequest{
		NewName:  name,
		Data:     data,
		TeamID:   plan.TeamID.ValueInt64(),
		FolderID: plan.FolderID.ValueInt64(),
		Mode:     "versionReplace",
	}

	story, err := r.client.ImportStory(&importRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Tines Story",
			"Could not import story, unexpected error: "+err.Error(),
		)
		return
	}

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
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.EntryAgentID = types.Int64Value(story.EntryAgentID)
	plan.ExitAgents, diags = types.ListValueFrom(ctx, types.Int64Type, story.ExitAgents)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.TeamID = types.Int64Value(story.TeamID)
	plan.Tags, diags = types.ListValueFrom(ctx, types.StringType, story.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
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
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

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
func (s *storyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*tines_cli.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *tines_cli.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	s.client = client
}
