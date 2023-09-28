package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"fmt"
)

// storyResource is the resource implementation.
type storyResource struct {
	client *Client
}

type storyResourceModel struct {
	Data          types.String `tfsdk:"data"`
	ID            types.Int64  `tfsdk:"id"`
	TinesApiToken types.String `tfsdk:"tines_api_token"`
	TenantUrl     types.String `tfsdk:"tenant_url"`
	TeamID        types.Int64  `tfsdk:"team_id"`
	FolderID      types.Int64  `tfsdk:"folder_id"`
}

type StoryImportRequest struct {
	NewName  string                 `json:"new_name"`
	TeamID   int64                  `json:"team_id,omitempty"`
	FolderID int64                  `json:"folder_id,omitempty"`
	Mode     string                 `json:"mode,omitempty"`
	Data     map[string]interface{} `json:"data"`
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
		Description: "Manage a Tines Story.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Tines Story identifier.",
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
				Description: "Tines team ID.",
				Optional:    true,
			},
			"folder_id": schema.Int64Attribute{
				Description: "Tines folder ID.",
				Optional:    true,
			},
		},
	}
}

// Create imports the json data as Tines Story and sets the initial Terraform state.
func (r *storyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Importing Story")

	var plan storyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.client.apiToken = plan.TinesApiToken.ValueString()
	r.client.tenantUrl = plan.TenantUrl.ValueString()

	sirB, err := r.client.buildStoryData(plan.Data.String(), plan.TeamID.ValueInt64(), plan.FolderID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing data",
			"Could not create story, unexpected error: "+err.Error(),
		)
		return
	}

	s, err := r.client.ImportStory(ctx, sirB)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating story",
			"Could not create story, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.Int64Value(s.ID)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Keeping this as a placeholder.
func (r *storyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	r.client.apiToken = plan.TinesApiToken.ValueString()
	r.client.tenantUrl = plan.TenantUrl.ValueString()

	sirB, err := r.client.buildStoryData(plan.Data.String(), plan.TeamID.ValueInt64(), plan.FolderID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing data",
			"Could not create story, unexpected error: "+err.Error(),
		)
		return
	}

	s, err := r.client.ImportStory(ctx, sirB)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating story",
			"Could not update story, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.Int64Value(s.ID)

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

	r.client.apiToken = state.TinesApiToken.ValueString()
	r.client.tenantUrl = state.TenantUrl.ValueString()

	// Delete existing story
	_, err := r.client.DeleteStory(ctx, state.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Tines Story",
			"Could not delete story, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (s *storyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *tines.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	s.client = client
}
