package util

import (
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
)

func CycleJSON(src any, tgt any) error {
	b, _ := jsoniter.Marshal(src)
	return FromJSON(b, tgt)
}

func ToJSON(data any) ([]byte, error) {
	jsoned, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return jsoned, nil
}

func FromJSON(msg json.RawMessage, tgt any) error {
	return jsoniter.Unmarshal(msg, tgt)
}

func ToJSONString(x interface{}) string {
	return string(ToJSONBytes(x, true))
}

func ToJSONCompact(x interface{}) string {
	return string(ToJSONBytes(x, false))
}

func ToJSONBytes(x interface{}, indent bool) []byte {
	if indent {
		b, _ := json.MarshalIndent(x, "", "  ")
		return b
	}
	b, _ := json.Marshal(x)
	return b
}
