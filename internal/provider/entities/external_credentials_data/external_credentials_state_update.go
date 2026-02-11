package external_credentials_data

import (
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	sdkExternalCredentials "github.com/control-monkey/controlmonkey-sdk-go/services/external_credentials"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateStateAfterRead(apiEntity *sdkExternalCredentials.ExternalCredentials, state *ResourceModel, diagnostics *diag.Diagnostics) {
	state.ID = types.StringValue(controlmonkey.StringValue(apiEntity.ID))
	state.Name = types.StringValue(controlmonkey.StringValue(apiEntity.Name))
}
