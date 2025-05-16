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
	_ resource.Resource                = &workspaceResource{}
	_ resource.ResourceWithConfigure   = &workspaceResource{}
	_ resource.ResourceWithImportState = &workspaceResource{}
)

// NewWorkspaceResource is a helper function to simplify the provider implementation.
func NewWorkspaceResource() resource.Resource {
	return &workspaceResource{}
}

// workspaceResource is the resource implementation.
type workspaceResource struct {
	client *adverity.Client
}

// workspaceResourceModel maps the resource schema data.
type workspaceResourceModel struct {
	ID          types.Int64   `tfsdk:"id"`
	Name        types.String  `tfsdk:"name"`
	Slug        types.String  `tfsdk:"slug"`
	DatalakeID  types.Int64   `tfsdk:"datalake_id"`
	ParentID    types.Int64   `tfsdk:"parent_id"`
	Parameters  types.Dynamic `tfsdk:"parameters"`
	LastUpdated types.String  `tfsdk:"last_updated"`
}

// Configure adds the provider configured client to the resource.
func (r *workspaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *workspaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

// Schema defines the schema for the resource.
func (r *workspaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a workspace.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the workspace.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the workspace.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the workspace.",
				Required:    true,
			},
			"slug": schema.StringAttribute{
				Description: "Slug of the workspace.",
				Computed:    true,
			},
			"datalake_id": schema.Int64Attribute{
				Description: "Numeric identifier of the datalake.",
				Optional:    true,
			},
			"parent_id": schema.Int64Attribute{
				Description: "Numeric identifier of the parent workspace.",
				Optional:    true,
			},
			"parameters": schema.DynamicAttribute{
				Description: "Additional workspace parameters.",
				Optional:    true,
			},
		},
	}
}

// Create a new resource.
func (r *workspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan workspaceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	payload := &adverity.WorkspaceConfig{
		Name:       plan.Name.ValueStringPointer(),
		DatalakeID: plan.DatalakeID.ValueInt64Pointer(),
		ParentID:   plan.ParentID.ValueInt64Pointer(),
	}

	if !plan.Parameters.IsNull() {
		parameters := utils.ExpandParameters(plan.Parameters.UnderlyingValue(), path.Root("parameters"), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		payload.Parameters = &parameters
	}

	// Create new workspace
	workspace, err := r.client.CreateWorkspace(payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating workspace",
			"Could not create workspace, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate computed attribute values
	plan.ID = types.Int64Value(workspace.ID)
	plan.Slug = types.StringValue(workspace.Slug)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *workspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state workspaceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed workspace value from Adverity
	workspace, err := r.client.ReadWorkspace(state.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Adverity workspace",
			"Could not read workspace, unexpected error: "+err.Error(),
		)
		return
	}

	// Overwrite state with refreshed attributes
	state.Name = types.StringValue(workspace.Name)
	state.Slug = types.StringValue(workspace.Slug)
	state.ParentID = types.Int64Value(workspace.ParentID)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *workspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan workspaceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	// Generate API request body from plan
	payload := &adverity.WorkspaceConfig{
		Name:       plan.Name.ValueStringPointer(),
		DatalakeID: plan.DatalakeID.ValueInt64Pointer(),
		ParentID:   plan.ParentID.ValueInt64Pointer(),
	}

	if !plan.Parameters.IsNull() {
		parameters := utils.ExpandParameters(plan.Parameters.UnderlyingValue(), path.Root("parameters"), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		payload.Parameters = &parameters
	}

	// Retrieve slug from state
	var slug types.String
	diags = req.State.GetAttribute(ctx, path.Root("slug"), &slug)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing workspace
	workspace, err := r.client.UpdateWorkspace(slug.ValueString(), payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Adverity workspace",
			"Could not update workspace, unexpected error: "+err.Error(),
		)
		return
	}

	// Update resource state with updated attributes and timestamp
	plan.Slug = types.StringValue(workspace.Slug)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *workspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state workspaceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing workspace
	_, err := r.client.DeleteWorkspace(state.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Adverity workspace",
			"Could not delete workspace, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *workspaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
