package commons

import (
	"github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"strings"
)

type ConverterType string

const (
	CreateConverter ConverterType = "create"
	UpdateConverter ConverterType = "update"
)

func IsNotFoundResponseError(err error) bool {
	retVal := strings.Contains(err.Error(), commons.ErrorCodeNotFound)
	return retVal
}
