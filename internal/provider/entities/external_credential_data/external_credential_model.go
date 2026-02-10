package external_credential_data

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	Name   types.String `tfsdk:"name"`
	Vendor types.String `tfsdk:"vendor"`
	ID     types.String `tfsdk:"id"`
}
