package flat

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestFlatten(t *testing.T) {
	tests := []struct {
		given   string
		want    map[string]interface{}
		options Options
	}{
		{
			`{"hello":{"world":"good morning"}}`,
			map[string]interface{}{
				"hello.world": "good morning",
			},
			Options{},
		},
	}
	for i, test := range tests {
		var given interface{}
		err := json.Unmarshal([]byte(test.given), &given)
		if err != nil {
			t.Errorf("%d: failed to unmarshal test: %v", i+1, err)
		}
		got, err := Flatten(given.(map[string]interface{}), test.options)
		if err != nil {
			t.Errorf("%d: failed to flatten: %v", i+1, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d: mismatch, got: %v want: %v", i+1, got, test.want)
		}
	}
}
