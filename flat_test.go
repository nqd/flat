package flat

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestFlatten(t *testing.T) {
	tests := []struct {
		given   bson.M
		options *Options
		want    bson.M
	}{
		// test with different primitives
		// String: 'world',
		// Number: 1234.99,
		// Boolean: true,
		// null: null,
		{
			bson.M{"hello": "world"},
			nil,
			bson.M{"hello": "world"},
		},
		{
			bson.M{"hello": 1234.99},
			nil,
			bson.M{"hello": 1234.99},
		},
		{
			bson.M{"hello": true},
			nil,
			bson.M{"hello": true},
		},
		{
			bson.M{"hello": nil},
			nil,
			bson.M{"hello": nil},
		},
		// nested once
		{
			bson.M{"hello": bson.M{}},
			nil,
			bson.M{"hello": bson.M{}},
		},
		{
			bson.M{"hello": bson.M{"world": "good morning"}},
			nil,
			bson.M{"hello.world": "good morning"},
		},
		{
			bson.M{"hello": bson.M{"world": 1234.99}},
			nil,
			bson.M{"hello.world": 1234.99},
		},
		{
			bson.M{"hello": bson.M{"world": true}},
			nil,
			bson.M{"hello.world": true},
		},
		{
			bson.M{"hello": bson.M{"world": nil}},
			nil,
			bson.M{"hello.world": nil},
		},
		// empty slice
		{
			bson.M{"hello": bson.M{"world": bson.A{}}},
			nil,
			bson.M{"hello.world": bson.A{}},
		},
		// slice
		{
			bson.M{"hello": bson.M{"world": bson.A{"one", "two"}}},
			nil,
			bson.M{
				"hello.world.0": "one",
				"hello.world.1": "two",
			},
		},
		// nested twice
		{
			bson.M{"hello": bson.M{"world": bson.M{"again": "good morning"}}},
			nil,
			bson.M{"hello.world.again": "good morning"},
		},
		// multiple keys
		{
			bson.M{
				"hello": bson.M{
					"lorem": bson.M{
						"ipsum": "again",
						"dolor": "sit",
					},
				},
				"world": bson.M{
					"lorem": bson.M{
						"ipsum": "again",
						"dolor": "sit",
					},
				},
			},
			nil,
			bson.M{
				"hello.lorem.ipsum": "again",
				"hello.lorem.dolor": "sit",
				"world.lorem.ipsum": "again",
				"world.lorem.dolor": "sit"},
		},
		// empty object
		{
			bson.M{"hello": bson.M{"empty": bson.M{"nested": bson.M{}}}},
			nil,
			bson.M{"hello.empty.nested": bson.M{}},
		},
		// custom delimiter
		{
			bson.M{"hello": bson.M{"world": bson.M{"again": "good morning"}}},
			&Options{
				Delimiter: ":",
				MaxDepth:  20,
			},
			bson.M{"hello:world:again": "good morning"},
		},
		// custom depth
		{
			bson.M{
				"hello": bson.M{
					"world": bson.M{
						"again": "good morning",
					},
				},
				"lorem": bson.M{
					"ipsum": bson.M{
						"dolor": "good evening",
					},
				},
			},
			&Options{
				MaxDepth:  2,
				Delimiter: ".",
			},
			bson.M{
				"hello.world": bson.M{"again": "good morning"},
				"lorem.ipsum": bson.M{"dolor": "good evening"},
			},
		},
		// custom safe = true
		{
			bson.M{"hello": bson.M{"world": bson.A{"one", "two"}}},
			&Options{
				Safe:      true,
				Delimiter: ".",
			},
			bson.M{
				"hello.world": bson.A{"one", "two"},
			},
		},
	}
	for i, test := range tests {
		got, err := Flatten(test.given, test.options)
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
		flat    bson.M
		options *Options
		want    bson.M
	}{
		{
			bson.M{"hello": "world"},
			nil,
			bson.M{"hello": "world"},
		},
		{
			bson.M{"hello": 1234.56},
			nil,
			bson.M{"hello": 1234.56},
		},
		{
			bson.M{"hello": true},
			nil,
			bson.M{"hello": true},
		},
		// nested twice
		{
			bson.M{"hello.world.again": "good morning"},
			nil,
			bson.M{
				"hello": bson.M{
					"world": bson.M{
						"again": "good morning",
					},
				},
			},
		},
		// multiple keys
		{
			bson.M{
				"hello.lorem.ipsum": "again",
				"hello.lorem.dolor": "sit",
				"world.lorem.ipsum": "again",
				"world.lorem.dolor": "sit",
				"world":             bson.M{"greet": "hello"},
			},
			nil,
			bson.M{
				"hello": bson.M{
					"lorem": bson.M{
						"ipsum": "again",
						"dolor": "sit",
					},
				},
				"world": bson.M{
					"greet": "hello",
					"lorem": bson.M{
						"ipsum": "again",
						"dolor": "sit",
					},
				},
			},
		},
		// nested objects do not clobber each other
		{
			bson.M{
				"foo.bar": bson.M{"t": 123},
				"foo":     bson.M{"k": 456},
			},
			nil,
			bson.M{
				"foo": bson.M{
					"bar": bson.M{
						"t": 123,
					},
					"k": 456,
				},
			},
		},
		// custom delimiter
		{
			bson.M{
				"hello world again": "good morning",
			},
			&Options{
				Delimiter: " ",
			},
			bson.M{
				"hello": bson.M{
					"world": bson.M{
						"again": "good morning",
					},
				},
			},
		},
		// do not overwrite
		{
			bson.M{
				"travis":           "true",
				"travis_build_dir": "/home/foo",
			},
			&Options{
				Delimiter: "_",
			},
			bson.M{
				"travis": "true",
			},
		},
		// keys with nil values
		{
			bson.M{
				"foo.bar": bson.M{"t": nil},
				"foo":     bson.M{"k": nil},
			},
			nil,
			bson.M{
				"foo": bson.M{
					"bar": bson.M{
						"t": nil,
					},
					"k": nil,
				},
			},
		},
		// todo
		// overwrite true
		// {
		// 	bson.M{
		// 		"travis":           "true",
		// 		"travis_build_dir": "/home/foo",
		// 	},
		// 	Options{
		// 		Delimiter: "_",
		// 		Overwrite: true,
		// 	},
		// 	bson.M{
		// 		"travis": bson.M{
		// 			"build": bson.M{
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
