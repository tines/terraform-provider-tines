package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/dynamicvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/tines/go-sdk/tines"
)

// tinesResource is the resource implementation. Semantically,
// This should probably be called a resourceResource but that
// gets very confusing, so we make an exception and refer to this
// as a tinesResource.

type tinesResource struct {
	client *tines.Client
}

type tinesResourceModel struct {
	Id           types.Int64   `tfsdk:"id"`
	Name         types.String  `tfsdk:"name"`
	Description  types.String  `tfsdk:"description"`
	Value        types.Dynamic `tfsdk:"value"`
	TeamId       types.Int64   `tfsdk:"team_id"`
	FolderId     types.Int64   `tfsdk:"folder_id"`
	UserId       types.Int64   `tfsdk:"user_id"`
	ReadAccess   types.String  `tfsdk:"read_access"`
	SharedTeams  types.List    `tfsdk:"shared_team_slugs"`
	Slug         types.String  `tfsdk:"slug"`
	TestEnabled  types.Bool    `tfsdk:"test_resource_enabled"`
	TestResource types.Dynamic `tfsdk:"test_resource"`
	TestValue    types.Dynamic `tfsdk:"test_value"`
	IsTest       types.Bool    `tfsdk:"is_test"`
	LiveResId    types.Int64   `tfsdk:"live_resource_id"`
	CreatedAt    types.String  `tfsdk:"created_at"`
	UpdatedAt    types.String  `tfsdk:"updated_at"`
	RefActions   types.List    `tfsdk:"referencing_action_ids"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &tinesResource{}
	_ resource.ResourceWithConfigure   = &tinesResource{}
	_ resource.ResourceWithImportState = &tinesResource{}
)

// NewTinesResource is a helper function to simplify the provider implementation.
func NewTinesResource() resource.Resource {
	return &tinesResource{}
}

// Metadata returns the resource type name.
func (r *tinesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource"
}

const TINES_RESOURCE_DESCRIPTION = `
Tines Resources can be created with strings, arrays, or JSON objects as values. In general, Tines Resources should be used
for non-secret values that are reused repeatedly across Stories or Actions: e.g. things like domain names, references to
lists of users or IP addresses, etc. For secret values such as passwords and API keys, you should use Tines Credentials instead.
Throughout this documentation, we attempt to distinguish Tines Resources (which are stored in the Tines platform) from Terraform
Resources (which are stored in Terraform state files), although Tines Resources can be managed as a Terraform Resource.`

func (r *tinesResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: TINES_RESOURCE_DESCRIPTION,
		Version:     0, // This needs to be incremented every time we change the schema, and accompanied by a schema migration.
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The Tines-generated identifier for this Tines Resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the Tines Resource.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A long-form description of the Tines Resource.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("Managed via Terraform"),
			},
			"value": schema.DynamicAttribute{
				Description: "Contents of the Tines Resource as a JSON array, object, or string.",
				Required:    true,
			},
			"team_id": schema.Int64Attribute{
				Description: "The ID of Tines Team where this Tines Resource will be located.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"folder_id": schema.Int64Attribute{
				Description: "The ID of folder where the Tines Resource will be located.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.Int64Attribute{
				Description: "The ID of user that created the Tines Resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"read_access": schema.StringAttribute{
				Description: "Controls who is allowed to use this Tines Resource (TEAM, GLOBAL, SPECIFIC_TEAMS). default: TEAM.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("TEAM", "GLOBAL", "SPECIFIC_TEAMS"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"shared_team_slugs": schema.ListAttribute{
				Description: "List of teams' slugs where this resource can be used. Required to set read_access to SPECIFIC_TEAMS.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []validator.List{
					listvalidator.AlsoRequires(path.MatchRoot("read_access")),
				},
			},
			"slug": schema.StringAttribute{
				Description: "An underscored representation of the Tines Resource name.",
				Computed:    true,
			},
			"test_resource_enabled": schema.BoolAttribute{
				Description: "A boolean value indicating whether the Tines Resource is enabled for using a test Tines Reesource value during non-production Story execution.",
				Optional:    true,
				Computed:    true,
			},
			"test_resource": schema.DynamicAttribute{
				Description: "A JSON block representing the test version of this Tines Resource.",
				Computed:    true,
			},
			"test_value": schema.DynamicAttribute{
				Description: "Contents of the test version of this Tines Resource as a JSON array, object, or string.",
				Optional:    true,
				Validators: []validator.Dynamic{
					dynamicvalidator.AlsoRequires(path.MatchRoot("test_resource_enabled")),
				},
			},
			"is_test": schema.BoolAttribute{
				Description: "Boolean flag indicating whether the Tines Resource production or test value should be updated.",
				Optional:    true,
			},
			"live_resource_id": schema.Int64Attribute{
				Description: "Optional when updating a test Tines Resource value.",
				Optional:    true,
			},
			"referencing_action_ids": schema.ListAttribute{
				Description: "A list of Action IDs in Tines Stories that reference this Tines Resource value. This Resource should not be removed if this is non-null.",
				ElementType: types.Int64Type,
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The ISO 8601 Timestamp representing date and time the Tines Resource was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "The ISO 8601 Timestamp representing date and time the Tines Resource was last updated.",
				Computed:    true,
			},
		},
	}
}

// Creates a new Tines Resource and sets the initial Terraform state.
func (r *tinesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Tines Resource")
	var plan tinesResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	val, diags := r.getUnderlyingDynamicValue(ctx, &plan.Value)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newResource := tines.Resource{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		TeamId:      int(plan.TeamId.ValueInt64()),
		Value:       val,
	}

	// Add optional attributes to the new Tines Resource if they have been set.
	if !plan.FolderId.IsNull() && !plan.FolderId.IsUnknown() {
		newResource.FolderId = int(plan.Id.ValueInt64())
	}

	if !plan.ReadAccess.IsNull() && !plan.ReadAccess.IsUnknown() {
		newResource.ReadAccess = plan.ReadAccess.ValueString()
	}

	if !plan.SharedTeams.IsNull() && !plan.SharedTeams.IsUnknown() {
		diags = plan.SharedTeams.ElementsAs(ctx, newResource.SharedTeams, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !plan.TestEnabled.IsNull() && !plan.TestEnabled.IsUnknown() {
		newResource.TestResEnabled = plan.TestEnabled.ValueBool()
	}

	// Create the Tines Resource. The API does not permit creating a Tines Resource with an associated
	// Test Tines Resource at the same time, so if there is a test value specified, we handle it in a
	// separate step.
	tr, err := r.client.CreateResource(ctx, &newResource)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Tines Resource",
			"Could not create resource, unexpected error: "+err.Error(),
		)
		return
	}

	// Check and see if there is a test value associated with this Tines Resource, and add it to the
	// parent resource if so.
	updateRequired := false
	updatedResource := tines.Resource{
		Id:           tr.Id,
		TestResource: true,
	}

	if !plan.TestEnabled.IsNull() && !plan.TestEnabled.IsUnknown() {
		updatedResource.TestResEnabled = plan.TestEnabled.ValueBool()
		updateRequired = true
	}

	if !plan.TestValue.IsNull() && !plan.TestValue.IsUnknown() {
		updatedResource.Value = plan.TestValue.UnderlyingValue()
		updateRequired = true
	}

	if updateRequired {
		tr, err = r.client.UpdateResource(ctx, updatedResource.Id, &updatedResource)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Adding Test Tines Resource",
				"Could not create test version of resource, unexpected error: "+err.Error(),
			)
			return
		}

	}

	// Convert populated resource values to Terraform types in the plan.
	diags = r.convertTinesResourceToPlan(ctx, &plan, tr)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data from the plan.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Retrieve the current infrastructure state.
func (r *tinesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var localState tinesResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &localState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	remoteState, err := r.client.GetResource(ctx, int(localState.Id.ValueInt64()))
	if err != nil {
		// Treat HTTP 404 Not Found status as a signal to recreate resource
		// and return early.
		if tinesErr, ok := err.(tines.Error); ok {
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

	diags := r.convertTinesResourceToPlan(ctx, &localState, remoteState)
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

// Update performs an in-place update of the Tines Resource and sets the updated Terraform state on success.
func (r *tinesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Tines Resource")
	var plan, state tinesResourceModel
	var resourceUpdate tines.Resource

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

	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		resourceUpdate.Id = int(plan.Id.ValueInt64())
	}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		resourceUpdate.Name = plan.Name.ValueString()
	}

	if !plan.Value.IsNull() && !plan.Value.IsUnderlyingValueUnknown() {
		val, diags := r.getUnderlyingDynamicValue(ctx, &plan.Value)
		if diags.HasError() {
			return
		}
		resourceUpdate.Value = val
	}

	newResource, err := r.client.UpdateResource(ctx, resourceUpdate.Id, &resourceUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Tines Resource",
			"Could not update Tines Resource, unexpected error: "+err.Error(),
		)
		return
	}

	// Populate all the computed values in the plan.
	diags = r.convertTinesResourceToPlan(ctx, &plan, newResource)
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

// Deletes the Tines Resource and removes the Terraform state on success.
func (r *tinesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Tines Resource")
	// Retrieve values from state
	var state tinesResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Tines Resource.
	err := r.client.DeleteResource(ctx, int(state.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Tines Resource",
			"Could not delete Tines Resource, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *tinesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing Tines Resource")
	// Retrieve import ID and save to id attribute
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Import Error",
			"Could not determine the ID of the Tines Resource, unexpected error: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// Configure adds the provider configured client to the resource.
func (r *tinesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*tines.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Tines Client Configure Type",
			fmt.Sprintf("Expected *tines.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *tinesResource) convertTinesResourceToPlan(ctx context.Context, plan *tinesResourceModel, tr *tines.Resource) (diags diag.Diagnostics) {
	// Note: we don't attempt to convert the TestValue from a Tines Resource to a plan value
	// because that construct doesn't exist outside of Terraform. It's a convenience method
	// to be able to know when an API call for the `value` is meant to update the test vs
	// the production value.
	plan.Id = types.Int64Value(int64(tr.Id))
	plan.Name = types.StringValue(tr.Name)
	plan.Value = types.DynamicValue(types.StringValue(tr.Value.(string)))
	plan.Description = types.StringValue(tr.Description)
	plan.TeamId = types.Int64Value(int64(tr.TeamId))
	plan.FolderId = types.Int64Value(int64(tr.FolderId))
	plan.UserId = types.Int64Value(int64(tr.UserId))
	plan.ReadAccess = types.StringValue(tr.ReadAccess)
	plan.SharedTeams, diags = types.ListValueFrom(ctx, types.StringType, tr.SharedTeams)
	if diags.HasError() {
		return diags
	}
	plan.Slug = types.StringValue(tr.Slug)
	plan.TestEnabled = types.BoolValue(tr.TestResEnabled)
	// Set the computed value of TestResource to nil if one is not enabled - the Tines
	// API won't return a value for test_resource if one is not enabled.
	if tr.TestResource != nil {
		plan.TestResource = types.DynamicValue(types.StringValue(tr.TestResource.(string)))
	} else {
		plan.TestResource = types.DynamicValue(types.StringValue(""))
	}
	plan.RefActions, diags = types.ListValueFrom(ctx, types.Int64Type, tr.RefActions)
	if diags.HasError() {
		return diags
	}
	plan.CreatedAt = types.StringValue(tr.CreatedAt)
	plan.UpdatedAt = types.StringValue(tr.UpdatedAt)

	return diags
}

func (r *tinesResource) getUnderlyingDynamicValue(ctx context.Context, res *types.Dynamic) (any, diag.Diagnostics) {
	// Handle the underlying value in the dynamic type.
	var diags diag.Diagnostics
	switch val := res.UnderlyingValue().(type) {
	case types.String:
		tflog.Info(ctx, "Dynamic value is a string")
		return val.ValueString(), nil
	case types.Number:
		tflog.Info(ctx, "Dynamic value is a number")
		return val.ToNumberValue(ctx)
	case types.Tuple:
		tflog.Info(ctx, "Dynamic value is a tuple")
		return nil, nil
	}

	tflog.Info(ctx, "Dynamic value is an unsupported type")
	diags.AddError("Dynamic value is an unsupported type.", "Only String, Number, and Tuple values are supported.")
	return nil, diags
}
