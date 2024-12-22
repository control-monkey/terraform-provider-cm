package disaster_recovery_configuration

import (
	apiDisasterRecovery "github.com/control-monkey/controlmonkey-sdk-go/services/disaster_recovery"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*apiDisasterRecovery.DisasterRecoveryConfiguration, bool) {
	var retVal *apiDisasterRecovery.DisasterRecoveryConfiguration

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(apiDisasterRecovery.DisasterRecoveryConfiguration)
	hasChanges := false

	if state == nil {
		state = new(ResourceModel) // dummy initialization
		hasChanges = true          // must have changes because before is null and after is not
	}

	if plan.Scope != state.Scope {
		retVal.SetScope(plan.Scope.ValueStringPointer())
		hasChanges = true
	}

	if plan.CloudAccountId != state.CloudAccountId {
		retVal.SetCloudAccountId(plan.CloudAccountId.ValueStringPointer())
		hasChanges = true
	}

	if bs, hasChanged := backupStrategyConverter(plan.BackupStrategy, state.BackupStrategy, converterType); hasChanged {
		retVal.SetBackupStrategy(bs)
		hasChanges = true
	}

	return retVal, hasChanges
}

func vcsInfoConverter(plan *VcsInfoModel, state *VcsInfoModel, converterType commons.ConverterType) (*apiDisasterRecovery.VcsInfo, bool) {
	var retVal *apiDisasterRecovery.VcsInfo

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(apiDisasterRecovery.VcsInfo)
	hasChanges := false

	if state == nil {
		state = new(VcsInfoModel) // dummy initialization
		hasChanges = true         // must have changes because before is null and after is not
	}

	if plan.ProviderId != state.ProviderId {
		retVal.SetProviderId(plan.ProviderId.ValueStringPointer())
		hasChanges = true
	}
	if plan.RepoName != state.RepoName {
		retVal.SetRepoName(plan.RepoName.ValueStringPointer())
		hasChanges = true
	}
	if plan.Branch != state.Branch {
		retVal.SetBranch(plan.Branch.ValueStringPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}

func backupStrategyConverter(plan *BackupStrategyModel, state *BackupStrategyModel, converterType commons.ConverterType) (*apiDisasterRecovery.BackupStrategy, bool) {
	var retVal *apiDisasterRecovery.BackupStrategy

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(apiDisasterRecovery.BackupStrategy)
	hasChanges := false

	if state == nil {
		state = new(BackupStrategyModel) // dummy initialization
		hasChanges = true                // must have changes because before is null and after is not
	}

	if plan.IncludeManagedResources != state.IncludeManagedResources {
		retVal.SetIncludeManagedResources(plan.IncludeManagedResources.ValueBoolPointer())
		hasChanges = true
	}

	if plan.Mode != state.Mode {
		retVal.SetMode(plan.Mode.ValueStringPointer())
		hasChanges = true
	}

	if vcsInfo, hasChanged := vcsInfoConverter(plan.VcsInfo, state.VcsInfo, converterType); hasChanged {
		retVal.SetVcsInfo(vcsInfo)
		hasChanges = true
	}

	if plan.Groups != state.Groups {
		var groupsList []*map[string]interface{}

		if plan.Groups.IsNull() {
			groupsList = nil
		} else {
			plan.Groups.Unmarshal(&groupsList)
		}

		retVal.SetGroups(groupsList)
		hasChanges = true
	}

	return retVal, hasChanges
}
