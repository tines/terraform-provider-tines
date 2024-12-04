package provider

import (
	"context"
	"os"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/tines/terraform-provider-tines/internal/tines_cli"
)

// Ensure TinesProvider satisfies various provider interfaces.
var _ provider.Provider = &TinesProvider{}

// TinesProvider defines the provider implementation.
type TinesProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// tinesProviderModel describes the provider data model.
type tinesProviderModel struct {
	Tenant types.String `tfsdk:"tenant"`
	ApiKey types.String `tfsdk:"api_key"`
}

func (p *TinesProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tines"
	resp.Version = p.version
}

// Schema describes the schema for Tines provider.
func (p *TinesProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"tenant": schema.StringAttribute{
				Optional:    true,
				Description: "If this value is not set in the configuration, you must set the TINES_TENANT environment variable instead.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^https:\/\/[a-zA-Z0-9-\.]+\.[a-zA-Z0-9-]+\.[a-zA-Z0-9-]+$`),
						"must be a valid hostname in the format https://example.tines.com",
					),
				},
			},
			"api_key": schema.StringAttribute{
				Optional:    true,
				Description: "If this value is not set in the configuration, you must set the TINES_API_KEY environment variable instead.",
				Sensitive:   true,
			},
		},
	}
}

// Configure initializes the Tines client.
func (p *TinesProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Tines client")

	var config tinesProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.Tenant.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("tenant"),
			"Uknown Tines Tenant",
			"The provider cannot create the Tines API client as there is an unknown configuration value for the Tines tenant",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Tines API Key",
			"The provider cannot create the Tines API client as there is an unknown configuration value for the Tines API key",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	tenant := os.Getenv("TINES_TENANT")
	apiKey := os.Getenv("TINES_API_KEY")

	if !config.Tenant.IsNull() {
		tenant = config.Tenant.ValueString()
	}

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	if tenant == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("tenant"),
			"Missing Tines Tenant",
			"The provider cannot create the Tines API client as the Tines tenant URL is missing",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Tines API Key",
			"The provider cannot create the Tines API client as the Tines API key is missing",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Tines client using the configuration values.
	c, err := tines_cli.NewClient(tenant, apiKey, p.version)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Tines API Client",
			"An unexpected error occurred when creating the Tines API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Tines Client Error: "+err.Error(),
		)
		return
	}

	// Make the Tines client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = c
	resp.ResourceData = c

	tflog.Info(ctx, "Configured Tines client", map[string]any{"success": true})
}

// Resources returns the available Tines API resources.
func (p *TinesProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewStoryResource,
	}
}

// DataSources returns the available data sources (currently none).
func (p *TinesProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TinesProvider{
			version: version,
		}
	}
}
