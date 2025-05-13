// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"os"
	"terraform-provider-adverity/internal/adverity"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure AdverityProvider satisfies various provider interfaces.
var _ provider.Provider = &AdverityProvider{}

// AdverityProvider defines the provider implementation.
type AdverityProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// AdverityProviderModel describes the provider data model.
type AdverityProviderModel struct {
	InstanceUrl types.String `tfsdk:"instance_url"`
	AuthToken   types.String `tfsdk:"auth_token"`
}

func (p *AdverityProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "adverity"
	resp.Version = p.version
}

func (p *AdverityProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Adverity.",
		Attributes: map[string]schema.Attribute{
			"instance_url": schema.StringAttribute{
				Description: "Instance URL pointing to Adverity API (e.g. <placeholder>.datatap.adverity.com). May also be provided via ADVERITY_INSTANCE_URL environment variable.",
				Optional:    true,
			},
			"auth_token": schema.StringAttribute{
				Description: "Authentication token for Adverity API. May also be provided via ADVERITY_AUTH_TOKEN environment variable.",
				Optional:    true,
			},
		},
	}
}

func (p *AdverityProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Adverity API client")

	// Retrieve provider data from configuration

	var config AdverityProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.InstanceUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("instance_url"),
			"Unknown Adverity instance URL",
			"The provider cannot create the Adverity API client as there is an unknown configuration value for the Adverity instance URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ADVERITY_INSTANCE_URL environment variable.",
		)
	}

	if config.AuthToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("auth_token"),
			"Unknown Adverity auth token",
			"The provider cannot create the Adverity API client as there is an unknown configuration value for the Adverity auth token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ADVERITY_AUTH_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	instanceUrl := os.Getenv("ADVERITY_INSTANCE_URL")
	authToken := os.Getenv("ADVERITY_AUTH_TOKEN")

	if !config.InstanceUrl.IsNull() {
		instanceUrl = config.InstanceUrl.ValueString()
	}

	if !config.AuthToken.IsNull() {
		authToken = config.AuthToken.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if instanceUrl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("instance_url"),
			"Missing Adverity instance URL",
			"The provider cannot create the Adverity API client as there is a missing or empty value for the Adverity instance URL. "+
				"Set the host value in the configuration or use the ADVERITY_INSTANCE_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if authToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("auh_token"),
			"Missing Adverity auth token",
			"The provider cannot create the Adverity API client as there is a missing or empty value for the Adverity auth token. "+
				"Set the username value in the configuration or use the ADVERITY_AUTH_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "adverity_instance_url", instanceUrl)
	ctx = tflog.SetField(ctx, "adverity_auth_token", authToken)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "adverity_auth_token")

	tflog.Debug(ctx, "Creating Adverity API client")

	// Create a new Adverity client using the configuration values
	client, err := adverity.NewClient(&instanceUrl, &authToken)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Adverity API client",
			"An unexpected error occurred when creating the Adverity API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Adverity client error: "+err.Error(),
		)
		return
	}

	// Make the Adverity client available during DataSource and Resource type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Adverity API client", map[string]any{"success": true})
}

// Resources defines the resources implemented in the provider.
func (p *AdverityProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
		NewWorkspaceResource,
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *AdverityProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AdverityProvider{
			version: version,
		}
	}
}
