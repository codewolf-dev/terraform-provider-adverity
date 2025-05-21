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
	_ resource.Resource                = &destinationResource{}
	_ resource.ResourceWithConfigure   = &destinationResource{}
	_ resource.ResourceWithImportState = &destinationResource{}
)

// NewDestinationResource is a helper function to simplify the provider implementation.
func NewDestinationResource() resource.Resource {
	return &destinationResource{}
}

// destinationResource is the resource implementation.
type destinationResource struct {
	client *adverity.Client
}

// destinationResourceModel maps the resource schema data.
type destinationResourceModel struct {
	DestinationTypeId types.Int64  `tfsdk:"destination_type_id"`
	ID                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	StackID           types.Int64  `tfsdk:"stack_id"`
	AuthID            types.Int64  `tfsdk:"auth_id"`
	//SchemaMapping          types.Bool    `tfsdk:"schema_mapping"`
	//ColumnNamesToLowerCase types.Bool    `tfsdk:"column_names_to_lowercase"`
	//ForceString            types.Bool    `tfsdk:"force_string"`
	//FormatHeaders          types.Bool    `tfsdk:"format_headers"`
	//HeadersFormatting      types.Int64   `tfsdk:"headers_formatting"`
	Parameters  types.Dynamic `tfsdk:"parameters"`
	LastUpdated types.String  `tfsdk:"last_updated"`
}

func (r *destinationResource) refreshState(destination *adverity.DestinationResponse, state *destinationResourceModel) {
	state.ID = types.Int64Value(destination.ID)
	state.Name = types.StringValue(destination.Name)
	state.StackID = types.Int64Value(destination.StackID)
	state.AuthID = types.Int64Value(destination.AuthID)
	//state.SchemaMapping = types.BoolValue(destination.SchemaMapping)
	//state.ColumnNamesToLowerCase = types.BoolValue(destination.ColumnNamesToLowerCase)
	//state.ForceString = types.BoolValue(destination.ForceString)
	//state.FormatHeaders = types.BoolValue(destination.FormatHeaders)
	//state.HeadersFormatting = types.Int64Value(destination.HeadersFormatting)
}

// Configure adds the provider configured client to the resource.
func (r *destinationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *destinationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination"
}

// Schema defines the schema for the resource.
func (r *destinationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a destination.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the destination.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the destination.",
				Computed:    true,
			},
			"destination_type_id": schema.Int64Attribute{
				Description: "Numeric identifier of the destination type.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the destination.",
				Required:    true,
			},
			"stack_id": schema.Int64Attribute{
				Description: "Numeric identifier of the workspace.",
				Optional:    true,
			},
			"auth_id": schema.Int64Attribute{
				Description: "Numeric identifier of the authentication.",
				Optional:    true,
			},
			"parameters": schema.DynamicAttribute{
				Description: "Additional destination parameters.",
				Optional:    true,
			},
		},
	}
}

// Create a new resource.
func (r *destinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan destinationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	payload := &adverity.DestinationConfig{
		Name:    plan.Name.ValueStringPointer(),
		StackID: plan.StackID.ValueInt64Pointer(),
		AuthID:  plan.AuthID.ValueInt64Pointer(),
		//SchemaMapping:          plan.SchemaMapping.ValueBoolPointer(),
		//ColumnNamesToLowerCase: plan.ColumnNamesToLowerCase.ValueBoolPointer(),
		//ForceString:            plan.ForceString.ValueBoolPointer(),
		//FormatHeaders:          plan.FormatHeaders.ValueBoolPointer(),
		//HeadersFormatting:      plan.HeadersFormatting.ValueInt64Pointer(),
	}

	if !plan.Parameters.IsNull() {
		parameters := utils.ExpandParameters(plan.Parameters.UnderlyingValue(), path.Root("parameters"), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		payload.Parameters = &parameters
	}

	// Create new destination
	destination, err := r.client.CreateDestination(int(plan.DestinationTypeId.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating destination",
			"Could not create destination, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate computed attribute values
	r.refreshState(destination, &plan)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *destinationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state destinationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed destination value from Adverity
	destination, err := r.client.ReadDestination(int(state.DestinationTypeId.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Adverity destination",
			"Could not read destination, unexpected error: "+err.Error(),
		)
		return
	}

	// Overwrite state with refreshed attributes
	r.refreshState(destination, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *destinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan destinationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	// Generate API request body from plan
	payload := &adverity.DestinationConfig{
		Name:    plan.Name.ValueStringPointer(),
		StackID: plan.StackID.ValueInt64Pointer(),
		AuthID:  plan.AuthID.ValueInt64Pointer(),
		//SchemaMapping:          plan.SchemaMapping.ValueBoolPointer(),
		//ColumnNamesToLowerCase: plan.ColumnNamesToLowerCase.ValueBoolPointer(),
		//ForceString:            plan.ForceString.ValueBoolPointer(),
		//FormatHeaders:          plan.FormatHeaders.ValueBoolPointer(),
		//HeadersFormatting:      plan.HeadersFormatting.ValueInt64Pointer(),
	}

	if !plan.Parameters.IsNull() {
		parameters := utils.ExpandParameters(plan.Parameters.UnderlyingValue(), path.Root("parameters"), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		payload.Parameters = &parameters
	}

	// Update existing destination
	destination, err := r.client.UpdateDestination(int(plan.DestinationTypeId.ValueInt64()), int(plan.ID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Adverity destination",
			"Could not update destination, unexpected error: "+err.Error(),
		)
		return
	}

	// Update resource state with updated attributes and timestamp
	r.refreshState(destination, &plan)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *destinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state destinationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing destination
	_, err := r.client.DeleteDestination(int(state.DestinationTypeId.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Adverity destination",
			"Could not delete destination, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *destinationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
