package notfication_endpoint

import (
	sdkNotification "github.com/control-monkey/controlmonkey-sdk-go/services/notification"
	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
)

func UpdateStateAfterRead(res *sdkNotification.Endpoint, state *ResourceModel) {
	state.Name = helpers.StringValueOrNull(res.Name)
	state.Protocol = helpers.StringValueOrNull(res.Protocol)
	state.Url = helpers.StringValueOrNull(res.Url)
}
