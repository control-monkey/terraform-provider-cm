package template_namespace_mappings

import (
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID         types.String      `tfsdk:"id"`
	TemplateId types.String      `tfsdk:"template_id"`
	Namespaces []*NamespaceModel `tfsdk:"namespaces"`
}

type NamespaceModel struct { //When new field is added consider Hash() function
	NamespaceId types.String `tfsdk:"namespace_id"`
}

func (e *NamespaceModel) Hash() string {
	return e.NamespaceId.ValueString()
}

func (e *NamespaceModel) GetBlockIdentifier() string {
	retVal := ""

	if helpers.IsKnown(e.NamespaceId) {
		retVal += e.Hash() // do not use e.Hash if another property is added to Model
	}

	return retVal
}
