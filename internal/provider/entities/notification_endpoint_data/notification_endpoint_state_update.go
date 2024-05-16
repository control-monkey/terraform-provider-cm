package notification_endpoint_data

import (
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	sdkNotification "github.com/control-monkey/controlmonkey-sdk-go/services/notification"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateStateAfterRead(apiEntity *sdkNotification.Endpoint, state *ResourceModel, diagnostics *diag.Diagnostics) {
	state.ID = types.StringValue(controlmonkey.StringValue(apiEntity.ID))
	state.Name = types.StringValue(controlmonkey.StringValue(apiEntity.Name))
}
