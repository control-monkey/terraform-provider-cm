package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	cmTypes "github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	sdkNotification "github.com/control-monkey/controlmonkey-sdk-go/services/notification"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/interfaces"
	tfEventsSubscriptions "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/events_subscriptions"
	cmStringValidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &EventsSubscriptionsResource{}

func NewEventsSubscriptionsResource() resource.Resource {
	return &EventsSubscriptionsResource{}
}

type EventsSubscriptionsResource struct {
	client *ControlMonkeyAPIClient
}

func (r *EventsSubscriptionsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_events_subscriptions"
}

func (r *EventsSubscriptionsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys events subscriptions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of this resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("Scope of the resource. Allowed values: %s.", helpers.EnumForDocs(cmTypes.EventSubscriptionScopeTypes)),
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(cmTypes.EventSubscriptionScopeTypes...),
				},
			},
			"scope_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the resource to which the subscriptions are attached.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					cmStringValidators.NotBlank(),
				},
			},
			"subscriptions": schema.SetNestedAttribute{
				MarkdownDescription: "Specifies a list of events subscriptions.",
				Optional:            true,
				Computed:            true,
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.ObjectType{
							AttrTypes: map[string]attr.Type{},
						},
						[]attr.Value{
							types.ObjectValueMust(
								map[string]attr.Type{}, map[string]attr.Value{}),
						},
					),
				),
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique ID of the subscription.",
							Computed:            true,
						},
						"event_type": schema.StringAttribute{
							MarkdownDescription: "The type of the event. Find supported types [here](https://docs.controlmonkey.io/controlmonkey-api/api-enumerations#event-types)",
							Required:            true,
						},
						"notification_endpoint_id": schema.StringAttribute{
							MarkdownDescription: "The unique ID of the endpoint to which the notification will be sent.",
							Required:            true,
							Validators: []validator.String{
								cmStringValidators.NotBlank(),
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *EventsSubscriptionsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EventsSubscriptionsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data tfEventsSubscriptions.ResourceModel

	if diags := req.Config.Get(ctx, &data); diags.HasError() {
		return
	}

	if helpers.IsKnown(data.Scope) {
		if data.Scope.ValueString() != cmTypes.OrganizationScope {
			if data.ScopeId.IsNull() {
				resp.Diagnostics.AddError(validationError, fmt.Sprintf("scope_id is required for scope '%s'", data.Scope.ValueString()))
			}
		} else {
			if helpers.IsKnown(data.ScopeId) {
				resp.Diagnostics.AddError(validationError, fmt.Sprintf("scope_id is cannot be set for scope '%s'", data.Scope.ValueString()))
			}
		}
	}

	if len(data.Subscriptions) > 0 {
		identifiers := interfaces.GetIdentifiers(data.Subscriptions)

		if helpers.IsUnique(identifiers) == false {
			duplicates := helpers.FindDuplicates(identifiers, false)
			for _, d := range duplicates {
				resp.Diagnostics.AddError(validationError, fmt.Sprintf("'%s' cannot be assigned multiple times", tfEventsSubscriptions.CleanIdentifier(d)))
			}
		}
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *EventsSubscriptionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//Get current state
	var state tfEventsSubscriptions.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID
	scope, scopeId := r.breakdownId(id)
	res, err := r.client.Client.notification.ListEventSubscriptions(ctx, scope, scopeId)

	if err != nil {
		resourceIdentifier := r.logIdentifier(scope, scopeId)

		if commons.IsNotFoundResponseError(err) {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(resourceNotFoundError, fmt.Sprintf("Resource of '%s' not found", resourceIdentifier))
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read events subscriptions for resource '%s'", resourceIdentifier), err.Error())
		return
	}

	tfEventsSubscriptions.UpdateStateAfterRead(res, &state, scope, scopeId)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *EventsSubscriptionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan tfEventsSubscriptions.ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeResult := tfEventsSubscriptions.Merge(&plan, nil, commons.CreateMerger)
	diags, newEntities := r.createEntities(ctx, mergeResult.EntitiesToCreate, plan.Scope.ValueString(), plan.ScopeId.ValueStringPointer())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.updateIdForTfSubscriptions(plan, newEntities)
	plan.ID = r.buildId(plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EventsSubscriptionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Retrieve values from plan
	var plan tfEventsSubscriptions.ResourceModel
	var state tfEventsSubscriptions.ResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	scope := plan.Scope.ValueString()
	scopeId := plan.ScopeId.ValueStringPointer()
	mergeResult := tfEventsSubscriptions.Merge(&plan, &state, commons.UpdateMerger)

	diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, scope, scopeId)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags, newEntities := r.createEntities(ctx, mergeResult.EntitiesToCreate, scope, scopeId)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.updateIdForTfSubscriptions(plan, newEntities)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EventsSubscriptionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state tfEventsSubscriptions.ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	scope := state.Scope.ValueString()
	scopeId := state.ScopeId.ValueStringPointer()
	mergeResult := tfEventsSubscriptions.Merge(nil, &state, commons.DeleteMerger)

	diags = r.deleteEntities(ctx, mergeResult.EntitiesToDelete, scope, scopeId)
	resp.Diagnostics.Append(diags...)
}

func (r *EventsSubscriptionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

//region Private Methods

func (r *EventsSubscriptionsResource) createEntities(ctx context.Context, entitiesToCreate []*sdkNotification.EventSubscription, scope string, scopeId *string) (diag.Diagnostics, []*sdkNotification.EventSubscription) {
	var diags diag.Diagnostics
	var newEntities []*sdkNotification.EventSubscription

	resourceIdentifier := r.logIdentifier(scope, scopeId)
	tflog.Info(ctx, fmt.Sprintf("Adding %d subscription to resource %s.", len(entitiesToCreate), resourceIdentifier))

	for _, e := range entitiesToCreate {
		newEntity, err := r.client.Client.notification.CreateEventSubscription(ctx, e)
		if err == nil {
			newEntities = append(newEntities, newEntity)
		} else {
			subscriptionIdentifier := r.logSubscriptionIdentifier(*e.EventType, *e.NotificationEndpointId)
			if commons.IsAlreadyExistResponseError(err) {
				diags = diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceAlreadyExists, fmt.Sprintf("Resource already has subscription %s. Import operation is required",
						subscriptionIdentifier)),
				}
			} else if commons.IsNotFoundResponseError(err) {
				diags = diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to add subscription %s. Error: %s",
						subscriptionIdentifier, err)),
				}
			} else {
				diags = diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to add subscription %s.", subscriptionIdentifier),
						err.Error()),
				}
			}
		}

		if diags.HasError() {
			return diags, nil
		}
	}

	return diags, newEntities
}

