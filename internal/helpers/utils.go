package helpers

import (
	"fmt"
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

func BoolValueOrNull(v *bool) types.Bool {
	var r types.Bool

	if v != nil {
		r = types.BoolValue(controlmonkey.BoolValue(v))
	} else {
		r = types.BoolNull()
	}

	return r
}

func StringValueOrNull(v *string) types.String {
	var r types.String

	if v != nil {
		r = types.StringValue(controlmonkey.StringValue(v))
	} else {
		r = types.StringNull()
	}

	return r
}

func EnumForDocs(stringArray []string) string {
	return fmt.Sprintf("[%s]", strings.Join(stringArray, ", "))
}
