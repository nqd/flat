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

out, err := flat.Flatten(in, nil)
// out = map[string]interface{}{
//     "a": "b",
//     "c.d": "e",
//     "c.f": "g",
//     "z.0": 2,
//     "z.1": 1.4567,
// }
```

### Unflatten

```{go}
in := map[string]interface{}{
    "foo.bar": map[string]interface{}{"t": 123},
    "foo":     map[string]interface{}{"k": 456},
}

out, err := flat.Unflatten(in, nil)
// out = map[string]interface{}{
//     "foo": map[string]interface{}{
//         "bar": map[string]interface{}{
//             "t": 123,
//         },
//         "k": 456,
//     },
// }
```

## Options

### Delimiter

### Safe

### MaxDepth