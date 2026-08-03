package provider

import (
	"context"
	"fmt"

	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	tfRunTask "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/run_task"
	cmStringValidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &RunTaskResource{}

func NewRunTaskResource() resource.Resource {
	return &RunTaskResource{}
}

type RunTaskResource struct {
	client *ControlMonkeyAPIClient
}

func (r *RunTaskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_run_task"
}

func (r *RunTaskResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys run tasks. For more information: [ControlMonkey Documentation](https://docs.controlmonkey.io/main-concepts/stack/run-tasks)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the run task.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the run task.",
				Required:            true,
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
				},
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL of the run task endpoint.",
				Required:            true,
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
				},
			},
			"is_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the run task is enabled.",
				Required:            true,
			},
			"hmac_key": schema.StringAttribute{
				MarkdownDescription: "The HMAC key for the run task. The API never returns this field, so it is stored in Terraform" +
					" state in plain text. Prefer `hmac_key_wo`, which is never written to state.",
				Optional:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("hmac_key_wo"), path.MatchRoot("hmac_key_wo_version")),
				},
			},
			"hmac_key_wo": schema.StringAttribute{
				MarkdownDescription: "The HMAC key for the run task, as a write-only argument. Unlike `hmac_key` the value is never" +
					" persisted to Terraform state or plan files. Available on Terraform / OpenTofu 1.11 and later." +
					" Because a write-only value is not stored, the provider cannot detect a change to it - increment" +
					" `hmac_key_wo_version` to send a new key.",
				Optional:  true,
				Sensitive: true,
				WriteOnly: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("hmac_key")),
					stringvalidator.AlsoRequires(path.MatchRoot("hmac_key_wo_version")),
				},
			},
			"hmac_key_wo_version": schema.Int64Attribute{
				MarkdownDescription: "Version of `hmac_key_wo`. Increment this value to make the provider send the current" +
					" `hmac_key_wo` to ControlMonkey. Changing `hmac_key_wo` on its own produces no change, because a" +
					" write-only value is not stored and therefore cannot be compared.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.MatchRoot("hmac_key")),
					int64validator.AlsoRequires(path.MatchRoot("hmac_key_wo")),
				},
			},
			"is_hmac_key_configured": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether an HMAC key has been configured for this run task.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *RunTaskResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ControlMonkeyAPIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ControlMonkeyAPIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *RunTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tfRunTask.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.runTask.ReadRunTask(ctx, id)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read run task '%s'", id), err.Error())
		return
	}

	tfRunTask.UpdateStateAfterRead(res, &state)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *RunTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tfRunTask.ResourceModel
	var config tfRunTask.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	// A write-only attribute is stripped from the plan, so hmac_key_wo is only readable from the config.
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, _ := tfRunTask.Converter(&plan, nil, commons.CreateConverter)
	if config.HmacKeyWo.IsNull() == false {
		body.SetHmacKey(config.HmacKeyWo.ValueStringPointer())
	}

	res, err := r.client.Client.runTask.CreateRunTask(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceCreationFailedError,
			fmt.Sprintf("failed to create run task, error: %s", err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(controlmonkey.StringValue(res.ID))
	plan.IsHmacKeyConfigured = helpers.BoolValueOrNull(res.IsHmacKeyConfigured)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *RunTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan tfRunTask.ResourceModel
	var state tfRunTask.ResourceModel
	var config tfRunTask.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	// A write-only attribute is stripped from the plan, so hmac_key_wo is only readable from the config.
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	body, _ := tfRunTask.Converter(&plan, &state, commons.UpdateConverter)
	if hmacKey := determineHmacKeyForUpdate(&plan, &state, &config); hmacKey != nil {
		body.SetHmacKey(hmacKey)
	}

	res, err := r.client.Client.runTask.UpdateRunTask(ctx, id, body)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.Diagnostics.AddError(resourceNotFoundError, fmt.Sprintf("Run task '%s' not found", id))
			return
		}

		resp.Diagnostics.AddError(
			resourceUpdateFailedError,
			fmt.Sprintf("failed to update run task %s, error: %s", id, err),
		)
		return
	}

	plan.IsHmacKeyConfigured = helpers.BoolValueOrNull(res.IsHmacKeyConfigured)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// determineHmacKeyForUpdate returns the write-only HMAC key to send, or nil to leave it alone.
// hmac_key_wo_version is the only trigger, since a write-only value never reaches state.
func determineHmacKeyForUpdate(plan, state, config *tfRunTask.ResourceModel) *string {
	if config.HmacKeyWo.IsNull() {
		return nil
	}

	if state.HmacKeyWoVersion.IsNull() || plan.HmacKeyWoVersion.Equal(state.HmacKeyWoVersion) == false {
		return config.HmacKeyWo.ValueStringPointer()
	}

	return nil
}

func (r *RunTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tfRunTask.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	_, err := r.client.Client.runTask.DeleteRunTask(ctx, id)
	if err != nil {
		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Run task deletion failed",
			fmt.Sprintf("Failed to delete run task %s, error: %s", id, err),
		)
	}
}

func (r *RunTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
