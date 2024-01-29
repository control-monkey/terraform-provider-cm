package helpers

import "github.com/hashicorp/terraform-plugin-framework/types"

func IsTrue(v types.Bool) bool {
	retVal := v.IsNull() == false && v.IsUnknown() == false && v.ValueBool()
	return retVal
}
