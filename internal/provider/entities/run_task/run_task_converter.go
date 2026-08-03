package run_task

import (
	sdkRunTask "github.com/control-monkey/controlmonkey-sdk-go/services/run_task"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons"
)

func Converter(plan *ResourceModel, state *ResourceModel, converterType commons.ConverterType) (*sdkRunTask.RunTask, bool) {
	var retVal *sdkRunTask.RunTask

	if plan == nil {
		if state == nil {
			return nil, false // both are the same, no changes
		} else {
			return nil, true // before had data, after update is null -> update to null
		}
	}

	retVal = new(sdkRunTask.RunTask)
	hasChanges := false

	if state == nil {
		state = new(ResourceModel) // dummy initialization
		hasChanges = true          // must have changes because before is null and after is not
	}

	if plan.Name != state.Name {
		retVal.SetName(plan.Name.ValueStringPointer())
		hasChanges = true
	}
	if plan.Url != state.Url {
		retVal.SetUrl(plan.Url.ValueStringPointer())
		hasChanges = true
	}
	if plan.IsEnabled != state.IsEnabled {
		retVal.SetIsEnabled(plan.IsEnabled.ValueBoolPointer())
		hasChanges = true
	}
	if !plan.HmacKey.IsNull() && !plan.HmacKey.IsUnknown() {
		retVal.SetHmacKey(plan.HmacKey.ValueStringPointer())
		hasChanges = true
	} else if plan.HmacKey.IsNull() && !state.HmacKey.IsNull() && plan.HmacKeyWoVersion.IsNull() {
		// Only clear the key when it is being removed outright. When hmac_key_wo_version is set the
		// caller is migrating to the write-only key, and the resource sends that value instead -
		// clearing here would put HmacKey in nullFields and the marshaller rejects both at once.
		retVal.SetHmacKey(nil)
		hasChanges = true
	}

	return retVal, hasChanges
}
