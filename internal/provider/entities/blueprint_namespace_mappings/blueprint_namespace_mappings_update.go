package blueprint_namespace_mappings

import (
	sdkBlueprint "github.com/control-monkey/controlmonkey-sdk-go/services/blueprint"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(res []*sdkBlueprint.BlueprintNamespaceMapping, state *ResourceModel) {
	blueprintNamespaces := res

	state.BlueprintId = state.ID

	if blueprintNamespaces != nil {
		ec := updateStateAfterReadBlueprintNamespaces(blueprintNamespaces)
		state.Namespaces = ec
	} else {
		state.Namespaces = nil
	}
}

func updateStateAfterReadBlueprintNamespaces(blueprintNamespaces []*sdkBlueprint.BlueprintNamespaceMapping) []*NamespaceModel {
	var retVal []*NamespaceModel

	if len(blueprintNamespaces) > 0 {
		retVal = make([]*NamespaceModel, len(blueprintNamespaces))

		for i, namespace := range blueprintNamespaces {
			u := updateStateAfterReadNamespace(namespace)
			retVal[i] = &u
		}
	}

	return retVal
}

func updateStateAfterReadNamespace(namespace *sdkBlueprint.BlueprintNamespaceMapping) NamespaceModel {
	var retVal NamespaceModel

	retVal.NamespaceId = helpers.StringValueOrNull(namespace.NamespaceId)

	return retVal
}
