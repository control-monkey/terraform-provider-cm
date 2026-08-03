package run_task

import "github.com/hashicorp/terraform-plugin-framework/types"

type ResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Url                 types.String `tfsdk:"url"`
	IsEnabled           types.Bool   `tfsdk:"is_enabled"`
	HmacKey             types.String `tfsdk:"hmac_key"`
	HmacKeyWo           types.String `tfsdk:"hmac_key_wo"`
	HmacKeyWoVersion    types.Int64  `tfsdk:"hmac_key_wo_version"`
	IsHmacKeyConfigured types.Bool   `tfsdk:"is_hmac_key_configured"`
}
