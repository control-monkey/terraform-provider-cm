package commons

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WithNullSetDefault sets a null set default, taking the element type from the attribute's own
// NestedObject. Null and not empty - UpdateStateAfterRead leaves these collections nil.
func WithNullSetDefault(a schema.SetNestedAttribute) schema.SetNestedAttribute {
	a.Default = setdefault.StaticValue(types.SetNull(a.NestedObject.Type()))
	return a
}
