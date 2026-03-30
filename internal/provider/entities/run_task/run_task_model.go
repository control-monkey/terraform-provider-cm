package run_task

import "github.com/hashicorp/terraform-plugin-framework/types"

type ResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Url                 types.String `tfsdk:"url"`
	IsEnabled           types.Bool   `tfsdk:"is_enabled"`
	HmacKey             types.String `tfsdk:"hmac_key"`
	IsHmacKeyConfigured types.Bool   `tfsdk:"is_hmac_key_configured"`
}
