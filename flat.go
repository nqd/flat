package flat

import "log"

type Options struct {
	Delimiter string
	Safe      bool
	Object    bool
	MaxDepth  int
}

func Flatten(nested map[string]interface{}, opts Options) (m map[string]interface{}, err error) {
	m, err = flatten("", nested, opts)
	return
}

func flatten(prefix string, nested interface{}, opts Options) (m map[string]interface{}, err error) {
	switch nested := nested.(type) {
	case map[string]interface{}:
		// map
		for k, v := range nested {
			newKey := k
			if prefix != "" {
				newKey = prefix + "." + k
			}
			switch v.(type) {
			case map[string]interface{}, []interface{}:
				m, err = flatten(newKey, v, opts)
			default:
				m = make(map[string]interface{})
				m[newKey] = v
			}
		}
	case []interface{}:
		log.Println("slice")
		// slice
	default:
		log.Println("error")
	}
	return
}
