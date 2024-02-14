package organization

import (
	"github.com/control-monkey/controlmonkey-sdk-go/services/organization"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"reflect"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*organization.OrgConfiguration, bool) {
	var retVal *organization.OrgConfiguration

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(organization.OrgConfiguration)
	hasChanges := false

	if state == nil {
		state = new(ResourceModel) // dummy initialization
		hasChanges = true          // must have changes because before is null and after is not
	}

	if iacConfig, hasChanged := iacConfigConverter(plan.IacConfig, state.IacConfig, converterType); hasChanged {
		retVal.SetIacConfig(iacConfig)
		hasChanges = true
	}

	if stateFilesLocations, hasChanged := s3StateFilesLocationsConverter(plan.S3StateFilesLocations, state.S3StateFilesLocations, converterType); hasChanged {
		retVal.SetTfStateFilesS3Locations(stateFilesLocations)
		hasChanges = true
	}

	if runnerConfig, hasChanged := runnerConfigConverter(plan.RunnerConfig, state.RunnerConfig, converterType); hasChanged {
		retVal.SetRunnerConfig(runnerConfig)
		hasChanges = true
	}

	if suppressedResources, hasChanged := suppressedResourcesConverter(plan.SuppressedResources, state.SuppressedResources, converterType); hasChanged {
		retVal.SetSuppressedResources(suppressedResources)
		hasChanges = true
	}

	if reportConfigurations, hasChanged := reportConfigurationsConverter(plan.ReportConfigurations, state.ReportConfigurations, converterType); hasChanged {
		retVal.SetReportConfigurations(reportConfigurations)
		hasChanges = true
	}

	return retVal, hasChanges
}

func iacConfigConverter(plan *IacConfigModel, state *IacConfigModel, converterType commons.ConverterType) (*organization.IacConfig, bool) {
	var retVal *organization.IacConfig

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(organization.IacConfig)
	hasChanges := false

	if state == nil {
		state = new(IacConfigModel) // dummy initialization
		hasChanges = true           // must have changes because before is null and after is not
	}

	if plan.TerraformVersion != state.TerraformVersion {
		retVal.SetTerraformVersion(plan.TerraformVersion.ValueStringPointer())
		hasChanges = true
	}
	if plan.TerragruntVersion != state.TerragruntVersion {
		retVal.SetTerragruntVersion(plan.TerragruntVersion.ValueStringPointer())
		hasChanges = true
	}
	if plan.OpentofuVersion != state.OpentofuVersion {
		retVal.SetOpentofuVersion(plan.OpentofuVersion.ValueStringPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}

func s3StateFilesLocationsConverter(plan []*S3StateFilesLocationModel, state []*S3StateFilesLocationModel, converterType commons.ConverterType) ([]*organization.S3StateFilesLocation, bool) {
	var retVal []*organization.S3StateFilesLocation
	hasChanged := false

	if reflect.DeepEqual(plan, state) == false {
		hasChanged = true

		if plan != nil {
			retVal = make([]*organization.S3StateFilesLocation, 0)

			for _, r := range plan {
				rule := s3StateFilesLocationConverter(r)
				retVal = append(retVal, rule)
			}
		}
	}

	return retVal, hasChanged
}

func s3StateFilesLocationConverter(plan *S3StateFilesLocationModel) *organization.S3StateFilesLocation {
	retVal := new(organization.S3StateFilesLocation)

	retVal.SetBucketName(plan.BucketName.ValueStringPointer())
	retVal.SetBucketRegion(plan.BucketRegion.ValueStringPointer())
	retVal.SetAwsAccountId(plan.AwsAccountId.ValueStringPointer())

	return retVal
}

func runnerConfigConverter(plan *RunnerConfigModel, state *RunnerConfigModel, converterType commons.ConverterType) (*organization.RunnerConfig, bool) {
	var retVal *organization.RunnerConfig

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(organization.RunnerConfig)
	hasChanges := false

	if state == nil {
		state = new(RunnerConfigModel) // dummy initialization
		hasChanges = true              // must have changes because before is null and after is not
	}

	if plan.Mode != state.Mode {
		retVal.SetMode(plan.Mode.ValueStringPointer())
		hasChanges = true
	}

	if innerProperty, hasInnerChanges := helpers.TfListStringConverter(plan.Groups, state.Groups); hasInnerChanges {
		retVal.SetGroups(innerProperty)
		hasChanges = true
	}

	if plan.IsOverridable != state.IsOverridable {
		retVal.SetIsOverridable(plan.IsOverridable.ValueBoolPointer())
		hasChanges = true
	}

	return retVal, hasChanges
}

func suppressedResourcesConverter(plan *SuppressedResourcesModel, state *SuppressedResourcesModel, converterType commons.ConverterType) (*organization.SuppressedResources, bool) {
	var retVal *organization.SuppressedResources

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(organization.SuppressedResources)
	hasChanges := false

	if state == nil {
		state = new(SuppressedResourcesModel) // dummy initialization
		hasChanges = true                     // must have changes because before is null and after is not
	}

	if innerProperty, hasInnerChanges := tagPropertiesConverter(plan.ManagedByTags, state.ManagedByTags, converterType); hasInnerChanges {
		retVal.SetManagedByTags(innerProperty)
		hasChanges = true
	}

	return retVal, hasChanges
}

func tagPropertiesConverter(plan []*TagPropertiesModel, state []*TagPropertiesModel, converterType commons.ConverterType) ([]*organization.TagProperties, bool) {
	var retVal []*organization.TagProperties
	hasChanged := false

	if reflect.DeepEqual(plan, state) == false {
		hasChanged = true

		if plan != nil {
			retVal = make([]*organization.TagProperties, 0)

			for _, tfTag := range plan {
				apiTag := tagConverter(tfTag)
				retVal = append(retVal, apiTag)
			}
		}
	}

	return retVal, hasChanged
}

func tagConverter(plan *TagPropertiesModel) *organization.TagProperties {
	retVal := new(organization.TagProperties)

	retVal.SetKey(plan.Key.ValueStringPointer())
	retVal.SetValue(plan.Value.ValueStringPointer())

	return retVal
}

func reportConfigurationsConverter(plan []*ReportConfigurationModel, state []*ReportConfigurationModel, converterType commons.ConverterType) ([]*organization.ReportConfiguration, bool) {
	var retVal []*organization.ReportConfiguration
	hasChanged := false

	if reflect.DeepEqual(plan, state) == false {
		hasChanged = true

		if plan != nil {
			retVal = make([]*organization.ReportConfiguration, 0)

			for _, r := range plan {
				rc := reportConfigurationConverter(r)
				retVal = append(retVal, rc)
			}
		}
	}

	return retVal, hasChanged
}

func reportConfigurationConverter(plan *ReportConfigurationModel) *organization.ReportConfiguration {
	retVal := new(organization.ReportConfiguration)

	retVal.SetType(plan.Type.ValueStringPointer())

	recipients := recipientsConverter(plan.Recipients)
	retVal.SetRecipients(recipients)

	retVal.SetEnabled(plan.Enabled.ValueBoolPointer())

	return retVal
}

func recipientsConverter(plan *ReportRecipientsModel) *organization.ReportRecipients {
	retVal := new(organization.ReportRecipients)

	retVal.SetAllAdmins(plan.AllAdmins.ValueBoolPointer())
	retVal.SetEmailAddresses(helpers.TfListToStringPointerSlice(plan.EmailAddresses))
	retVal.SetEmailAddressesToExclude(helpers.TfListToStringPointerSlice(plan.EmailAddressesToExclude))

	return retVal
}
