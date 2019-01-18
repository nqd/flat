package flat

import (
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

	if opts.MaxDepth == 0 {
		opts.MaxDepth = 20
	}
	m, err = f("", 0, nested, opts)

	return
}

func f(prefix string, depth int, nested interface{}, opts Options) (flatmap map[string]interface{}, err error) {
	flatmap = make(map[string]interface{})

	switch nested := nested.(type) {
	case map[string]interface{}:
		if depth >= opts.MaxDepth {
			flatmap[prefix] = nested
			return
		}
		if reflect.DeepEqual(nested, map[string]interface{}{}) {
			flatmap[prefix] = nested
			return
		}
		for k, v := range nested {
			// create new key
			newKey := k
			if prefix != "" {
				newKey = prefix + opts.Delimiter + newKey
			}
			fm1, fe := f(newKey, depth+1, v, opts)
			if fe != nil {
				err = fe
				return
			}
			update(flatmap, fm1)
		}
	case []interface{}:
		if opts.Safe {
			flatmap[prefix] = nested
			return
		}
		for i, v := range nested {
			newKey := strconv.Itoa(i)
			if prefix != "" {
				newKey = prefix + opts.Delimiter + newKey
			}
			fm1, fe := f(newKey, depth+1, v, opts)
			if fe != nil {
				err = fe
				return
			}
			update(flatmap, fm1)
		}
	default:
		flatmap[prefix] = nested
	}
	return
}

// merge is the function that update to map with from and key
// example: from is a map
// to = {"hi": "there"}
// from = {"foo": "bar"}
// new to = {"hi": "there", "foo": "bar"}
// example: from is an empty map
// to = {"hi": "there"}
// from = {}
// key = "world"
// key = {"hi": "there", "world": {}}
func update(to map[string]interface{}, from map[string]interface{}) {
	for kt, vt := range from {
		to[kt] = vt
	}
}
