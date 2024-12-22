package helpers

import (
	"encoding/json"
	"strings"
)

func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

func NormalizeJsonArrayString(s string) string {
	var retVal string
	var jsonObjectStructure []map[string]interface{}

	err := json.Unmarshal([]byte(s), &jsonObjectStructure)
	if err != nil {
		return retVal
	}

	if ret, err := json.Marshal(jsonObjectStructure); err == nil {
		retVal = string(ret)
	}

	return retVal
}

func NormalizeJsonString(s string) string {
	var retVal string
	var jsonObjectStructure map[string]interface{}

	err := json.Unmarshal([]byte(s), &jsonObjectStructure)
	if err != nil {
		return retVal
	}

	if ret, err := json.Marshal(jsonObjectStructure); err == nil {
		retVal = string(ret)
	}

	return retVal
}
