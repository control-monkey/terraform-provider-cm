package helpers

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mpvl/unique"
	"reflect"
	"strconv"
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

func StringPointerSliceToTfList(vs []*string) types.List {
	var retVal types.List

	if vs != nil {
		var values []attr.Value

		for _, v := range vs {
			values = append(values, types.StringValue(*v))
		}

		retVal = types.ListValueMust(types.StringType, values)
	} else {
		retVal = types.ListNull(types.StringType)
	}

	return retVal
}

func TfListToStringPointerSlice(vs types.List) []*string {
	var retVal []*string

	if vs.IsNull() == false {
		retVal = make([]*string, 0)

		for _, pattern := range vs.Elements() {
			val := TrimDoubleQuotesIfPresent(pattern.String())
			retVal = append(retVal, &val)
		}
	}

	return retVal
}

func TrimDoubleQuotesIfPresent(s string) string {
	retVal := s

	if len(retVal) > 0 && retVal[0] == '"' {
		retVal = retVal[1:]
	}
	if len(retVal) > 0 && retVal[len(retVal)-1] == '"' {
		retVal = retVal[:len(retVal)-1]
	}

	return retVal
}

func TfListStringConverter(plan types.List, state types.List) ([]*string, bool) {
	var retVal []*string
	hasChanged := false

	if reflect.DeepEqual(plan.Elements(), state.Elements()) == false {
		elements := divideTfListToTfElements(plan)
		retVal = stringPointerSliceOrNull(elements)
		hasChanged = true
	}

	return retVal, hasChanged
}

func DoesTfListContainsEmptyValue(tfValues types.List) bool {
	retVal := false

	elements := divideTfListToTfElements(tfValues)

	for _, v := range elements {
		if v.IsNull() || IsBlank(v.ValueString()) {
			retVal = true
			break
		}
	}

	return retVal
}

func IsTfStringSliceUnique(tfList types.List) bool {
	var retVal bool

	elements := divideTfListToTfElements(tfList)
	retVal = unique.StringsAreUnique(stringValuesSliceFromTfSlice(elements))

	return retVal
}

func CheckAndGetIfNumericString(s string) (bool, float64) {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return false, i
	} else {
		return true, i
	}
}

func EnumForDocs(stringArray []string) string {
	return fmt.Sprintf("[%s]", strings.Join(stringArray, ", "))
}

//region Unexported

func divideTfListToTfElements(tfValues types.List) []types.String {
	retVal := make([]types.String, 0, len(tfValues.Elements()))
	tfValues.ElementsAs(nil, &retVal, false) // ctx is not used in the inner logic

	return retVal
}

// stringValuesSliceFromTfSlice
// Note - types.String.ValueString() can throw exception in case of having Null
func stringValuesSliceFromTfSlice(vs []types.String) []string {
	var retVal []string

	if vs != nil {
		retVal = make([]string, 0)

		for _, v := range vs {
			retVal = append(retVal, v.ValueString())
		}
	}

	return retVal
}

func stringPointerSliceOrNull(vs []types.String) []*string {
	var retVal []*string

	if vs != nil {
		retVal = make([]*string, 0)

		for _, pattern := range vs {
			retVal = append(retVal, pattern.ValueStringPointer())
		}
	}

	return retVal
}

//endregion
