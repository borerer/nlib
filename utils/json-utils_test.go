package utils

import "testing"

func TestToJSONString(t *testing.T) {
	v := map[string]interface{}{
		"a": "b",
	}
	s := ToJSONString(v)
	if s != "{\"a\":\"b\"}" {
		t.Fatal(s)
	}
}
