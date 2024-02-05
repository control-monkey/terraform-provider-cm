package template

import (
	sdkTemplate "github.com/control-monkey/controlmonkey-sdk-go/services/template"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(res *sdkTemplate.Template, state *ResourceModel) {
	template := res

	state.IacType = helpers.StringValueOrNull(template.IacType)
	state.Name = helpers.StringValueOrNull(template.Name)
	state.Description = helpers.StringValueOrNull(template.Description)

	if template.VcsInfo != nil {
		vcs := updateStateAfterReadVcsInfo(template.VcsInfo)
		state.VcsInfo = &vcs
	} else {
		state.VcsInfo = nil
	}

	if res.Policy != nil {
		policy := updateStateAfterReadPolicy(res.Policy)
		state.Policy = &policy
	} else {
		state.Policy = nil
	}

	state.SkipStateRefreshOnDestroy = helpers.BoolValueOrNull(template.SkipStateRefreshOnDestroy)
}

func updateStateAfterReadVcsInfo(vcsInfo *sdkTemplate.VcsInfo) VcsInfoModel {
	var retVal VcsInfoModel

	retVal.ProviderId = helpers.StringValueOrNull(vcsInfo.ProviderId)
	retVal.RepoName = helpers.StringValueOrNull(vcsInfo.RepoName)
	retVal.Path = helpers.StringValueOrNull(vcsInfo.Path)
	retVal.Branch = helpers.StringValueOrNull(vcsInfo.Branch)

	return retVal
}

func updateStateAfterReadPolicy(policy *sdkTemplate.Policy) PolicyModel {
	var retVal PolicyModel

	if policy.TtlConfig != nil {
		ttlConfig := updateStateAfterReadTtlConfig(policy.TtlConfig)
		retVal.TtlConfig = &ttlConfig
	} else {
		retVal.TtlConfig = nil
	}

	return retVal
}

func updateStateAfterReadTtlConfig(ttlConfig *sdkTemplate.TtlConfig) TtlConfigModel {
	var retVal TtlConfigModel

	if ttlConfig.MaxTtl != nil {
		ttl := updateStateAfterReadTtlDefinition(ttlConfig.MaxTtl)
		retVal.MaxTtl = &ttl
	} else {
		retVal.MaxTtl = nil
	}
	if ttlConfig.DefaultTtl != nil {
		ttl := updateStateAfterReadTtlDefinition(ttlConfig.DefaultTtl)
		retVal.DefaultTtl = &ttl
	} else {
		retVal.DefaultTtl = nil
	}

	return retVal
}

func updateStateAfterReadTtlDefinition(ttl *sdkTemplate.TtlDefinition) TtlDefinitionModel {
	var retVal TtlDefinitionModel

	retVal.Type = helpers.StringValueOrNull(ttl.Type)
	retVal.Value = helpers.Int64ValueOrNull(ttl.Value)

	return retVal
}
