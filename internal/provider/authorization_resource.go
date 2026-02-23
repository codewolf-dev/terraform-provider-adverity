// Copyright codewolf.dev 2025, 0
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
	_ resource.Resource                = &authorizationResource{}
	_ resource.ResourceWithConfigure   = &authorizationResource{}
	_ resource.ResourceWithImportState = &authorizationResource{}
)

// NewAuthorizationResource is a helper function to simplify the provider implementation.
func NewAuthorizationResource() resource.Resource {
	return &authorizationResource{}
}

// authorizationResource is the resource implementation.
type authorizationResource struct {
	client *adverity.Client
}

// authorizationResourceModel maps the resource schema data.
type authorizationResourceModel struct {
	AuthorizationTypeId types.Int64   `tfsdk:"authorization_type_id"`
	ID                  types.Int64   `tfsdk:"id"`
	Name                types.String  `tfsdk:"name"`
	StackID             types.Int64   `tfsdk:"stack_id"`
	IsAuthorized        types.Bool    `tfsdk:"is_authorized"`
	Parameters          types.Dynamic `tfsdk:"parameters"`
	LastUpdated         types.String  `tfsdk:"last_updated"`
}

func (r *authorizationResource) refreshState(authorization *adverity.AuthorizationResponse, state *authorizationResourceModel) {
	state.ID = types.Int64Value(authorization.ID)
	state.Name = types.StringValue(authorization.Name)
	state.StackID = types.Int64Value(authorization.StackID)
	state.IsAuthorized = types.BoolValue(authorization.IsAuthorized)
}

// Configure adds the provider configured client to the resource.
func (r *authorizationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *authorizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorization"
}

// Schema defines the schema for the resource.
func (r *authorizationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an authorization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the authorization.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the authorization.",
				Computed:    true,
			},
			"authorization_type_id": schema.Int64Attribute{
				Description: "Numeric identifier of the authorization type.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the authorization.",
				Required:    true,
			},
			"stack_id": schema.Int64Attribute{
				Description: "Numeric identifier of the workspace.",
				Optional:    true,
			},
			"is_authorized": schema.BoolAttribute{
				Description: "Whether the authorization is authorized.",
				Computed:    true,
			},
			"parameters": schema.DynamicAttribute{
				Description: "Additional authorization parameters.",
				Optional:    true,
			},
		},
	}
}

// Create a new resource.
func (r *authorizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan authorizationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	payload := &adverity.AuthorizationConfig{
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

	// Create new authorization
	authorization, err := r.client.CreateAuthorization(int(plan.AuthorizationTypeId.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating authorization",
			"Could not create authorization, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate computed attribute values
	r.refreshState(authorization, &plan)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *authorizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state authorizationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed authorization value from Adverity
	authorization, err := r.client.ReadAuthorization(int(state.AuthorizationTypeId.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Adverity authorization",
			"Could not read authorization, unexpected error: "+err.Error(),
		)
		return
	}

	// Overwrite state with refreshed attributes
	r.refreshState(authorization, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *authorizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan authorizationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	// Generate API request body from plan
	payload := &adverity.AuthorizationConfig{
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

	// Update existing authorization
	authorization, err := r.client.UpdateAuthorization(int(plan.AuthorizationTypeId.ValueInt64()), int(plan.ID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Adverity authorization",
			"Could not update authorization, unexpected error: "+err.Error(),
		)
		return
	}

	// Update resource state with updated attributes and timestamp
	r.refreshState(authorization, &plan)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *authorizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state authorizationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing authorization
	_, err := r.client.DeleteAuthorization(int(state.AuthorizationTypeId.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Adverity authorization",
			"Could not delete authorization, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *authorizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
