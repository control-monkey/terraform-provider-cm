package cm_listvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.List = sizeExactly{}

// sizeExactly validates that list contains at least n elements
// and at most max elements.
type sizeExactly struct {
	n int
}

// Description describes the validation in plain text formatting.
func (v sizeExactly) Description(_ context.Context) string {
	return fmt.Sprintf("list must contain exactly %d elements", v.n)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v sizeExactly) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v sizeExactly) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elems := req.ConfigValue.Elements()

	if len(elems) != v.n {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", len(elems)),
		))
	}
}

// SizeExactly returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a List.
//   - Contains at exactly number of elements.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func SizeExactly(n int) validator.List {
	return sizeExactly{
		n: n,
	}
}
