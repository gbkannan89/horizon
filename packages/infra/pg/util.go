package pg

import (
	"encoding/json"
	"fmt"
)

func ToJSON(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal json: %w", err)
	}
	return data, nil
}

func FromJSON(data []byte, dest interface{}) error {
	if len(data) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}
	return nil
}

func ToJSONString(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func Ptr[T any](v T) *T { return &v }

func SafeStr(s *string) string {
	if s == nil { return "" }
	return *s
}

func SafeInt64(i *int64) int64 {
	if i == nil { return 0 }
	return *i
}

func SafeFloat64(f *float64) float64 {
	if f == nil { return 0 }
	return *f
}
