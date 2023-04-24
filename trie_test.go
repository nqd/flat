package flat

import (
	"reflect"
	"strings"
	"testing"
)

func TestTrie(t *testing.T) {
	tests := []struct {
		data map[string]interface{}
		want map[string]interface{}
	}{
		{
			map[string]interface{}{"hello": "world"},
			map[string]interface{}{"hello": "world"},
		},
		{
			map[string]interface{}{"hello.world.again": "good morning"},
			map[string]interface{}{
				"hello": map[string]interface{}{
					"world": map[string]interface{}{
						"again": "good morning",
					},
				},
			},
		},
		{
			map[string]interface{}{"a.0.0": 1, "a.0.1": 2, "a.1.0": 21, "a.1.1": 22},
			map[string]interface{}{
				"a": []interface{}{
					[]interface{}{1, 2},
					[]interface{}{21, 22},
				},
			},
		},
		{
			map[string]interface{}{"a.0.0": "1", "a.0.1": "2", "a.1.0": "21", "a.1.1": "22"},
			map[string]interface{}{
				"a": []interface{}{
					[]interface{}{"1", "2"},
					[]interface{}{"21", "22"},
				},
			},
		},
		{
			map[string]interface{}{"a.0": "1", "a.1": "2", "b": "21"},
			map[string]interface{}{"a": []interface{}{"1", "2"}, "b": "21"},
		},
		{
			map[string]interface{}{"a.b.0": "1", "a.b.1": "2", "c": "21"},
			map[string]interface{}{"a": map[string]interface{}{"b": []interface{}{"1", "2"}}, "c": "21"},
		},
		{
			map[string]interface{}{"a.0.b.0": "1", "a.0.b.1": "2", "a.1.b.0": "3", "a.1.b.1": "4", "c": "21"},
			map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{
						"b": []interface{}{"1", "2"},
					},
					map[string]interface{}{
						"b": []interface{}{"3", "4"},
					},
				},
				"c": "21",
			},
		},
		{
			map[string]interface{}{
				"a.0.b.0.d": "1", "a.0.b.0.e": "2",
				"a.0.b.1.d": "3", "a.0.b.1.e": "4",
				"a.1.b.0.d": "11", "a.1.b.0.e": "12",
				"a.1.b.0.f": "13", "a.1.b.0.g": "14",
				"a.1.b.1.d": "15", "a.1.b.1.e": "16",
				"a.1.b.1.f": "17", "a.1.b.1.g": "18",
				"c": "21"},
			map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{
						"b": []interface{}{
							map[string]interface{}{"d": "1", "e": "2"},
							map[string]interface{}{"d": "3", "e": "4"},
						},
					},
					map[string]interface{}{
						"b": []interface{}{
							map[string]interface{}{"d": "11", "e": "12", "f": "13", "g": "14"},
							map[string]interface{}{"d": "15", "e": "16", "f": "17", "g": "18"},
						},
					},
				},
				"c": "21",
			},
		},
	}

	for i, test := range tests {
		root := &TrieNode{}
		for k, v := range test.data {
			root.insert(strings.Split(k, "."), v)
		}

		got := root.unflatten()
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d: mismatch, got: %v want: %v", i+1, got, test.want)
		}
	}
}
