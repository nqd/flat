# flat [![Build Status](https://secure.travis-ci.org/nqd/flat.png?branch=master)](http://travis-ci.org/nqd/flat)

Take a golang map and flatten it or unfatten a map with delimited key.

This work inspired by the [nodejs flat package](https://github.com/hughsk/flat/)

## Method

### Flatten

```{go}
in := map[string]interface{}{
   "a": "b",
   "c": map[string]interface{}{
       "d": "e",
       "f": "g",
   },
   "z": [2, 1.4567],
}

out, err := flat.Flatten(in, flat.Options{})
// out = map[string]interface{}{
//     "a": "b",
//     "c.d": "e",
//     "c.f": "g",
//     "z.0": 2,
//     "z.1": 1.4567,
// }
```

### Unflatten

## Options

### Delimiter

### Safe

### MaxDepth