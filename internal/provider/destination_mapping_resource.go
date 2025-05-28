// Copyright (c) codewolf.dev
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"time"

	"terraform-provider-adverity/internal/adverity"
	"terraform-provider-adverity/internal/provider/utils"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &destinationMappingResource{}
	_ resource.ResourceWithConfigure   = &destinationMappingResource{}
	_ resource.ResourceWithImportState = &destinationMappingResource{}
)

// NewDestinationMappingResource is a helper function to simplify the provider implementation.
func NewDestinationMappingResource() resource.Resource {
	return &destinationMappingResource{}
}

// destinationMappingResource is the resource implementation.
type destinationMappingResource struct {
	client *adverity.Client
}

// destinationMappingResourceModel maps the resource schema data.
type destinationMappingResourceModel struct {
	DestinationTypeId types.Int64   `tfsdk:"destination_type_id"`
	DestinationId     types.Int64   `tfsdk:"destination_id"`
	DatastreamId      types.Int64   `tfsdk:"datastream_id"`
	ID                types.Int64   `tfsdk:"id"`
	Enabled           types.Bool    `tfsdk:"enabled"`
	TableName         types.String  `tfsdk:"table_name"`
	Parameters        types.Dynamic `tfsdk:"parameters"`
	LastUpdated       types.String  `tfsdk:"last_updated"`
}

func (r *destinationMappingResource) refreshState(destinationMapping *adverity.DestinationMappingResponse, state *destinationMappingResourceModel) {
	state.ID = types.Int64Value(destinationMapping.ID)
	state.DatastreamId = types.Int64Value(destinationMapping.DatastreamID)
	state.Enabled = types.BoolValue(destinationMapping.Enabled)
	state.TableName = types.StringValue(destinationMapping.TableName)
}

// Configure adds the provider configured client to the resource.
func (r *destinationMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

// Metadata returns the resource type name.
func (r *destinationMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination_mapping"
}

// Schema defines the schema for the resource.
func (r *destinationMappingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a destination mapping.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the destination mapping.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the destination mapping.",
				Computed:    true,
			},
			"destination_type_id": schema.Int64Attribute{
				Description: "Numeric identifier of the destination mapping type.",
				Required:    true,
			},
			"destination_id": schema.Int64Attribute{
				Description: "Numeric identifier of the destination mapping type.",
				Required:    true,
			},
			"datastream_id": schema.Int64Attribute{
				Description: "Numeric identifier of the datastream.",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Name of the destination mapping.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"table_name": schema.StringAttribute{
				Description: "Name of the target table.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"parameters": schema.DynamicAttribute{
				Description: "Additional destination mapping parameters.",
				Optional:    true,
			},
		},
	}
}

// Create a new resource.
func (r *destinationMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan destinationMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	payload := &adverity.DestinationMappingConfig{
		DatastreamId: plan.DatastreamId.ValueInt64Pointer(),
		Enabled:      plan.Enabled.ValueBoolPointer(),
	}

	if !plan.TableName.IsUnknown() {
		payload.TableName = plan.TableName.ValueStringPointer()
	}

	if !plan.Parameters.IsNull() {
		parameters := utils.ExpandParameters(plan.Parameters.UnderlyingValue(), path.Root("parameters"), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		payload.Parameters = &parameters
	}

	// Create new destination mapping
	destinationMapping, err := r.client.CreateDestinationMapping(int(plan.DestinationTypeId.ValueInt64()), int(plan.DestinationId.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating destination mapping",
			"Could not create destinationMapping, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate computed attribute values
	r.refreshState(destinationMapping, &plan)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *destinationMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state destinationMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed destination mapping value from Adverity
	destinationMapping, err := r.client.ReadDestinationMapping(int(state.DestinationTypeId.ValueInt64()), int(state.DestinationId.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Adverity destination mapping",
			"Could not read destinationMapping, unexpected error: "+err.Error(),
		)
		return
	}

	// Overwrite state with refreshed attributes
	r.refreshState(destinationMapping, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *destinationMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan destinationMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	// Generate API request body from plan
	payload := &adverity.DestinationMappingConfig{
		DatastreamId: plan.DatastreamId.ValueInt64Pointer(),
		Enabled:      plan.Enabled.ValueBoolPointer(),
	}

	if !plan.TableName.IsUnknown() {
		payload.TableName = plan.TableName.ValueStringPointer()
	}

	if !plan.Parameters.IsNull() {
		parameters := utils.ExpandParameters(plan.Parameters.UnderlyingValue(), path.Root("parameters"), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		payload.Parameters = &parameters
	}

	// Update existing destination mapping
	destinationMapping, err := r.client.UpdateDestinationMapping(int(plan.DestinationTypeId.ValueInt64()), int(plan.DestinationId.ValueInt64()), int(plan.ID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Adverity destination mapping",
			"Could not update destinationMapping, unexpected error: "+err.Error(),
		)
		return
	}

	// Update resource state with updated attributes and timestamp
	r.refreshState(destinationMapping, &plan)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *destinationMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state destinationMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing destination mapping
	_, err := r.client.DeleteDestinationMapping(int(state.DestinationTypeId.ValueInt64()), int(state.DestinationId.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Adverity destination mapping",
			"Could not delete destinationMapping, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *destinationMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
