package commons

import (
	"github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"strings"
)

type ConverterType string

const (
	CreateConverter ConverterType = "create"
	UpdateConverter ConverterType = "update"

	CreateMerger ConverterType = "createMerger"
	UpdateMerger ConverterType = "updateMerger"
	DeleteMerger ConverterType = "deleteMerger"
)

func IsNotFoundResponseError(err error) bool {
	retVal := strings.Contains(err.Error(), commons.ErrorCodeNotFound)
	return retVal
}

func IsAlreadyExistResponseError(err error) bool {
	retVal := strings.Contains(err.Error(), commons.ErrorCodeAlreadyExist)
	return retVal
}

func DoesErrorContains(err error, s string) bool {
	retVal := strings.Contains(err.Error(), s)
	return retVal
}
