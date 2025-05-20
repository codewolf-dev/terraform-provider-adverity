// Copyright (c) HashiCorp, Inc.
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
	_ resource.Resource                = &datastreamResource{}
	_ resource.ResourceWithConfigure   = &datastreamResource{}
	_ resource.ResourceWithImportState = &datastreamResource{}
)

// NewDatastreamResource is a helper function to simplify the provider implementation.
func NewDatastreamResource() resource.Resource {
	return &datastreamResource{}
}

// datastreamResource is the resource implementation.
type datastreamResource struct {
	client *adverity.Client
}

type datastreamScheduleModel struct {
	CronPreset         types.String `tfsdk:"cron_preset"`
	CronType           types.String `tfsdk:"cron_type"`
	CronInterval       types.Int64  `tfsdk:"cron_interval"`
	CronIntervalStart  types.Int64  `tfsdk:"cron_interval_start"`
	CronStartOfDay     types.String `tfsdk:"cron_start_of_day"`
	TimeRangePreset    types.Int64  `tfsdk:"time_range_preset"`
	DeltaType          types.Int64  `tfsdk:"delta_type"`
	DeltaInterval      types.Int64  `tfsdk:"delta_interval"`
	DeltaIntervalStart types.Int64  `tfsdk:"delta_interval_start"`
	DeltaStartOfDay    types.String `tfsdk:"delta_start_of_day"`
	FixedStart         types.String `tfsdk:"fixed_start"`
	FixedEnd           types.String `tfsdk:"fixed_end"`
	OffsetDays         types.Int64  `tfsdk:"offset_days"`
	NotBeforeDate      types.String `tfsdk:"not_before_date"`
	NotBeforeTime      types.String `tfsdk:"not_before_time"`
}

// datastreamResourceModel maps the resource schema data.
type datastreamResourceModel struct {
	DatastreamTypeId    types.Int64               `tfsdk:"datastream_type_id"`
	ID                  types.Int64               `tfsdk:"id"`
	Name                types.String              `tfsdk:"name"`
	Description         types.String              `tfsdk:"description"`
	StackID             types.Int64               `tfsdk:"stack_id"`
	AuthID              types.Int64               `tfsdk:"auth_id"`
	Schedules           []datastreamScheduleModel `tfsdk:"schedule"`
	Enabled             types.Bool                `tfsdk:"enabled"`
	DataType            types.String              `tfsdk:"datatype"`
	RetentionType       types.Int64               `tfsdk:"retention_type"`
	RetentionNumber     types.Int64               `tfsdk:"retention_number"`
	ManageExtractNames  types.Bool                `tfsdk:"manage_extract_names"`
	ExtractNameKeys     types.String              `tfsdk:"extract_name_keys"`
	IsInsightsMediaplan types.Bool                `tfsdk:"is_insights_mediaplan"`
	Parameters          types.Dynamic             `tfsdk:"parameters"`
	LastUpdated         types.String              `tfsdk:"last_updated"`
}

func (r *datastreamResource) refreshState(datastream *adverity.DatastreamResponse, state *datastreamResourceModel) {
	state.ID = types.Int64Value(datastream.ID)
	state.Name = types.StringValue(datastream.Name)
	state.Description = types.StringValue(datastream.Description)
	state.StackID = types.Int64Value(datastream.StackID)
	state.AuthID = types.Int64Value(datastream.AuthID)
	state.DataType = types.StringValue(datastream.DataType)
	state.RetentionType = types.Int64Value(datastream.RetentionType)
	state.RetentionNumber = types.Int64Value(datastream.RetentionNumber)
	state.IsInsightsMediaplan = types.BoolValue(datastream.IsInsightsMediaplan)
	state.ManageExtractNames = types.BoolValue(datastream.ManageExtractNames)
	state.ExtractNameKeys = types.StringValue(datastream.ExtractNameKeys)
	state.Enabled = types.BoolValue(datastream.Enabled)

	var schedules []datastreamScheduleModel
	for _, schedule := range datastream.Schedules {
		schedules = append(schedules, datastreamScheduleModel{
			CronPreset:         types.StringValue(schedule.CronPreset),
			CronType:           types.StringValue(schedule.CronType),
			CronInterval:       types.Int64Value(schedule.CronInterval),
			CronIntervalStart:  types.Int64Value(schedule.CronIntervalStart),
			CronStartOfDay:     types.StringValue(schedule.CronStartOfDay),
			TimeRangePreset:    types.Int64Value(schedule.TimeRangePreset),
			DeltaType:          types.Int64Value(schedule.DeltaType),
			DeltaInterval:      types.Int64Value(schedule.DeltaInterval),
			DeltaIntervalStart: types.Int64Value(schedule.DeltaIntervalStart),
			DeltaStartOfDay:    types.StringValue(schedule.DeltaStartOfDay),
			FixedStart:         types.StringValue(schedule.FixedStart),
			FixedEnd:           types.StringValue(schedule.FixedEnd),
			OffsetDays:         types.Int64Value(schedule.OffsetDays),
			NotBeforeDate:      types.StringValue(schedule.NotBeforeDate),
			NotBeforeTime:      types.StringValue(schedule.NotBeforeTime),
		})
	}
	state.Schedules = schedules
}

