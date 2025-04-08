package jsonhelper

import (
	"encoding/json"
)

func ToString(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
