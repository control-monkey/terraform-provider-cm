package organization

import (
	sdkOrganization "github.com/control-monkey/controlmonkey-sdk-go/services/organization"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(res *sdkOrganization.OrgConfiguration, state *ResourceModel) {

	if res.IacConfig != nil {
		ic := updateStateAfterReadIacConfig(res.IacConfig)
		state.IacConfig = &ic
	} else {
		state.IacConfig = nil
	}

	if res.S3StateFilesLocations != nil {
		locations := updateStateAfterReadS3StateFilesLocations(res.S3StateFilesLocations)
		state.S3StateFilesLocations = locations
	} else {
		state.S3StateFilesLocations = nil
	}

	if res.RunnerConfig != nil {
		rc := updateStateAfterReadRunnerConfig(res.RunnerConfig)
		state.RunnerConfig = &rc
	} else {
		state.RunnerConfig = nil
	}

	if res.SuppressedResources != nil {
		sr := updateStateAfterReadSuppressedResources(res.SuppressedResources)
		state.SuppressedResources = &sr
	} else {
		state.SuppressedResources = nil
	}

	if res.ReportConfigurations != nil {
		rc := updateStateAfterReadReportConfigurations(res.ReportConfigurations)
		state.ReportConfigurations = rc
	} else {
		state.ReportConfigurations = nil
	}

}

func updateStateAfterReadS3StateFilesLocations(apiEntities []*sdkOrganization.S3StateFilesLocation) []*S3StateFilesLocationModel {
	var retVal []*S3StateFilesLocationModel

	if apiEntities != nil {
		retVal = make([]*S3StateFilesLocationModel, 0)

		for _, apiEntity := range apiEntities {
			tfEntity := updateStateAfterReadLocations(apiEntity)
			retVal = append(retVal, &tfEntity)
		}
	}

	return retVal
}

func updateStateAfterReadLocations(apiEntity *sdkOrganization.S3StateFilesLocation) S3StateFilesLocationModel {
	var retVal S3StateFilesLocationModel

	retVal.BucketName = helpers.StringValueOrNull(apiEntity.BucketName)
	retVal.BucketRegion = helpers.StringValueOrNull(apiEntity.BucketRegion)
	retVal.AwsAccountId = helpers.StringValueOrNull(apiEntity.AwsAccountId)

	return retVal
}

func updateStateAfterReadIacConfig(iacConfig *sdkOrganization.IacConfig) IacConfigModel {
	var retVal IacConfigModel

	retVal.TerraformVersion = helpers.StringValueOrNull(iacConfig.TerraformVersion)
	retVal.TerragruntVersion = helpers.StringValueOrNull(iacConfig.TerragruntVersion)
	retVal.OpentofuVersion = helpers.StringValueOrNull(iacConfig.OpentofuVersion)

	return retVal
}

func updateStateAfterReadRunnerConfig(rc *sdkOrganization.RunnerConfig) RunnerConfigModel {
	var retVal RunnerConfigModel

	retVal.Mode = helpers.StringValueOrNull(rc.Mode)
	retVal.Groups = helpers.StringPointerSliceToTfList(rc.Groups)
	retVal.IsOverridable = helpers.BoolValueOrNull(rc.IsOverridable)

	return retVal
}

func updateStateAfterReadSuppressedResources(apiEntity *sdkOrganization.SuppressedResources) SuppressedResourcesModel {
	var retVal SuppressedResourcesModel

	if apiEntity.ManagedByTags != nil {
		ts := updateStateAfterReadManagedByTags(apiEntity.ManagedByTags)
		retVal.ManagedByTags = ts
	} else {
		retVal.ManagedByTags = nil
	}

	return retVal
}

func updateStateAfterReadManagedByTags(apiEntities []*sdkOrganization.TagProperties) []*TagPropertiesModel {
	var retVal []*TagPropertiesModel

	if apiEntities != nil {
		retVal = make([]*TagPropertiesModel, 0)

		for _, apiEntity := range apiEntities {
			tfEntity := updateStateAfterReadTag(apiEntity)
			retVal = append(retVal, &tfEntity)
		}
	}

	return retVal
}

func updateStateAfterReadTag(apiEntity *sdkOrganization.TagProperties) TagPropertiesModel {
	var retVal TagPropertiesModel

	retVal.Key = helpers.StringValueOrNull(apiEntity.Key)
	retVal.Value = helpers.StringValueOrNull(apiEntity.Value)

	return retVal
}

func updateStateAfterReadReportConfigurations(apiEntities []*sdkOrganization.ReportConfiguration) []*ReportConfigurationModel {
	var retVal []*ReportConfigurationModel

	if apiEntities != nil {
		retVal = make([]*ReportConfigurationModel, 0)

		for _, apiEntity := range apiEntities {
			tfEntity := updateStateAfterReadReportConfiguration(apiEntity)
			retVal = append(retVal, &tfEntity)
		}
	}

	return retVal
}

func updateStateAfterReadReportConfiguration(apiEntity *sdkOrganization.ReportConfiguration) ReportConfigurationModel {
	var retVal ReportConfigurationModel

	retVal.Type = helpers.StringValueOrNull(apiEntity.Type)

	if apiEntity.Recipients != nil {
		tfEntity := updateStateAfterReadReportRecipients(apiEntity.Recipients)
		retVal.Recipients = &tfEntity
	} else {
		retVal.Recipients = nil
	}

	retVal.Enabled = helpers.BoolValueOrNull(apiEntity.Enabled)

	return retVal
}

func updateStateAfterReadReportRecipients(apiEntity *sdkOrganization.ReportRecipients) ReportRecipientsModel {
	var retVal ReportRecipientsModel

	retVal.AllAdmins = helpers.BoolValueOrNull(apiEntity.AllAdmins)
	retVal.EmailAddresses = helpers.StringPointerSliceToTfList(apiEntity.EmailAddresses)
	retVal.EmailAddressesToExclude = helpers.StringPointerSliceToTfList(apiEntity.EmailAddressesToExclude)

	return retVal
}
