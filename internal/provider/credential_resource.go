package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/tines/go-sdk/tines"
)

// credentialResource is the resource implementation for Tines Credentials.
type credentialResource struct {
	client *tines.Client
}

type credentialResourceModel struct {
	Id             types.Int64  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Mode           types.String `tfsdk:"mode"`
	Value          types.String `tfsdk:"value"`
	ValueWo        types.String `tfsdk:"value_wo"`
	ValueWoVersion types.Int64  `tfsdk:"value_wo_version"`
	TeamId         types.Int64  `tfsdk:"team_id"`
	FolderId       types.Int64  `tfsdk:"folder_id"`
	Description    types.String `tfsdk:"description"`
	ReadAccess     types.String `tfsdk:"read_access"`
	SharedTeams    types.List   `tfsdk:"shared_team_slugs"`
	Slug           types.String `tfsdk:"slug"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &credentialResource{}
	_ resource.ResourceWithConfigure   = &credentialResource{}
	_ resource.ResourceWithImportState = &credentialResource{}
)

// NewCredentialResource is a helper function to simplify the provider implementation.
func NewCredentialResource() resource.Resource {
	return &credentialResource{}
}

// Metadata returns the resource type name.
func (r *credentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential"
}

const TINES_CREDENTIAL_DESCRIPTION = `
Tines Credentials store secret values such as API keys, passwords, and tokens that are referenced by Stories and Actions
without exposing the secret value itself. This resource currently supports managing TEXT mode Credentials, which makes it
well suited to automated secret rotation: when an upstream token is rotated, the new value can be written to the associated
Tines Credential as part of the same Terraform apply. For non-secret values that are reused across Stories, use a Tines
Resource (` + "`tines_resource`" + `) instead.`

func (r *credentialResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: TINES_CREDENTIAL_DESCRIPTION,
		Version:     0, // This needs to be incremented every time we change the schema, and accompanied by a schema migration.
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The Tines-generated identifier for this Tines Credential.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the Tines Credential.",
				Required:    true,
			},
			"mode": schema.StringAttribute{
				Description: "The type of the Tines Credential. Only TEXT is currently supported by this resource.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("TEXT"),
				Validators: []validator.String{
					stringvalidator.OneOf("TEXT"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				Description: "The secret value of the Tines Credential. The Tines API never returns this value, so it cannot be " +
					"populated on import and drift is detected from configuration only. This attribute is stored in Terraform " +
					"state (marked sensitive); for secrets that must never be persisted to state, use `value_wo` instead. " +
					"Exactly one of `value` or `value_wo` must be set. Rotate a secret by updating this value.",
				Optional:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("value"),
						path.MatchRoot("value_wo"),
					),
				},
			},
			"value_wo": schema.StringAttribute{
				Description: "The secret value of the Tines Credential, supplied as a write-only argument so that it is never " +
					"persisted to Terraform state (requires Terraform >= 1.11). Because write-only values are not stored, " +
					"changes are not detected automatically: bump `value_wo_version` to signal that the secret should be " +
					"rewritten to Tines (e.g. during secret rotation). Exactly one of `value` or `value_wo` must be set.",
				Optional:  true,
				Sensitive: true,
				WriteOnly: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.AlsoRequires(path.MatchRoot("value_wo_version")),
				},
			},
			"value_wo_version": schema.Int64Attribute{
				Description: "A user-managed version counter for `value_wo`. Increment this whenever the write-only secret " +
					"changes so that Terraform triggers an in-place update (secret rotation). Has no effect unless " +
					"`value_wo` is set.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.MatchRoot("value_wo")),
				},
			},
			"team_id": schema.Int64Attribute{
				Description: "The ID of the Tines Team where this Tines Credential will be located.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"folder_id": schema.Int64Attribute{
				Description: "The ID of the folder where the Tines Credential will be located.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A long-form description of the Tines Credential.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("Managed via Terraform"),
			},
			"read_access": schema.StringAttribute{
				Description: "Controls who is allowed to use this Tines Credential (TEAM, GLOBAL, SPECIFIC_TEAMS). default: TEAM.",
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
				Description: "List of teams' slugs where this Credential can be used. Required to set read_access to SPECIFIC_TEAMS.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []validator.List{
					listvalidator.AlsoRequires(path.MatchRoot("read_access")),
				},
			},
			"slug": schema.StringAttribute{
				Description: "An underscored representation of the Tines Credential name.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The ISO 8601 Timestamp representing date and time the Tines Credential was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "The ISO 8601 Timestamp representing date and time the Tines Credential was last updated.",
				Computed:    true,
			},
		},
	}
}

// Create creates a new Tines Credential and sets the initial Terraform state.
func (r *credentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Tines Credential")
	var plan credentialResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Write-only attributes are not present in the plan, so read the effective
	// secret value from configuration.
	secret, diags := r.resolveSecretValue(ctx, req.Config, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newCred := tines.Credential{
		Name:        plan.Name.ValueString(),
		Mode:        tines.CredentialType(plan.Mode.ValueString()),
		Description: plan.Description.ValueString(),
		TeamId:      int(plan.TeamId.ValueInt64()),
		CredentialPayload: tines.CredentialPayload{
			TextValue: secret,
		},
	}

	// Add optional attributes to the new Tines Credential if they have been set.
	if !plan.FolderId.IsNull() && !plan.FolderId.IsUnknown() {
		newCred.FolderId = int(plan.FolderId.ValueInt64())
	}

	if !plan.ReadAccess.IsNull() && !plan.ReadAccess.IsUnknown() {
		newCred.ReadAccess = plan.ReadAccess.ValueString()
	}

	if !plan.SharedTeams.IsNull() && !plan.SharedTeams.IsUnknown() {
		diags = plan.SharedTeams.ElementsAs(ctx, &newCred.SharedTeams, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	cred, err := r.client.CreateCredential(ctx, &newCred)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Tines Credential",
			"Could not create credential, unexpected error: "+err.Error(),
		)
		return
	}

	// Convert populated credential values to Terraform types in the plan. The
	// secret value is preserved from the plan because the API never returns it.
	diags = r.convertCredentialToPlan(ctx, &plan, cred)
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

// Read retrieves the current infrastructure state.
func (r *credentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var localState credentialResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &localState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	remoteState, err := r.client.GetCredential(ctx, int(localState.Id.ValueInt64()))
	if err != nil {
		// Treat HTTP 404 Not Found status as a signal to recreate the resource
		// and return early.
		//
		// NOTE: as of go-sdk v0.2.2/v0.3.0 this branch cannot currently fire.
		// GetCredential re-wraps the underlying request error in a fresh
		// tines.Error without copying StatusCode (and the numeric code is not
		// recoverable from the message string), so tinesErr.StatusCode is
		// always 0 here. GetResource has the same limitation; GetStory returns
		// the error unwrapped, which is why the identical pattern works for
		// stories. The correct fix is in go-sdk (propagate StatusCode when
		// wrapping, or return the error unwrapped like GetStory). The check is
		// retained so that "recreate on 404" begins working automatically once
		// the SDK is fixed, without a further change here. See PR #85 review.
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

	// This resource only manages TEXT mode Credentials. If an existing
	// Credential of another mode is imported (ImportState only sets the ID and
	// cannot enforce the mode validator), refuse to reconcile it rather than
	// writing a non-TEXT mode into state, which would otherwise cause the next
	// plan to propose a destroy-and-recreate of a Credential this resource
	// cannot recreate.
	if remoteState.Mode != tines.CredentialTypeText {
		resp.Diagnostics.AddError(
			"Unsupported Tines Credential Mode",
			fmt.Sprintf(
				"Tines Credential %d has mode %q, but the tines_credential resource only supports %q mode "+
					"Credentials. Managing or importing non-TEXT Credentials is not supported.",
				remoteState.Id, remoteState.Mode, tines.CredentialTypeText,
			),
		)
		return
	}

	// The Tines API does not return the secret value, so we retain whatever is
	// already in state. convertCredentialToPlan leaves the Value field untouched.
	diags := r.convertCredentialToPlan(ctx, &localState, remoteState)
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

// Update performs an in-place update of the Tines Credential and sets the updated Terraform state on success.
func (r *credentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Tines Credential")
	var plan, state credentialResourceModel

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

	// Write-only attributes are not present in the plan, so read the effective
	// secret value from configuration.
	secret, diags := r.resolveSecretValue(ctx, req.Config, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The SDK requires both the Id and Mode to be set when updating a Credential.
	credUpdate := tines.Credential{
		Id:   int(state.Id.ValueInt64()),
		Mode: tines.CredentialType(plan.Mode.ValueString()),
		CredentialPayload: tines.CredentialPayload{
			TextValue: secret,
		},
	}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		credUpdate.Name = plan.Name.ValueString()
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		credUpdate.Description = plan.Description.ValueString()
	}

	if !plan.FolderId.IsNull() && !plan.FolderId.IsUnknown() {
		credUpdate.FolderId = int(plan.FolderId.ValueInt64())
	}

	if !plan.ReadAccess.IsNull() && !plan.ReadAccess.IsUnknown() {
		credUpdate.ReadAccess = plan.ReadAccess.ValueString()
	}

	if !plan.SharedTeams.IsNull() && !plan.SharedTeams.IsUnknown() {
		diags = plan.SharedTeams.ElementsAs(ctx, &credUpdate.SharedTeams, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	cred, err := r.client.UpdateCredential(ctx, credUpdate.Id, &credUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Tines Credential",
			"Could not update Tines Credential, unexpected error: "+err.Error(),
		)
		return
	}

	// Populate all the computed values in the plan.
	diags = r.convertCredentialToPlan(ctx, &plan, cred)
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

// Delete deletes the Tines Credential and removes the Terraform state on success.
func (r *credentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Tines Credential")
	// Retrieve values from state.
	var state credentialResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Tines Credential.
	err := r.client.DeleteCredential(ctx, int(state.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Tines Credential",
			"Could not delete Tines Credential, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *credentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing Tines Credential")
	// Retrieve import ID and save to id attribute.
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Import Error",
			"Could not determine the ID of the Tines Credential, unexpected error: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// Configure adds the provider configured client to the resource.
func (r *credentialResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// resolveSecretValue returns the effective secret to send to the Tines API,
// reading the write-only `value_wo` attribute from configuration when it is set
// and otherwise falling back to the state-persisted `value` attribute. The
// schema guarantees (via ExactlyOneOf) that exactly one of the two is supplied.
func (r *credentialResource) resolveSecretValue(ctx context.Context, config tfsdk.Config, plan *credentialResourceModel) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Write-only attributes are never present in the plan or state; they can
	// only be read from configuration.
	var valueWo types.String
	diags.Append(config.GetAttribute(ctx, path.Root("value_wo"), &valueWo)...)
	if diags.HasError() {
		return "", diags
	}

	if !valueWo.IsNull() && !valueWo.IsUnknown() {
		return valueWo.ValueString(), diags
	}

	return plan.Value.ValueString(), diags
}

// convertCredentialToPlan maps API response values onto the Terraform plan. The
// secret `value`/`value_wo` attributes are intentionally not modified here
// because the Tines API never returns the secret; the configured value (or, for
// write-only, nothing) is preserved instead.
func (r *credentialResource) convertCredentialToPlan(ctx context.Context, plan *credentialResourceModel, cred *tines.Credential) (diags diag.Diagnostics) {
	plan.Id = types.Int64Value(int64(cred.Id))
	plan.Name = types.StringValue(cred.Name)
	plan.Mode = types.StringValue(string(cred.Mode))
	plan.Description = types.StringValue(cred.Description)
	plan.TeamId = types.Int64Value(int64(cred.TeamId))
	plan.FolderId = types.Int64Value(int64(cred.FolderId))
	plan.ReadAccess = types.StringValue(cred.ReadAccess)
	plan.SharedTeams, diags = types.ListValueFrom(ctx, types.StringType, cred.SharedTeams)
	if diags.HasError() {
		return diags
	}
	plan.Slug = types.StringValue(cred.Slug)
	plan.CreatedAt = types.StringValue(cred.CreatedAt)
	plan.UpdatedAt = types.StringValue(cred.UpdatedAt)

	return diags
}
