# flat (BSON variant)
**forked from [nqd/flat](https://github.com/nqd/flat)**

Take a mongo bson.M and flatten it or unflatten a map with delimited key.

Forked from 
This work inspired by the [nodejs flat package](https://github.com/hughsk/flat/)

## Method

### Flatten

Flatten given map, returns a map one level deep.

```{go}
in := bson.M{
    "a": "b",
    "c": bson.M{
        "d": "e",
        "f": "g",
    },
    "z": bson.A{2, 1.4567},
}

out, err := flat.Flatten(in, nil)
// out = bson.M{
//     "a": "b",
//     "c.d": "e",
//     "c.f": "g",
//     "z.0": 2,
//     "z.1": 1.4567,
// }
```

### Unflatten

Since there is flatten, flat should have unflatten.

```{go}
in := bson.M{
    "foo.bar": bson.M{"t": 123},
    "foo":     bson.M{"k": 456},
}

out, err := flat.Unflatten(in, nil)
// out = bson.M{
//     "foo": bson.M{
//         "bar": bson.M{
//             "t": 123,
//         },
//         "k": 456,
//     },
// }
```

## Options

### Delimiter

Use a custom delimiter for flattening/unflattening your objects. Default value is `.`.

```{go}
in := bson.M{
   "hello": bson.M{
       "world": bson.M{
           "again": "good morning",
        }
    },
}

out, err := flat.Flatten(in, &flat.Options{
    Delimiter: ":",
})
// out = bson.M{
//     "hello:world:again": "good morning",
// }
```

### Safe

<!-- When Safe is true, both fatten and unflatten will preserve bson.As and their contents. Default Safe value is `false`. -->
When Safe is true, fatten will preserve arrays and their contents. Default Safe value is `false`.

```{go}
in := bson.M{
    "hello": bson.M{
        "world": bson.A{
            "one",
            "two",
        }
   },
}

out, err := flat.Flatten(in, &flat.Options{
    Delimiter: ".",
    Safe:      true,
})
// out = bson.M{
//     "hello.world": bson.A{"one", "two"},
// }
```

<!-- Example of Unflatten goes here -->

### MaxDepth

MaxDepth is the maximum number of nested objects to flatten. MaxDepth can be any integer number. MaxDepth = 0 means no limit.

Default MaxDepth value is `0`.

```{go}
in := bson.M{
    "hello": bson.M{
        "world": bson.M{
            "again": "good morning",
        }
   },
}

out, err := flat.Flatten(in, &flat.Options{
    Delimiter: ".",
    MaxDepth:  2,
})
// out = bson.M{
//     "hello.world": bson.M{"again": "good morning"},
// }
```

## Todos

- [ ] Safe option for Unflatten
- [ ] Overwrite