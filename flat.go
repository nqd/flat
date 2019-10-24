package flat

import (
	"errors"
	"strings"

	"github.com/imdario/mergo"
)

// isObject checks wether the given input is an object
func isObject(in interface{}) bool {
	_, ok := in.(map[string]interface{})
	return ok
}

// isArray checks wether the given input is an array
func isArray(in interface{}) bool {
	_, ok := in.([]interface{})
	return ok
}

// recursivelyUnflattenArray recusively unflattens an array
func recursivelyUnflattenArray(in []interface{}, opts *Options) ([]interface{}, error) {
	out := make([]interface{}, len(in))

	err := errors.New("")
	for key, value := range in {
		if isObject(value) {
			object, _ := value.(map[string]interface{})

			out[key], err = recursivelyUnflattenObject(object, opts)
			if err != nil {
				return nil, err
			}
		} else if isArray(value) {
			array, _ := value.([]interface{})

			out[key], err = recursivelyUnflattenArray(array, opts)
			if err != nil {
				return nil, err
			}
		} else {
			out[key] = value
		}
	}

	return out, nil
}

// recursivelyUnflattenObject recursively callten the Unflatten function
func recursivelyUnflattenObject(in interface{}, opts *Options) (map[string]interface{}, error) {
	inputMap, ok := in.(map[string]interface{})
	if !ok {
		return nil, errors.New("provided input was not an object")
	}

	out, err := unflatten(inputMap, opts)
	if err != nil {
		return nil, err
	}

	for key, value := range out {
		if isObject(value) {
			out[key], err = recursivelyUnflattenObject(value, opts)
			if err != nil {
				return nil, err
			}
		} else if isArray(value) {
			array, _ := value.([]interface{})

			out[key], err = recursivelyUnflattenArray(array, opts)
			if err != nil {
				return nil, err
			}
		} else {
			out[key] = value
		}
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
