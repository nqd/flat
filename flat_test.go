package flat

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestFlatten(t *testing.T) {
	tests := []struct {
		given   string
		options *Options
		want    map[string]interface{}
	}{
		// test with different primitives
		// String: 'world',
		// Number: 1234.99,
		// Boolean: true,
		// null: null,
		{
			`{"hello": "world"}`,
			nil,
			map[string]interface{}{"hello": "world"},
		},
		{
			`{"hello": 1234.99}`,
			nil,
			map[string]interface{}{"hello": 1234.99},
		},
		{
			`{"hello": true}`,
			nil,
			map[string]interface{}{"hello": true},
		},
		{
			`{"hello": null}`,
			nil,
			map[string]interface{}{"hello": nil},
		},
		// nested once
		{
			`{"hello":{}}`,
			nil,
			map[string]interface{}{"hello": map[string]interface{}{}},
		},
		{
			`{"hello":{"world":"good morning"}}`,
			nil,
			map[string]interface{}{"hello.world": "good morning"},
		},
		{
			`{"hello":{"world":1234.99}}`,
			nil,
			map[string]interface{}{"hello.world": 1234.99},
		},
		{
			`{"hello":{"world":true}}`,
			nil,
			map[string]interface{}{"hello.world": true},
		},
		{
			`{"hello":{"world":null}}`,
			nil,
			map[string]interface{}{"hello.world": nil},
		},
		// empty slice
		{
			`{"hello":{"world":[]}}`,
			nil,
			map[string]interface{}{"hello.world": []interface{}{}},
		},
		// slice
		{
			`{"hello":{"world":["one","two"]}}`,
			nil,
			map[string]interface{}{
				"hello.world.0": "one",
				"hello.world.1": "two",
			},
		},
		// nested twice
		{
			`{"hello":{"world":{"again":"good morning"}}}`,
			nil,
			map[string]interface{}{"hello.world.again": "good morning"},
		},
		// multiple keys
		{
			`{
				"hello": {
					"lorem": {
						"ipsum":"again",
						"dolor":"sit"
					}
				},
				"world": {
					"lorem": {
						"ipsum":"again",
						"dolor":"sit"
					}
				}
			}`,
			nil,
			map[string]interface{}{
				"hello.lorem.ipsum": "again",
				"hello.lorem.dolor": "sit",
				"world.lorem.ipsum": "again",
				"world.lorem.dolor": "sit"},
		},
		// empty object
		{
			`{"hello":{"empty":{"nested":{}}}}`,
			nil,
			map[string]interface{}{"hello.empty.nested": map[string]interface{}{}},
		},
		// custom delimiter
		{
			`{"hello":{"world":{"again":"good morning"}}}`,
			&Options{
				Delimiter: ":",
				MaxDepth:  20,
			},
			map[string]interface{}{"hello:world:again": "good morning"},
		},
		// custom depth
		{
			`{
				"hello": {
					"world": {
						"again": "good morning"
					}
				},
				"lorem": {
					"ipsum": {
						"dolor": "good evening"
					}
				}
			}
			`,
			&Options{
				MaxDepth:  2,
				Delimiter: ".",
			},
			map[string]interface{}{
				"hello.world": map[string]interface{}{"again": "good morning"},
				"lorem.ipsum": map[string]interface{}{"dolor": "good evening"},
			},
		},
		// custom safe = true
		{
			`{"hello":{"world":["one","two"]}}`,
			&Options{
				Safe:      true,
				Delimiter: ".",
			},
			map[string]interface{}{
				"hello.world": []interface{}{"one", "two"},
			},
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

func TestUnflatten(t *testing.T) {
	tests := []struct {
		flat    map[string]interface{}
		options *Options
		want    map[string]interface{}
	}{
		{
			map[string]interface{}{"hello": "world"},
			nil,
			map[string]interface{}{"hello": "world"},
		},
		{
			map[string]interface{}{"hello": 1234.56},
			nil,
			map[string]interface{}{"hello": 1234.56},
		},
		{
			map[string]interface{}{"hello": true},
			nil,
			map[string]interface{}{"hello": true},
		},
		// nested twice
		{
			map[string]interface{}{"hello.world.again": "good morning"},
			nil,
			map[string]interface{}{
				"hello": map[string]interface{}{
					"world": map[string]interface{}{
						"again": "good morning",
					},
				},
			},
		},
		// multiple keys
		{
			map[string]interface{}{
				"hello.lorem.ipsum": "again",
				"hello.lorem.dolor": "sit",
				"world.lorem.ipsum": "again",
				"world.lorem.dolor": "sit",
				"world":             map[string]interface{}{"greet": "hello"},
			},
			nil,
			map[string]interface{}{
				"hello": map[string]interface{}{
					"lorem": map[string]interface{}{
						"ipsum": "again",
						"dolor": "sit",
					},
				},
				"world": map[string]interface{}{
					"greet": "hello",
					"lorem": map[string]interface{}{
						"ipsum": "again",
						"dolor": "sit",
					},
				},
			},
		},
		// nested objects do not clobber each other
		{
			map[string]interface{}{
				"foo.bar": map[string]interface{}{"t": 123},
				"foo":     map[string]interface{}{"k": 456},
			},
			nil,
			map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": map[string]interface{}{
						"t": 123,
					},
					"k": 456,
				},
			},
		},
		// custom delimiter
		{
			map[string]interface{}{
				"hello world again": "good morning",
			},
			&Options{
				Delimiter: " ",
			},
			map[string]interface{}{
				"hello": map[string]interface{}{
					"world": map[string]interface{}{
						"again": "good morning",
					},
				},
			},
		},
		// do not overwrite
		{
			map[string]interface{}{
				"travis":           "true",
				"travis_build_dir": "/home/foo",
			},
			&Options{
				Delimiter: "_",
			},
			map[string]interface{}{
				"travis": "true",
			},
		},
		// todo
		// overwrite true
		// {
		// 	map[string]interface{}{
		// 		"travis":           "true",
		// 		"travis_build_dir": "/home/foo",
		// 	},
		// 	Options{
		// 		Delimiter: "_",
		// 		Overwrite: true,
		// 	},
		// 	map[string]interface{}{
		// 		"travis": map[string]interface{}{
		// 			"build": map[string]interface{}{
		// 				"dir": "/home/foo",
		// 			},
		// 		},
		// 	},
		// },
	}
	for i, test := range tests {
		got, err := Unflatten(test.flat, test.options)
		if err != nil {
			t.Errorf("%d: failed to unflatten: %v", i+1, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d: mismatch, got: %v want: %v", i+1, got, test.want)
		}
	}
}
