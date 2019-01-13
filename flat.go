package flat

import (
	"log"
	"reflect"
	"strconv"
)

type Options struct {
	Delimiter string
	Safe      bool
	Object    bool
	MaxDepth  int
}

func Flatten(nested map[string]interface{}, opts Options) (m map[string]interface{}, err error) {
	// construct default value
	if opts.Delimiter == "" {
		opts.Delimiter = "."
	}
	m, err = flatten("", nested, opts)
	return
}

func flatten(prefix string, nested interface{}, opts Options) (m map[string]interface{}, err error) {
	m = make(map[string]interface{})
	switch nested := nested.(type) {
	case map[string]interface{}:
		log.Println("map", nested)
		// map
		for k, v := range nested {
			// create new key
			newKey := k
			if prefix != "" {
				newKey = prefix + opts.Delimiter + k
			}
			switch v.(type) {
			case map[string]interface{}, []interface{}:
				temp, fErr := flatten(newKey, v, opts)
				if fErr != nil {
					err = fErr
					return
				}
				// empty map {}
				if reflect.DeepEqual(temp, map[string]interface{}{}) {
					m[newKey] = temp
				} else {
					// merge temp to m
					for kt, vt := range temp {
						m[kt] = vt
					}
				}
			default:
				m[newKey] = v
			}
		}
	case []interface{}:
		log.Println("slice", nested)
		// slice
		for i, v := range nested {
			newKey := strconv.Itoa(i)
			if prefix != "" {
				newKey = prefix + opts.Delimiter + newKey
			}
			switch v.(type) {
			case map[string]interface{}, []interface{}:
				temp, fErr := flatten(newKey, v, opts)
				if fErr != nil {
					err = fErr
					return
				}
				// empty map {}
				if reflect.DeepEqual(temp, map[string]interface{}{}) {
					m[newKey] = temp
				} else {
					// merge temp to m
					for kt, vt := range temp {
						m[kt] = vt
					}
				}
			default:
				m[newKey] = v
			}
		}
	default:
		log.Println("error")
	}
	return
}
