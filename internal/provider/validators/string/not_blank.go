package cm_stringvalidator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strings"
)

var _ validator.String = notBlankValidator{}

// notBlankValidator validates that the value does not match one of the values.
type notBlankValidator struct {
}

func (v notBlankValidator) Description(_ context.Context) string {
	return "value must not be empty"
}

func (v notBlankValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v notBlankValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	if value.IsNull() || strings.TrimSpace(value.ValueString()) == "" {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value.String(),
		))
	}
}

// NotBlank checks that the String is not blank
func NotBlank() validator.String {
	return notBlankValidator{}
}
