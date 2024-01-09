package commons

import (
	cm_stringvalidators "github.com/control-monkey/terraform-provider-cm/internal/provider/validators/string"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ValidateUniqueNotEmptyListWithNoBlankValues() []validator.List {
	return []validator.List{listvalidator.SizeAtLeast(1), listvalidator.UniqueValues(), listvalidator.ValueStringsAre(cm_stringvalidators.NotBlank())}
}
