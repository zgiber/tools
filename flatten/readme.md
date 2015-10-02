# tools
--
    import "github.com/zgiber/tools/flatten"


## Usage

#### func  Flatten

```go
func Flatten(src interface{}, prefix, separator string) (map[string]interface{}, error)
```
Flatten takes a source, prefix and separator, and produces a flat
map[string]interface{} as a result, using prefix and separator for creating 
the keys in it. 
The src can be of type []interface{}, map[string]interface{} or
valid JSON encoded string and []byte. If the source has bigger than 1 depth, the
separator is used to produce keys in the result. If the source contains a
conflicting key, it will be prefixed (again) in the result If the source is a
slice, the result will have the index of the item as the key, as a string type.
