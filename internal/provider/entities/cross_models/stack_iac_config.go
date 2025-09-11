package cross_models

import (
	"github.com/control-monkey/controlmonkey-sdk-go/services/cross_models"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//region Model

type IacConfigModel struct {
	TerraformVersion   types.String `tfsdk:"terraform_version"`
	TerragruntVersion  types.String `tfsdk:"terragrunt_version"`
	OpentofuVersion    types.String `tfsdk:"opentofu_version"`
	IsTerragruntRunAll types.Bool   `tfsdk:"is_terragrunt_run_all"`
	VarFiles           types.List   `tfsdk:"var_files"`
}

//endregion

//region Create/Update Converter

func IacConfigConverter(plan *IacConfigModel, state *IacConfigModel, converterType commons.ConverterType) (*cross_models.IacConfig, bool) {
	var retVal *cross_models.IacConfig

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(cross_models.IacConfig)
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
	if plan.IsTerragruntRunAll != state.IsTerragruntRunAll {
		retVal.SetIsTerragruntRunAll(plan.IsTerragruntRunAll.ValueBoolPointer())
		hasChanges = true
	}
	if innerProperty, hasInnerChanges := helpers.TfListStringConverter(plan.VarFiles, state.VarFiles); hasInnerChanges {
		retVal.SetVarFiles(innerProperty)
		hasChanges = true
	}

	return retVal, hasChanges
}

//endregion

//region Update State After Read

func UpdateStateAfterReadIacConfig(iacConfig *cross_models.IacConfig) IacConfigModel {
	var retVal IacConfigModel

	retVal.TerraformVersion = helpers.StringValueOrNull(iacConfig.TerraformVersion)
	retVal.TerragruntVersion = helpers.StringValueOrNull(iacConfig.TerragruntVersion)
	retVal.OpentofuVersion = helpers.StringValueOrNull(iacConfig.OpentofuVersion)
	retVal.IsTerragruntRunAll = helpers.BoolValueOrNull(iacConfig.IsTerragruntRunAll)
	retVal.VarFiles = helpers.StringPointerSliceToTfList(iacConfig.VarFiles)

	return retVal
}

//endregion
