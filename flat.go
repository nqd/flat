package flat

import (
	"errors"
	"strings"

	"github.com/imdario/mergo"
)

// castObject tries to cast to a generic object
func castObject(in interface{}) (map[string]interface{}, bool) {
	casted, ok := in.(map[string]interface{})
	return casted, ok
}

// castArray tries to cast to a generic array
func castArray(in interface{}) ([]interface{}, bool) {
	casted, ok := in.([]interface{})
	return casted, ok
}

// recursivelyUnflattenArray recusively unflattens an array
func recursivelyUnflattenArray(in []interface{}, opts *Options) ([]interface{}, error) {
	out := make([]interface{}, len(in))

	err := errors.New("")
	for key, value := range in {
		if object, ok := castObject(value); ok {

			out[key], err = recursivelyUnflattenObject(object, opts)
			if err != nil {
				return nil, err
			}

			continue
		}
		if array, ok := castArray(value); ok {
			out[key], err = recursivelyUnflattenArray(array, opts)
			if err != nil {
				return nil, err
			}

			continue

		}
		out[key] = value
	}

	return out, nil
}

// recursivelyUnflattenObject recursively callten the Unflatten function
func recursivelyUnflattenObject(in interface{}, opts *Options) (map[string]interface{}, error) {
	inputMap, ok := in.(map[string]interface{})
	if !ok {
		return nil, errors.New("provided input was not an object")
	}

	// flatten current depth
	out, err := unflatten(inputMap, opts)
	if err != nil {
		return nil, err
	}

	// check whether some fields if they need to be recursed (type of array or "objects")
	// and recurse them if neccessary
	for key, value := range out {
		if _, ok := castObject(value); ok {
			out[key], err = recursivelyUnflattenObject(value, opts)
			if err != nil {
				return nil, err
			}

			continue
		}

		if array, ok := castArray(value); ok {
			out[key], err = recursivelyUnflattenArray(array, opts)
			if err != nil {
				return nil, err
			}

			continue
		}

		out[key] = value
	}

	return out, nil
}

// Options the flatten options.
// By default: Demiliter = "."
type Options struct {
	Delimiter string
	Safe      bool
	MaxDepth  int
}

// Unflatten the map, it returns a nested map of a map
// By default, the flatten has Delimiter = "."
func Unflatten(flat map[string]interface{}, opts *Options) (nested map[string]interface{}, err error) {
	if opts == nil {
		opts = &Options{
			Delimiter: ".",
		}
	}
	nested, err = recursivelyUnflattenObject(flat, opts)
	return
}

func unflatten(flat map[string]interface{}, opts *Options) (nested map[string]interface{}, err error) {
	nested = make(map[string]interface{})

	for k, v := range flat {
		temp := uf(k, v, opts).(map[string]interface{})
		err = mergo.Merge(&nested, temp)
		if err != nil {
			return
		}
	}

	return
}

func uf(k string, v interface{}, opts *Options) (n interface{}) {
	n = v

	keys := strings.Split(k, opts.Delimiter)

	for i := len(keys) - 1; i >= 0; i-- {
		temp := make(map[string]interface{})
		temp[keys[i]] = n
		n = temp
	}

	return
}
