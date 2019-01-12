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
	switch nested.(type) {
	case map[string]interface{}:
		// map
		log.Println("flatten a map")
		for k, v := range nested.(map[string]interface{}) {
			newKey := prefix + "." + k
			switch v.(type) {
			case map[string]interface{}, []interface{}:
				m, err = flatten(newKey, v, opts)
			default:
				log.Println("default")
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
