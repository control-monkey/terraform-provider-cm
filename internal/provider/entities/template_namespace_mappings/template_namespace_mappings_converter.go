package template_namespace_mappings

import (
	"github.com/control-monkey/controlmonkey-sdk-go/services/template"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/interfaces"
	"github.com/hashicorp/go-set/v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MergedEntities struct {
	EntitiesToCreate []*template.TemplateNamespaceMapping
	EntitiesToUpdate []*template.TemplateNamespaceMapping
	EntitiesToDelete []*template.TemplateNamespaceMapping
}

func Merge(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) *MergedEntities {
	retVal := new(MergedEntities)

	if plan == nil {
		plan = new(ResourceModel) // delete merger
	}

	if state == nil {
		state = new(ResourceModel) // create merger
	}

	var templateId types.String
	if plan.TemplateId.IsNull() == false {
		templateId = plan.TemplateId
	} else {
		templateId = state.TemplateId
	}

	mergeResult := interfaces.MergeEntities(plan.Namespaces, state.Namespaces)
	retVal.EntitiesToCreate = convertEntities(mergeResult.EntitiesToCreate, templateId)
	retVal.EntitiesToUpdate = convertEntities(mergeResult.EntitiesToUpdate, templateId)
	retVal.EntitiesToDelete = convertEntities(mergeResult.EntitiesToDelete, templateId)

	return retVal
}

func convertEntities(entities set.Collection[*NamespaceModel], templateId types.String) []*template.TemplateNamespaceMapping {
	retVal := make([]*template.TemplateNamespaceMapping, entities.Size())

	for i, u := range entities.Slice() {
		tu := new(template.TemplateNamespaceMapping)
		tu.SetTemplateId(templateId.ValueStringPointer())
		tu.SetNamespaceId(u.NamespaceId.ValueStringPointer())

		retVal[i] = tu
	}

	return retVal
}
