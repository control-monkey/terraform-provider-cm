package external_credential_data

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	Vendor types.String `tfsdk:"vendor"`
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
}
