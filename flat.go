package flat

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/imdario/mergo"
	"go.mongodb.org/mongo-driver/bson"
)

// Options the flatten options.
// By default: Demiliter = "."
type Options struct {
	Delimiter string
	Safe      bool
	MaxDepth  int
}

// Flatten the map, it returns a map one level deep
// regardless of how nested the original map was.
// By default, the flatten has Delimiter = ".", and
// no limitation of MaxDepth
func Flatten(nested bson.M, opts *Options) (m bson.M, err error) {
	if opts == nil {
		opts = &Options{
			Delimiter: ".",
		}
	}

	m, err = flatten("", 0, nested, opts)

	return
}

func flatten(prefix string, depth int, nested interface{}, opts *Options) (flatmap bson.M, err error) {
	flatmap = bson.M{}

	switch nested := nested.(type) {
	case bson.M:
		if opts.MaxDepth != 0 && depth >= opts.MaxDepth {
			flatmap[prefix] = nested
			return
		}
		if reflect.DeepEqual(nested, bson.M{}) {
			flatmap[prefix] = nested
			return
		}
		for k, v := range nested {
			// create new key
			newKey := k
			if prefix != "" {
				newKey = prefix + opts.Delimiter + newKey
			}
			fm1, fe := flatten(newKey, depth+1, v, opts)
			if fe != nil {
				err = fe
				return
			}
			update(flatmap, fm1)
		}
	case bson.A:
		if opts.Safe {
			flatmap[prefix] = nested
			return
		}
		if reflect.DeepEqual(nested, bson.A{}) {
			flatmap[prefix] = nested
			return
		}
		for i, v := range nested {
			newKey := strconv.Itoa(i)
			if prefix != "" {
				newKey = prefix + opts.Delimiter + newKey
			}
			fm1, fe := flatten(newKey, depth+1, v, opts)
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

// update is the function that update to map with from
// example:
// to = {"hi": "there"}
// from = {"foo": "bar"}
// then, to = {"hi": "there", "foo": "bar"}
func update(to bson.M, from bson.M) {
	for kt, vt := range from {
		to[kt] = vt
	}
}

// Unflatten the map, it returns a nested map of a map
// By default, the flatten has Delimiter = "."
func Unflatten(flat bson.M, opts *Options) (nested bson.M, err error) {
	if opts == nil {
		opts = &Options{
			Delimiter: ".",
		}
	}
	nested, err = unflatten(flat, opts)
	return
}

func unflatten(flat bson.M, opts *Options) (nested bson.M, err error) {
	nested = make(bson.M)

	for k, v := range flat {
		temp := uf(k, v, opts).(bson.M)
		err = mergo.Merge(&nested, temp, func(c *mergo.Config) { c.Overwrite = true })
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
		temp := make(bson.M)
		temp[keys[i]] = n
		n = temp
	}

	return
}
