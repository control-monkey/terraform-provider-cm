package template_namespace_mappings

import (
	sdkTemplate "github.com/control-monkey/controlmonkey-sdk-go/services/template"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(res []*sdkTemplate.TemplateNamespaceMapping, state *ResourceModel) {
	templateNamespaces := res

	state.TemplateId = state.ID

	if templateNamespaces != nil {
		ec := updateStateAfterReadTemplateNamespaces(templateNamespaces)
		state.Namespaces = ec
	} else {
		state.Namespaces = nil
	}
}

func updateStateAfterReadTemplateNamespaces(templateNamespaces []*sdkTemplate.TemplateNamespaceMapping) []*NamespaceModel {
	var retVal []*NamespaceModel

	if len(templateNamespaces) > 0 {
		retVal = make([]*NamespaceModel, len(templateNamespaces))

		for i, namespace := range templateNamespaces {
			u := updateStateAfterReadNamespace(namespace)
			retVal[i] = &u
		}
	}

	return retVal
}

func updateStateAfterReadNamespace(namespace *sdkTemplate.TemplateNamespaceMapping) NamespaceModel {
	var retVal NamespaceModel

	retVal.NamespaceId = helpers.StringValueOrNull(namespace.NamespaceId)

	return retVal
}
