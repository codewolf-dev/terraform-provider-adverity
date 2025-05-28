// Copyright (c) codewolf.dev
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"

	"terraform-provider-adverity/internal/adverity"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &datastreamTypeDataSource{}
	_ datasource.DataSourceWithConfigure = &datastreamTypeDataSource{}
)

// NewDatastreamTypeDataSource is a helper function to simplify the provider implementation.
func NewDatastreamTypeDataSource() datasource.DataSource {
	return &datastreamTypeDataSource{}
}

// datastreamTypeDataSource is the data source implementation.
type datastreamTypeDataSource struct {
	client *adverity.Client
}

// datastreamTypeDataSourceModel maps the data source schema data.
type datastreamTypeDataSourceModel struct {
	SearchTerm types.String `tfsdk:"search_term"`
	Results    types.Map    `tfsdk:"results"`
}

// Configure adds the provider configured client to the data source.
func (d *datastreamTypeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*adverity.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *adverity.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *datastreamTypeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_datastream_type"
}

// Schema defines the schema for the data source.
func (d *datastreamTypeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of datastream types.",
		Attributes: map[string]schema.Attribute{
			"search_term": schema.StringAttribute{
				Description: "Search term to filter on.",
				Required:    true,
			},
			"results": schema.MapAttribute{
				Description: "Results containing slug to id mapping.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *datastreamTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Retrieve values from config
	var data datastreamTypeDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	datastreamTypes, err := d.client.QueryDatastreamTypes(data.SearchTerm.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying Adverity datastream types",
			"Could not query datastream types, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to model
	results := make(map[string]attr.Value, len(datastreamTypes))
	for _, datastreamType := range datastreamTypes {
		results[datastreamType.Slug] = types.Int64Value(datastreamType.ID)
	}

	// Convert results to types.Map
	data.Results, diags = types.MapValue(types.Int64Type, results)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
