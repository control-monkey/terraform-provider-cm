package helpers

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
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

func StringSliceOrNull(vs *[]*string) []types.String {
	var retVal []types.String

	if vs != nil {
		retVal = make([]types.String, 0)

		for _, v := range *vs {
			retVal = append(retVal, StringValueOrNull(v))
		}
	}

	return retVal
}

func StringPointerSliceOrNull(vs []types.String) []*string {
	var retVal []*string

	if vs != nil {
		retVal = make([]*string, 0)

		for _, pattern := range vs {
			retVal = append(retVal, pattern.ValueStringPointer())
		}
	}

	return retVal
}

func TfStringSliceConverter(plan []types.String, state []types.String) (*[]*string, bool) {
	var retVal []*string
	hasChanged := false

	if reflect.DeepEqual(plan, state) == false {
		retVal = StringPointerSliceOrNull(plan)
		hasChanged = true
	}

	return &retVal, hasChanged
}

func EnumForDocs(stringArray []string) string {
	return fmt.Sprintf("[%s]", strings.Join(stringArray, ", "))
}