// Configure adds the provider configured client to the resource.
func (r *datastreamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *datastreamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_datastream"
}

// Schema defines the schema for the resource.
func (r *datastreamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a datastream.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the datastream.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the datastream.",
				Computed:    true,
			},
			"datastream_type_id": schema.Int64Attribute{
				Description: "Numeric identifier of the datastream type.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the datastream.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the datastream.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"stack_id": schema.Int64Attribute{
				Description: "Numeric identifier of the workspace.",
				Optional:    true,
			},
			"auth_id": schema.Int64Attribute{
				Description: "Numeric identifier of the connection.",
				Optional:    true,
			},
			"datatype": schema.StringAttribute{
				Description: "Type of the datastream ('Live' or 'Staging').",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether to enable the datastream.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"retention_type": schema.Int64Attribute{
				Description: "Numeric identifier of the retention type.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"retention_number": schema.Int64Attribute{
				Description: "Number of fetches/extracts/days to retain.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"manage_extract_names": schema.BoolAttribute{
				Description: "Whether to manage extract names.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"extract_name_keys": schema.StringAttribute{
				Description: "Date column to use for managing extract names.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_insights_mediaplan": schema.BoolAttribute{
				Description: "Whether to treat extracts as insights mediaplans.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"parameters": schema.DynamicAttribute{
				Description: "Additional datastream parameters.",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"schedule": schema.ListNestedBlock{
				Description: "Schedule the datastream.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"cron_preset": schema.StringAttribute{
							Description: "Cron preset.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"cron_type": schema.StringAttribute{
							Description: "Cron type.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"cron_interval": schema.Int64Attribute{
							Description: "Cron interval.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"cron_interval_start": schema.Int64Attribute{
							Description: "Cron interval start.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"cron_start_of_day": schema.StringAttribute{
							Description: "Cron start of day.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"time_range_preset": schema.Int64Attribute{
							Description: "Time range preset.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"delta_type": schema.Int64Attribute{
							Description: "Delta type.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"delta_interval": schema.Int64Attribute{
							Description: "Delta interval.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"delta_interval_start": schema.Int64Attribute{
							Description: "Delta interval start.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"delta_start_of_day": schema.StringAttribute{
							Description: "Delta start of day.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"fixed_start": schema.StringAttribute{
							Description: "Fixed start.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"fixed_end": schema.StringAttribute{
							Description: "Fixed end.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"offset_days": schema.Int64Attribute{
							Description: "Offset days.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"not_before_date": schema.StringAttribute{
							Description: "Not before date.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"not_before_time": schema.StringAttribute{
							Description: "Not before time.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
		},
	}
}

// Create a new resource.
func (r *datastreamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan datastreamResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	payload := &adverity.DatastreamCreateConfig{
		Name:                plan.Name.ValueStringPointer(),
		Description:         plan.Description.ValueStringPointer(),
		StackID:             plan.StackID.ValueInt64Pointer(),
		AuthID:              plan.AuthID.ValueInt64Pointer(),
		DataType:            plan.DataType.ValueStringPointer(),
		IsInsightsMediaplan: plan.IsInsightsMediaplan.ValueBoolPointer(),
		ManageExtractNames:  plan.ManageExtractNames.ValueBoolPointer(),
		ExtractNameKeys:     plan.ExtractNameKeys.ValueStringPointer(),
		Enabled:             plan.Enabled.ValueBoolPointer(),
	}
	if !plan.DataType.IsUnknown() {
		payload.DataType = plan.DataType.ValueStringPointer()
	}
	if !plan.RetentionType.IsUnknown() {
		payload.RetentionType = plan.RetentionType.ValueInt64Pointer()
	}
	if !plan.RetentionNumber.IsUnknown() {
		payload.RetentionNumber = plan.RetentionNumber.ValueInt64Pointer()
	}

	if !plan.Parameters.IsNull() {
		parameters := utils.ExpandParameters(plan.Parameters.UnderlyingValue(), path.Root("parameters"), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		payload.Parameters = &parameters
	}

	var schedules []adverity.ScheduleConfig
	for _, schedule := range plan.Schedules {
		config := adverity.ScheduleConfig{}
		if !schedule.CronPreset.IsUnknown() {
			config.CronPreset = schedule.CronPreset.ValueStringPointer()
		}
		if !schedule.CronType.IsUnknown() {
			config.CronType = schedule.CronType.ValueStringPointer()
		}
		if !schedule.CronInterval.IsUnknown() {
			config.CronInterval = schedule.CronInterval.ValueInt64Pointer()
		}
		if !schedule.CronIntervalStart.IsUnknown() {
			config.CronIntervalStart = schedule.CronIntervalStart.ValueInt64Pointer()
		}
		if !schedule.CronStartOfDay.IsUnknown() {
			config.CronStartOfDay = schedule.CronStartOfDay.ValueStringPointer()
		}
		if !schedule.TimeRangePreset.IsUnknown() {
			config.TimeRangePreset = schedule.TimeRangePreset.ValueInt64Pointer()
		}
		if !schedule.DeltaType.IsUnknown() {
			config.DeltaType = schedule.DeltaType.ValueInt64Pointer()
		}
		if !schedule.DeltaInterval.IsUnknown() {
			config.DeltaInterval = schedule.DeltaInterval.ValueInt64Pointer()
		}
		if !schedule.DeltaIntervalStart.IsUnknown() {
			config.DeltaIntervalStart = schedule.DeltaIntervalStart.ValueInt64Pointer()
		}
		if !schedule.DeltaStartOfDay.IsUnknown() {
			config.DeltaStartOfDay = schedule.DeltaStartOfDay.ValueStringPointer()
		}
		if !schedule.FixedStart.IsUnknown() {
			config.FixedStart = schedule.FixedStart.ValueStringPointer()
		}
		if !schedule.FixedEnd.IsUnknown() {
			config.FixedEnd = schedule.FixedEnd.ValueStringPointer()
		}
		if !schedule.OffsetDays.IsUnknown() {
			config.OffsetDays = schedule.OffsetDays.ValueInt64Pointer()
		}
		if !schedule.NotBeforeDate.IsUnknown() {
			config.NotBeforeDate = schedule.NotBeforeDate.ValueStringPointer()
		}
		if !schedule.NotBeforeTime.IsUnknown() {
			config.NotBeforeTime = schedule.NotBeforeTime.ValueStringPointer()
		}
		schedules = append(schedules, config)
	}
	payload.Schedules = &schedules

	// Create new datastream
	datastream, err := r.client.CreateDatastream(int(plan.DatastreamTypeId.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating datastream",
			"Could not create datastream, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate computed attribute values
	r.refreshState(datastream, &plan)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *datastreamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state datastreamResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed datastream value from Adverity
	datastream, err := r.client.ReadDatastream(int(state.DatastreamTypeId.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Adverity datastream",
			"Could not read datastream, unexpected error: "+err.Error(),
		)
		return
	}

	// Overwrite state with refreshed attributes
	r.refreshState(datastream, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *datastreamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan datastreamResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	// Generate API request body from plan
	payload := &adverity.DatastreamUpdateConfig{
		Name:                plan.Name.ValueStringPointer(),
		Description:         plan.Description.ValueStringPointer(),
		StackID:             plan.StackID.ValueInt64Pointer(),
		AuthID:              plan.AuthID.ValueInt64Pointer(),
		DataType:            plan.DataType.ValueStringPointer(),
		RetentionType:       plan.RetentionType.ValueInt64Pointer(),
		RetentionNumber:     plan.RetentionNumber.ValueInt64Pointer(),
		IsInsightsMediaplan: plan.IsInsightsMediaplan.ValueBoolPointer(),
		ManageExtractNames:  plan.ManageExtractNames.ValueBoolPointer(),
		ExtractNameKeys:     plan.ExtractNameKeys.ValueStringPointer(),
	}
	if !plan.DataType.IsUnknown() {
		payload.DataType = plan.DataType.ValueStringPointer()
	}
	if !plan.RetentionType.IsUnknown() {
		payload.RetentionType = plan.RetentionType.ValueInt64Pointer()
	}
	if !plan.RetentionNumber.IsUnknown() {
		payload.RetentionNumber = plan.RetentionNumber.ValueInt64Pointer()
	}

	if !plan.Parameters.IsNull() {
		parameters := utils.ExpandParameters(plan.Parameters.UnderlyingValue(), path.Root("parameters"), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		payload.Parameters = &parameters
	}

	schedulePayload := &adverity.DatastreamScheduleConfig{
		Enabled: plan.Enabled.ValueBoolPointer(),
	}

	var schedules []adverity.ScheduleConfig
	for _, schedule := range plan.Schedules {
		config := adverity.ScheduleConfig{}
		if !schedule.CronPreset.IsUnknown() {
			config.CronPreset = schedule.CronPreset.ValueStringPointer()
		}
		if !schedule.CronType.IsUnknown() {
			config.CronType = schedule.CronType.ValueStringPointer()
		}
		if !schedule.CronInterval.IsUnknown() {
			config.CronInterval = schedule.CronInterval.ValueInt64Pointer()
		}
		if !schedule.CronIntervalStart.IsUnknown() {
			config.CronIntervalStart = schedule.CronIntervalStart.ValueInt64Pointer()
		}
		if !schedule.CronStartOfDay.IsUnknown() {
			config.CronStartOfDay = schedule.CronStartOfDay.ValueStringPointer()
		}
		if !schedule.TimeRangePreset.IsUnknown() {
			config.TimeRangePreset = schedule.TimeRangePreset.ValueInt64Pointer()
		}
		if !schedule.DeltaType.IsUnknown() {
			config.DeltaType = schedule.DeltaType.ValueInt64Pointer()
		}
		if !schedule.DeltaInterval.IsUnknown() {
			config.DeltaInterval = schedule.DeltaInterval.ValueInt64Pointer()
		}
		if !schedule.DeltaIntervalStart.IsUnknown() {
			config.DeltaIntervalStart = schedule.DeltaIntervalStart.ValueInt64Pointer()
		}
		if !schedule.DeltaStartOfDay.IsUnknown() {
			config.DeltaStartOfDay = schedule.DeltaStartOfDay.ValueStringPointer()
		}
		if !schedule.FixedStart.IsUnknown() {
			config.FixedStart = schedule.FixedStart.ValueStringPointer()
		}
		if !schedule.FixedEnd.IsUnknown() {
			config.FixedEnd = schedule.FixedEnd.ValueStringPointer()
		}
		if !schedule.OffsetDays.IsUnknown() {
			config.OffsetDays = schedule.OffsetDays.ValueInt64Pointer()
		}
		if !schedule.NotBeforeDate.IsUnknown() {
			config.NotBeforeDate = schedule.NotBeforeDate.ValueStringPointer()
		}
		if !schedule.NotBeforeTime.IsUnknown() {
			config.NotBeforeTime = schedule.NotBeforeTime.ValueStringPointer()
		}
		schedules = append(schedules, config)
	}
	schedulePayload.Schedules = &schedules

	// Update existing datastream
	datastream, err := r.client.UpdateDatastream(int(plan.DatastreamTypeId.ValueInt64()), int(plan.ID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Adverity datastream",
			"Could not update datastream, unexpected error: "+err.Error(),
		)
		return
	}

	// Update existing datastream schedule
	datastream, err = r.client.UpdateDatastreamSchedule(int(plan.ID.ValueInt64()), schedulePayload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Adverity datastream schedule",
			"Could not update datastream schedule, unexpected error: "+err.Error(),
		)
		return
	}

	// Update resource state with updated attributes and timestamp
	r.refreshState(datastream, &plan)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *datastreamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state datastreamResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing datastream
	_, err := r.client.DeleteDatastream(int(state.DatastreamTypeId.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Adverity datastream",
			"Could not delete datastream, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *datastreamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