func (r *EventsSubscriptionsResource) deleteEntities(ctx context.Context, entitiesToDelete []*sdkNotification.EventSubscription, scope string, scopeId *string) diag.Diagnostics {
	var retVal diag.Diagnostics

	resourceIdentifier := r.logIdentifier(scope, scopeId)
	tflog.Info(ctx, fmt.Sprintf("Deleting %d subscriptions from resource %s.", len(entitiesToDelete), resourceIdentifier))

	for _, e := range entitiesToDelete {
		_, err := r.client.Client.notification.DeleteEventSubscription(ctx, *e.ID)

		if err != nil {
			if commons.IsNotFoundResponseError(err) {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(resourceNotFoundError, fmt.Sprintf("Failed to delete subscription id '%s' from resource %s. Error: %s",
						*e.ID, resourceIdentifier, err)),
				}
			} else {
				return diag.Diagnostics{
					diag.NewErrorDiagnostic(fmt.Sprintf("Failed to delete subscription id '%s' from resource %s.", *e.ID, resourceIdentifier),
						err.Error()),
				}
			}
		}
	}

	return retVal
}

func (r *EventsSubscriptionsResource) logIdentifier(scope string, scopeId *string) string {
	retVal := fmt.Sprintf("scope '%s'", scope)
	if scopeId != nil {
		retVal += fmt.Sprintf(" with scope_id '%s'", *scopeId)
	}
	return retVal
}

func (r *EventsSubscriptionsResource) logSubscriptionIdentifier(eventType string, notificationEndpointId string) string {
	retVal := fmt.Sprintf("event_type '%s' to notification_endpoint_id '%s'", eventType, notificationEndpointId)
	return retVal
}

func (r *EventsSubscriptionsResource) buildId(plan tfEventsSubscriptions.ResourceModel) types.String {
	var retVal types.String
	var id string

	scope := plan.Scope.ValueString()

	if plan.ScopeId.IsNull() {
		id = scope
	} else {
		scopeId := plan.ScopeId.ValueString()
		id = fmt.Sprintf("%s/%s", scope, scopeId)
	}

	retVal = helpers.StringValueOrNull(&id)

	return retVal
}

func (r *EventsSubscriptionsResource) breakdownId(id types.String) (string, *string) {
	var scope string
	var scopeId *string

	idVal := id.ValueString()

	if strings.Contains(idVal, "/") {
		split := strings.Split(idVal, "/")
		scope = split[0]
		scopeId = &split[1]
	} else {
		scope = idVal
		scopeId = nil
	}

	return scope, scopeId
}

// When creating new subscriptions, the state file should have the new id of the subscription. The id of the subscription is used for deleteEntities
func (r *EventsSubscriptionsResource) updateIdForTfSubscriptions(plan tfEventsSubscriptions.ResourceModel, newEntities []*sdkNotification.EventSubscription) {
	for _, planEntity := range plan.Subscriptions {
		for _, newEntity := range newEntities {
			if planEntity.NotificationEndpointId.ValueString() == controlmonkey.StringValue(newEntity.NotificationEndpointId) &&
				planEntity.EventType.ValueString() == controlmonkey.StringValue(newEntity.EventType) {
				planEntity.ID = helpers.StringValueOrNull(newEntity.ID)
			}
		}
	}
}

//endregion
