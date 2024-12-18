package blueprint

import (
	apiBlueprint "github.com/control-monkey/controlmonkey-sdk-go/services/blueprint"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/entities/cross_models"
)

func UpdateStateAfterRead(res *apiBlueprint.Blueprint, state *ResourceModel) {
	blueprint := res

	state.Name = helpers.StringValueOrNull(blueprint.Name)
	state.Description = helpers.StringValueIfNotEqual(blueprint.Description, "")

	if blueprint.BlueprintVcsInfo != nil {
		bvi := updateStateAfterReadBlueprintVcsInfo(blueprint.BlueprintVcsInfo)
		state.BlueprintVcsInfo = &bvi
	} else {
		state.BlueprintVcsInfo = nil
	}

	if blueprint.StackConfiguration != nil {
		sc := updateStateAfterReadStackConfiguration(blueprint.StackConfiguration)
		state.StackConfiguration = &sc
	} else {
		state.StackConfiguration = nil
	}

	if blueprint.SubstituteParameters != nil {
		ec := updateStateAfterReadSubstituteParameters(blueprint.SubstituteParameters)
		state.SubstituteParameters = ec
	} else {
		state.SubstituteParameters = nil
	}

	state.SkipPlanOnStackInitialization = helpers.BoolValueOrNull(blueprint.SkipPlanOnStackInitialization)
	state.AutoApproveApplyOnInitialization = helpers.BoolValueOrNull(blueprint.AutoApproveApplyOnInitialization)
}

func updateStateAfterReadBlueprintVcsInfo(vi *apiBlueprint.VcsInfo) VcsInfoModel {
	var retVal VcsInfoModel

	retVal.ProviderId = helpers.StringValueOrNull(vi.ProviderId)
	retVal.RepoName = helpers.StringValueOrNull(vi.RepoName)
	retVal.Path = helpers.StringValueOrNull(vi.Path)
	retVal.Branch = helpers.StringValueOrNull(vi.Branch)

	return retVal
}

func updateStateAfterReadStackConfiguration(sc *apiBlueprint.StackConfiguration) StackConfigurationModel {
	var retVal StackConfigurationModel

	retVal.NamePattern = helpers.StringValueOrNull(sc.NamePattern)
	retVal.IacType = helpers.StringValueOrNull(sc.IacType)

	if sc.VcsInfoWithPatterns != nil {
		vp := updateStateAfterReadVcsInfoWithPatterns(sc.VcsInfoWithPatterns)
		retVal.VcsInfoWithPatterns = &vp
	} else {
		retVal.VcsInfoWithPatterns = nil
	}

	if sc.DeploymentApprovalPolicy != nil {
		dap := cross_models.UpdateStateAfterReadDeploymentApprovalPolicy(sc.DeploymentApprovalPolicy)
		retVal.DeploymentApprovalPolicy = &dap
	} else {
		sc.DeploymentApprovalPolicy = nil
	}

	return retVal
}

func updateStateAfterReadVcsInfoWithPatterns(vp *apiBlueprint.StackVcsInfoWithPatterns) StackVcsInfoWithPatternsModel {
	var retVal StackVcsInfoWithPatternsModel

	retVal.ProviderId = helpers.StringValueOrNull(vp.ProviderId)
	retVal.RepoName = helpers.StringValueOrNull(vp.RepoName)
	retVal.PathPattern = helpers.StringValueOrNull(vp.PathPattern)
	retVal.BranchPattern = helpers.StringValueOrNull(vp.BranchPattern)

	return retVal
}

func updateStateAfterReadSubstituteParameters(sps []*apiBlueprint.SubstituteParameter) []*SubstituteParameterModel {
	var retVal []*SubstituteParameterModel

	if sps != nil {
		retVal = make([]*SubstituteParameterModel, 0)

		for _, sp := range sps {
			spm := updateStateAfterReadParameter(sp)
			retVal = append(retVal, &spm)
		}
	}

	return retVal
}

func updateStateAfterReadParameter(credentials *apiBlueprint.SubstituteParameter) SubstituteParameterModel {
	var retVal SubstituteParameterModel

	retVal.Key = helpers.StringValueOrNull(credentials.Key)
	retVal.Description = helpers.StringValueOrNull(credentials.Description)
	retVal.ValueConditions = cross_models.UpdateStateAfterReadValueConditions(credentials.ValueConditions)

	return retVal
}
