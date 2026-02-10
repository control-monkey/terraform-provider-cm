package external_credential_data

import (
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	sdkExternalCredential "github.com/control-monkey/controlmonkey-sdk-go/services/external_credential"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateStateAfterRead(apiEntity *sdkExternalCredential.ExternalCredential, state *ResourceModel, diagnostics *diag.Diagnostics) {
	state.ID = types.StringValue(controlmonkey.StringValue(apiEntity.ID))
	state.Name = types.StringValue(controlmonkey.StringValue(apiEntity.Name))
}
