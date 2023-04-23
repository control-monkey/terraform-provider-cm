package provider

import (
	"context"
	"fmt"

	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	"github.com/control-monkey/controlmonkey-sdk-go/services/namespace"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
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
var _ resource.Resource = &NamespaceResource{}

func NewNamespaceResource() resource.Resource {
	return &NamespaceResource{}
}

type NamespaceResource struct {
	client *ControlMonkeyAPIClient
}

type NamespaceResourceModel struct {
	ID                  types.String                 `tfsdk:"id"`
	Name                types.String                 `tfsdk:"name"`
	Description         types.String                 `tfsdk:"description"`
	ExternalCredentials *[]*externalCredentialsModel `tfsdk:"external_credentials"`
	Policy              *namespacePolicyModel        `tfsdk:"policy"`
}

type externalCredentialsModel struct {
	Type                  types.String `tfsdk:"type"`
	ExternalCredentialsId types.String `tfsdk:"external_credentials_id"`
}

type namespacePolicyModel struct {
	TtlConfig *namespaceTtlConfigModel `tfsdk:"ttl_config"`
}

type namespaceTtlConfigModel struct {
	MaxTtl     *namespaceTtlDefinitionModel `tfsdk:"max_ttl"`
	DefaultTtl *namespaceTtlDefinitionModel `tfsdk:"default_ttl"`
}

type namespaceTtlDefinitionModel struct {
	Type  types.String `tfsdk:"type"`
	Value types.Int64  `tfsdk:"value"`
}

var externalCredentialTypes = []string{"awsAssumeRole", "gcpServiceAccount", "azureServicePrincipal"}
var namespaceTtlTypes = []string{"hours", "days"}

func (r *NamespaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace"
}

func (r *NamespaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates, updates and destroys namespaces.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the namespace.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the namespace.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the namespace.",
				Optional:            true,
			},
			"external_credentials": schema.ListNestedAttribute{
				MarkdownDescription: "List of cloud credentials attached to the namespace.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: fmt.Sprintf("The type of the credentials. Allowed values: %s.", helpers.EnumForDocs(externalCredentialTypes)),
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(externalCredentialTypes...),
							},
						},
						"external_credentials_id": schema.StringAttribute{
							MarkdownDescription: "The Control Monkey unique ID of the credentials.",
							Required:            true,
						},
					},
				},
			},
			"policy": schema.SingleNestedAttribute{
				MarkdownDescription: "The policy of the namespace.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"ttl_config": schema.SingleNestedAttribute{
						MarkdownDescription: "The time to live config of the namespace policy regarding to its stacks.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"max_ttl": schema.SingleNestedAttribute{
								MarkdownDescription: "The max time to live for new stacks in the namespace.",
								Required:            true,
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										MarkdownDescription: fmt.Sprintf("The type of the ttl. Allowed values: %s.", helpers.EnumForDocs(namespaceTtlTypes)),
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf(namespaceTtlTypes...),
										},
									},
									"value": schema.Int64Attribute{
										MarkdownDescription: "The value that corresponds the type",
										Required:            true,
									},
								},
							},
							"default_ttl": schema.SingleNestedAttribute{
								MarkdownDescription: "The default time to live for new stacks in the namespace.",
								Required:            true,
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										MarkdownDescription: fmt.Sprintf("The type of the ttl. Allowed values: %s.", helpers.EnumForDocs(namespaceTtlTypes)),
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf(namespaceTtlTypes...),
										},
									},
									"value": schema.Int64Attribute{
										MarkdownDescription: "The value that corresponds the type",
										Required:            true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *NamespaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the latest data.
func (r *NamespaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//Get current state
	var state NamespaceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	res, err := r.client.Client.namespace.ReadNamespace(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read namespace %s", id), fmt.Sprintf("%s", err))
		return
	}

	namespace_ := res.Namespace

	state.Name = helpers.StringValueOrNull(namespace_.Name)
	state.Description = helpers.StringValueOrNull(namespace_.Description)

	if namespace_.ExternalCredentials != nil {
		var creds []*externalCredentialsModel

		for _, v := range *namespace_.ExternalCredentials {
			var ec externalCredentialsModel
			ec.Type = helpers.StringValueOrNull(v.Type)
			ec.ExternalCredentialsId = helpers.StringValueOrNull(v.ExternalCredentialsId)
			creds = append(creds, &ec)
		}
		state.ExternalCredentials = &creds
	} else {
		state.ExternalCredentials = nil
	}

	var policy namespacePolicyModel
	if namespace_.Policy != nil {
		if namespace_.Policy.TtlConfig != nil {
			var ttlc namespaceTtlConfigModel
			var maxTtl namespaceTtlDefinitionModel
			var defTtl namespaceTtlDefinitionModel

			maxTtl.Type = helpers.StringValueOrNull(namespace_.Policy.TtlConfig.MaxTtl.Type)
			maxTtl.Value = helpers.Int64ValueOrNull(namespace_.Policy.TtlConfig.MaxTtl.Value)
			defTtl.Type = helpers.StringValueOrNull(namespace_.Policy.TtlConfig.DefaultTtl.Type)
			defTtl.Value = helpers.Int64ValueOrNull(namespace_.Policy.TtlConfig.DefaultTtl.Value)

			ttlc.MaxTtl = &maxTtl
			ttlc.DefaultTtl = &defTtl
			policy.TtlConfig = &ttlc
		} else {
			policy.TtlConfig = nil
		}
		state.Policy = &policy
	} else {
		state.Policy = nil
	}
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *NamespaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan NamespaceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := cmNamespaceFromPlanConverter(plan)

	res, err := r.client.Client.namespace.CreateNamespace(ctx, &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Namespace creation failed",
			fmt.Sprintf("failed to create namespace, error: %s", err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(controlmonkey.StringValue(res.Namespace.ID))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *NamespaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan NamespaceResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	body := cmNamespaceFromPlanConverter(plan)

	_, err := r.client.Client.namespace.UpdateNamespace(ctx, id, &body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Namespace update failed",
			fmt.Sprintf("failed to update namespace %s, error: %s", id, err.Error()),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *NamespaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state NamespaceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	_, err := r.client.Client.namespace.DeleteNamespace(ctx, id)

	if err != nil {
		errMsg := err.Error()
		resp.Diagnostics.AddError(
			"Namespace deletion failed",
			fmt.Sprintf("Failed to delete namespace %s, error: %s", id, errMsg),
		)
		return
	}
}

func (r *NamespaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

//region Private

func cmNamespaceFromPlanConverter(plan NamespaceResourceModel) namespace.Namespace {
	var body namespace.Namespace

	body.SetName(plan.Name.ValueStringPointer())
	body.SetDescription(plan.Description.ValueStringPointer())

	if plan.ExternalCredentials != nil {
		var creds []*namespace.ExternalCredentials

		for _, extCreds := range *plan.ExternalCredentials {
			var ec namespace.ExternalCredentials
			ec.Type = extCreds.Type.ValueStringPointer()
			ec.ExternalCredentialsId = extCreds.ExternalCredentialsId.ValueStringPointer()
			creds = append(creds, &ec)
		}
		body.ExternalCredentials = &creds
	} else {
		body.ExternalCredentials = nil
	}

	var policy namespace.Policy
	if plan.Policy != nil {
		if plan.Policy.TtlConfig != nil {
			var ttlConfig namespace.TtlConfig
			var maxTtl namespace.TtlDefinition
			var defTtl namespace.TtlDefinition

			maxTtl.SetType(plan.Policy.TtlConfig.MaxTtl.Type.ValueStringPointer())
			maxTtl.SetValue(controlmonkey.Int(int(plan.Policy.TtlConfig.MaxTtl.Value.ValueInt64())))
			defTtl.SetType(plan.Policy.TtlConfig.DefaultTtl.Type.ValueStringPointer())
			defTtl.SetValue(controlmonkey.Int(int(plan.Policy.TtlConfig.DefaultTtl.Value.ValueInt64())))

			ttlConfig.SetMaxTtl(&maxTtl)
			ttlConfig.SetDefaultTtl(&defTtl)
			policy.SetTtlConfig(&ttlConfig)
		} else {
			policy.SetTtlConfig(nil)
		}
		body.SetPolicy(&policy)
	} else {
		body.SetPolicy(nil)
	}

	return body
}

//endregion
