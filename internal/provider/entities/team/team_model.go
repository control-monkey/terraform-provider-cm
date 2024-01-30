package team

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	CustomIdpId types.String `tfsdk:"custom_idp_id"`
}
