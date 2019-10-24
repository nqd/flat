package flat

import (
	"reflect"
	"testing"
)

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
		// use flat syntax in nested object
		{
			map[string]interface{}{
				"pew": map[string]interface{}{
					"woop.party": "rainbows!",
				},
			},
			nil,
			map[string]interface{}{
				"pew": map[string]interface{}{
					"woop": map[string]interface{}{
						"party": "rainbows!",
					},
				},
			},
		},
		// use flat syntax in nested array objects
		{
			map[string]interface{}{
				"pew": []interface{}{
					map[string]interface{}{"woop.party": "rainbows!"},
				},
			},
			nil,
			map[string]interface{}{
				"pew": []interface{}{
					map[string]interface{}{
						"woop": map[string]interface{}{
							"party": "rainbows!",
						},
					},
				},
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
