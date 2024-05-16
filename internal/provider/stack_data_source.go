package provider

import (
	"context"
	"fmt"
	tfStack "github.com/control-monkey/terraform-provider-cm/internal/provider/entities/stack_data"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &StackDataSource{}

func NewStackDataSource() datasource.DataSource {
	return &StackDataSource{}
}

type StackDataSource struct {
	client *ControlMonkeyAPIClient
}

func (r *StackDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stack"
}

func (r *StackDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the stack.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("id"), path.MatchRoot("name")),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the stack.",
				Optional:            true,
			},
			"namespace_id": schema.StringAttribute{
				MarkdownDescription: "The namespace ID where the stack is located.",
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *StackDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (r *StackDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	//Get current state
	var state tfStack.ResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stackId := state.ID.ValueStringPointer()
	stackName := state.Name.ValueStringPointer()
	namespaceId := state.NamespaceId.ValueStringPointer()
	res, err := r.client.Client.stack.ListStacks(ctx, stackId, stackName, namespaceId)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read stack"), fmt.Sprintf("%s", err))
		return
	} else if len(res) == 0 {
		resp.Diagnostics.AddError(fmt.Sprintf(resourceNotFoundError), fmt.Sprintf(stackNotFoundError))
		return
	} else if len(res) > 1 {
		errMsg := fmt.Sprintf("Found multiple stacks with name '%s'", *stackName)
		if namespaceId != nil {
			errMsg += fmt.Sprintf(" in namespace id '%s'; use additional constraints to reduce matches to a single match", *namespaceId)
		}
		errMsg += "."
		resp.Diagnostics.AddError(multipleEntitiesError, errMsg)
		return
	}

	tfStack.UpdateStateAfterRead(res[0], &state, &resp.Diagnostics)

	// Set refreshed state
	// Save data into Terraform state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
