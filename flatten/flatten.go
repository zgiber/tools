package tools

import (
	"strconv"
	"strings"
)

import (
	"encoding/json"
	"fmt"
)

type flattener struct {
	prefix    string
	separator string
	dst       map[string]interface{}
}

// Flatten takes a source, prefix and separator, and produces a flat
// map[string]interface{} from the source as a result, using prefix for
// the keys in the result. The src can be of type []interface{},
// map[string]interface{} or valid JSON encoded string and []byte.
// If the source has bigger than 1 depth, the separator is used to
// produce keys in the result. If the source contains a conflicting
// key, it will be prefixed (again) in the result
// If the source is a slice, the result will have the index of the item
// as the key, as a string type.
func Flatten(src interface{}, prefix, separator string) (map[string]interface{}, error) {

	f := &flattener{
		prefix:    prefix,
		separator: separator,
		dst:       map[string]interface{}{},
	}

	switch v := src.(type) {
	case map[string]interface{}:
		f.flattenMap(v, f.prefix)
		return f.dst, nil

	case []interface{}:
		f.flattenSlice(v, f.prefix)
		return f.dst, nil

	case []byte:
		unmarshaled := map[string]interface{}{}
		err := json.Unmarshal(v, &unmarshaled)
		if err != nil {
			return nil, err
		}

		f.flattenMap(unmarshaled, f.prefix)
		return f.dst, nil

	case string:

		unmarshaled := map[string]interface{}{}
		err := json.Unmarshal([]byte(v), &unmarshaled)
		if err != nil {
			return nil, err
		}

		f.flattenMap(unmarshaled, f.prefix)
		return f.dst, nil
	default:
		return nil, fmt.Errorf("Flatten: Unsupported type.")
	}
}

func (f *flattener) flattenSlice(src []interface{}, prefix string) {
	var currentPrefix string
	for i, elem := range src {
		if prefix != f.prefix {
			currentPrefix = strings.Join([]string{prefix, strconv.Itoa(i)}, f.separator)
		} else {
			currentPrefix = fmt.Sprintf("%s%v", prefix, i)
		}
		switch t := elem.(type) {
		case map[string]interface{}:
			f.flattenMap(t, currentPrefix)
		case []interface{}:
			f.flattenSlice(t, currentPrefix)
		default:
			f.dst[currentPrefix] = t
		}
	}
}

func (f *flattener) flattenMap(src map[string]interface{}, prefix string) {
	var currentPrefix string
	for k, v := range src {
		if prefix != f.prefix {
			currentPrefix = strings.Join([]string{prefix, k}, f.separator) // fmt.Sprintf("%s_%s", prefix, k)
		} else {
			currentPrefix = strings.Join([]string{prefix, k}, "")
		}
		switch t := v.(type) {
		case map[string]interface{}:
			f.flattenMap(t, currentPrefix)
		case []interface{}:
			f.flattenSlice(t, currentPrefix)

		default:
			f.dst[currentPrefix] = t
		}
	}
}
