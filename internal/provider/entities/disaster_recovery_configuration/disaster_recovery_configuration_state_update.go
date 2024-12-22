package disaster_recovery_configuration

import (
	"encoding/json"
	apiDisasterRecovery "github.com/control-monkey/controlmonkey-sdk-go/services/disaster_recovery"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
)

func UpdateStateAfterRead(apiEntity *apiDisasterRecovery.DisasterRecoveryConfiguration, state *ResourceModel) {
	state.Scope = helpers.StringValueOrNull(apiEntity.Scope)
	state.CloudAccountId = helpers.StringValueOrNull(apiEntity.CloudAccountId)

	if apiEntity.BackupStrategy != nil {
		bs := updateStateAfterReadBackupStrategy(apiEntity.BackupStrategy)
		state.BackupStrategy = &bs
	} else {
		state.BackupStrategy = nil
	}
}

func updateStateAfterReadVcsInfo(apiEntity *apiDisasterRecovery.VcsInfo) VcsInfoModel {
	var retVal VcsInfoModel

	retVal.ProviderId = helpers.StringValueOrNull(apiEntity.ProviderId)
	retVal.RepoName = helpers.StringValueOrNull(apiEntity.RepoName)
	retVal.Branch = helpers.StringValueOrNull(apiEntity.Branch)

	return retVal
}

func updateStateAfterReadBackupStrategy(apiEntity *apiDisasterRecovery.BackupStrategy) BackupStrategyModel {
	var retVal BackupStrategyModel

	retVal.IncludeManagedResources = helpers.BoolValueOrNull(apiEntity.IncludeManagedResources)
	retVal.Mode = helpers.StringValueOrNull(apiEntity.Mode)

	if apiEntity.VcsInfo != nil {
		vi := updateStateAfterReadVcsInfo(apiEntity.VcsInfo)
		retVal.VcsInfo = &vi
	} else {
		retVal.VcsInfo = nil
	}

	if apiEntity.Groups != nil {
		groupsString, err := json.Marshal(apiEntity.Groups)

		if err != nil {
			retVal.GroupsJson = jsontypes.NewNormalizedNull()
		} else {
			retVal.GroupsJson = jsontypes.NewNormalizedValue(string(groupsString))
		}
	} else {
		retVal.GroupsJson = jsontypes.NewNormalizedNull()
	}
	return retVal
}
