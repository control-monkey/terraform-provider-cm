package helpers

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

func BoolValueOrNull(v *bool) types.Bool {
	var r types.Bool

	if v != nil {
		r = types.BoolValue(*v)
	} else {
		r = types.BoolNull()
	}

	return r
}

func Int64ValueOrNull(v *int) types.Int64 {
	var r types.Int64

	if v != nil {
		r = types.Int64Value(int64(*v))
	} else {
		r = types.Int64Null()
	}

	return r
}

func StringValueOrNull(v *string) types.String {
	var r types.String

	if v != nil {
		r = types.StringValue(*v)
	} else {
		r = types.StringNull()
	}

	return r
}

func StringSlice(vs []*string) []types.String {
	var arr []types.String

	for _, v := range vs {
		arr = append(arr, StringValueOrNull(v))
	}

	return arr
}

func EnumForDocs(stringArray []string) string {
	return fmt.Sprintf("[%s]", strings.Join(stringArray, ", "))
}
