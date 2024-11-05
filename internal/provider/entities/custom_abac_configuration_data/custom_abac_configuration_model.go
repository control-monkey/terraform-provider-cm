package custom_abac_configuration_data

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}
