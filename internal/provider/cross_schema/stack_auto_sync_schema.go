package cross_schema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var AutoSyncSchema = schema.SingleNestedAttribute{
	MarkdownDescription: "Set up auto sync configurations.",
	Optional:            true,
	Attributes: map[string]schema.Attribute{
		"deploy_when_drift_detected": schema.BoolAttribute{
			MarkdownDescription: "If set to `true`, a deployment will start automatically upon detecting a drift or multiple drifts",
			Optional:            true,
		},
	},
}
