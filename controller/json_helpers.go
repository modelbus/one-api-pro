package controller

import "encoding/json"

// jsonMarshalString is a tiny helper that marshals v to a JSON string
// and returns an empty string on error. Used in places that need a
// plain string for storage (system_settings.value) without dragging
// in fmt everywhere.
func jsonMarshalString(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// jsonUnmarshalString unmarshals a JSON string into v.
func jsonUnmarshalString(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}
