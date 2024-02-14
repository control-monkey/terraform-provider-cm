package notfication_endpoint

import "github.com/hashicorp/terraform-plugin-framework/types"

type ResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Protocol types.String `tfsdk:"protocol"`
	Url      types.String `tfsdk:"url"`
}
