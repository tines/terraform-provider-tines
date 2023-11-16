package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
type tinesProviderModel struct{}

func (p *TinesProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tines"
	resp.Version = p.version
}

// Schema describes the schema for Tines provider.
func (p *TinesProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

// Configure initializes the Tines client.
func (p *TinesProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Tines client")

	var data tinesProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Tines client using the configuration values.
	client, err := NewClient()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Tines API Client",
			"An unexpected error occurred when creating the Tines API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Tines Client Error: "+err.Error(),
		)
		return
	}
	client.version = p.version

	// Make the Tines client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Tines client", map[string]any{"success": true})
}

// Resources returns the available Story Resource.
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
