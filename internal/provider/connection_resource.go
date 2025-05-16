// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
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
	_ resource.Resource                = &connectionResource{}
	_ resource.ResourceWithConfigure   = &connectionResource{}
	_ resource.ResourceWithImportState = &connectionResource{}
)

// NewConnectionResource is a helper function to simplify the provider implementation.
func NewConnectionResource() resource.Resource {
	return &connectionResource{}
}

// connectionResource is the resource implementation.
type connectionResource struct {
	client *adverity.Client
}

// connectionResourceModel maps the resource schema data.
type connectionResourceModel struct {
	ConnectionTypeId types.Int64   `tfsdk:"connection_type_id"`
	ID               types.Int64   `tfsdk:"id"`
	Name             types.String  `tfsdk:"name"`
	StackID          types.Int64   `tfsdk:"stack_id"`
	IsAuthorized     types.Bool    `tfsdk:"is_authorized"`
	Parameters       types.Dynamic `tfsdk:"parameters"`
	LastUpdated      types.String  `tfsdk:"last_updated"`
}

// Configure adds the provider configured client to the resource.
func (r *connectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *connectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connection"
}

// Schema defines the schema for the resource.
func (r *connectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a connection.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the connection.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the connection.",
				Computed:    true,
			},
			"connection_type_id": schema.Int64Attribute{
				Description: "Numeric identifier of the connection type.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the connection.",
				Required:    true,
			},
			"stack_id": schema.Int64Attribute{
				Description: "Numeric identifier of the workspace.",
				Optional:    true,
			},
			"is_authorized": schema.BoolAttribute{
				Description: "Whether the connection is authorized.",
				Computed:    true,
			},
			"parameters": schema.DynamicAttribute{
				Description: "Additional connection parameters.",
				Optional:    true,
			},
		},
	}
}

// Create a new resource.
func (r *connectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan connectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	payload := &adverity.ConnectionConfig{
		Name:    plan.Name.ValueStringPointer(),
		StackID: plan.StackID.ValueInt64Pointer(),
	}

	if !plan.Parameters.IsNull() {
		parameters := utils.ExpandParameters(plan.Parameters.UnderlyingValue(), path.Root("parameters"), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		payload.Parameters = &parameters
	}

	// Create new connection
	connection, err := r.client.CreateConnection(int(plan.ConnectionTypeId.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating connection",
			"Could not create connection, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate computed attribute values
	plan.ID = types.Int64Value(connection.ID)
	plan.IsAuthorized = types.BoolValue(connection.IsAuthorized)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *connectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state connectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed connection value from Adverity
	connection, err := r.client.ReadConnection(int(state.ConnectionTypeId.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Adverity connection",
			"Could not read connection, unexpected error: "+err.Error(),
		)
		return
	}

	// Overwrite state with refreshed attributes
	state.Name = types.StringValue(connection.Name)
	state.StackID = types.Int64Value(connection.StackID)
	state.IsAuthorized = types.BoolValue(connection.IsAuthorized)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *connectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan connectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	// Generate API request body from plan
	payload := &adverity.ConnectionConfig{
		Name:    plan.Name.ValueStringPointer(),
		StackID: plan.StackID.ValueInt64Pointer(),
	}

	if !plan.Parameters.IsNull() {
		parameters := utils.ExpandParameters(plan.Parameters.UnderlyingValue(), path.Root("parameters"), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		payload.Parameters = &parameters
	}

	// Update existing connection
	connection, err := r.client.UpdateConnection(int(plan.ConnectionTypeId.ValueInt64()), int(plan.ID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Adverity connection",
			"Could not update connection, unexpected error: "+err.Error(),
		)
		return
	}

	// Update resource state with updated attributes and timestamp
	plan.IsAuthorized = types.BoolValue(connection.IsAuthorized)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *connectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state connectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing connection
	_, err := r.client.DeleteConnection(int(state.ConnectionTypeId.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Adverity connection",
			"Could not delete connection, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *connectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
